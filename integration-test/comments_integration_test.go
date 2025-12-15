package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

/*const (
	host           = "app" // контейнер или локальный хост сервиса
	requestTimeout = 5 * time.Second
	httpURL        = "http://" + host + ":8080"
	basePathV1     = httpURL + "/v1"
	attempts       = 20
)

var errHealthCheck = fmt.Errorf("url %s is not available", httpURL+"/healthz")

func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func healthCheck(attempts int) error {
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, httpURL+"/healthz", nil)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", httpURL+"/healthz", attempts)
		time.Sleep(time.Second)
		attempts--
	}
	return errHealthCheck
}

func TestMain(m *testing.M) {
	if err := healthCheck(attempts); err != nil {
		log.Fatalf("Integration tests: service not available: %s", err)
	}
	os.Exit(m.Run())
}*/

// -----------------------------
// POST /v1/comments
func TestHTTPCreateComment(t *testing.T) {
	tests := []struct {
		description string
		body        string
		expected    int
	}{
		{"Valid comment", `{"entity_id":1,"text":"Hello world"}`, http.StatusCreated},
		{"Missing text", `{"entity_id":1}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			url := basePathV1 + "/comments"

			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
			defer cancel()

			resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, url, bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expected {
				t.Errorf("Expected status %d, got %d", tt.expected, resp.StatusCode)
			}
		})
	}
}

// -----------------------------
// GET /v1/comments/:id
func TestHTTPGetComment(t *testing.T) {
	commentID := 1 // подготовь в БД заранее
	url := fmt.Sprintf("%s/comments/%d", basePathV1, commentID)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var comment struct {
		ID       int    `json:"id"`
		EntityID int    `json:"entity_id"`
		Text     string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if comment.ID != commentID {
		t.Errorf("Expected ID %d, got %d", commentID, comment.ID)
	}
}

// -----------------------------
// GET /v1/comments?entity_id=123
func TestHTTPListComments(t *testing.T) {
	entityID := 1
	url := fmt.Sprintf("%s/comments?entity_id=%d&page=1&limit=20", basePathV1, entityID)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var body struct {
		Comments []struct {
			ID       int    `json:"id"`
			EntityID int    `json:"entity_id"`
			Text     string `json:"text"`
		} `json:"comments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(body.Comments) == 0 {
		t.Error("Expected at least one comment")
	}
}

// -----------------------------
// PUT /v1/comments/:id
func TestHTTPUpdateComment(t *testing.T) {
	commentID := 1
	url := fmt.Sprintf("%s/comments/%d", basePathV1, commentID)
	body := `{"text":"Updated comment"}`

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodPut, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// -----------------------------
// DELETE /v1/comments/:id
func TestHTTPDeleteComment(t *testing.T) {
	commentID := 1
	url := fmt.Sprintf("%s/comments/%d", basePathV1, commentID)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusNoContent, resp.StatusCode)
	}
}
