package ConfigManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// writeData writes the data to the JSON file.
//
// Parameters:
//   - elem: The data to write.
//
// Returns:
//   - error: An error if one occurred while writing the data.
func (jm *JSONManager[T]) writeData(elem T) error {
	dir := filepath.Dir(jm.loc)

	err := os.MkdirAll(dir, jm.dirPerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(jm.loc, os.O_WRONLY|os.O_CREATE, jm.filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(elem, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// openFile opens the JSON file and returns the data.
//
// Parameters:
//   - create: Whether to create the file if it doesn't exist.
//
// Returns:
//   - *T: The data from the JSON file.
//   - error: An error if one occurred while opening the file.
func (jm *JSONManager[T]) openFile(create bool) (T, error) {
	dir := filepath.Dir(jm.loc)

	res := jm.data.Default().(T)

	err := os.MkdirAll(dir, jm.dirPerm)
	if err != nil {
		return res, err
	}

	file, err := os.OpenFile(jm.loc, os.O_RDONLY|os.O_CREATE, jm.filePerm)
	if err != nil {
		return res, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return res, fmt.Errorf("could not get file info: %w", err)
	}

	size := fileInfo.Size()

	if size == 0 {
		if !create {
			return res, fmt.Errorf("file is empty")
		}

		err = jm.writeData(res)
		if err != nil {
			return res, fmt.Errorf("could not write data: %w", err)
		}

		return res, nil
	}

	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		return res, fmt.Errorf("could not read file: %w", err)
	}

	err = json.Unmarshal(data, res)
	if err != nil {
		return res, fmt.Errorf("could not unmarshal data: %w", err)
	}

	return res, nil
}
