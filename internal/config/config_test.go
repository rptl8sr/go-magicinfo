package config

import (
	"regexp"
	"testing"
)

func TestParseMac(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid MAC address with dashes",
			input:    "ab-Cd-eF-00-11-22",
			expected: "ab-cd-ef-00-11-22",
		},
		{
			name:     "Valid MAC address with colon",
			input:    "ab:Cd:eF:00:11:22",
			expected: "ab-cd-ef-00-11-22",
		},
		{
			name:     "Valid MAC address without separator",
			input:    "abCdeF001122",
			expected: "ab-cd-ef-00-11-22",
		},
		{
			name:     "Invalid MAC address short",
			input:    "abCdeF00112",
			expected: "",
		},
		{
			name:     "Invalid MAC address empty",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid MAC address other symbols",
			input:    "1234567890gh",
			expected: "",
		},
	}

	re = regexp.MustCompile(`^([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})$`)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseMac(tc.input, re)
			if got != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, got)
			}
		})
	}
}
