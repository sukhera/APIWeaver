package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetFileExtension returns the file extension (without the dot)
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		return ext[1:] // Remove the dot
	}
	return ""
}

// IsMarkdownFile checks if a file is a Markdown file
func IsMarkdownFile(filename string) bool {
	ext := strings.ToLower(GetFileExtension(filename))
	return ext == "md" || ext == "markdown"
}

// IsYAMLFile checks if a file is a YAML file
func IsYAMLFile(filename string) bool {
	ext := strings.ToLower(GetFileExtension(filename))
	return ext == "yaml" || ext == "yml"
}

// IsJSONFile checks if a file is a JSON file
func IsJSONFile(filename string) bool {
	ext := strings.ToLower(GetFileExtension(filename))
	return ext == "json"
}

// ReadFileWithLimit reads a file with a size limit
func ReadFileWithLimit(filename string, maxSize int64) ([]byte, error) {
	// Check file size first
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if info.Size() > maxSize {
		return nil, fmt.Errorf("file size %d exceeds limit %d", info.Size(), maxSize)
	}

	return os.ReadFile(filename)
}

// WriteFileAtomic writes a file atomically by writing to a temp file first
func WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	// Create temp file in the same directory
	dir := filepath.Dir(filename)
	tempFile, err := os.CreateTemp(dir, "temp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	tempName := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(tempName) // Clean up if something goes wrong
	}()

	// Write data to temp file
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Sync to ensure data is written
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// Close temp file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Set permissions
	if err := os.Chmod(tempName, perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Atomically rename temp file to target
	if err := os.Rename(tempName, filename); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// SafeFileName sanitizes a filename by removing/replacing unsafe characters
func SafeFileName(filename string) string {
	// Replace unsafe characters with underscore
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// Remove leading/trailing spaces and dots
	result = strings.Trim(result, " .")
	
	// Ensure filename is not empty
	if result == "" {
		result = "unnamed"
	}
	
	return result
}

// GetRelativePath gets the relative path from base to target
func GetRelativePath(base, target string) (string, error) {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	
	return filepath.Rel(absBase, absTarget)
}

// ListFiles lists files in a directory with optional extension filter
func ListFiles(dir string, extensions []string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			if len(extensions) == 0 {
				files = append(files, path)
			} else {
				ext := strings.ToLower(GetFileExtension(path))
				for _, allowedExt := range extensions {
					if ext == strings.ToLower(allowedExt) {
						files = append(files, path)
						break
					}
				}
			}
		}
		
		return nil
	})
	
	return files, err
}

// FileSize returns the size of a file
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// IsExecutable checks if a file is executable
func IsExecutable(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	
	return info.Mode()&0111 != 0
}