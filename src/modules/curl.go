package modules

import (
	"net/http"
	"time"
)

func IgCurl(endPoint string, authInfo AuthInfo, method string) (error, *http.Response) {
	client := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(method, endPoint, nil)
	if err != nil {
		return err, nil
	}
	req.Header.Set("x-ig-app-id", authInfo.AppId)
	for _, cookie := range authInfo.CookieMap {
		netCookie := http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Path:    cookie.Path,
			Domain:  cookie.Domain,
			Expires: time.Unix(int64(cookie.Expiry), 0),
			Secure:  cookie.Secure,
		}
		req.AddCookie(&netCookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	return nil, resp
}
