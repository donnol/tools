package jwt

import (
	"math/rand"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	rand.Seed(time.Now().Unix())

	secret := []byte("Xadfdfoere2324212afasf34wraf090uadfafdIEJF039038")
	token := New(secret)
	id := rand.Int()
	r, err := token.Sign(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)

	userID, err := token.Verify(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userID)

	if userID != id {
		t.Fatalf("Bad userID, got %d\n", userID)
	}

	t.Run("VerifyEmpty", func(t *testing.T) {
		r, err := token.Verify("")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r)
	})
}
