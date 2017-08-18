package oauth2

import (
	"fmt"
)

type OAuth2Scope string

const (
	OAUTH2_SCOPE_BASE     OAuth2Scope = "snsapi_base"
	OAUTH2_SCOPE_USERINFO OAuth2Scope = "snsapi_userinfo"
	//OAUTH2_SCOPE_LOGIN    OAuth2Scope = "snsapi_login"
)

// Common part of OAuth2 API result.
type OAuth2ResultBase struct {
	// ErrCode default to 0, which is good since no errcode means success
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Is ErrCode == 0 ?
func (r *OAuth2ResultBase) OK() bool {

	return r.ErrCode == 0

}

// Return error if not OK().
func (r *OAuth2ResultBase) Error() error {

	if r.OK() {

		return nil

	}

	return fmt.Errorf("errcode=%v errmsg=%+q", r.ErrCode, r.ErrMsg)

}
