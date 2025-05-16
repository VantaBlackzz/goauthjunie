package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "securePassword123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Long password",
			password: "veryLongPasswordThatIsMoreThan72CharactersWhichIsBcryptLimitButShouldStillWork123456789012345678901234567890",
			wantErr:  true, // bcrypt has a 72-byte limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPassword() returned empty hash for password: %v", tt.password)
			}
			if !tt.wantErr && hash == tt.password {
				t.Errorf("HashPassword() returned unhashed password: %v", tt.password)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, create some hashed passwords for testing
	password1 := "securePassword123"
	password2 := "differentPassword456"

	hash1, err := HashPassword(password1)
	if err != nil {
		t.Fatalf("Failed to hash password for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: password2,
			hash:     hash1,
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalid_hash",
			want:     false,
		},
		{
			name:     "Empty hash",
			password: password1,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHashAndCheckIntegration tests the integration between HashPassword and CheckPasswordHash
func TestHashAndCheckIntegration(t *testing.T) {
	passwords := []string{
		"simple",
		"complex!@#$%^&*()",
		"with spaces and symbols !@#",
		"",
		"1234567890",
	}

	for _, password := range passwords {
		t.Run("Integration test for: "+password, func(t *testing.T) {
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Verify the correct password matches
			if !CheckPasswordHash(password, hash) {
				t.Errorf("CheckPasswordHash() failed to verify correct password: %v", password)
			}

			// Verify a different password doesn't match
			if password != "" && CheckPasswordHash(password+"different", hash) {
				t.Errorf("CheckPasswordHash() incorrectly verified wrong password")
			}
		})
	}
}
