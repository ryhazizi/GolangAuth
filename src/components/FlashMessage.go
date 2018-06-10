package components

import (
	"encoding/base64"
	"net/http"
	"time"
)

func SetFlashMessage(w http.ResponseWriter, name string, value []byte) {
	cookie := &http.Cookie{Name: name, Value: encode(value)}

	http.SetCookie(w, cookie)
}

func GetFlashMessage(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}

	value, err := decode(cookie.Value)
	if err != nil {
		return nil, err
	}

	domainCookie := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}

	http.SetCookie(w, domainCookie)
	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func CheckMessage(s string, b []byte) bool {
	if len(s) != len(b) {
		return false
	}
	for i, x := range b {
		if x != s[i] {
			return false
		}
	}
	return true
}
