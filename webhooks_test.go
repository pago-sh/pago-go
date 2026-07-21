package pago

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// testWebhookSecret is a Pago endpoint secret: the "whsec_" prefix followed by
// the base64-encoded 32 byte signing key.
const testWebhookSecret = "whsec_AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8="

// testWebhookRotatedSecret is a second, unrelated secret used to exercise
// rotation, where two signatures ride on the same header.
const testWebhookRotatedSecret = "whsec_ZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXp7fH1+f4CBgoM="

var testEventTypes = map[string]struct{}{"dummy.event": {}}

// signWebhookWith produces Standard Webhooks headers for a body, the same way
// the Pago API signs its deliveries: the HMAC key is the base64 decoding of
// the secret's payload, never the secret string itself.
func signWebhookWith(t *testing.T, secret, webhookID, body string, timestamp time.Time) http.Header {
	t.Helper()

	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	if err != nil {
		t.Fatalf("test secret is not a valid whsec_ secret: %v", err)
	}

	signedContent := fmt.Sprintf("%s.%d.%s", webhookID, timestamp.Unix(), body)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(signedContent))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	headers := http.Header{}
	headers.Set("Webhook-Id", webhookID)
	headers.Set("Webhook-Timestamp", fmt.Sprintf("%d", timestamp.Unix()))
	headers.Set("Webhook-Signature", "v1,"+signature)
	return headers
}

// signWebhook signs a body with the canonical test secret and webhook id.
func signWebhook(t *testing.T, body string, timestamp time.Time) http.Header {
	t.Helper()
	return signWebhookWith(t, testWebhookSecret, "test-webhook", body, timestamp)
}

