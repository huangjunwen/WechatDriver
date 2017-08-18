package pay

import (
	"fmt"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"net/http"
)

// Communicate to Wechat's payment service.
type Pay struct {
	// The app configuration.
	config *wx.AppConfig

	// The client to do http API call.
	client *http.Client

	// MD5 or HMAC-SHA256.
	DefaultSignType SignType

	// Max size when reading incoming result. Default to 8k.
	MaxResultSize int
}

// Create Pay instance from app config and optinal a HTTP client. NOTE:
// some Pay APIs need tls client cert verification, if you want to use a custom
// client, remember to set its tls config, see: AppConfig.PayClientTLSConfig
func NewPay(config *wx.AppConfig, client *http.Client) (*Pay, error) {

	if config == nil || config.AppID == "" {

		return nil, fmt.Errorf("config's AppID missing")

	}

	tls_config, err := config.PayClientTLSConfig()

	if err != nil {

		return nil, err

	}

	if client == nil {

		// Make a client using config.PayClientTLSConfig
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tls_config,
			},
			Timeout: wx.DefaultAPITimeout,
		}

	}

	return &Pay{
		config:          config,
		client:          client,
		DefaultSignType: SIGN_TYPE_MD5,
	}, nil

}
