package ConfigManager

import (
	"encoding/json"
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Status is a type that represents the status of a set of keys.
type Status[T uc.Enumer] struct {
	// status is the status of the keys.
	status map[T]bool

	// isModified is a flag that indicates if the status has been modified.
	isModified bool

	// mu is a mutex that protects the status and isModified fields.
	mu sync.RWMutex
}

// MarshalJSON implements the json.Marshaler interface.
func (u *Status[T]) MarshalJSON() ([]byte, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	u.isModified = false

	return json.MarshalIndent(u.status, "", "  ")
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *Status[T]) UnmarshalJSON(data []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	var status map[T]bool

	err := json.Unmarshal(data, &status)
	if err != nil {
		return err
	}

	u.status = status
	u.isModified = false

	return nil
}

// NewStatus creates a new Status object.
//
// Parameters:
//   - keys: The keys of the status.
//
// Returns:
//   - *Status: The new Status object.
func NewStatus[T uc.Enumer](keys []T) *Status[T] {
	status := make(map[T]bool, len(keys))
	for _, key := range keys {
		status[key] = false
	}

	return &Status[T]{
		status:     status,
		isModified: false,
	}
}

// Change changes the status of a key.
//
// Parameters:
//   - key: The key to change.
//   - value: The new value of the key.
//
// Behaviors:
//   - If the key does not exist, the function does nothing.
func (u *Status[T]) Change(key T, value bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	prev, ok := u.status[key]
	if !ok || prev == value {
		return
	}

	u.status[key] = value
	u.isModified = true
}

// GetStatus returns the status of the keys.
//
// Returns:
//   - map[T]bool: The status of the keys.
func (u *Status[T]) GetStatus() map[T]bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.status
}

// GetValue returns the value of a key.
//
// Parameters:
//   - key: The key to get the value of.
//
// Returns:
//   - bool: The value of the key.
//
// Behaviors:
//   - If the key does not exist, the function returns false.
func (u *Status[T]) GetValue(key T) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	val, ok := u.status[key]
	if !ok {
		return false
	}

	return val
}

// IsModified returns true if the status has been modified.
//
// Returns:
//   - bool: True if the status has been modified, false otherwise.
func (u *Status[T]) IsModified() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.isModified
}
