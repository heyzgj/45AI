//go:build ignore
// +build ignore

// This is a small helper script to generate a JWT for manual testing.
// Run it with `go run backend/create_token.go` if needed. It is excluded from
// regular `go vet` and build by the `ignore` build tag above.

package main; import ("fmt"; "time"; "github.com/golang-jwt/jwt/v5"); func main() { token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "1", "userID": 1, "exp": time.Now().Add(time.Hour * 24).Unix()}); tokenString, _ := token.SignedString([]byte("45ai_super_secret_jwt_key_for_development_only_change_in_production")); fmt.Print(tokenString) }
