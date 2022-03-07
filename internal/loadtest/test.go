package loadtest

import (
	"fmt"
	root "github.com/Thorin0ak/mercure-test/internal"
	"github.com/Thorin0ak/mercure-test/internal/token"
	"log"
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
	req        *http.Request
	config     *root.Config
	tokenMaker *token.Maker
}

func generateMockSseData(topicUri string, evtType string) url.Values {
	data := url.Values{}
	data.Set("topic", topicUri)
	data.Set("data", "yo")
	data.Set("type", evtType)
	data.Set("private", "on")
	return data
}

func (t *Test) Run() error {
	client := http.Client{}
	resp, err := client.Do(t.req)
	if err != nil {
		fmt.Printf("could not POST to Mercure Hub: %v", err)
		return fmt.Errorf("got error %s", err.Error())
	}

	if resp.StatusCode > 299 {
		resp.Body.Close()
		return fmt.Errorf("POST request to Mercure Hub received error: %d", resp.StatusCode)
	}

	return nil
}

func NewTest(config *root.Config, headers http.Header) (*Test, error) {
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
	jwtToken, err := m.CreateToken("123456", fmt.Sprintf("%s/%s", config.Hermes.TopicUri, config.Hermes.EventType), time.Minute*15)
	if err != nil {
		log.Fatalln(err)
	}

	payload := generateMockSseData(config.Hermes.TopicUri, config.Hermes.EventType)
	encodedPayload := payload.Encode()
	req, err := http.NewRequest(http.MethodPost, hubUrl, strings.NewReader(encodedPayload))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(authorizationHeader, fmt.Sprintf(authorizationHeaderFormat, jwtToken))
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, strings.Join(v[:], ","))
		}
	}

	return &Test{req: req, config: config, tokenMaker: &m}, nil
}
