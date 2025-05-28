package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	validToken, _ := MakeJWT(userId, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserId  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserId:  userId,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "Invalid.token.string",
			tokenSecret: "secret",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong secret",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserId, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotUserId != tt.wantUserId {
				t.Errorf("ValidateJWT() gotUserId = %v, want %v", gotUserId, tt.wantUserId)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "Valid Authorization",
			header:  http.Header{"Authorization": []string{"Bearer VALID"}},
			token:   "VALID",
			wantErr: false,
		}, {
			name:    "Marformed Authorization Header",
			header:  http.Header{"Authorization": []string{"InvalidBearer VALID"}},
			token:   "",
			wantErr: true,
		}, {
			name:    "Missing Authorization Header",
			header:  http.Header{},
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, want = %v", err, tt.wantErr)
			}
			if token != tt.token {
				t.Errorf("GetBearerToken() gotToken = %v, wantToken = %v", token, tt.token)
			}
		})
	}
}
