// Package push delivers newly created docket cards to a phone without an app
// (JOURNAL_SPEC §8.2): ntfy or plain webhook, one config line either way.
package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/maceip/clew/internal/config"
)

// Send pushes one item. Every pushed item names why it blocks (I8) — the
// body must carry the reason (composed by the differ's alert text).
func Send(p config.Push, title, body string) (bool, error) {
	if p.URL == "" {
		return false, nil // push not configured: docket + GitHub-mobile remain the read path
	}
	cl := &http.Client{Timeout: 10 * time.Second}
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := sendOnce(cl, p, title, body)
		if err != nil {
			return false, err
		}
		resp.Body.Close()
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
			return true, nil
		}
		if attempt < 2 && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) {
			time.Sleep(retryDelay(resp.Header.Get("Retry-After"), attempt))
			continue
		}
		return false, fmt.Errorf("push endpoint returned %s", resp.Status)
	}
	return false, fmt.Errorf("push endpoint did not accept the card")
}

func sendOnce(cl *http.Client, p config.Push, title, body string) (*http.Response, error) {
	var req *http.Request
	var err error
	switch p.Kind {
	case "webhook":
		payload, _ := json.Marshal(map[string]string{"title": title, "body": body})
		req, err = http.NewRequest("POST", p.URL, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	default: // ntfy
		req, err = http.NewRequest("POST", p.URL, bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Title", title)
		req.Header.Set("Priority", "high")
	}
	return cl.Do(req)
}

func retryDelay(header string, attempt int) time.Duration {
	header = strings.TrimSpace(header)
	if seconds, err := strconv.Atoi(header); err == nil && seconds >= 0 {
		delay := time.Duration(seconds) * time.Second
		if delay > 5*time.Second {
			return 5 * time.Second
		}
		return delay
	}
	if at, err := http.ParseTime(header); err == nil {
		delay := time.Until(at)
		if delay < 0 {
			return 0
		}
		if delay > 5*time.Second {
			return 5 * time.Second
		}
		return delay
	}
	return time.Duration(attempt+1) * 250 * time.Millisecond
}
