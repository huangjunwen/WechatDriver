package oauth2

import (
	"context"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"net/url"
	"strings"
)

// The result of OAuth2.AccessToken and OAuth2.RefreshAccessToken
type AccessTokenResult struct {
	OAuth2ResultBase

	AccessToken  string `json:"access_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"` // "," seperated
	// NOTE: UnionID is an undocumented field.
	UnionID string `json:"unionid"`
}

// Scopes return the authorized scopes of the access token.
func (r *AccessTokenResult) Scopes() (ret []OAuth2Scope) {
	for _, scope := range strings.Split(r.Scope, ",") {
		switch scope {
		case string(OAUTH2_SCOPE_BASE):
			ret = append(ret, OAUTH2_SCOPE_BASE)
		case string(OAUTH2_SCOPE_USERINFO):
			ret = append(ret, OAUTH2_SCOPE_USERINFO)
		default:
		}
	}
	return
}

// HasScope check a scope's existence.
func (r *AccessTokenResult) HasScope(scope OAuth2Scope) bool {
	for _, scp := range r.Scopes() {
		if scp == scope {
			return true
		}
	}
	return false
}

// Get access token (and refresh token) from code. Return error only when there
// is error in transport. Caller should also check r.Error().
func (o *OAuth2) AccessToken(ctx context.Context, code string, l wx.Logger) (
	*AccessTokenResult, error) {

	r := &AccessTokenResult{}

	if err := o.callOAuth2API(ctx,
		"https://api.weixin.qq.com/sns/oauth2/access_token?"+url.Values{
			"appid":      []string{o.config.AppID},
			"secret":     []string{o.config.AppSecret},
			"code":       []string{code},
			"grant_type": []string{"authorization_code"},
		}.Encode(), r, l); err != nil {
		return nil, err
	}

	return r, nil

}

// Refresh access token (and refresh token) from refresh token. Return error only when there
// is error in transport. Caller should also check r.Error().
func (o *OAuth2) RefreshAccessToken(ctx context.Context, refresh_token string, l wx.Logger) (
	*AccessTokenResult, error) {

	r := &AccessTokenResult{}

	if err := o.callOAuth2API(ctx,
		"https://api.weixin.qq.com/sns/oauth2/refresh_token?"+url.Values{
			"appid":         []string{o.config.AppID},
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refresh_token},
		}.Encode(), r, l); err != nil {
		return nil, err
	}

	return r, nil

}
