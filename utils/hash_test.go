package utils

import "testing"

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string
		want bool
		got  func(*testing.T) bool
	}{
		{
			name: "it returns hashed password",
			want: true,
			got: func(t *testing.T) bool {
				hashedPassword, err := HashPassword("hashPassword")

				if err != nil {
					t.Fatal(err)
				}

				return CheckPassword("hashPassword", hashedPassword)
			},
		},
		{
			name: "it returns an emtpy string and error when hash function returns an error",
			want: false,
			got: func(t *testing.T) bool {
				hashedPassword, err := HashPassword(`01234567890123456789012345678901234567890123456789012345622211111111111111111111111111111111111:`)

				return hashedPassword == "" && err.Error() == "password length exceeds 72 bytes"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			want := test.want
			got := test.got(t)

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name string
		want bool
		got  func() bool
	}{
		{
			name: "It return false for non hashed password",
			want: false,
			got: func() bool {
				return CheckPassword("Test", "passwrod")
			},
		},
		{
			name: "It return false when password is match",
			want: false,
			got: func() bool {
				hashedPassword, _ := HashPassword("password")

				return CheckPassword("Test", hashedPassword)
			},
		},
		{
			name: "It return true when password and hash match",
			want: true,
			got: func() bool {
				hashedPassword, _ := HashPassword("password")

				return CheckPassword("password", hashedPassword)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			want := test.want
			got := test.got()

			if want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}
