package utils

import (
	"testing"
)

func TestMD5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"hello"}, "5d41402abc4b2a76b9719d911017c592"},
		{"test2", args{"world"}, "7d793037a0760186574b0282f2f435e7"},
		{"test3", args{""}, "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.str); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	// Test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "aGVsbG8="},
		{"world", "d29ybGQ="},
		{"", ""},
	}

	// Iterate over test cases
	for _, test := range tests {
		// Call the function being tested
		output := Base64Encode(test.input)

		// Check if the output matches the expected value
		if output != test.expected {
			t.Errorf("Test case failed: input=%s, expected=%s, output=%s", test.input, test.expected, output)
		}
	}
}
