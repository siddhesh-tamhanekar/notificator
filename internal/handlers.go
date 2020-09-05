package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"siddhesh-tamhanekar/notificator/pkg/pool"
	"sync"
	"time"
)

var notificationRequestPool *sync.Pool

// SetHandlers set application handlers
func SetHandlers(s *http.ServeMux) {
	notificationRequestPool = &sync.Pool{
		New: func() interface{} {
			return new(sendManyRequest)
		},
	}
	s.HandleFunc("/json", jsonHandler)
	s.HandleFunc("/send", sendHandler)
	s.HandleFunc("/sendMany", sendManyHandler)
	s.HandleFunc("/home", homeHandler)

}
func homeHandler(w http.ResponseWriter, req *http.Request) {
	id := rand.Intn(99999)
	job := pool.Job{
		Id:        1,
		CreatedAt: time.Now(),
		Run: func() {
			fmt.Println(id, "dummy job started")
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			fmt.Println(id, "dummy job completed")
		},
	}
	pool.GetInstance().AddJob(job)
	fmt.Fprintf(w, "Hello world")
}

func jsonResponse(w http.ResponseWriter, jsonData interface{}) {
	jsonBytes, error := json.Marshal(jsonData)
	if error != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Response is not in valid json format")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

// parseJSONRequest parses the req.body and returns struct filled.
func parseJSONRequest(r io.ReadCloser, structure interface{}) error {
	defer r.Close()
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	// json decode
	if err = json.Unmarshal(body, structure); err != nil {
		return err
	}
	return nil

}
func jsonHandler(w http.ResponseWriter, req *http.Request) {
	jsonResponse(w, pool.GetInstance().Stats())
}

type sendManyRequest struct {
	Notifications []Notification
}

func sendManyHandler(w http.ResponseWriter, req *http.Request) {

	bodyArr := notificationRequestPool.Get().(*sendManyRequest)
	if err := parseJSONRequest(req.Body, bodyArr); err != nil {
		http.Error(w, err.Error(), 500)
	}
	results := make([]map[string]string, 0)

	for _, notification := range bodyArr.Notifications {
		result := make(map[string]string)

		err := notification.Validate()
		if err == nil {
			err = notification.Send()
		}
		if err != nil {
			result["error"] = err.Error()
		} else {
			result["success"] = "job queued successfully"
		}
		results = append(results, result)

	}
	notificationRequestPool.Put(bodyArr)
	jsonResponse(w, bodyArr)
}

func sendHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1024 * 1024 * 20)

	medium := req.Form.Get("type")
	var err error
	var not *Notification
	switch medium {
	case "email":
		// initialize email sender
		to := req.Form.Get("to")
		from := req.Form.Get("to")
		body := req.Form.Get("body")
		subject := req.Form.Get("subject")
		not = &Notification{
			NotificationType: "email",
			To:               to,
			From:             from,
			Body:             body,
			Subject:          subject,
		}
	case "sms":
		message := req.Form.Get("message")
		number := req.Form.Get("number")
		smsType := req.Form.Get("sms_type")
		not = &Notification{
			NotificationType: "sms",
			Message:          message,
			Number:           number,
			SmsType:          smsType,
		}
		SetSMSNotificationJob(message, number, smsType)
	case "fcm":
		data := req.Form.Get("data")
		not = &Notification{
			NotificationType: "fcm",
			Data:             data,
		}
	}
	if not == nil {
		jsonResponse(w, map[string]string{
			"msg": fmt.Sprintf("type %s not valid", medium),
		})
	}
	if err = not.Validate(); err == nil {
		err = not.Send()
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	jsonResponse(w, map[string]string{
		"msg": "Notification Queued Successfully",
	})
}
