package pay

import (
	"bytes"
	"context"
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"io"
	"strconv"
)

type OrderQueryParam struct {
	PayParam

	// --- Required one of the two
	TransactionID string     `wx_pay:"transaction_id"`
	OutTradeNO    string     `wx_pay:"out_trade_no"`
	Version       APIVersion `wx_pay:"version"`
}

type OrderQueryResult struct {
	PayResult

	DeviceInfo         string     `wx_pay:"device_info"`
	OpenID             string     `wx_pay:"openid"`
	IsSubscribe        YN         `wx_pay:"is_subscribe"`
	TradeType          TradeType  `wx_pay:"trade_type"`
	TradeState         TradeState `wx_pay:"trade_state"`
	BankType           string     `wx_pay:"bank_type"`
	TotalFee           uint32     `wx_pay:"total_fee"`
	SettlementTotalFee uint32     `wx_pay:"settlement_total_fee"` // Do not use this
	FeeType            string     `wx_pay:"fee_type"`
	CashFee            uint32     `wx_pay:"cash_fee"`
	CashFeeType        string     `wx_pay:"cash_fee_type"`
	TransactionID      string     `wx_pay:"transaction_id"`
	OutTradeNO         string     `wx_pay:"out_trade_no"`
	Attach             string     `wx_pay:"attach"`
	TimeEnd            Datetime   `wx_pay:"time_end"`
	TradeStateDesc     string     `wx_pay:"trade_state_desc"`

	// Version notes: recommend to use API_VERSION_1_0
	//
	// For API_VERSION_DEFAULT (default value), the following fields are presented:
	//   coupon_fee: the sum fee of coupons used in this order
	//   coupon_count: the number of coupons used in this order
	//   coupon_id_$n (dynamic): the id of each coupon
	//   coupon_type_$n (dynamic): the type of each coupon -> CASH/NO_CASH
	//   coupon_fee_$n (dynamic): the fee of each coupon
	// The last three has dynamic field names so they are stored in Extra
	//
	// For API_VERSION_1_0, the following fields are NOT presented:
	//   coupon_fee/coupon_count/coupon_id_$n/coupon_type_$n/coupon_fee_$n
	// Instead, these information is encapsulated in a JSON field promotion_detail. See PromotionDetailSt.
	//
	Version         APIVersion          `wx_pay:"version"`
	CouponFee       uint32              `wx_pay:"coupon_fee"`
	CouponCount     uint32              `wx_pay:"coupon_count"`
	PromotionDetail PromotionDetailInfo `wx_pay:"promotion_detail"`
	Extra           map[string]string   `wx_pay:"*"`

	UnifiedPromotionDetail *PromotionDetailInfo
}

func (r *OrderQueryResult) unifyPromotionDetail() error {

	switch r.Version {

	default:

		return fmt.Errorf("Unknown version %+q", r.Version)

	case API_VERSION_1_0:

		r.UnifiedPromotionDetail = &r.PromotionDetail

		return nil

	case API_VERSION_DEFAULT:

		p := &PromotionDetailInfo{}

		if r.CouponCount == 0 {

			r.UnifiedPromotionDetail = p

			return nil

		}

		p.Items = make([]PromotionDetailItem, 0, int(r.CouponCount))

		for i := 0; i < int(r.CouponCount); i++ {

			var ok bool

			var coupon_fee, coupon_type, coupon_id string

			if coupon_fee, ok = r.Extra[fmt.Sprintf("coupon_fee_%d", i)]; !ok {

				return fmt.Errorf("coupon_fee_%d not found", i)

			}

			amount, err := strconv.ParseUint(coupon_fee, 10, 32)

			if err != nil {

				return err

			}

			if coupon_type, ok = r.Extra[fmt.Sprintf("coupon_type_%d", i)]; !ok {

				return fmt.Errorf("coupon_type_%d not found", i)

			}

			switch coupon_type {

			case "CASH":

				coupon_type = "COUPON"

			case "NO_CASH":

				coupon_type = "DISCOUNT"

			default:

				return fmt.Errorf("Unknown coupon type %+q", coupon_type)

			}

			coupon_id, _ = r.Extra[fmt.Sprintf("coupon_id_%d", i)]

			p.Items = append(p.Items, PromotionDetailItem{
				PromotionID: coupon_id,
				Type:        coupon_type,
				Amount:      uint32(amount),
			})

		}

		r.UnifiedPromotionDetail = p

		return nil

	}

}

func (pay *Pay) OrderQuery(ctx context.Context, p *OrderQueryParam, l wx.Logger) (
	r *OrderQueryResult, err error) {

	if p.TransactionID == "" && p.OutTradeNO == "" {

		return nil, fmt.Errorf("OrderQuery: require one of transaction_id/out_trade_no")

	}

	p.PayParam.fillFrom(pay)

	r = &OrderQueryResult{}

	err = pay.callPayPAI(ctx, "https://api.mch.weixin.qq.com/pay/orderquery", p, r, l)

	if err != nil {

		return

	}

	err = r.unifyPromotionDetail()

	return

}

func (pay *Pay) PaymentNotify(r io.Reader, l wx.Logger) (*OrderQueryResult, error) {

	var body bytes.Buffer

	err := wx.LimitRead(r, &body, int64(pay.maxResultSize()))

	if err != nil {

		return nil, err

	}

	if l != nil {

		l.Printf("notify_body=%+q\n", body.Bytes())

	}

	result := &OrderQueryResult{}

	if err := pay.decodeResult(&body, result, SIGN_TYPE_MD5); err != nil {

		return nil, err

	}

	err = result.unifyPromotionDetail()

	if err != nil {

		return nil, err

	}

	return result, nil

}
