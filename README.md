# Notificator
This notificator is written in golang which uses worker pool implementation.
You can send email,sms and fcm push notification using this server.
- SMS providers are vary region to region so please change sms implementation as per your requirement.
- Email are sent using smtp.

# Installaction
- Copy config.yaml.example to config.yaml
- from terminal go to root directory
- build the project `go build ./cmd/server`
- run the project using `./serve` 


