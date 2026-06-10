package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_decodeHashParams(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		encoded_hash string
		want         *argon2id_params
		want2        []byte
		want3        []byte
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "kek",
			encoded_hash: "$argon2id$v=19$m=69000,t=300,p=2$" + base64.RawStdEncoding.EncodeToString([]byte("sallty")) + "$" + base64.RawStdEncoding.EncodeToString([]byte("haaashy")),
			want: &argon2id_params{
				memory:      69000,
				iterations:  300,
				parallelism: 2,
				saltLength:  6,
				keyLength:   7,
			},
			want2:   []byte("sallty"),
			want3:   []byte("haaashy"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, got3, gotErr := decodeHashParams(tt.encoded_hash)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("decodeHashParams() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("decodeHashParams() succeeded unexpectedly")
			}
			if *got != *tt.want {
				t.Errorf("decodeHashParams() = %v, want %v", got, tt.want)
			}
			if !bytes.Equal(got2, tt.want2) {
				t.Errorf("decodeHashParams() = %v, want %v", got2, tt.want2)
			}
			if !bytes.Equal(got3, tt.want3) {
				t.Errorf("decodeHashParams() = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_authservice_HashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     string
	}{
		{
			name: "just to print default password hash",
			password: "password123",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s authservice
			s.hashParams = argon2id_params{
				memory:      256 * 1024,
				iterations:  3,
				parallelism: 1,
				saltLength:  16,
				keyLength:   32,
			}
			got := s.HashPassword(tt.password)
			fmt.Printf("t: %v\n", got)
		})
	}
}
