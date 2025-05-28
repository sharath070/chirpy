package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "wrongPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct Password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Wrong Password",
			password: password2,
			hash:     hash2,
			wantErr:  false,
		},
		{
			name:     "Password doesn't match hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid Hash",
			password: password1,
			hash:     "invalid hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
