package notificator

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Host string `yaml:host`
	}

	Pool struct {
		MaxWorkers         int `yaml:"MaxWorkers"`
		IdleWorkers        int `yaml:"IdleWorkers"`
		JobQueueCapacity   int `yaml:"JobQueueCapacity"`
		WorkerIdleTimeSecs int `yaml:"WorkerIdleTimeSecs"`
	}

	Sms struct {
		Username string `yaml:"Username"`
		Pwd      string `yaml:"Pwd"`
		ApiKey   string `yaml:"ApiKey"`
		Sender   string `yaml:"Sender"`
	}

	Fcm struct {
		ApiKey string `yaml:"ApiKey"`
	}

	Email struct {
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	}
}

var cfg *config

func ConfigInstance() *config {
	if cfg == nil {
		cfg = &config{}
		cfg.parse()
	}
	return cfg
}
func (c *config) parse() {

	bytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("There is no config file found")
	}
	// fmt.Println(string(bytes))
	yaml.Unmarshal(bytes, c)
}
