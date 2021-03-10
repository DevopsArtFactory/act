package runner

import (
	"fmt"
	"testing"
)

func TestIsValidAddress(t *testing.T) {
	testData := []struct {
		input    string
		expected error
	}{
		{
			input:    "10.10.10.10/30",
			expected: nil,
		},
		{
			input:    "10.10.10.10/33",
			expected: fmt.Errorf("cidr base should be between 0 and 32: %s", "10.10.10.10/33"),
		},
		{
			input:    "123.251.129/30",
			expected: fmt.Errorf("wrong IP address: %s", "123.251.129/30"),
		},
		{
			input:    "123.421.0.0/32",
			expected: fmt.Errorf("each class number should be between 0 and 255: %s", "123.421.0.0/32"),
		},
	}

	for _, td := range testData {
		if err := IsValidAddress(td.input); (err != nil && err.Error() != td.expected.Error()) || (err == nil && err != td.expected) {
			t.Errorf(err.Error())
			t.Errorf("validation error: %s", td.input)
		}
	}
}
