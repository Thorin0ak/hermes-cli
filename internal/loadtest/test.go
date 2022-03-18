package loadtest

import (
	"fmt"
	root "github.com/Thorin0ak/mercure-test/internal"
	"github.com/Thorin0ak/mercure-test/internal/token"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	authorizationHeader       = "Authorization"
	authorizationHeaderFormat = "Bearer %s"
)

type Tester interface {
	Run() error
}

type Test struct {
	config     *root.Config
	tokenMaker token.Maker
	hubUrl     string
}

func generateMockSseData(topicUri string, evtType string) url.Values {
	data := url.Values{}
	data.Set("topic", topicUri)
	data.Set("data", "yo")
	data.Set("type", evtType)
	data.Set("private", "on")
	return data
}

func publish(client *http.Client, url string, payload string, headers http.Header) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Could not create new HTP request: %v", err)
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("could not POST to Mercure Hub: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("Mercure returned error code %v", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not parse response body: %v\n", err)
	}

	log.Printf("Mercure ACK: %v\n", string(body))
}

func (t *Test) Run(headers http.Header) {
	client := http.Client{}
	durationStream := make(chan time.Duration)
	var err error

	token, err := t.tokenMaker.CreateToken("123456", t.config.Hermes.TopicUri, time.Minute*15)
	if err != nil {
		log.Fatalln(err)
	}

	h := http.Header{}
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	h.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, token))
	// override headers
	if len(headers) > 0 {
		for k, v := range headers {
			h.Set(k, strings.Join(v[:], ","))
		}
	}

	payload := generateMockSseData(t.config.Hermes.TopicUri, t.config.Hermes.EventType)
	encodedPayload := payload.Encode()

	go func() {
		defer close(durationStream)
		for i := 0; i < t.config.Hermes.NumEvents; i++ {
			begin := time.Now()
			w := rand.Intn(t.config.Hermes.MaxWaitTimes-t.config.Hermes.MinWaitTimes) + t.config.Hermes.MinWaitTimes
			log.Printf("interval: %d", w)
			time.Sleep(time.Millisecond * time.Duration(w))
			publish(&client, t.hubUrl, encodedPayload, h)
			since := time.Since(begin)
			durationStream <- since
		}
	}()

	for duration := range durationStream {
		fmt.Printf("%v taken to publish update\n", duration)
	}
}

func NewTest(config *root.Config) (*Test, error) {
	// TODO: pass context.Context and use req.WithContext(ctx)
	env := config.Hermes.ActiveEnv
	var hubUrl, secret string
	for i := 0; i < len(config.Mercure.Envs); i++ {
		if config.Mercure.Envs[i].Name == env {
			hubUrl = config.Mercure.Envs[i].HubUrl
			secret = config.Mercure.Envs[i].JwtSecret
		}
	}

	if len(hubUrl) == 0 {
		errMsg := fmt.Sprintf("no config found for active env %s", env)
		log.Fatal(errMsg)
	}

	m, err := token.NewJWTMaker(secret)
	if err != nil {
		log.Fatalln(err)
	}

	return &Test{config: config, tokenMaker: m, hubUrl: hubUrl}, nil
}
