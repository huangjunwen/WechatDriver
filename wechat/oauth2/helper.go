package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	wx "github.com/huangjunwen/WechatDriver/wechat"
	"net/http"
)

func (o *OAuth2) maxResultSize() int {

	max_result_size := o.MaxResultSize

	if max_result_size <= 0 {

		max_result_size = 4 * 1024

	}

	return max_result_size

}

// Prepare oauth2 GET request
func (o *OAuth2) prepareOAuth2GetRequest(URL string, l wx.Logger) (*http.Request, error) {

	if l != nil {

		l.Printf("method=\"GET\" url=%+q\n", URL)

	}

	return http.NewRequest("GET", URL, nil)

}

func (o *OAuth2) parseOAuth2Response(result interface{}, resp *http.Response, l wx.Logger) (err error) {

	var body bytes.Buffer

	if err = wx.LimitRead(resp.Body, &body, int64(o.maxResultSize())); err != nil {

		return

	}

	if l != nil {

		l.Printf("status=%+q proto=%+q body=%+q\n", resp.Status, resp.Proto, body.Bytes())

	}

	return json.NewDecoder(&body).Decode(result)

}

func (o *OAuth2) callOAuth2API(ctx context.Context, URL string, result interface{}, l wx.Logger) (err error) {

	var (
		req  *http.Request
		resp *http.Response
	)

	if req, err = o.prepareOAuth2GetRequest(URL, l); err != nil {

		return

	}

	req = req.WithContext(ctx)

	if resp, err = o.client.Do(req); err != nil {

		return

	}

	if err = o.parseOAuth2Response(result, resp, l); err != nil {

		return

	}

	return

}
