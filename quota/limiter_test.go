package quota

import (
	"testing"
	"time"
)

func TestTokenBucketRefills(t *testing.T) {
	l := New(2, 1)
	now := time.Unix(0, 0)
	l.now = func() time.Time { return now }
	if !l.Allow("tenant", 1).Allowed || !l.Allow("tenant", 1).Allowed || l.Allow("tenant", 1).Allowed {
		t.Fatal("capacity enforcement failed")
	}
	now = now.Add(time.Second)
	if !l.Allow("tenant", 1).Allowed {
		t.Fatal("expected refill")
	}
}
