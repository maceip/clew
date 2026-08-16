// Package push delivers human-blocking inbox items to a phone without an app
// (JOURNAL_SPEC §8.2): ntfy or plain webhook, one config line either way.
package push

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"clew/internal/config"
)

// Send pushes one item. Every pushed item names why it blocks (I8) — the
// body must carry the reason (composed by the differ's alert text).
func Send(p config.Push, title, body string) error {
	if p.URL == "" {
		return nil // push not configured: inbox + GitHub-mobile remain the read path
	}
	cl := &http.Client{Timeout: 10 * time.Second}
	switch p.Kind {
	case "webhook":
		payload, _ := json.Marshal(map[string]string{"title": title, "body": body})
		req, err := http.NewRequest("POST", p.URL, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := cl.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		return nil
	default: // ntfy
		req, err := http.NewRequest("POST", p.URL, bytes.NewReader([]byte(body)))
		if err != nil {
			return err
		}
		req.Header.Set("Title", title)
		req.Header.Set("Priority", "high")
		resp, err := cl.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		return nil
	}
}
