package binance

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Client struct {
	client *http.Client
}

func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 240 * time.Second,
		},
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest(GET, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			log.Printf("request failed (attempt %d/%d): %v\n", attempt, maxRetries, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			defer res.Body.Close()
			lastErr = err
			log.Printf("failed to read response body (attempt %d/%d): %v\n", attempt, maxRetries, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			break
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			lastErr = errors.New(string(body))
			log.Printf("request returned status %d (attempt %d/%d): %s\n", res.StatusCode, attempt, maxRetries, string(body))

			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		return body, nil
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}

func (c *Client) Post(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(POST, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.New(string(body))
	}

	return body, err
}

func (c *Client) Delete(url string) ([]byte, error) {
	req, err := http.NewRequest(DELETE, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New(string(resBody))
	}

	return resBody, nil
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
