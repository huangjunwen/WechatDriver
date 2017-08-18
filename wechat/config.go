package wechat

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

// A structure containing a Wechat app's parameter configuration.
type AppConfig struct {
	// The app id. (required)
	AppID string `json:"app_id"`

	// App secret.
	AppSecret string `json:"app_secret"`

	// Wechat payment merchandiser id.
	PayMchID string `json:"pay_mch_id"`

	// Wechat payment secret key.
	PayKey string `json:"pay_key"`

	// Wechat payment tls client certificate (PEM format).
	PayClientTLSCertPEM string `json:"pay_client_tls_cert"`

	// Wechat payment tls client private key (PEM format).
	PayClientTLSKeyPEM string `json:"pay_client_tls_key"`

	// List of root CA certificates payment client should trust (PEM format).
	PayRootCAPEM []string `json:"pay_root_ca"`
}

// Return tls.Config which can be used in constructing HTTP client. Example:
//
//   client := &http.Client{
//       Transport: &http.Transport{
//           TLSClientConfig: config.PayClientTLSConfig(),
//       },
//   }
func (config *AppConfig) PayClientTLSConfig() (*tls.Config, error) {

	if len(config.PayClientTLSCertPEM) == 0 || len(config.PayClientTLSKeyPEM) == 0 {

		return nil, fmt.Errorf("missing PayClientTLSCertPEM/PayClientTLSKeyPEM")

	}

	var (
		cert tls.Certificate
		err  error
	)

	// Load client cert/key
	cert, err = tls.X509KeyPair([]byte(config.PayClientTLSCertPEM), []byte(config.PayClientTLSKeyPEM))

	if err != nil {

		return nil, err

	}

	// Load root ca pool
	pool := x509.NewCertPool()

	for _, pem := range config.PayRootCAPEM {

		if !pool.AppendCertsFromPEM([]byte(pem)) {

			return nil, fmt.Errorf("Failed to load CA cert")

		}

	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: false,
	}, nil

}

var DefaultAPITimeout time.Duration = 30 * time.Second