func TestValidateWebhookAcceptsAKnownEvent(t *testing.T) {
	body := `{"type":"dummy.event","value":"payload"}`

	eventType, err := ValidateWebhook(
		[]byte(body),
		signWebhook(t, body, time.Now()),
		testWebhookSecret,
		testEventTypes,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eventType != "dummy.event" {
		t.Fatalf("expected the event type, got %q", eventType)
	}
}

// The HMAC key is the decoded secret, so a signature produced with the raw
// secret string as key — the pre-fix behaviour — must not verify.
func TestValidateWebhookDerivesTheKeyByBase64DecodingTheSecret(t *testing.T) {
	body := `{"type":"dummy.event"}`
	timestamp := time.Now()

	mac := hmac.New(sha256.New, []byte(testWebhookSecret))
	mac.Write([]byte(fmt.Sprintf("test-webhook.%d.%s", timestamp.Unix(), body)))

	headers := signWebhook(t, body, timestamp)
	headers.Set("Webhook-Signature", "v1,"+base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError for a raw-secret-keyed signature, got %v", err)
	}
}

// A third-party Standard Webhooks library derives the same key from the same
// secret, so a signature it produces must verify here.
func TestValidateWebhookMatchesAStandardWebhooksSignature(t *testing.T) {
	body := `{"type":"dummy.event","value":"payload"}`
	timestamp := time.Now()

	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(testWebhookSecret, "whsec_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(fmt.Sprintf("msg_2b.%d.%s", timestamp.Unix(), body)))

	headers := http.Header{}
	headers.Set("Webhook-Id", "msg_2b")
	headers.Set("Webhook-Timestamp", fmt.Sprintf("%d", timestamp.Unix()))
	headers.Set("Webhook-Signature", "v1,"+base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	if _, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// During a rotation both the old and the new secret sign the delivery; a
// receiver holding either one must accept it.
func TestValidateWebhookAcceptsAnySignatureDuringRotation(t *testing.T) {
	body := `{"type":"dummy.event"}`
	timestamp := time.Now()

	old := signWebhookWith(t, testWebhookSecret, "test-webhook", body, timestamp)
	rotated := signWebhookWith(t, testWebhookRotatedSecret, "test-webhook", body, timestamp)

	headers := old.Clone()
	headers.Set(
		"Webhook-Signature",
		old.Get("Webhook-Signature")+" "+rotated.Get("Webhook-Signature"),
	)

	for name, secret := range map[string]string{
		"old secret":     testWebhookSecret,
		"rotated secret": testWebhookRotatedSecret,
	} {
		if _, err := ValidateWebhook([]byte(body), headers, secret, testEventTypes); err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
	}
}

func TestValidateWebhookAcceptsMultipleSignatures(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhook(t, body, time.Now())
	headers.Set("Webhook-Signature", "v1,bm90LWEtc2lnbmF0dXJl "+headers.Get("Webhook-Signature"))

	if _, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWebhookRejectsAForgedSignature(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhook(t, body, time.Now())
	headers.Set("Webhook-Signature", "v1,"+base64.StdEncoding.EncodeToString(make([]byte, 32)))

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

// A signature made with a secret the receiver doesn't hold must be rejected.
func TestValidateWebhookRejectsASignatureFromAnotherSecret(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhookWith(t, testWebhookRotatedSecret, "test-webhook", body, time.Now())

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

func TestValidateWebhookRejectsATamperedBody(t *testing.T) {
	body := `{"type":"dummy.event","amount":100}`
	headers := signWebhook(t, body, time.Now())
	tampered := `{"type":"dummy.event","amount":900}`

	_, err := ValidateWebhook([]byte(tampered), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

// The webhook id is part of the signed content, so swapping it must invalidate
// the signature — otherwise a delivery could be replayed under a new id.
func TestValidateWebhookRejectsATamperedID(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhook(t, body, time.Now())
	headers.Set("Webhook-Id", "attacker-chosen-id")

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

// The timestamp is signed too, so a captured delivery can't be re-dated to slip
// back inside the tolerance window.
func TestValidateWebhookRejectsAReplayedDeliveryWithARefreshedTimestamp(t *testing.T) {
	body := `{"type":"dummy.event"}`
	captured := signWebhook(t, body, time.Now().Add(-30*time.Minute))

	headers := captured.Clone()
	headers.Set("Webhook-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

func TestValidateWebhookRejectsMalformedSignatureEncoding(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhook(t, body, time.Now())
	headers.Set("Webhook-Signature", "v1,not-base64!")

	_, err := ValidateWebhook([]byte(body), headers, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

func TestValidateWebhookRejectsMissingHeaders(t *testing.T) {
	_, err := ValidateWebhook([]byte("{}"), http.Header{}, testWebhookSecret, testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

func TestValidateWebhookRejectsAnEmptySecret(t *testing.T) {
	body := `{"type":"dummy.event"}`

	_, err := ValidateWebhook([]byte(body), signWebhook(t, body, time.Now()), "", testEventTypes)

	var typed *WebhookVerificationError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookVerificationError, got %v", err)
	}
}

// A malformed secret is a configuration mistake, and it must say so instead of
// failing as an ordinary signature mismatch.
func TestValidateWebhookRejectsAMalformedSecretExplicitly(t *testing.T) {
	body := `{"type":"dummy.event"}`
	headers := signWebhook(t, body, time.Now())

	for name, testCase := range map[string]struct{ secret, wants string }{
		"missing prefix": {
			secret: "AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8=",
			wants:  "must start with",
		},
		"raw random string": {
			secret: "test-secret",
			wants:  "must start with",
		},
		"prefixed but not base64": {
			secret: "whsec_not base64!!",
			wants:  "base64-encoded signing key",
		},
		"prefix only": {
			secret: "whsec_",
			wants:  "no signing key",
		},
	} {
		_, err := ValidateWebhook([]byte(body), headers, testCase.secret, testEventTypes)

		var typed *WebhookVerificationError
		if !errors.As(err, &typed) {
			t.Fatalf("%s: expected a *WebhookVerificationError, got %v", name, err)
		}
		if !strings.Contains(err.Error(), testCase.wants) {
			t.Fatalf("%s: expected the error to explain the malformed secret (%q), got %q",
				name, testCase.wants, err.Error())
		}
		if strings.Contains(err.Error(), "no matching signature found") {
			t.Fatalf("%s: a malformed secret must not be reported as a signature mismatch", name)
		}
	}
}

func TestValidateWebhookRejectsStaleAndFutureTimestamps(t *testing.T) {
	body := `{"type":"dummy.event"}`
	for name, timestamp := range map[string]time.Time{
		"stale":  time.Now().Add(-6 * time.Minute),
		"future": time.Now().Add(6 * time.Minute),
	} {
		_, err := ValidateWebhook(
			[]byte(body),
			signWebhook(t, body, timestamp),
			testWebhookSecret,
			testEventTypes,
		)
		var typed *WebhookVerificationError
		if !errors.As(err, &typed) {
			t.Fatalf("%s: expected a *WebhookVerificationError, got %v", name, err)
		}
	}
}

func TestValidateWebhookRejectsAnUnknownEventType(t *testing.T) {
	body := `{"type":"future.event"}`

	_, err := ValidateWebhook(
		[]byte(body),
		signWebhook(t, body, time.Now()),
		testWebhookSecret,
		testEventTypes,
	)

	var typed *WebhookUnknownTypeError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookUnknownTypeError, got %v", err)
	}
	if typed.EventType != "future.event" {
		t.Fatalf("expected the unknown event type, got %q", typed.EventType)
	}
}

func TestValidateWebhookRejectsAMalformedPayload(t *testing.T) {
	body := "{"

	_, err := ValidateWebhook(
		[]byte(body),
		signWebhook(t, body, time.Now()),
		testWebhookSecret,
		testEventTypes,
	)

	var typed *WebhookError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *WebhookError, got %v", err)
	}
}

func TestWebhookErrorHierarchy(t *testing.T) {
	var base *WebhookError
	if !errors.As(error(NewWebhookVerificationError("boom")), &base) {
		t.Fatal("expected a verification error to unwrap to *WebhookError")
	}
	if !errors.As(error(NewWebhookUnknownTypeError("x")), &base) {
		t.Fatal("expected an unknown type error to unwrap to *WebhookError")
	}
}
