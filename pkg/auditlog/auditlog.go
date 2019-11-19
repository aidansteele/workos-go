// Package auditlog is a package to send audit logs events to WorkOS.
//
// Example:
//   func main() {
//       auditlog.SetAPIKey("my_api_key")
//
//       // Wherever you need to publish an audit log event:
//       err := auditlog.Publish(context.Background(), auditlog.Event{
//           Action:     "document.viewed",
//           ActionType: auditlog.Create,
//           ActorName:  "Jairo Kunde",
//           ActorID:    "user_01DGZ0FAXN978HCET66Q98QMTQ",
//           Group:      "abstract.com",
//           Location:   "55.27.223.26",
//           OccurredAt: time.Now(),
//           TargetName: "central.class",
//           TargetID:   "doc_01DGZ0FAXP4HA4X0BVFKS0ZH4Y",
//       })
//       if err != nil {
//           // Handle error.
//       }
//   }
package auditlog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	// DefaultClient is the client used by SetAPIKey and Publish functions.
	DefaultClient = &Client{
		APIKey:   os.Getenv("WORKOS_API_KEY"),
		Endpoint: "https://api.workos.com/events",
	}

	// GlobalMetadata are metadata that are injected in every audit log events.
	GlobalMetadata Metadata

	errTooMuchMetadataKeys = errors.New("too much metadata key")
)

// SetAPIKey sets the WorkOS API key to use when using Publish.
func SetAPIKey(k string) {
	DefaultClient.APIKey = k
}

// Publish publishes the given event.
func Publish(ctx context.Context, e Event) error {
	return DefaultClient.Publish(ctx, e)
}

// Event represents an Audit Log event.
type Event struct {
	Action     string     `json:"action"`
	ActionType ActionType `json:"action_type"`
	ActorName  string     `json:"actor_name"`
	ActorID    string     `json:"actor_id"`
	Group      string     `json:"group"`

	// A key that ensures that an same event is not processed multiple time.
	// Once the event is sent for the first time, the lock on the key expires
	// after 24 hours.
	//
	// An idempotency key is automatically generated if not set.
	IdempotencyKey string `json:"-"`

	// An ip address that locates where the audit log occured.
	Location string `json:"location"`

	// The event metadata. It can't contain more than 50 keys. A key can't
	// exeed 40 characters.
	Metadata Metadata `json:"metadata,omitempty"`

	// The time when the audit log occured.
	//
	// Defaults to time.Now().
	OccurredAt time.Time `json:"occurred_at"`

	TargetName string `json:"target_name"`
	TargetID   string `json:"target_id"`
}

// ActionType is the type that holds the CRUD action used for the WorkOS Audit
// Log.
type ActionType string

// Constants that enumerate the different action types.
const (
	Create ActionType = "C"
	Read   ActionType = "R"
	Update ActionType = "U"
	Delete ActionType = "D"
)

// Metadata represents metadata to be attached to an audit log event.
type Metadata map[string]interface{}

// Merges the given metadata. Values from m are not overridden by the ones from
// other.
func (m Metadata) merge(other Metadata) {
	for k, v := range other {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}
}

func (m Metadata) validate() error {
	if len(m) > 50 {
		return errTooMuchMetadataKeys
	}

	for k := range m {
		if l := len(k); l > 40 {
			return fmt.Errorf("metadata key %q exceed 40 characters: %d", k, l)
		}
	}

	return nil
}

func defaultTime(t time.Time) time.Time {
	if t == (time.Time{}) {
		t = time.Now().UTC()
	}
	return t
}

func defaultIdempotencyKey(key string) string {
	if key == "" {
		return uuid.New().String()
	}
	return key
}