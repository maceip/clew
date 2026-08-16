package push

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"clew/internal/config"
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
