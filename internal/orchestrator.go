package internal

import (
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	net "github.com/subchord/go-sse"
	"go.uber.org/zap"
	"io/ioutil"
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
	logger     *zap.SugaredLogger
}

func generateMockSseData(topicUri string, evtType string) url.Values {
	data := url.Values{}
	data.Set("topic", topicUri)
	data.Set("data", "mock")
	data.Set("type", evtType)
	data.Set("private", "on")
	return data
}

func publish(client *http.Client, url string, payload string, headers http.Header, logger *zap.SugaredLogger) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		logger.Fatalf("Could not create new HTTP request: %v", err)
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("could not POST to Mercure Hub: %v\n", err)
		return
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode > 299 {
		logger.Errorf("Mercure returned error code %v", resp.StatusCode)
		return
	}

	// this seems to be required to re-use the http client
	ioutil.ReadAll(resp.Body) //nolint:errcheck
}

func (o *Orchestrator) subscribe(wg *sync.WaitGroup, stdoutBuffer *bytes.Buffer, headers http.Header) {
	token, err := o.tokenMaker.CreateToken("123456", o.config.Hermes.TopicUri, time.Minute*15, false)
	if err != nil {
		o.logger.Fatal(err)
	}

	h := http.Header{}
	h.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, token))
	if len(headers) > 0 {
		for k, v := range headers {
			h.Set(k, strings.Join(v[:], ","))
		}
	}

	evtUrl := fmt.Sprintf("%v?topic=%v", o.hubUrl, o.config.Hermes.TopicUri)
	stream, err := net.ConnectWithSSEFeed(evtUrl, h)
	if err != nil {
		o.logger.Fatal(err)
	}
	defer stream.Close()

	sub, err := stream.Subscribe(o.config.Hermes.EventType)
	if err != nil {
		o.logger.Errorf("error with subscription: %v", err)
		return
	}
	defer sub.Close()
	for {
		select {
		case evt := <-sub.Feed():
			fmt.Fprintf(stdoutBuffer, "Mercure Event: %v\n", evt)
			wg.Done()
		case err := <-sub.ErrFeed():
			o.logger.Error(err)
			return
		}
	}
}

func (o *Orchestrator) Run(pubHeaders http.Header, subHeaders http.Header) {
	client := http.Client{}
	durationStream := make(chan time.Duration)
	// Progress bar
	bar := pb.StartNew(o.config.Hermes.NumEvents)
	bar.SetWriter(os.Stdout)
	// in-memory buffer to avoid writing to stdout while the progress bar is there
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout) //nolint:errcheck

	var wg sync.WaitGroup
	if !o.config.Hermes.PublishOnly {
		wg.Add(o.config.Hermes.NumEvents)
		go o.subscribe(&wg, &stdoutBuff, subHeaders)
	}

	token, err := o.tokenMaker.CreateToken("123456", o.config.Hermes.TopicUri, time.Minute*15, true)
	if err != nil {
		o.logger.Fatal(err)
	}

	h := http.Header{}
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	h.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, token))
	if len(pubHeaders) > 0 {
		for k, v := range pubHeaders {
			h.Set(k, strings.Join(v[:], ","))
		}
	}

	payload := generateMockSseData(o.config.Hermes.TopicUri, o.config.Hermes.EventType)
	encodedPayload := payload.Encode()

	go func() {
		defer close(durationStream)
		for i := 0; i < o.config.Hermes.NumEvents; i++ {
			bar.Increment()
			begin := time.Now()
			w := rand.Intn(o.config.Hermes.MaxWaitTimes-o.config.Hermes.MinWaitTimes) + o.config.Hermes.MinWaitTimes
			time.Sleep(time.Millisecond * time.Duration(w))
			publish(&client, o.hubUrl, encodedPayload, h, o.logger)
			since := time.Since(begin)
			durationStream <- since
		}
	}()

	for range durationStream {
		// TODO: collect data
	}
	bar.Finish()
	if !o.config.Hermes.PublishOnly {
		wg.Wait()
	}
}

func GetDummy(config *Config, hubUrl string, secret string, logger *zap.SugaredLogger) *Orchestrator {
	m, _ := NewJWTMaker(secret, logger)
	return &Orchestrator{config, m, hubUrl, logger}
}

func NewOrchestrator(config *Config, logger *zap.SugaredLogger) (*Orchestrator, error) {
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
		return nil, fmt.Errorf("no config found for active env %s", env)
	}

	m, err := NewJWTMaker(secret, logger)
	if err != nil {
		return nil, err
	}

	return &Orchestrator{config, m, hubUrl, logger}, nil
}
