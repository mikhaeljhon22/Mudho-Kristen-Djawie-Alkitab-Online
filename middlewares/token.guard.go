package middlewares

import (
	"net/http"
	"strings"

	"github.com/kataras/jwt"
)

var sharedKey = []byte("sercrethatmaycontainch@r$32chars")

type TokenClaims struct {
	TokenClaims string `json:"tokenClaims"`
	UserID      int    `json:"userID"`
}

func JWTVerif(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Pastikan format "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Ambil token-nya
		jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
		jwtToken = strings.TrimSpace(jwtToken)
		if jwtToken == "" {
			http.Error(w, "Empty JWT token", http.StatusUnauthorized)
			return
		}

		// Verifikasi token
		_, err := jwt.Verify(jwt.HS256, sharedKey, []byte(jwtToken))
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
