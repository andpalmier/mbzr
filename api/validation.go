package api

import (
	"fmt"
	"regexp"
)

var (
	// SHA256 regex: exactly 64 hexadecimal characters
	sha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	// SHA1 regex: exactly 40 hexadecimal characters
	sha1Regex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
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

// ValidateSHA1 checks if a string is a valid SHA1 hash
func ValidateSHA1(hash string) error {
	if !sha1Regex.MatchString(hash) {
		return fmt.Errorf("invalid SHA1 hash format: must be 40 hexadecimal characters")
	}
	return nil
}

// ValidateHash checks if a string is a hash accepted by the get_info endpoint,
// which takes SHA256, SHA1 or MD5.
func ValidateHash(hash string) error {
	if sha256Regex.MatchString(hash) || sha1Regex.MatchString(hash) || md5Regex.MatchString(hash) {
		return nil
	}
	return fmt.Errorf("invalid hash format: must be SHA256 (64 hex), SHA1 (40 hex) or MD5 (32 hex)")
}

// ValidateTag checks if a tag is safe (alphanumeric, dash, underscore only)
func ValidateTag(tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("tag cannot be empty")
	}
	if len(tag) > 100 {
		return fmt.Errorf("tag too long: maximum 100 characters")
	}
	// Deliberately permissive. The API documents [A-Za-z0-9.- ] for tags you
	// submit, but says nothing about tags you query; rejecting a tag here that
	// the API would have accepted turns a working query into a client error.
	// Malformed tags come back as illegal_tag, which is reported to the user.
	for _, r := range tag {
		if !isValidTagChar(r) {
			return fmt.Errorf("tag contains invalid characters: only alphanumeric, dash, underscore, dot, and space allowed")
		}
	}
	return nil
}

func isValidTagChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_' || r == ' ' || r == '.'
}
