package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashedPassword("password")
	if err != nil {
		t.Errorf("hata hashing şifresi: %v", err)
	}
	if hash == "" {
		t.Error("hash'in boş olmaması bekleniyor")
	}
	if hash == "password" {
		t.Error("beklenen hash'in paroladan farklı olması")
	}
}

func TestComparePasswords (t *testing.T) {
	hash, err := HashedPassword("password")
	if err != nil {
		t.Errorf("hata hashing şifresi: %v", err)
	}
	if !ComparePasswords(hash, []byte("password")) {
		t.Errorf("beklenen parolanın hash ile eşleşmesi")
	}
	if ComparePasswords(hash, []byte("notpassword")) {
		t.Errorf("beklenen parolanın hash ile eşleşmemesi")
	}
}