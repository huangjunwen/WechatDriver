package oauth2

import (
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"net/http"
)

// Communicate to Wechat's OAuth2 service.
type OAuth2 struct {
	// The app configuration.
	config *wx.AppConfig

	// The client to do http API call, can be nil.
	client *http.Client

	// Max size when reading incoming result. Default to 4k.
	MaxResultSize int
}

// Create OAuth2 instance from app config (and optional a HTTP client). The config
// should contain valid AppID and AppSecret.
func NewOAuth2(config *wx.AppConfig, client *http.Client) (*OAuth2, error) {

	if config == nil || config.AppID == "" || config.AppSecret == "" {

		return nil, fmt.Errorf("NewOAuth2: config's AppID/AppSecret missing")

	}

	if client == nil {

		client = http.DefaultClient

	}

	return &OAuth2{
		config: config,
		client: client,
	}, nil

}
