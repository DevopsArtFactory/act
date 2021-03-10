package aws

import "testing"

func TestParseWebACLID(t *testing.T) {
	testData := []struct {
		input    string
		expected string
	}{
		{
			input:    "test-web-acl / 1234-1234-1234",
			expected: "1234-1234-1234",
		},
	}

	for _, td := range testData {
		if ParseWebACLID(td.input) != td.expected {
			t.Errorf("expected: %s, output: %s", td.expected, ParseWebACLID(td.input))
		}
	}
}
