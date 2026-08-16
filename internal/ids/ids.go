// Package ids issues ULID-based identifiers for journal entries and events.
// Entry ids are prefixed "e", event ids "v", per JOURNAL_SPEC §3.2.
package ids

import (
	"crypto/rand"
	"regexp"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	mu      sync.Mutex
	entropy = ulid.Monotonic(rand.Reader, 0)

	entryRe = regexp.MustCompile(`^e[0-9A-HJKMNP-TV-Z]{26}$`)
	eventRe = regexp.MustCompile(`^v[0-9A-HJKMNP-TV-Z]{26}$`)
)

func newULID(t time.Time) string {
	mu.Lock()
	defer mu.Unlock()
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// NewEntry returns a new entry id. t stamps the ULID time component so that
// backfilled entries sort at their original position.
func NewEntry(t time.Time) string { return "e" + newULID(t) }

// NewEvent returns a new event id.
func NewEvent(t time.Time) string { return "v" + newULID(t) }

func ValidEntry(id string) bool { return entryRe.MatchString(id) }
func ValidEvent(id string) bool { return eventRe.MatchString(id) }

// Time extracts the timestamp encoded in an id, or zero if unparseable.
func Time(id string) time.Time {
	if len(id) != 27 {
		return time.Time{}
	}
	u, err := ulid.ParseStrict(id[1:])
	if err != nil {
		return time.Time{}
	}
	return time.UnixMilli(int64(u.Time())).UTC()
}
