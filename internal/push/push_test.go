package push

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/maceip/clew/internal/config"
)

func TestSendDistinguishesDisabledDeliveredAndRejected(t *testing.T) {
	if sent, err := Send(config.Push{}, "title", "body"); err != nil || sent {
		t.Fatalf("disabled Send() = %v, %v; want false, nil", sent, err)
	}

	ok := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer ok.Close()
	if sent, err := Send(config.Push{URL: ok.URL}, "title", "body"); err != nil || !sent {
		t.Fatalf("successful Send() = %v, %v; want true, nil", sent, err)
	}

	rejected := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "no", http.StatusForbidden)
	}))
	defer rejected.Close()
	if sent, err := Send(config.Push{URL: rejected.URL}, "title", "body"); err == nil || sent {
		t.Fatalf("rejected Send() = %v, %v; want false, error", sent, err)
	}
}

func TestSendRetriesRateLimitAndServerFailures(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		switch calls.Add(1) {
		case 1:
			w.Header().Set("Retry-After", "0")
			http.Error(w, "slow down", http.StatusTooManyRequests)
		case 2:
			http.Error(w, "temporary", http.StatusBadGateway)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()
	sent, err := Send(config.Push{URL: server.URL}, "title", "body")
	if err != nil || !sent {
		t.Fatalf("retried Send() = %v, %v", sent, err)
	}
	if calls.Load() != 3 {
		t.Fatalf("calls=%d want 3", calls.Load())
	}
}
