package pago

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WebhookToleranceSeconds is how far a webhook timestamp may drift from the
// current time before the message is rejected.
const WebhookToleranceSeconds = 5 * 60

// WebhookSecretPrefix is the prefix every Pago endpoint secret carries. What
// follows it is the base64-encoded HMAC signing key.
const WebhookSecretPrefix = "whsec_"

// WebhookError is the base error raised while processing a Pago webhook.
type WebhookError struct {
	baseError
	Message string
	Err     error
}

func (e *WebhookError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("pago: %s: %v", e.Message, e.Err)
	}
	return "pago: " + e.Message
}

func (e *WebhookError) Unwrap() error { return e.Err }

// NewWebhookError builds a WebhookError.
func NewWebhookError(message string, err error) *WebhookError {
	return &WebhookError{Message: message, Err: err}
}

// WebhookVerificationError is raised when a webhook signature cannot be
// verified.
type WebhookVerificationError struct {
	WebhookError
}

// NewWebhookVerificationError builds a WebhookVerificationError.
func NewWebhookVerificationError(message string) *WebhookVerificationError {
	return &WebhookVerificationError{WebhookError{Message: message}}
}

func (e *WebhookVerificationError) Unwrap() error { return &e.WebhookError }

// WebhookUnknownTypeError is raised when a verified webhook has a type this
// SDK version does not know about.
type WebhookUnknownTypeError struct {
	WebhookError
	EventType string
}

// NewWebhookUnknownTypeError builds a WebhookUnknownTypeError.
func NewWebhookUnknownTypeError(eventType string) *WebhookUnknownTypeError {
	return &WebhookUnknownTypeError{
		WebhookError: WebhookError{
			Message: fmt.Sprintf("unknown webhook event type: %q", eventType),
		},
		EventType: eventType,
	}
}

func (e *WebhookUnknownTypeError) Unwrap() error { return &e.WebhookError }

// ValidateWebhook verifies the signature of a raw webhook request and returns
// its event type.
//
// It implements the Standard Webhooks specification: the signed content is
// "<webhook-id>.<webhook-timestamp>.<body>", authenticated with HMAC-SHA256
// and encoded as standard base64. The signing key is not the secret string
// itself — it is the base64 decoding of whatever follows the "whsec_" prefix.
// Headers are read through http.Header, so passing (*http.Request).Header
// works directly.
func ValidateWebhook(
	body []byte,
	headers http.Header,
	secret string,
	eventTypes map[string]struct{},
) (string, error) {
	if err := VerifyWebhookSignature(body, headers, secret); err != nil {
		return "", err
	}

	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "", NewWebhookError("failed to parse webhook payload", err)
	}
	if _, known := eventTypes[envelope.Type]; !known {
		return "", NewWebhookUnknownTypeError(envelope.Type)
	}
	return envelope.Type, nil
}

// webhookSigningKey derives the HMAC key from an endpoint secret.
//
// A Pago secret is "whsec_" followed by the base64-encoded key, and the key is
// what signs the message. Anything that doesn't fit that shape is reported
// instead of being silently used as a key, which would fail verification later
// with a misleading "no matching signature found".
func webhookSigningKey(secret string) ([]byte, error) {
	if secret == "" {
		return nil, NewWebhookVerificationError("secret can't be empty")
	}

	encodedKey, hasPrefix := strings.CutPrefix(secret, WebhookSecretPrefix)
	if !hasPrefix {
		return nil, NewWebhookVerificationError(fmt.Sprintf(
			"secret must start with %q; pass the endpoint secret exactly as issued",
			WebhookSecretPrefix,
		))
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		key, err = base64.RawStdEncoding.DecodeString(encodedKey)
	}
	if err != nil {
		return nil, NewWebhookVerificationError(fmt.Sprintf(
			"secret must be %q followed by a base64-encoded signing key",
			WebhookSecretPrefix,
		))
	}
	if len(key) == 0 {
		return nil, NewWebhookVerificationError(fmt.Sprintf(
			"secret carries no signing key after the %q prefix", WebhookSecretPrefix,
		))
	}

	return key, nil
}

// VerifyWebhookSignature verifies the Standard Webhooks signature of a raw
// request body.
func VerifyWebhookSignature(body []byte, headers http.Header, secret string) error {
	key, err := webhookSigningKey(secret)
	if err != nil {
		return err
	}

	webhookID := headers.Get("webhook-id")
	webhookTimestamp := headers.Get("webhook-timestamp")
	webhookSignature := headers.Get("webhook-signature")
	if webhookID == "" || webhookTimestamp == "" || webhookSignature == "" {
		return NewWebhookVerificationError("missing required headers")
	}

	timestamp, err := strconv.ParseFloat(webhookTimestamp, 64)
	if err != nil || math.IsInf(timestamp, 0) || math.IsNaN(timestamp) {
		return NewWebhookVerificationError("invalid signature headers")
	}

	now := float64(time.Now().Unix())
	if timestamp < now-WebhookToleranceSeconds {
		return NewWebhookVerificationError("message timestamp too old")
	}
	if timestamp > now+WebhookToleranceSeconds {
		return NewWebhookVerificationError("message timestamp too new")
	}

	signedContent := fmt.Sprintf("%s.%d.%s", webhookID, int64(math.Floor(timestamp)), body)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(signedContent))
	expected := mac.Sum(nil)

	// The header may carry several "v1,<base64>" entries separated by spaces so
	// a secret can be rotated without dropping deliveries: any one of them
	// matching accepts the message. Every candidate is compared, and the result
	// is accumulated rather than returned early, so neither the outcome nor the
	// position of the matching entry is observable through timing.
	matched := 0
	for _, versionedSignature := range strings.Fields(webhookSignature) {
		version, signature, found := strings.Cut(versionedSignature, ",")
		if !found || version != "v1" {
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			continue
		}
		matched |= subtle.ConstantTimeCompare(expected, decoded)
	}
	if matched == 1 {
		return nil
	}

	return NewWebhookVerificationError("no matching signature found")
}
