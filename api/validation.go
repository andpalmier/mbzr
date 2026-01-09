package api

import (
	"fmt"
	"regexp"
)

var (
	// SHA256 regex: exactly 64 hexadecimal characters
	sha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	// MD5 regex: exactly 32 hexadecimal characters
	md5Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
)

// ValidateSHA256 checks if a string is a valid SHA256 hash
func ValidateSHA256(hash string) error {
	if !sha256Regex.MatchString(hash) {
		return fmt.Errorf("invalid SHA256 hash format: must be 64 hexadecimal characters")
	}
	return nil
}

// ValidateMD5 checks if a string is a valid MD5 hash
func ValidateMD5(hash string) error {
	if !md5Regex.MatchString(hash) {
		return fmt.Errorf("invalid MD5 hash format: must be 32 hexadecimal characters")
	}
	return nil
}

// ValidateTag checks if a tag is safe (alphanumeric, dash, underscore only)
func ValidateTag(tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("tag cannot be empty")
	}
	if len(tag) > 100 {
		return fmt.Errorf("tag too long: maximum 100 characters")
	}
	// Allow alphanumeric, dash, underscore, and space
	for _, r := range tag {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ' || r == '.') {
			return fmt.Errorf("tag contains invalid characters: only alphanumeric, dash, underscore, dot, and space allowed")
		}
	}
	return nil
}
