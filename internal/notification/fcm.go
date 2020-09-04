package notification

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FCMSender struct {
	ApiKey string
}

var sender *FCMSender

func (f *FCMSender) Send(data string) error {

	url := "https://fcm.googleapis.com/fcm/send"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "key="+f.ApiKey)

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respbytes))
	return nil
}
