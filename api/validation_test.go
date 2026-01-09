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
