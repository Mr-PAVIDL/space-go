package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net/http"
	"space-go/internal/model"
	"time"
)

type DatsEdenSpaceClient struct {
	BaseURL     string
	AuthToken   string
	HTTPClient  *http.Client
	RateLimiter *rate.Limiter
}

func NewClient(baseURL, authToken string) *DatsEdenSpaceClient {
	return &DatsEdenSpaceClient{
		BaseURL:     baseURL,
		AuthToken:   authToken,
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
		RateLimiter: rate.NewLimiter(rate.Every(time.Second/4), 1),
	}
}

func (c *DatsEdenSpaceClient) createRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, endpoint), buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.AuthToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *DatsEdenSpaceClient) doRequest(ctx context.Context, req *http.Request, v interface{}) error {
	err := c.RateLimiter.Wait(ctx)
	if err != nil {
		return fmt.Errorf("rate limit error: %w", err)
	}

tryAgain:
	req = req.WithContext(ctx)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("Rate limit exceeded, waiting for 2 seconds")
			time.Sleep(2 * time.Second)
			goto tryAgain
		} else if resp.StatusCode == http.StatusBadGateway {
			log.Println("Bad gateway, waiting for 2 seconds")
			time.Sleep(2 * time.Second)
			goto tryAgain
		} else if resp.StatusCode == http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			return fmt.Errorf("API request error: status code %d, response: %s", resp.StatusCode, bodyString)
		} else {
			log.Println("Hz what happened, waiting for 2 seconds")
			time.Sleep(2 * time.Second)
			goto tryAgain
		}
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}

func (c *DatsEdenSpaceClient) GetUniverse(ctx context.Context) (*model.UniverseResponse, error) {
	req, err := c.createRequest("GET", "/player/universe", nil)
	if err != nil {
		return nil, err
	}

	var player model.UniverseResponse
	if err := c.doRequest(ctx, req, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (c *DatsEdenSpaceClient) Travel(ctx context.Context, travelRequest model.TravelRequest) (*model.TravelResponse, error) {
	req, err := c.createRequest("POST", "/player/travel", travelRequest)
	if err != nil {
		return nil, err
	}

	var travelResponse model.TravelResponse
	if err := c.doRequest(ctx, req, &travelResponse); err != nil {
		return nil, err
	}

	return &travelResponse, nil
}

func (c *DatsEdenSpaceClient) CollectGarbage(ctx context.Context, collectRequest model.CollectRequest) (*model.CollectResponse, error) {
	req, err := c.createRequest("POST", "/player/collect", collectRequest)
	if err != nil {
		return nil, err
	}

	var collectResponse model.CollectResponse
	if err := c.doRequest(ctx, req, &collectResponse); err != nil {
		return nil, err
	}

	return &collectResponse, nil
}

func (c *DatsEdenSpaceClient) ResetGameState(ctx context.Context) (*model.ResetResponse, error) {
	req, err := c.createRequest("DELETE", "/player/reset", nil)
	if err != nil {
		return nil, err
	}

	var acceptedResponse model.ResetResponse
	if err := c.doRequest(ctx, req, &acceptedResponse); err != nil {
		return nil, err
	}

	return &acceptedResponse, nil
}

func (c *DatsEdenSpaceClient) GetRounds(ctx context.Context) (*model.RoundsResponse, error) {
	req, err := c.createRequest("GET", "/player/rounds", nil)
	if err != nil {
		return nil, err
	}

	var roundList model.RoundsResponse
	if err := c.doRequest(ctx, req, &roundList); err != nil {
		return nil, err
	}

	return &roundList, nil
}
