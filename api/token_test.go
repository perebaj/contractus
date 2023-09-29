package api

import "testing"

func TestAuthGenerateToken(t *testing.T) {
	a := Auth{
		JWTSecretKey: "secret",
	}
	got, err := a.GenerateToken("test@gmail.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZ21haWwuY29tIn0.Lk-U6IdaIHjtI3JJ6fbyBgc0xdjU5YeWGlThsye6Fss"

	assert(t, got, want)
}
