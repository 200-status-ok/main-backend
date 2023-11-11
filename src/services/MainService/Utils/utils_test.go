package Utils

import (
	"fmt"
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

func TestGetTime(t *testing.T) {
	tests := []struct {
		name          string
		location      string
		expectedError error
	}{
		{
			name:          "valid location",
			location:      "Asia/Tehran",
			expectedError: nil,
		},
		{
			name:          "invalid location",
			location:      "invalid_location",
			expectedError: fmt.Errorf("unknown time zone %s", "invalid_location"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetTime(test.location)
			if err != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", test.expectedError, err)
				}
			}
		})
	}
}
