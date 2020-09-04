package notification

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SMSSender struct {
	Username, Pwd, ApiKey, Sender string
}

// SendOTPMessage
func (s *SMSSender) SendOTPMessage(message, number string) error {
	message = url.QueryEscape(message)
	url := fmt.Sprintf("https://onlinecloudbox.com/v2/sms/sthapi.php?login=%s&pword=%s&api_key=%s&msg=%s&sender=%s&mobnum=%s&route_id=8", s.Username, s.Pwd, s.ApiKey, message, s.Sender, number)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Printf(string(resp))
	res.Body.Close()
	return nil

}

// SendTransactionMessage used to send transactional message
func (s *SMSSender) SendTransactionMessage(message, number string) error {
	message = url.QueryEscape(message)

	url := fmt.Sprintf("https://onlinecloudbox.com/v2/sms/stapi.php?login=%s&pword=%s&api_key=%s&msg=%s&sender=%s&mobnum=%s&route_id=3", s.Username, s.Pwd, s.ApiKey, message, s.Sender, number)
	_, err := http.Get(url)
	if err != nil {
		return err
	}
	return nil
}
