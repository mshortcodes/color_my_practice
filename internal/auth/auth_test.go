package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password1 := "abc123"
	password2 := "def456"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	cases := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "incorrectpw",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "password doesn't match hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
