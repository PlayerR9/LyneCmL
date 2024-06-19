package ConfigManager

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// JSONEncoder is an interface for encoding and decoding JSON data.
type JSONEncoder interface {
	json.Marshaler
	json.Unmarshaler

	// Default returns a new instance of the JSONEncoder.
	//
	// Returns:
	//   - JSONEncoder: A new instance of the JSONEncoder.
	Default() JSONEncoder
}

// JSONManager is a manager for saving and loading data to and from a JSON file.
type JSONManager[T JSONEncoder] struct {
	// data is the data to save and load.
	data T

	// loc is the location of the JSON file.
	loc string

	// dirPerm is the permissions of the directory.
	dirPerm os.FileMode

	// filePerm is the permissions of the file.
	filePerm os.FileMode

	// mu is a mutex for thread safety.
	mu sync.RWMutex
}

// Create creates a new JSON file at the location of the JSONManager[T].
//
// Returns:
//   - error: An error if one occurred while creating the file.
func (m *JSONManager[T]) Create() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	empty := m.data.Default().(T)

	err := m.writeData(empty)
	if err != nil {
		return err
	}

	return nil
}

// NewJSONManager creates a new JSONManager.
//
// Parameters:
//   - loc: The path to the JSON file.
//
// Returns:
//   - *JSONManager[T]: The new JSONManager.
func NewJSONManager[T JSONEncoder](loc string) *JSONManager[T] {
	return &JSONManager[T]{
		data:     (*new(T)).Default().(T),
		loc:      loc,
		dirPerm:  0777,
		filePerm: 0666,
	}
}

// SetDirPermissions sets the permissions of the directory.
//
// Parameters:
//   - perm: The permissions of the directory.
func (m *JSONManager[T]) SetDirPermissions(perm os.FileMode) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dirPerm = perm
}

// SetFilePermissions sets the permissions of the file.
//
// Parameters:
//   - perm: The permissions of the file.
func (m *JSONManager[T]) SetFilePermissions(perm os.FileMode) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dirPerm = perm
}

// ChangeData changes the data of the JSONManager[T].
//
// Parameters:
//   - data: The new data to save and load.
//
// Behaviors:
//   - If data is nil, a new instance of T is created.
func (m *JSONManager[T]) ChangeData(data T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = data
}

// ChangePath changes the path of the JSON file.
//
// Parameters:
//   - path: The new path to the JSON file.
func (m *JSONManager[T]) ChangePath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loc = path
}

// Load loads the data from the JSON file.
//
// The data is stored in the JSONManager[T] and can be accessed with GetData.
//
// Parameters:
//   - create: Whether to create the file if it doesn't exist.
//
// Returns:
//   - error: An error if there was an issue loading the data.
func (m *JSONManager[T]) Load(create bool) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, err := m.openFile(create)
	if err != nil {
		return *new(T), err
	}

	m.data = res

	return m.data, nil
}

// GetData returns the data of the JSONManager[T].
//
// Returns:
//   - *T: The data of the JSONManager[T].
//
// Behaviors:
//   - Never returns nil.
func (m *JSONManager[T]) GetData() T {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

// Save saves the data to the JSON file. It will overwrite the file if it already exists.
//
// Returns:
//   - error: An error if there was an issue saving the data.
func (m *JSONManager[T]) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.writeData(m.data)
	if err != nil {
		return err
	}

	return nil
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (m *JSONManager[T]) GetLocation() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.loc
}

// GetPermissions returns the permissions of the file.
//
// Returns:
//   - [2]os.FileMode: An array of os.FileMode representing the permissions of the directory and file.
func (m *JSONManager[T]) GetPermissions() [2]os.FileMode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return [2]os.FileMode{m.dirPerm, m.filePerm}
}

// Exists checks if a file exists at the location of the JSONManager[T].
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (m *JSONManager[T]) Exists() (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, err := os.Stat(m.loc)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// Delete deletes the file at the location.
//
// Returns:
//   - error: An error if one occurred while deleting the file.
func (m *JSONManager[T]) Delete() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return os.Remove(m.loc)
}
