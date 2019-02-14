package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type Config struct {
	Filename    string `toml:"filename"`
	Password    string `toml:"password"`
	DeviceToken string `toml:"device_token"`
	Topic       string `toml:"topic"`
	Mode        string `toml:"mode"`
}

func main() {
	var (
		message = flag.String("message", "Hello!", "")
	)
	flag.Parse()

	var c Config
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}

	cert, err := certificate.FromP12File(c.Filename, c.Password)
	if err != nil {
		log.Fatal(err)
	}

	var client *apns2.Client
	switch c.Mode {
	case "production":
		client = apns2.NewClient(cert).Production()
	default:
		client = apns2.NewClient(cert).Development()
	}

	res, err := client.Push(&apns2.Notification{
		CollapseID:  "",
		DeviceToken: c.DeviceToken,
		Topic:       c.Topic,
		Payload:     `{"aps":{"alert":"` + *message + `"}}`,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
