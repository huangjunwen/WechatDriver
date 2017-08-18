package pay

import (
	"context"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
)

type UnifiedOrderParam struct {
	PayParam

	// --- Required
	TradeType      TradeType `wx_pay:"trade_type"`
	Body           string    `wx_pay:"body"`
	OutTradeNO     string    `wx_pay:"out_trade_no"`
	TotalFee       uint32    `wx_pay:"total_fee"`
	SpbillCreateIP string    `wx_pay:"spbill_create_ip"`
	NotifyURL      string    `wx_pay:"notify_url"`

	// --- Required in some cases
	ProductID string `wx_pay:"product_id"`
	OpenID    string `wx_pay:"openid"`

	// --- Optional
	DeviceInfo string   `wx_pay:"device_info"`
	Detail     string   `wx_pay:"detail"`
	Attach     string   `wx_pay:"attach"`
	FeeType    string   `wx_pay:"fee_type"`
	TimeStart  Datetime `wx_pay:"time_start"`
	TimeExpire Datetime `wx_pay:"time_expire"`
	GoodsTag   string   `wx_pay:"goods_tag"`
	LimitPay   string   `wx_pay:"limit_pay"`
}

type UnifiedOrderResult struct {
	PayResult

	TradeType TradeType `wx_pay:"trade_type"`
	PrepayID  string    `wx_pay:"prepay_id"`
	CodeURL   string    `wx_pay:"code_url"`
}

func (pay *Pay) UnifiedOrder(ctx context.Context, p *UnifiedOrderParam, l wx.Logger) (
	r *UnifiedOrderResult, err error) {

	if p.TradeType == TRADE_TYPE_JSAPI && p.OpenID == "" {

		return nil, fmt.Errorf("UnifiedOrder: OpenID is required when trade type is JSAPI")

	}

	if p.TradeType == TRADE_TYPE_NATIVE && p.ProductID == "" {

		return nil, fmt.Errorf("UnifiedOrder: ProductID is required when trade type is NATIVE")

	}

	p.PayParam.fillFrom(pay)

	r = &UnifiedOrderResult{}

	err = pay.callPayPAI(ctx, "https://api.mch.weixin.qq.com/pay/unifiedorder", p, r, l)

	return

}
