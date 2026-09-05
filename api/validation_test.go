package api

import "testing"

func TestValidateSHA256(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{"valid lowercase", "88d862aeb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1", false},
		{"valid uppercase", "88D862AEB067278155C67A6D4E5BE927B36F08149C950D75A3A419EB20560AA1", false},
		{"too short", "88d8", true},
		{"not hex", "zzzzzzzzb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSHA256(tt.hash); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSHA256() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{"valid", "Emotet", false},
		{"valid complex", "Win.Emotet_1", false},
		{"empty", "", true},
		{"too long", "ThisTagIsWayTooLongAndShouldDefinitelyFailBecauseItExceedsTheLimitOfSixtyFourCharactersWhichIsArbitraryButGood", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTag(tt.tag); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	tests := []struct {
		name  string
		hash  string
		valid bool
	}{
		// get_info accepts all three, confirmed against the live API.
		{"sha256", "cd24da0069d1d15c28538a0e9cb9610d0a47d6e7b36f00ac380ddfd0362afb93", true},
		{"sha1", "bab94357d255c22ec55e60dc55745d58b4d7ef12", true},
		{"md5", "56589e5d713295415379b4622c1e74e2", true},
		{"39 hex chars", "bab94357d255c22ec55e60dc55745d58b4d7ef1", false},
		{"41 hex chars", "bab94357d255c22ec55e60dc55745d58b4d7ef123", false},
		{"non hex", "zzz94357d255c22ec55e60dc55745d58b4d7ef12", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHash(tt.hash)
			if tt.valid && err != nil {
				t.Errorf("ValidateHash(%q) = %v, want nil", tt.hash, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("ValidateHash(%q) = nil, want an error", tt.hash)
			}
		})
	}
}

func TestValidateSHA1(t *testing.T) {
	if err := ValidateSHA1("bab94357d255c22ec55e60dc55745d58b4d7ef12"); err != nil {
		t.Errorf("valid SHA1 rejected: %v", err)
	}
	if err := ValidateSHA1("bab94357"); err == nil {
		t.Error("short hash accepted as SHA1")
	}
}
