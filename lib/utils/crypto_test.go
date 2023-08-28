package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
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

func TestGenerateCerts(t *testing.T) {
	type args struct {
		domain   string
		basePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test1",
			args{"test1.com", "/tmp"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateCerts(tt.args.domain, tt.args.basePath); (err != nil) != tt.wantErr {
				t.Errorf("GenerateCerts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	// Test case 1: Valid password
	password := "myStrongPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the hashed password is a valid bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		t.Errorf("Hashed password does not match the original password: %v", err)
	}

	// Test case 2: Empty password
	emptyPassword := ""
	hashedEmptyPassword, err := HashPassword(emptyPassword)
	if err != nil {
		t.Errorf("Expected no error for empty password, but got: %v", err)
	}
	if hashedEmptyPassword == "" {
		t.Errorf("Hashed empty password should not be empty")
	}
}

func TestCompareHashWithPassword(t *testing.T) {
	// Generate a hash from a password
	password := "myStrongPassword123"
	hashedPassword := SHA256(password)

	// Test case 1: Hash matches the original password
	if !CompareHashWithPassword(hashedPassword, password) {
		t.Error("Hash and password should match")
	}

	// Test case 2: Hash does not match the original password
	if CompareHashWithPassword(hashedPassword, "wrongPassword") {
		t.Error("Hash and wrong password should not match")
	}

	// Test case 3: Empty password
	if CompareHashWithPassword(hashedPassword, "") {
		t.Error("Hash and empty password should not match")
	}
}

func TestEncodePassword(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"hello", "world"}, "6672216e1a15e4d0e489a04864ed16b0cd1f7c1fbbc72543e996114e00f405eb7ed910252fcad3d7e3c5d608f2c9f59bb988"},
		{"test2", args{"world", "hello"}, "6bfd3de829e82922ce057ab178af48c0daaa1966d025e1a9c06989e2d45a355ecc87c1c052565ad58b31f7ed5898fa92fbe0"},
		{"test3", args{"", ""}, "0d60c695f7e95db4f9557f5b2410e77b296801a818bd99ffb8d3dd4972b28c5088341941719ec8689f071e9f4452427c6ad7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodePassword(tt.args.password, tt.args.salt); got != tt.want {
				t.Errorf("EncodePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"hello"}, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"test2", args{"world"}, "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{"test3", args{""}, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256(tt.args.text); got != tt.want {
				t.Errorf("SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}
