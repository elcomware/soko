package payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

// VerifyStripeSignature verifies the exact raw payload and a five-minute replay
// window. Callers must still validate event identity/amount and durably deduplicate.
func VerifyStripeSignature(payload []byte, header, secret string, now time.Time) error {
	if secret == "" {
		return ErrUnconfigured
	}
	if len(payload) > 1024*1024 || len(header) > 4096 {
		return errors.New("webhook exceeds limits")
	}
	var timestamp string
	var sigs []string
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		if kv[0] == "t" {
			if timestamp != "" {
				return errors.New("ambiguous timestamp")
			}
			timestamp = kv[1]
		}
		if kv[0] == "v1" {
			sigs = append(sigs, kv[1])
		}
	}
	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || now.Sub(time.Unix(unix, 0)) > 5*time.Minute || time.Unix(unix, 0).Sub(now) > 5*time.Minute {
		return errors.New("webhook timestamp outside tolerance")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(payload)
	want := mac.Sum(nil)
	for _, s := range sigs {
		got, err := hex.DecodeString(s)
		if err == nil && hmac.Equal(got, want) {
			return nil
		}
	}
	return errors.New("invalid webhook signature")
}
