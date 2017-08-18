package pay

import (
	"context"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
)

type CloseOrderParam struct {
	PayParam

	// --- Required
	OutTradeNO string `wx_pay:"out_trade_no"`
}

type CloseOrderResult struct {
	PayResult
}

func (pay *Pay) CloseOrder(ctx context.Context, p *CloseOrderParam, l wx.Logger) (
	r *CloseOrderResult, err error) {

	if p.OutTradeNO == "" {

		return nil, fmt.Errorf("CloseOrder: out_trade_no missing")

	}

	p.PayParam.fillFrom(pay)

	r = &CloseOrderResult{}

	err = pay.callPayPAI(ctx, "https://api.mch.weixin.qq.com/pay/closeorder", p, r, l)

	return

}
