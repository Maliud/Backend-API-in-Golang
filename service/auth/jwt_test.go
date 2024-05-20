package auth

import "testing"

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Errorf("JWT oluşturulurken hata %v", err)
	}
	if token == "" {
		t.Error("token'ın boş olmaması Gerekiyor.")
	}
}