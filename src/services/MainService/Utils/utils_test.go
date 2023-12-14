package Utils

import (
	"testing"
)

func TestUsernameValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected int
	}{
		{
			name:     "valid email",
			username: "alifakhary622@gmail.com",
			expected: 0,
		},
		{
			name:     "valid phone number",
			username: "+989123456789",
			expected: 2,
		},
		{
			name:     "valid phone number",
			username: "09123456789",
			expected: 4,
		},
		{
			name:     "invalid username",
			username: "invalid_username",
			expected: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := UsernameValidation(test.username)
			if actual != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestEmailRandomGenerator(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "valid email",
			expected: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := EmailRandomGenerator()
			if actual == "" {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		})
	}
}
