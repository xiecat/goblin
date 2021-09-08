package replace

import (
	"net/http"

	"goblin/pkg/utils"
)

func (rqRule *Request) Request(r *http.Request) error {
	if rqRule == nil {
		return nil
	}

	if !utils.EleInArray(r.Method, rqRule.Method) {
		return nil
	}
	if rqRule.Header != nil {
		for hkey, hvalue := range rqRule.Header {
			if hvalue == "" {
				r.Header.Del(hkey)
			} else {
				r.Header.Set(hkey, hvalue)
			}
		}
	}
	return nil
}
