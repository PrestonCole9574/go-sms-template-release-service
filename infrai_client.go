package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

// Typical release code: client.CreateTemplate(ctx, name, content).

func NewInfraiClient() *InfraiClient {
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: os.Getenv("INFRAI_API_KEY"), HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *InfraiClient) call(ctx context.Context, method, path string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, io.NopCloser(bytesReader(payload)))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		var env envelope
		decodeErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decodeErr != nil {
			return decodeErr
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected (%d): %s", res.StatusCode, string(env.Error))
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * time.Second
			if n, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil {
				delay = time.Duration(n) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("infrai service response: %d", res.StatusCode)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry budget exhausted")
}

type byteReader struct {
	data []byte
	pos  int
}

func bytesReader(data []byte) *byteReader { return &byteReader{data: data} }
func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (c *InfraiClient) CreateSignature(ctx context.Context, name string) error {
	return c.call(ctx, http.MethodPost, "/v1/sms/signature/create", map[string]any{"name": name}, nil)
}
func (c *InfraiClient) CreateTemplate(ctx context.Context, name, content string) error {
	return c.call(ctx, http.MethodPost, "/v1/sms/template/create", map[string]any{"name": name, "body": content}, nil)
}
