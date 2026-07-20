// SPDX-FileCopyrightText: 2025 HolyBits LLC
//
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"testing"
)

const (
	clientIDKey = "client_id"
	valueKey    = "value"
)

// TestAppConnectionDetails verifies that appConnectionDetails correctly
// extracts the terraform client_id attribute and maps it to clientId
// in the connection details secret.
func TestAppConnectionDetails(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want map[string]string
	}{
		{
			name: "extracts client_id as clientId",
			attr: map[string]any{
				clientIDKey: "abc123-456-def",
			},
			want: map[string]string{
				"clientId": "abc123-456-def",
			},
		},
		{
			name: "empty map returns nothing",
			attr: map[string]any{},
			want: map[string]string{},
		},
		{
			name: "nil map returns nothing",
			attr: nil,
			want: map[string]string{},
		},
		{
			name: "client_id with wrong type (int) is ignored",
			attr: map[string]any{
				clientIDKey: 42,
			},
			want: map[string]string{},
		},
		{
			name: "client_id with wrong type (bool) is ignored",
			attr: map[string]any{
				clientIDKey: true,
			},
			want: map[string]string{},
		},
		{
			name: "extra keys are ignored, only client_id matters",
			attr: map[string]any{
				clientIDKey:    "abc123",
				"display_name": "test-app",
				"object_id":    "obj-456",
			},
			want: map[string]string{
				"clientId": "abc123",
			},
		},
		{
			name: "empty string client_id is still published",
			attr: map[string]any{
				clientIDKey: "",
			},
			want: map[string]string{
				"clientId": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appConnectionDetails(tt.attr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d keys, want %d keys", len(got), len(tt.want))
			}
			for k, wantV := range tt.want {
				if string(got[k]) != wantV {
					t.Errorf("key %q: got %q, want %q", k, string(got[k]), wantV)
				}
			}
		})
	}
}

// TestPwdConnectionDetails verifies that pwdConnectionDetails correctly
// extracts the terraform value attribute and maps it as "value" in the
// connection details secret.
func TestPwdConnectionDetails(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want map[string]string
	}{
		{
			name: "extracts value",
			attr: map[string]any{
				valueKey: "super-secret-password",
			},
			want: map[string]string{
				valueKey: "super-secret-password",
			},
		},
		{
			name: "empty map returns nothing",
			attr: map[string]any{},
			want: map[string]string{},
		},
		{
			name: "nil map returns nothing",
			attr: nil,
			want: map[string]string{},
		},
		{
			name: "value with wrong type (int) is ignored",
			attr: map[string]any{
				valueKey: 12345,
			},
			want: map[string]string{},
		},
		{
			name: "extra keys are ignored, only value matters",
			attr: map[string]any{
				valueKey:       "secret123",
				"display_name": "test-pwd",
				"key_id":       "key-789",
			},
			want: map[string]string{
				valueKey: "secret123",
			},
		},
		{
			name: "empty string value is still published",
			attr: map[string]any{
				valueKey: "",
			},
			want: map[string]string{
				valueKey: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pwdConnectionDetails(tt.attr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d keys, want %d keys", len(got), len(tt.want))
			}
			for k, wantV := range tt.want {
				if string(got[k]) != wantV {
					t.Errorf("key %q: got %q, want %q", k, string(got[k]), wantV)
				}
			}
		})
	}
}
