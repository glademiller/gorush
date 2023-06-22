package web

import (
	"encoding/json"
	"errors"
	"fmt"
	webpush "github.com/SherClockHolmes/webpush-go"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	httpClientTimeout = 30 * time.Second
)

// Client type
type Client struct {
	HTTPClient      *http.Client
	vapidPublicKey  string
	vapidPrivateKey string
}

// Response type
type Response struct {
	StatusCode int
	Body       string
}

// // Subscription type
type Subscription struct {
	Endpoint string `json:"endpoint,omitempty"`
	Key      string `json:"key,omitempty"`
	Auth     string `json:"auth,omitempty"`
}

// Notification type
type Notification struct {
	Subscription *Subscription           `json:"subscription,omitempty"`
	Payload      *map[string]interface{} `json:"payload,omitempty"`
	TimeToLive   *uint                   `json:"time_to_live,omitempty"`
}

// Browser type
type Browser struct {
	Name     string
	ReDetect regexp.Regexp
	ReError  regexp.Regexp
}

// Browsers available
var Browsers = [...]Browser{
	{"Chrome", *regexp.MustCompile("https://android.googleapis.com/gcm/send/"), *regexp.MustCompile("<TITLE>(.*)</TITLE>")},
	{"Firefox", *regexp.MustCompile("https://updates.push.services.mozilla.com/wpush"), *regexp.MustCompile("\\\"errno\\\":\\s(\\d+)")},
}

// NewClient returns a new web.Client
func NewClient(vapidPrivateKey string, vapidPublicKey string) *Client {
	// vapidPrivateKey, vapidPublicKey, err := webpush.GenerateVAPIDKeys()
	// if err != nil {
	// 	panic(err)
	// }

	return &Client{
		HTTPClient: &http.Client{
			Timeout: httpClientTimeout,
		},
		vapidPrivateKey: vapidPrivateKey,
		vapidPublicKey:  vapidPublicKey,
	}
}

// Push sends a notification using Webpush
func (c *Client) Push(n *Notification) (*Response, error) {
	jsonBuffer, _ := json.Marshal(n.Payload)
	var timeToLive int
	if n.TimeToLive != nil {
		timeToLive = int(*n.TimeToLive)
	} else {
		timeToLive = 2419200
	}

	subscription := webpush.Subscription{
		Endpoint: n.Subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: n.Subscription.Key,
			Auth:   n.Subscription.Auth,
		},
	}

	resp, err := webpush.SendNotification(jsonBuffer, &subscription, &webpush.Options{
		Subscriber:      "example@example.com", // Do not include "mailto:"
		VAPIDPublicKey:  c.vapidPublicKey,
		VAPIDPrivateKey: c.vapidPrivateKey,
		TTL:             timeToLive,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var pushResponse Response
	fmt.Println("Status Code: " + strconv.Itoa(resp.StatusCode))
	pushResponse.StatusCode = resp.StatusCode
	pushResponse.Body = string(body)
	if resp.StatusCode != 201 {
		return &pushResponse, errors.New("Push endpoint returned incorrect status code: " + strconv.Itoa(resp.StatusCode))
	}

	return &pushResponse, nil
}
