package file_manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcstr "github.com/PlayerR9/go-commons/strings"
)

// FileExists checks if a file exists.
//
// Parameters:
//   - path: The path of the file.
//   - opts: The file settings options.
//
// Returns:
//   - bool: True if the file exists, false otherwise.
//   - error: The error if any.
//
// By default, it will allow directories and files with any extension.
func FileExists(path string, opts ...FileSettingsOption) (bool, error) {
	settings := FileSettings{
		allow_dir:    true,
		allow_file:   true,
		allowed_exts: nil,
	}

	for _, opt := range opts {
		opt(&settings)
	}

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return settings.allow_dir, nil
	}

	if !settings.allow_file {
		return false, nil
	}

	if len(settings.allowed_exts) == 0 {
		return true, nil
	}

	name := stat.Name()
	ext := filepath.Ext(name)

	_, ok := slices.BinarySearch(settings.allowed_exts, ext)
	return ok, nil
}

// AddSuffixToFileName adds a suffix to a file name.
//
// Parameters:
//   - filename: The name of the file.
//   - new_suffix: The new suffix to add.
//   - ext: The extension of the file. If not provided, it will be inferred
//     from the filename.
//
// Returns:
//   - string: The new file name with the suffix added.
//
// This function returns an empty string if the filename is empty.
func AddSuffixToFileName(filename, new_suffix string, ext string) string {
	if filename == "" || new_suffix == "" {
		return filename
	}

	if ext == "" {
		ext = filepath.Ext(filename)

		if ext == "" {
			return filename + new_suffix
		}
	}

	loc := strings.TrimSuffix(filename, ext)

	var builder strings.Builder

	builder.WriteString(loc)
	builder.WriteString(new_suffix)
	builder.WriteString(ext)

	return builder.String()
}

// ErrIfInvalidExt returns an error if the file name does not have
// one of the given extensions.
//
// Parameters:
//   - file_name: The name of the file.
//   - exts: The extensions to check against.
//
// Returns:
//   - error: An error if the file name does not have one of the given extensions.
func ErrIfInvalidExt(file_name string, exts ...string) error {
	if file_name == "" {
		return gcers.NewErrInvalidParameter("file_name", errors.New("no file name provided"))
	} else if len(exts) == 0 {
		return gcers.NewErrInvalidParameter("exts", errors.New("no extensions provided"))
	}

	ext := filepath.Ext(file_name)

	if ext == "" {
		return errors.New("expected file, got directory instead")
	}

	for _, e := range exts {
		if ext == e {
			return nil
		}
	}

	gcstr.QuoteStrings(exts)

	return fmt.Errorf("invalid file extension: %s", gcstr.OrString(exts, true))
}

// ModifyPath modifies a path based on the given suffix and sub_directories.
//
// Parameters:
//   - path: The path to modify.
//   - suffix: The suffix to add.
//   - sub_directories: The sub directories to add.
//
// Returns:
//   - string: The modified path.
//   - error: An error if the path is not a file.
//
// This function returns an empty string if the path is empty.
func ModifyPath(path, suffix string, sub_directories ...string) (string, error) {
	if path == "" {
		return "", nil
	}

	if len(sub_directories) > 0 {
		dir, file := filepath.Split(path)
		path = filepath.Join(dir, filepath.Join(sub_directories...), file)
	}

	if suffix != "" {
		ext := filepath.Ext(path)

		if ext == "" {
			return "", errors.New("expected file, got directory instead")
		}

		path = strings.TrimSuffix(path, ext)
		path += suffix + ext
	}

	return path, nil
}

// ReadDir reads a directory and returns a list of file paths in the directory; includes
// sub directories.
//
// Parameters:
//   - loc: The location of the directory to read.
//
// Returns:
//   - []string: The list of file paths in the directory.
//   - error: An error if the directory could not be read.
func ReadDir(loc string) ([]string, error) {
	if loc == "" {
		return nil, gcers.NewErrInvalidParameter("loc", gcers.NewErrEmpty(loc))
	}

	var sols []string

	dirs, err := os.ReadDir(loc)
	if err != nil {
		return sols, err
	}

	for len(dirs) > 0 {
		dir := dirs[0]
		dirs = dirs[1:]

		dir_name := dir.Name()

		full_loc := filepath.Join(loc, dir_name)

		if !dir.IsDir() {
			sols = append(sols, full_loc)
		} else {
			new_dirs, err := os.ReadDir(full_loc)
			if err != nil {
				return sols, err
			}

			dirs = append(dirs, new_dirs...)
		}
	}

	return sols, nil
}

// fix_exts fixes the extension list by sorting and removing duplicates/invalid extensions.
//
// Parameters:
//   - exts: The extension list.
//
// Returns:
//   - []string: The fixed extension list.
func fix_exts(exts []string) []string {
	if len(exts) == 0 {
		return nil
	}

	new_exts := make([]string, 0, len(exts))

	for _, ext := range exts {
		if ext == "" || !strings.HasPrefix(ext, ".") {
			continue
		}

		pos, ok := slices.BinarySearch(new_exts, ext)
		if !ok {
			new_exts = slices.Insert(new_exts, pos, ext)
		}
	}

	return new_exts[:len(new_exts):len(new_exts)]
}

// FilterFiles returns a list of files that have one of the given extensions.
//
// Parameters:
//   - files: The list of files to filter.
//   - exts: The extensions to filter for.
//
// Returns:
//   - []string: The filtered list of files.
//
// If no extensions are valid or provided, then the function returns all the directories in the list.
func FilterFiles(files []string, exts ...string) []string {
	if len(files) == 0 {
		return nil
	}

	var top int

	exts = fix_exts(exts)
	if len(exts) == 0 {

		for i := 0; i < len(files); i++ {
			ext := filepath.Ext(files[i])

			if ext == "" {
				files[top] = files[i]
				top++
			}
		}
	} else {
		for i := 0; i < len(files); i++ {
			ext := filepath.Ext(files[i])

			_, ok := slices.BinarySearch(exts, ext)
			if ok {
				files[top] = files[i]
				top++
			}
		}
	}

	return files[:top:top]
}
