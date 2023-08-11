package utils

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Test case 1: Generate a string of length 10
	n := uint(10)
	result := RandomString(n)
	if len(result) != int(n) {
		t.Errorf("Generated string length is not as expected: %d", n)
	}

	// Test case 2: Generate multiple strings and check for uniqueness
	iterations := 1000
	uniqueStrings := make(map[string]bool)
	for i := 0; i < iterations; i++ {
		result := RandomString(n)
		if len(result) != int(n) {
			t.Errorf("Generated string length is not as expected: %d", n)
		}
		// Check for uniqueness
		if _, exists := uniqueStrings[result]; exists {
			t.Errorf("Generated string is not unique: %s", result)
		}
		uniqueStrings[result] = true
	}

	// Test case 3: Generate a string of length 0
	n = uint(0)
	result = RandomString(n)
	if len(result) != int(n) {
		t.Errorf("Generated string length is not as expected: %d", n)
	}
}

func TestRandomChars(t *testing.T) {
	// Test generating a string with 10 random characters
	n := 10
	result, err := randomChars(n)
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}
	if len(result) != n {
		t.Errorf("Test case failed: expected string length %d, got %d", n, len(result))
	}

	// Test generating a string with 100 random characters
	n = 100
	result, err = randomChars(n)
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}
	if len(result) != n {
		t.Errorf("Test case failed: expected string length %d, got %d", n, len(result))
	}

	// Test generating a string with 0 random characters
	n = 0
	result, err = randomChars(n)
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}
	if len(result) != n {
		t.Errorf("Test case failed: expected string length %d, got %d", n, len(result))
	}

	// Test generating a string with a negative number of random characters, which should return an empty string
	n = -10
	result, err = randomChars(n)
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Test case failed: expected string length %d, got %d", 0, len(result))
	}
}
