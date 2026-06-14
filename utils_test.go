package blindindexstore

import "testing"

func Test_isHex(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected bool
	}{
		{
			name: "valid lowercase hex",
			input: "abcdef0123456789",
			expected: true,
		},
		{
			name: "valid uppercase hex",
			input: "ABCDEF0123456789",
			expected: true,
		},
		{
			name: "valid mixed case hex",
			input: "AbCdEf0123456789",
			expected: true,
		},
		{
			name: "valid 64 character hex (SHA256)",
			input: "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9",
			expected: true,
		},
		{
			name: "empty string",
			input: "",
			expected: true,
		},
		{
			name: "invalid characters - spaces",
			input: "abc def",
			expected: false,
		},
		{
			name: "invalid characters - special chars",
			input: "abc@def",
			expected: false,
		},
		{
			name: "invalid characters - letters g-z",
			input: "ghijklmnop",
			expected: false,
		},
		{
			name: "plain text",
			input: "SearchValue01",
			expected: false,
		},
		{
			name: "single valid hex char",
			input: "a",
			expected: true,
		},
		{
			name: "single invalid char",
			input: "g",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHex(tt.input)
			if result != tt.expected {
				t.Errorf("isHex(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
