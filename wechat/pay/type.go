package pay

import (
	"encoding/json"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"time"
)

// Common part of Pay parameters.
type PayParam struct {
	AppID    string `wx_pay:"appid"`
	MchID    string `wx_pay:"mch_id"`
	NonceStr string `wx_pay:"nonce_str"`
}

func (param *PayParam) fillFrom(pay *Pay) {

	param.AppID = pay.config.AppID

	param.MchID = pay.config.PayMchID

	param.NonceStr = wx.HexCryptoRandString(32)

}

// Common part of Pay results.
type PayResult struct {
	ReturnCode string `wx_pay:"return_code"`
	ReturnMsg  string `wx_pay:"return_msg"`
	ResultCode string `wx_pay:"result_code"`
	ResultMsg  string `wx_pay:"result_msg"`
	ErrCode    string `wx_pay:"err_code"`
	ErrCodeDes string `wx_pay:"err_code_des"`
	Sign       string `wx_pay:"sign"`
	AppID      string `wx_pay:"appid"`
	MchID      string `wx_pay:"mch_id"`
	NonceStr   string `wx_pay:"nonce_str"`
}

type APIVersion string

const (
	API_VERSION_DEFAULT APIVersion = ""
	API_VERSION_1_0     APIVersion = "1.0"
)

func (v *APIVersion) marshal() (string, error) {

	return string(*v), nil

}

func (v *APIVersion) unmarshal(s string) error {

	r := APIVersion(s)

	switch r {

	case API_VERSION_DEFAULT, API_VERSION_1_0:

		*v = r

		return nil

	default:

		return fmt.Errorf("Unknown version %+q", s)

	}

}

type SignType string

const (
	SIGN_TYPE_MD5         SignType = "MD5"
	SIGN_TYPE_HMAC_SHA256 SignType = "HMAC-SHA256"
)

type TradeType string

const (
	TRADE_TYPE_JSAPI  TradeType = "JSAPI"
	TRADE_TYPE_NATIVE TradeType = "NATIVE"
	TRADE_TYPE_APP    TradeType = "APP"
)

func (tt *TradeType) marshal() (string, error) {

	return string(*tt), nil

}

func (tt *TradeType) unmarshal(s string) error {

	v := TradeType(s)

	switch v {

	case TRADE_TYPE_JSAPI, TRADE_TYPE_NATIVE, TRADE_TYPE_APP:

		*tt = v

		return nil

	default:

		return fmt.Errorf("Unknown trade type %+q", s)

	}

}

type TradeState string

const (
	TRADE_STATE_SUCCESS    TradeState = "SUCCESS"
	TRADE_STATE_REFUND     TradeState = "REFUND"
	TRADE_STATE_NOTPAY     TradeState = "NOTPAY"
	TRADE_STATE_CLOSED     TradeState = "CLOSED"
	TRADE_STATE_REVOKED    TradeState = "REVOKED"
	TRADE_STATE_USERPAYING TradeState = "USERPAYING"
	TRADE_STATE_PAYERROR   TradeState = "PAYERROR"
)

func (ts *TradeState) marshal() (string, error) {

	return string(*ts), nil

}

func (ts *TradeState) unmarshal(s string) error {

	v := TradeState(s)

	switch v {

	case TRADE_STATE_SUCCESS, TRADE_STATE_REFUND, TRADE_STATE_NOTPAY, TRADE_STATE_CLOSED,
		TRADE_STATE_REVOKED, TRADE_STATE_USERPAYING, TRADE_STATE_PAYERROR:

		*ts = v

		return nil

	default:

		return fmt.Errorf("Unknown trade state %+q", s)

	}

}

// Datetime of format "yyyymmddHHMMSS".
type Datetime time.Time

const datetimeFmt string = "20060102150405"

func (dt *Datetime) marshal() (string, error) {

	return (*time.Time)(dt).Format(datetimeFmt), nil

}

func (dt *Datetime) unmarshal(s string) error {

	t, err := time.Parse(datetimeFmt, s)

	if err != nil {

		return err

	}

	*dt = Datetime(t)

	return nil

}

// "Y" (yes) or "N" (no)
type YN bool

func (yn *YN) marshal() (string, error) {

	if bool(*yn) {

		return "Y", nil

	} else {

		return "N", nil

	}

}

func (yn *YN) unmarshal(s string) error {

	switch s {

	default:

		return fmt.Errorf("Unknow YN value %+q", s)

	case "Y":

		*yn = YN(true)

		return nil

	case "N":

		*yn = YN(false)

		return nil

	}

}

type PromotionDetailInfo struct {
	Items []PromotionDetailItem `json:"promotion_detail"`
}

type PromotionDetailItem struct {
	ActivityID         string            `json:"activity_id"`         // 微信商户后台配置的批次 ID
	PromotionID        string            `json:"promotion_id"`        // 券或者立减优惠 ID
	Name               string            `json:"name"`                // 优惠名称
	Scope              string            `json:"scope"`               // GLOBAL- 全场代金券/SINGLE- 单品优惠
	Type               string            `json:"type"`                // COUPON- 预充代金券（走结算资金）/DISCOUNT- 免充优惠券
	Amount             uint32            `json:"amount"`              // 金额 = 微信出资金额 + 商家出资金额 + 其他出资方金额
	WXPayContribute    uint32            `json:"wxpay_contribute"`    // 微信出资金额
	MerchantContribute uint32            `json:"merchant_contribute"` // 商家出资金额
	OtherContribute    uint32            `json:"other_contribute"`    // 其他出资方金额
	GoodsDetail        []GoodsDetailItem // 单品信息
}

type GoodsDetailItem struct {
	GoodsID        string `json:"goods_id"`        // 商品编码
	Quantity       string `json:"quantity"`        // 商品数量
	Price          uint32 `json:"price"`           // 商品单价
	DiscountAmount uint32 `json:"discount_amount"` // 商品优惠金额
}

func (p *PromotionDetailInfo) marshal() (string, error) {

	if r, err := json.Marshal(p); err != nil {

		return "", err

	} else {

		return string(r), nil

	}

}

func (p *PromotionDetailInfo) unmarshal(s string) (err error) {

	err = json.Unmarshal([]byte(s), p)

	return

}
