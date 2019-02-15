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

type Client struct {
	apnsClient *apns2.Client
	topic      string
}

func NewClient(certPath, password, topic, mode string) (*Client, error) {
	cert, err := certificate.FromP12File(certPath, password)
	if err != nil {
		return nil, err
	}
	var client *apns2.Client
	switch mode {
	case "production":
		client = apns2.NewClient(cert).Production()
	default:
		client = apns2.NewClient(cert).Development()
	}
	return &Client{
		apnsClient: client,
		topic:      topic,
	}, err
}

func (c *Client) Push(message, deviceToken string) (*apns2.Response, error) {
	payload := fmt.Sprintf(`{"aps":{"alert":"%s"}}`, message)
	return c.apnsClient.Push(&apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       c.topic,
		Payload:     payload,
	})
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

	client, err := NewClient(c.Filename, c.Password, c.Topic, c.Mode)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Push(*message, c.DeviceToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %v %v\n", resp.StatusCode, resp.ApnsID, resp.Reason)
}
