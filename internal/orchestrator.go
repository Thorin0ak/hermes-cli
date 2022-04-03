package internal

import (
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	net "github.com/subchord/go-sse"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	authorizationHeader       = "Authorization"
	authorizationHeaderFormat = "Bearer %s"
)

type Orchestrator struct {
	config     *Config
	tokenMaker Maker
	hubUrl     string
}

func generateMockSseData(topicUri string, evtType string) url.Values {
	data := url.Values{}
	data.Set("topic", topicUri)
	data.Set("data", "mock")
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
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode > 299 {
		log.Printf("Mercure returned error code %v", resp.StatusCode)
		return
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not parse response body: %v\n", err)
	}
}

func (t *Orchestrator) subscribe(wg *sync.WaitGroup, stdoutBuffer *bytes.Buffer, headers http.Header) {
	token, err := t.tokenMaker.CreateToken("123456", t.config.Hermes.TopicUri, time.Minute*15, false)
	if err != nil {
		log.Fatalln(err)
	}

	h := http.Header{}
	h.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, token))
	if len(headers) > 0 {
		for k, v := range headers {
			h.Set(k, strings.Join(v[:], ","))
		}
	}

	evtUrl := fmt.Sprintf("%v?topic=%v", t.hubUrl, t.config.Hermes.TopicUri)
	stream, err := net.ConnectWithSSEFeed(evtUrl, h)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer stream.Close()

	sub, err := stream.Subscribe(t.config.Hermes.EventType)
	if err != nil {
		log.Printf("error with subscription: %v", err)
		return
	}
	defer sub.Close()
	for {
		select {
		case evt := <-sub.Feed():
			fmt.Fprintf(stdoutBuffer, "Mercure ACK: %v\n", evt)
			wg.Done()
		case err := <-sub.ErrFeed():
			log.Fatal(err)
			return
		}
	}
}

func (t *Orchestrator) Run(pubHeaders http.Header, subHeaders http.Header) {
	client := http.Client{}
	durationStream := make(chan time.Duration)
	// Progress bar
	bar := pb.StartNew(t.config.Hermes.NumEvents)
	bar.SetWriter(os.Stdout)
	// in-memory buffer to avoid writing to stdout while the progress bar is there
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout) //nolint:errcheck

	var wg sync.WaitGroup
	wg.Add(t.config.Hermes.NumEvents)
	go t.subscribe(&wg, &stdoutBuff, subHeaders)

	token, err := t.tokenMaker.CreateToken("123456", t.config.Hermes.TopicUri, time.Minute*15, true)
	if err != nil {
		log.Fatalln(err)
	}

	h := http.Header{}
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	h.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, token))
	if len(pubHeaders) > 0 {
		for k, v := range pubHeaders {
			h.Set(k, strings.Join(v[:], ","))
		}
	}

	payload := generateMockSseData(t.config.Hermes.TopicUri, t.config.Hermes.EventType)
	encodedPayload := payload.Encode()

	go func() {
		defer close(durationStream)
		for i := 0; i < t.config.Hermes.NumEvents; i++ {
			bar.Increment()
			begin := time.Now()
			w := rand.Intn(t.config.Hermes.MaxWaitTimes-t.config.Hermes.MinWaitTimes) + t.config.Hermes.MinWaitTimes
			time.Sleep(time.Millisecond * time.Duration(w))
			publish(&client, t.hubUrl, encodedPayload, h)
			since := time.Since(begin)
			durationStream <- since
		}
	}()

	for range durationStream {
		// TODO: collect data
	}
	bar.Finish()
	wg.Wait()
}

func NewOrchestrator(config *Config) (*Orchestrator, error) {
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

	m, err := NewJWTMaker(secret)
	if err != nil {
		log.Fatalln(err)
	}

	return &Orchestrator{config, m, hubUrl}, nil
}
