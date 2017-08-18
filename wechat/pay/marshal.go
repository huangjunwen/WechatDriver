package pay

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
)

// marshaler is used to encode/decode object into/from string.
type marshaler interface {
	marshal() (string, error)
	unmarshal(string) error
}

// Encoding object into string. r can be pointer to intXX/floatXX/string/[]byte
// or any marshaler.
func marshal(r interface{}) (string, error) {

	switch v := r.(type) {

	default:

		m, ok := v.(marshaler)

		if !ok {

			return "", fmt.Errorf("Unable to marshal %T", v)

		}

		return m.marshal()

	case *int, *uint, *int8, *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64,
		*float32, *float64:

		return fmt.Sprint(reflect.ValueOf(v).Elem()), nil

	case *string:

		return *v, nil

	case *[]byte:

		return string(*v), nil

	}

}

// Decode object from string. r can be pointer to intXX/floatXX/string/[]byte
// or any marshaler.
func unmarshal(s string, r interface{}) error {

	switch v := r.(type) {

	default:

		m, ok := v.(marshaler)

		if !ok {

			return fmt.Errorf("Unable to unmarshal %T", v)

		}

		return m.unmarshal(s)

	case *int, *uint, *int8, *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64,
		*float32, *float64:

		_, err := fmt.Sscan(s, v)

		return err

	case *string:

		*v = s

		return nil

	case *[]byte:

		*v = []byte(s)

		return nil

	}

}

// typeInfo stores type(struct) fields information for converting from/to dict.
type typeInfo struct {
	t           reflect.Type
	fields_info []fieldInfo
}

type fieldInfo struct {
	field_index []int  // struct field index (used by FieldByIndex)
	key         string // which key this field is mapped to/from dict or "*"
}

var dictType reflect.Type = reflect.TypeOf(map[string]string{})

func newTypeInfo(T reflect.Type) *typeInfo {

	ret := &typeInfo{
		t:           T,
		fields_info: nil,
	}

	var traverse func(reflect.Type, []int) error

	traverse = func(t reflect.Type, parent_index []int) error {

		if t.Kind() != reflect.Struct {

			return nil

		}

		for i := 0; i < t.NumField(); i++ {

			index := parent_index[:]

			index = append(index, i)

			f := t.Field(i)

			tag := f.Tag.Get("wx_pay")

			if tag == "" {

				if f.Type.Kind() == reflect.Struct {

					if err := traverse(f.Type, index); err != nil {

						return err

					}

				}

				continue

			}

			if tag == "*" {

				if f.Type != dictType {

					return fmt.Errorf(
						"Expect map[string]string for `wx_pay:\"*\"` but got %v", f.Type.String())

				}

			}

			ret.fields_info = append(ret.fields_info, fieldInfo{
				field_index: index,
				key:         tag,
			})

		}

		return nil

	}

	if err := traverse(T, []int{}); err != nil {

		panic(err)

	}

	return ret

}

// Cache.
var typeInfos map[reflect.Type]*typeInfo = make(map[reflect.Type]*typeInfo)

func getTypeInfo(v reflect.Value) *typeInfo {

	t := v.Type()

	if ret, ok := typeInfos[t]; ok {

		return ret

	}

	ret := newTypeInfo(t)

	typeInfos[t] = ret

	return ret

}

// Convert (ptr to) struct to dict
func structToDict(val interface{}) (map[string]string, error) {

	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {

		return nil, fmt.Errorf("Expect ptr to struct but got %T", val)

	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {

		return nil, fmt.Errorf("Expect ptr to struct but got %T", val)

	}

	ret := make(map[string]string)

	for _, field_info := range getTypeInfo(v).fields_info {

		f := v.FieldByIndex(field_info.field_index)

		// Skip zero fields
		if reflect.Zero(f.Type()).Interface() == f.Interface() {

			continue

		}

		if field_info.key == "*" {

			for k, v := range *f.Addr().Interface().(*map[string]string) {

				ret[k] = v

			}

			continue

		}

		s, err := marshal(f.Addr().Interface())

		if err != nil {

			return nil, err

		}

		ret[field_info.key] = s

	}

	return ret, nil

}

// Reverse operation of structToDict.
func structFromDict(dict map[string]string, val interface{}) error {

	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {

		return fmt.Errorf("Expect ptr to struct but got %T", val)

	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {

		return fmt.Errorf("Expect ptr to struct but got %T", val)

	}

	for _, field_info := range getTypeInfo(v).fields_info {

		if field_info.key == "*" {

			f := v.FieldByIndex(field_info.field_index).Addr().Interface().(*map[string]string)

			*f = dict

			continue

		}

		s, ok := dict[field_info.key]

		// No value, skipped.
		if !ok {

			continue

		}

		err := unmarshal(s, v.FieldByIndex(field_info.field_index).Addr().Interface())

		if err != nil {

			return err

		}

		delete(dict, field_info.key)

	}

	return nil

}

// Pay XML acts like a dict, it's simply a one level xml:
//   <xml>
//     <field1>value1</field1>
//     <field2>value2</field2>
//     ...
//   </xml>
type payXML struct {
	XMLName xml.Name      `xml:"xml"`
	Fields  []payXMLField `xml:",any"`
}

type payXMLField struct {
	XMLName xml.Name
	Text    string `xml:",chardata"`
}

// <field1>value1</field1> -> map[string]string{"field1": "value1"}
func (px *payXML) ToDict() map[string]string {

	ret := make(map[string]string, len(px.Fields))

	for _, f := range px.Fields {

		ret[f.XMLName.Local] = f.Text

	}

	return ret

}

// Reverse operation of ToDict.
func (px *payXML) FromDict(m map[string]string) {

	px.Fields = make([]payXMLField, 0, len(m))

	for k, v := range m {

		px.Fields = append(px.Fields, payXMLField{
			XMLName: xml.Name{
				Local: k,
			},
			Text: v,
		})

	}

}

// Decode from XML.
func (px *payXML) Decode(r io.Reader) error {

	return xml.NewDecoder(r).Decode(px)

}

// Encode to XML.
func (px *payXML) Encode() (*bytes.Buffer, error) {

	buf := new(bytes.Buffer)

	if err := xml.NewEncoder(buf).Encode(px); err != nil {

		return nil, err

	}

	return buf, nil

}
