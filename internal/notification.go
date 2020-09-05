package internal

import (
	"errors"
	"log"
	"siddhesh-tamhanekar/notificator"
	"siddhesh-tamhanekar/notificator/internal/notification"
	"siddhesh-tamhanekar/notificator/pkg/pool"
	"time"
)

type Notification struct {
	NotificationType                               string `json:"type"`
	SmsType                                        string `json:"sms_type"`
	Message, Number, From, To, Subject, Body, Data string
}

var emailSender *notification.SMTPEmailSender
var smsSender *notification.SMSSender
var fcmSender *notification.FCMSender

func (not *Notification) Validate() error {
	switch not.NotificationType {
	case "sms":
		if not.Number == "" || not.Message == "" {
			return errors.New("number and message are required parameters")
		}
		break
	case "email":
		if not.To == "" || not.From == "" || not.Subject == "" || not.Body == "" {
			return errors.New("to,from,body,subject are required parameters")
		}
		break
	case "fcm":
		if not.Data == "" {
			return errors.New("Data is required parameter")
		}
	}
	return nil
}
func (not *Notification) Send() error {
	var err error
	switch not.NotificationType {
	case "sms":
		err = SetSMSNotificationJob(not.Message, not.Number, not.SmsType)
		break
	case "email":
		err = SetEmailNotificationJob(not.To, not.From, not.Body, not.Subject)
		break
	case "fcm":
		err = SetFCMNotificationJob(not.Data)

	}
	return err
}

// SetEmailNotificationJob send email notification
func SetEmailNotificationJob(to, from, body, subject string) error {
	if to == "" || from == "" || body == "" || subject == "" {
		return errors.New("to,from,body,subject for email medium parameters are required")
	}
	if emailSender == nil {

		emailSender = &notification.SMTPEmailSender{
			Host:     notificator.ConfigInstance().Email.Host,
			Port:     notificator.ConfigInstance().Email.Port,
			Username: notificator.ConfigInstance().Email.Username,
			Password: notificator.ConfigInstance().Email.Password,
		}
	}

	job := pool.Job{
		Id:        1,
		CreatedAt: time.Now(),
		Run: func() {
			_, err := emailSender.Send(to, from, subject, body)
			if err != nil {
				log.Println("message sending failed err=" + err.Error())
			} else {
				log.Println("message sent successfully")
			}
		},
	}
	// job.Run()
	if err := pool.GetInstance().AddJob(job); err != nil {
		return err
	}
	return nil

}

//send sms notification
func SetSMSNotificationJob(message, number, smsType string) error {
	if smsSender == nil {

		smsSender = &notification.SMSSender{
			Username: notificator.ConfigInstance().Sms.Username,
			Pwd:      notificator.ConfigInstance().Sms.Pwd,
			ApiKey:   notificator.ConfigInstance().Sms.ApiKey,
			Sender:   notificator.ConfigInstance().Sms.Sender,
		}
	}
	if message == "" || number == "" {
		return errors.New("message,number are required parameters")
	}

	job := pool.Job{
		Id: 1,
		Run: func() {
			if smsType == "otp" {
				smsSender.SendOTPMessage(message, number)
			} else {
				smsSender.SendTransactionMessage(message, number)
			}

			log.Println("message sent successfully")

		},
	}
	if err := pool.GetInstance().AddJob(job); err != nil {
		return err
	}
	return nil
}

func SetFCMNotificationJob(data string) error {
	if fcmSender == nil {
		fcmSender = &notification.FCMSender{
			ApiKey: notificator.ConfigInstance().Fcm.ApiKey,
		}
	}
	job := pool.Job{
		Id: 1,
		Run: func() {
			fcmSender.Send(data)
			log.Println("notificaiton sent successfully")
		},
	}
	if err := pool.GetInstance().AddJob(job); err != nil {
		return err
	}
	return nil
}
