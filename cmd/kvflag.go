package cmd

import (
	"fmt"
	"sort"
	"strings"
)

// kvFlag collects repeatable "key=value" flags into key -> values.
// Repeating the same key appends, so a sample can carry several references
// of the same kind.
type kvFlag struct {
	allowed map[string]bool
	label   string
	pairs   map[string][]string
}

func newKVFlag(label string, allowed []string) *kvFlag {
	set := make(map[string]bool, len(allowed))
	for _, a := range allowed {
		set[a] = true
	}
	return &kvFlag{allowed: set, label: label, pairs: map[string][]string{}}
}

func (k *kvFlag) String() string {
	if len(k.pairs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(k.pairs))
	for key, vals := range k.pairs {
		for _, v := range vals {
			parts = append(parts, key+"="+v)
		}
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (k *kvFlag) Set(raw string) error {
	key, value, found := strings.Cut(raw, "=")
	if !found || key == "" || value == "" {
		return fmt.Errorf("expected %s in key=value form, got %q", k.label, raw)
	}
	if !k.allowed[key] {
		return fmt.Errorf("unknown %s key %q (allowed: %s)", k.label, key, strings.Join(k.keys(), ", "))
	}
	k.pairs[key] = append(k.pairs[key], value)
	return nil
}

func (k *kvFlag) keys() []string {
	out := make([]string, 0, len(k.allowed))
	for a := range k.allowed {
		out = append(out, a)
	}
	sort.Strings(out)
	return out
}

func (k *kvFlag) empty() bool { return len(k.pairs) == 0 }
