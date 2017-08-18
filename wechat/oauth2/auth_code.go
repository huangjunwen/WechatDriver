package oauth2

import (
	"fmt"
	"net/url"
)

// Generating a URL to start Oauth2 process in Wechat's builtin
// browser.
//
//   param scope: snsapi_base/snsapi_userinfo
//   param callback_url: the url to jump back after authorization
//   param csrf: CSRF token
//   return: URL to redirect
func (o *OAuth2) AuthCodeURL(scope OAuth2Scope, callback_url string, csrf string) string {

	switch scope {

	default:

		panic(fmt.Errorf("AuthURL: Not supported scope %v", scope))

	case OAUTH2_SCOPE_BASE, OAUTH2_SCOPE_USERINFO:

		break

	}

	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize"+
		"?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		url.QueryEscape(o.config.AppID),
		url.QueryEscape(callback_url),
		string(scope),
		url.QueryEscape(csrf),
	)

}

// Generating a URL to start Oauth2 process using QRcode in PC browser.
//
//   param callback_url: the url to jump back after authorization
//   param csrf: CSRF token
//   return: URL to redirect
func (o *OAuth2) AuthCodeURLQR(callback_url string, csrf string) string {

	return fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect"+
		"?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		url.QueryEscape(o.config.AppID),
		url.QueryEscape(callback_url),
		"snsapi_login",
		url.QueryEscape(csrf),
	)

}
