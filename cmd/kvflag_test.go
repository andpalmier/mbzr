package cmd

import "testing"

func TestKVFlagSet(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		wantErr bool
	}{
		{"single pair", []string{"any_run=https://app.any.run/tasks/1"}, false},
		{"repeated key appends", []string{"any_run=https://a", "any_run=https://b"}, false},
		{"value containing =", []string{"links=https://x.tld/?a=b"}, false},
		{"unknown key", []string{"nope=1"}, true},
		{"missing =", []string{"any_run"}, true},
		{"empty value", []string{"any_run="}, true},
		{"empty key", []string{"=value"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newKVFlag("reference", referenceKeys)
			var err error
			for _, in := range tt.inputs {
				if err = k.Set(in); err != nil {
					break
				}
			}
			if tt.wantErr && err == nil {
				t.Fatalf("expected an error for %v", tt.inputs)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestKVFlagRepeatedKeyAppends(t *testing.T) {
	k := newKVFlag("reference", referenceKeys)
	for _, v := range []string{"any_run=https://a", "any_run=https://b", "twitter=https://t"} {
		if err := k.Set(v); err != nil {
			t.Fatal(err)
		}
	}
	if got := len(k.pairs["any_run"]); got != 2 {
		t.Errorf("any_run has %d values, want 2", got)
	}
	if got := len(k.pairs["twitter"]); got != 1 {
		t.Errorf("twitter has %d values, want 1", got)
	}
}

// contextValues sends "comment" as a string and everything else as a list,
// matching the shape the API documents.
func TestContextValuesShape(t *testing.T) {
	k := newKVFlag("context", contextKeys)
	for _, v := range []string{"comment=very nasty", "dropped_by_malware=Gozi", "dropped_by_malware=Emotet"} {
		if err := k.Set(v); err != nil {
			t.Fatal(err)
		}
	}
	got := contextValues(k)
	if _, ok := got["comment"].(string); !ok {
		t.Errorf("comment should be a string, got %T", got["comment"])
	}
	list, ok := got["dropped_by_malware"].([]string)
	if !ok {
		t.Fatalf("dropped_by_malware should be a list, got %T", got["dropped_by_malware"])
	}
	if len(list) != 2 {
		t.Errorf("dropped_by_malware has %d values, want 2", len(list))
	}
}

func TestIsDeliveryMethod(t *testing.T) {
	if !isDeliveryMethod("email_attachment") {
		t.Error("email_attachment should be valid")
	}
	if isDeliveryMethod("carrier_pigeon") {
		t.Error("carrier_pigeon should not be valid")
	}
}
