package pay

import (
	"bytes"
	"context"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"hash"
	"io"
	"net/http"
	"strings"
)

func (pay *Pay) maxResultSize() int {

	max_result_size := pay.MaxResultSize

	if max_result_size <= 0 {

		max_result_size = 8 * 1024

	}

	return max_result_size

}

func (pay *Pay) normalizeSignType(sign_type SignType) SignType {

	for i := 0; i < 2; i++ {

		switch sign_type {

		case SIGN_TYPE_MD5, SIGN_TYPE_HMAC_SHA256:

			return sign_type

		default:

			sign_type = pay.DefaultSignType

		}

	}

	return SIGN_TYPE_MD5

}

func (pay *Pay) signFunction(sign_type SignType) func(map[string]string) string {

	var new_hash func() hash.Hash

	key := pay.config.PayKey

	switch pay.normalizeSignType(sign_type) {

	default:

		panic(fmt.Errorf("Unknown sign type %v", string(sign_type)))

	case SIGN_TYPE_MD5:

		new_hash = newMD5()

	case SIGN_TYPE_HMAC_SHA256:

		new_hash = newHMACSHA256(key)

	}

	return func(dict map[string]string) string {

		return signDict(dict, new_hash, key)

	}

}

// Low level method to sign and encode pay parameters (ptr to struct) into buffer.
func (pay *Pay) encodeParam(param interface{}, sign_type SignType) (buf *bytes.Buffer, err error) {

	buf = nil

	// struct -> dict.
	var dict map[string]string

	if dict, err = structToDict(param); err != nil {

		return

	}

	// Ensure no "sign_type"/"sign" then sign the dict.
	delete(dict, "sign_type")

	delete(dict, "sign")

	sign_type = pay.normalizeSignType(sign_type)

	dict["sign_type"] = string(sign_type)

	dict["sign"] = pay.signFunction(sign_type)(dict)

	// dict -> payXML
	pay_xml := &payXML{}

	pay_xml.FromDict(dict)

	// payXML -> io.Reader.
	buf, err = pay_xml.Encode()

	return

}

// Low level method to decode and verifiy pay result into ptr to struct.
func (pay *Pay) decodeResult(r io.Reader, result interface{}, sign_type SignType) (err error) {

	// io.Reader -> payXML.
	pay_xml := &payXML{}

	if err = pay_xml.Decode(r); err != nil {

		return

	}

	// payXML -> dict.
	dict := pay_xml.ToDict()

	// Verify return/result code and sign.
	sign, has_sign := dict["sign"]

	delete(dict, "sign")

	sign_verified := false

	if sign != "" && strings.ToLower(sign) == pay.signFunction(sign_type)(dict) {

		sign_verified = true

	}

	if has_sign {

		dict["sign"] = sign

	}

	return_code, _ := dict["return_code"]

	return_msg, _ := dict["return_msg"]

	result_code, _ := dict["result_code"]

	err_code, _ := dict["err_code"]

	err_code_des, _ := dict["err_code_des"]

	if return_code != "SUCCESS" || result_code != "SUCCESS" || !sign_verified {

		err = fmt.Errorf(
			"return_code=%+q return_msg=%+q result_code=%+q err_code=%+q err_code_des=%+q sign_verified=%v",
			return_code, return_msg, result_code, err_code, err_code_des, sign_verified,
		)

		return

	}

	// dict -> struct.
	err = structFromDict(dict, result)

	return

}

func (pay *Pay) preparePayRequest(param interface{}, URL string, sign_type SignType, l wx.Logger) (
	req *http.Request, err error) {

	var body *bytes.Buffer

	if body, err = pay.encodeParam(param, sign_type); err != nil {

		return

	}

	if l != nil {

		l.Printf("method=\"POST\" url=%+q body=%+q\n", URL, body.Bytes())

	}

	req, err = http.NewRequest("POST", URL, body)

	return

}

func (pay *Pay) parsePayResponse(result interface{}, resp *http.Response, sign_type SignType,
	l wx.Logger) (err error) {

	var body bytes.Buffer

	if err = wx.LimitRead(resp.Body, &body, int64(pay.maxResultSize())); err != nil {

		return

	}

	if l != nil {

		l.Printf("status=%+q proto=%+q body=%+q\n", resp.Status, resp.Proto, body.Bytes())

	}

	err = pay.decodeResult(&body, result, sign_type)

	return

}

func (pay *Pay) callPayPAI(ctx context.Context, URL string, param interface{},
	result interface{}, l wx.Logger) (err error) {

	var (
		req  *http.Request
		resp *http.Response
	)

	sign_type := pay.normalizeSignType(pay.DefaultSignType)

	if req, err = pay.preparePayRequest(param, URL, sign_type, l); err != nil {

		return

	}

	req = req.WithContext(ctx)

	if resp, err = pay.client.Do(req); err != nil {

		return

	}

	if err = pay.parsePayResponse(result, resp, sign_type, l); err != nil {

		return

	}

	return

}
