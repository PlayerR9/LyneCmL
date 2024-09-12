package file_manager

import (
	"slices"
	"strings"
)

// FileSettingsOption is a type that defines a file settings option.
//
// Parameters:
//   - fs: The file settings.
type FileSettingsOption func(fs *FileSettings)

// WithDir allows directories.
//
// Parameters:
//   - allow_dir: Whether to allow directories.
//
// Returns:
//   - FileSettingsOption: The file settings option.
func WithDir(allow_dir bool) FileSettingsOption {
	return func(fs *FileSettings) {
		fs.allow_dir = allow_dir
	}
}

// WithoutFile disallows files.
//
// Returns:
//   - FileSettingsOption: The file settings option.
func WithoutFile() FileSettingsOption {
	return func(fs *FileSettings) {
		fs.allow_file = false
	}
}

// WithFile allows files. Empty extensions are ignored. Leave empty
// to allow all extensions.
//
// Parameters:
//   - exts: The allowed extensions.
//
// Returns:
//   - FileSettingsOption: The file settings option.
func WithFileExts(exts ...string) FileSettingsOption {
	if len(exts) == 0 {
		return func(fs *FileSettings) {
			fs.allow_file = true
			fs.allowed_exts = nil
		}
	}

	unique := make([]string, 0, len(exts))

	for _, ext := range exts {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}

		pos, ok := slices.BinarySearch(unique, ext)
		if !ok {
			unique = slices.Insert(unique, pos, ext)
		}
	}

	unique = unique[:len(unique):len(unique)]

	return func(fs *FileSettings) {
		fs.allow_file = true
		fs.allowed_exts = unique
	}
}

// FileSettings is the settings for the file manager.
type FileSettings struct {
	// allow_dir is true if directories are allowed.
	allow_dir bool

	// allow_file is true if files are allowed.
	allow_file bool

	// allowed_exts is the list of allowed extensions.
	allowed_exts []string
}
