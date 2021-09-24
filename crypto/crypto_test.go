package crypto_test

import (
	"testing"

	"github.com/donnol/tools/crypto"
)

func TestCrypto(t *testing.T) {
	for _, cas := range []struct {
		key   string
		value string
	}{
		{"example key 1234", "10.00"},
		{"example key 1234", "10.50"},
		{"example key 1234", "10.05"},
		{"example key 123456789000", "10.00"},
		{"example key 123456789000", "10.50"},
		{"example key 123456789000", "10.05"},
		{"example key 123456789000 abcdefg", "10.00"},
		{"example key 123456789000 abcdefg", "10.50"},
		{"example key 123456789000 abcdefg", "10.05"},
		{"example key 123456789000 abcdefg", ""},
	} {
		c, err := crypto.NewCrypto(cas.key)
		if err != nil {
			t.Fatal(err)
		}
		r, err := c.Encrypt(cas.value)
		if err != nil {
			t.Fatal(err)
		}

		money, err := c.Decrypt(r)
		if err != nil {
			t.Fatal(err)
		}
		if money != cas.value {
			t.Fatal("Bad decode\n")
		}
	}
}
