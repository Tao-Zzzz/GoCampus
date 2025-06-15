package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_Struct(t *testing.T) {
	user := &User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Nickname:  "TestUser",
		Avatar:    "http://example.com/avatar.png",
		CreatedAt: time.Now().Truncate(time.Millisecond),
	}

	tests := []struct {
		name     string
		user     *User
		wantJSON string
		wantErr  bool
	}{
		{
			name: "Valid user JSON serialization",
			user: user,
			wantJSON: `{
				"id":"user123",
				"email":"test@example.com",
				"password":"hashedpassword",
				"nickname":"TestUser",
				"avatar":"http://example.com/avatar.png",
				"created_at":"` + user.CreatedAt.Format(time.RFC3339Nano) + `"
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON serialization
			data, err := json.Marshal(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// Normalize JSON strings for comparison (remove whitespace)
				var got, want map[string]interface{}
				if err := json.Unmarshal(data, &got); err != nil {
					t.Errorf("json.Unmarshal(got) error = %v", err)
				}
				if err := json.Unmarshal([]byte(tt.wantJSON), &want); err != nil {
					t.Errorf("json.Unmarshal(want) error = %v", err)
				}
				for key, wantValue := range want {
					if gotValue, exists := got[key]; !exists || gotValue != wantValue {
						t.Errorf("json.Marshal() field %s = %v, want %v", key, gotValue, wantValue)
					}
				}
			}
		})
	}
}