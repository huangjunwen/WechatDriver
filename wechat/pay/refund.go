package pay

import (
	"context"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
)

type RefundParam struct {
	PayParam

	// --- Required one of the two
	TransactionID string `wx_pay:"transaction_id"`
	OutTradeNO    string `wx_pay:"out_trade_no"`

	// --- Required
	OutRefundNO string `wx_pay:"out_refund_no"`
	TotalFee    uint32 `wx_pay:"total_fee"`
	RefundFee   uint32 `wx_pay:"refund_fee"`

	// --- Optional
	RefundFeeType string `wx_pay:"refund_fee_type"`
	RefundAccount string `wx_pay:"refund_account"`
	OpUserID      string `wx_pay:"op_user_id"`
	DeviceInfo    string `wx_pay:"device_info"`
}

type RefundResult struct {
	PayResult

	TransactionID       string `wx_pay:"transaction_id"`
	OutTradeNO          string `wx_pay:"out_trade_no"`
	OutRefundNO         string `wx_pay:"out_refund_no"`
	RefundID            string `wx_pay:"refund_id"`
	RefundFee           uint32 `wx_pay:"refund_fee"`
	SettlementRefundFee uint32 `wx_pay:"settlement_refund_fee"`
	TotalFee            uint32 `wx_pay:"total_fee"`
	SettlementTotalFee  uint32 `wx_pay:"settlement_total_fee"`
	FeeType             string `wx_pay:"fee_type"`
	CashFee             uint32 `wx_pay:"cash_fee"`
	CashFeeType         string `wx_pay:"cash_fee_type"`
	CashRefundFee       uint32 `wx_pay:"cash_refund_fee"`
	CouponRefundFee     uint32 `wx_pay:"coupon_refund_fee"`
	CouponRefundCount   uint32 `wx_pay:"coupon_refund_count"`
	Extra               map[string]string
}

func (pay *Pay) Refund(ctx context.Context, p *RefundParam, l wx.Logger) (
	r *RefundResult, err error) {

	if p.TransactionID == "" && p.OutTradeNO == "" {

		return nil, fmt.Errorf("Refund: at least one of transaction_id/out_trade_no required")

	}

	if p.OutRefundNO == "" || p.TotalFee == 0 || p.RefundFee == 0 {

		return nil, fmt.Errorf("Refund: out_refund_no/total_fee/refund_fee are required")

	}

	p.PayParam.fillFrom(pay)

	if p.OpUserID == "" {

		p.OpUserID = p.PayParam.MchID

	}

	r = &RefundResult{}

	err = pay.callPayPAI(ctx, "https://api.mch.weixin.qq.com/secapi/pay/refund", p, r, l)

	return

}
