package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

type contextKey string

// UserContextKey digunakan untuk menyimpan data klaim JWT ke konteks HTTP Request
const UserContextKey contextKey = "user"

// JWTClaims menampung klaim esensial dari Supabase JWT Token
type JWTClaims struct {
	Sub   string `json:"sub"` // UUID Pengguna di Supabase Auth
	Email string `json:"email"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

// AuthMiddleware mengamankan rute admin dengan memverifikasi JWT dari Supabase Auth
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Unauthorized: Authorization header tidak ditemukan"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error": "Unauthorized: Format header Authorization harus 'Bearer <token>'"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := VerifySupabaseJWT(token)
		if err != nil {
			http.Error(w, `{"error": "Unauthorized: `+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		// Menyimpan klaim user ke dalam context request
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifySupabaseJWT mendekode dan memverifikasi token JWT HS256 dari Supabase
func VerifySupabaseJWT(token string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("server error: SUPABASE_JWT_SECRET tidak disetel di server backend")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("format token tidak valid")
	}

	// 1. Verifikasi Tanda Tangan HS256 (HMAC-SHA256)
	message := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("gagal mendekode tanda tangan token")
	}

	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(message))
	expectedMac := mac.Sum(nil)

	if !hmac.Equal(sig, expectedMac) {
		return nil, errors.New("tanda tangan token tidak valid")
	}

	// 2. Dekode Payload JSON
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("gagal mendekode payload token")
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("gagal membaca isi klaim token")
	}

	// 3. Validasi Waktu Kedaluwarsa (exp)
	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token sudah kedaluwarsa (expired)")
	}

	return &claims, nil
}
