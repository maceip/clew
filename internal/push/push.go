// Package push delivers human-blocking inbox items to a phone without an app
// (JOURNAL_SPEC §8.2): ntfy or plain webhook, one config line either way.
package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"clew/internal/config"
)

// Send pushes one item. Every pushed item names why it blocks (I8) — the
// body must carry the reason (composed by the differ's alert text).
func Send(p config.Push, title, body string) (bool, error) {
	if p.URL == "" {
		return false, nil // push not configured: docket + GitHub-mobile remain the read path
	}
	cl := &http.Client{Timeout: 10 * time.Second}
	var resp *http.Response
	var err error
	switch p.Kind {
	case "webhook":
		payload, _ := json.Marshal(map[string]string{"title": title, "body": body})
		var req *http.Request
		req, err = http.NewRequest("POST", p.URL, bytes.NewReader(payload))
		if err != nil {
			return false, err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err = cl.Do(req)
	default: // ntfy
		var req *http.Request
		req, err = http.NewRequest("POST", p.URL, bytes.NewReader([]byte(body)))
		if err != nil {
			return false, err
		}
		req.Header.Set("Title", title)
		req.Header.Set("Priority", "high")
		resp, err = cl.Do(req)
	}
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false, fmt.Errorf("push endpoint returned %s", resp.Status)
	}
	return true, nil
}
