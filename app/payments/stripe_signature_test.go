package payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	now := time.Unix(1700000000, 0)
	p := []byte(`{"id":"evt_1"}`)
	stamp := strconv.FormatInt(now.Unix(), 10)
	m := hmac.New(sha256.New, []byte("test_secret"))
	m.Write([]byte(stamp + "."))
	m.Write(p)
	h := "t=" + stamp + ",v1=" + hex.EncodeToString(m.Sum(nil))
	if VerifyStripeSignature(p, h, "test_secret", now) != nil {
		t.Fatal("valid signature rejected")
	}
	for _, tt := range []struct {
		p    []byte
		h, s string
		now  time.Time
	}{{[]byte("changed"), h, "test_secret", now}, {p, h, "wrong", now}, {p, h, "test_secret", now.Add(6 * time.Minute)}, {p, h, "", now}, {p, "t=1,t=2,v1=bad", "test_secret", now}} {
		if VerifyStripeSignature(tt.p, tt.h, tt.s, tt.now) == nil {
			t.Fatal("invalid signature accepted")
		}
	}
}
