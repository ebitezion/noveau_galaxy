package middleware_test

import (
	"testing"

	"github.com/ebitezion/backend-framework/internal/middleware"
)

func TestEncryptMessageArmored(t *testing.T) {
	type args struct {
		publicKey string
		message   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid Encryption",
			args: args{
				publicKey: `-----BEGIN PGP PUBLIC KEY BLOCK-----
...
-----END PGP PUBLIC KEY BLOCK-----`,
				message: "Hello, this is a test message.",
			},
			wantErr: false,
		},
		{
			name: "Empty Message",
			args: args{
				publicKey: `-----BEGIN PGP PUBLIC KEY BLOCK-----
...
-----END PGP PUBLIC KEY BLOCK-----`,
				message: "",
			},
			wantErr: false, // Empty message is valid
		},
		{
			name: "Invalid Public Key",
			args: args{
				publicKey: `-----BEGIN PGP PUBLIC KEY BLOCK-----
				...
				-----END PGP PUBLIC KEY BLOCK-----`,
				message: "Hello, this is a test message.",
			},
			wantErr: true, // Invalid public key should return an error
		},
		{
			name: "Invalid Message",
			args: args{
				publicKey: `-----BEGIN PGP PUBLIC KEY BLOCK-----
				...
				-----END PGP PUBLIC KEY BLOCK-----`,
				message: "-----BEGIN PGP PUBLIC KEY BLOCK-----", // Invalid message
			},
			wantErr: true, // Invalid message should return an error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := middleware.EncryptMessageArmored(tt.args.publicKey, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptMessageArmored() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncryptMessageArmored() = %v, want %v", got, tt.want)
			}
		})
	}
}
