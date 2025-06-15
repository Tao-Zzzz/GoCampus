package jwt

import (
	"context"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

func TestJWTUtil_GenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret")
	userID := "user123"

	token, err := jwtUtil.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if claims["user_id"] != userID {
			t.Errorf("GenerateToken() user_id = %v, want %v", claims["user_id"], userID)
		}
	} else {
		t.Error("Invalid token claims")
	}
}

func TestJWTUtil_ValidateTokenFromContext(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret")
	userID := "user123"

	// Generate a valid token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	// Create context with valid token
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+tokenString))

	tests := []struct {
		name    string
		ctx     context.Context
		wantID  string
		wantErr bool
	}{
		{
			name:    "Valid token",
			ctx:     ctx,
			wantID:  userID,
			wantErr: false,
		},
		{
			name:    "Missing metadata",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:    "Invalid token format",
			ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Invalid")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := jwtUtil.ValidateTokenFromContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTokenFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id != tt.wantID {
				t.Errorf("ValidateTokenFromContext() id = %v, want %v", id, tt.wantID)
			}
		})
	}
}