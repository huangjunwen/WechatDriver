package oauth2

import (
	"context"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"net/url"
)

// The result of OAuth2.UserInfo
type UserInfoResult struct {
	OAuth2ResultBase

	OpenID   string `json:"openid"`
	UnionID  string `json:"unionid"`
	Nickname string `json:"nickname"`
	// ....sex...
	Sex        int      `json:"sex"`
	City       string   `json:"city"`
	Province   string   `json:"province"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
}

// Get user information from access_token. Return error only when there
// is error in transport. Caller should also check r.Error().
func (o *OAuth2) UserInfo(ctx context.Context, access_token, openid, lang string, l wx.Logger) (
	*UserInfoResult, error) {

	r := &UserInfoResult{}

	if err := o.callOAuth2API(ctx,
		"https://api.weixin.qq.com/sns/userinfo?"+url.Values{
			"access_token": []string{access_token},
			"openid":       []string{openid},
			"lang":         []string{lang},
		}.Encode(), r, l); err != nil {
		return nil, err
	}
	return r, nil

}
