package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// helper untuk membuat token JWT HS256 dummy
func generateTestJWT(secret string, claims JWTClaims) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadJSON, _ := json.Marshal(claims)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	message := headerEncoded + "." + payloadEncoded

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature, nil
}

func TestVerifySupabaseJWT(t *testing.T) {
	secret := "my-super-secret-test-key"
	_ = os.Setenv("SUPABASE_JWT_SECRET", secret)
	defer func() {
		_ = os.Unsetenv("SUPABASE_JWT_SECRET")
	}()

	t.Run("Valid Token", func(t *testing.T) {
		claims := JWTClaims{
			Sub:   "user-123",
			Email: "admin@klinikcepat.com",
			Role:  "authenticated",
			Exp:   time.Now().Add(1 * time.Hour).Unix(),
		}

		token, err := generateTestJWT(secret, claims)
		if err != nil {
			t.Fatalf("Gagal generate token: %v", err)
		}

		decodedClaims, err := VerifySupabaseJWT(token)
		if err != nil {
			t.Fatalf("Gagal memverifikasi token valid: %v", err)
		}

		if decodedClaims.Sub != claims.Sub {
			t.Errorf("Sub = %v, want %v", decodedClaims.Sub, claims.Sub)
		}
		if decodedClaims.Email != claims.Email {
			t.Errorf("Email = %v, want %v", decodedClaims.Email, claims.Email)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		claims := JWTClaims{
			Sub: "user-123",
			Exp: time.Now().Add(-1 * time.Hour).Unix(), // expired
		}

		token, err := generateTestJWT(secret, claims)
		if err != nil {
			t.Fatalf("Gagal generate token: %v", err)
		}

		_, err = VerifySupabaseJWT(token)
		if err == nil || !strings.Contains(err.Error(), "kedaluwarsa") {
			t.Errorf("Harus mengembalikan error kedaluwarsa, got: %v", err)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		claims := JWTClaims{
			Sub: "user-123",
			Exp: time.Now().Add(1 * time.Hour).Unix(),
		}

		token, err := generateTestJWT("wrong-secret-key", claims)
		if err != nil {
			t.Fatalf("Gagal generate token: %v", err)
		}

		_, err = VerifySupabaseJWT(token)
		if err == nil || !strings.Contains(err.Error(), "tanda tangan token tidak valid") {
			t.Errorf("Harus mengembalikan error tanda tangan tidak valid, got: %v", err)
		}
	})

	t.Run("Invalid Format", func(t *testing.T) {
		_, err := VerifySupabaseJWT("invalid.tokenformat")
		if err == nil || !strings.Contains(err.Error(), "format token tidak valid") {
			t.Errorf("Harus mengembalikan error format tidak valid, got: %v", err)
		}
	})
}

func TestAuthMiddleware(t *testing.T) {
	secret := "my-super-secret-test-key"
	_ = os.Setenv("SUPABASE_JWT_SECRET", secret)
	defer func() {
		_ = os.Unsetenv("SUPABASE_JWT_SECRET")
	}()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dapatkan klaim dari context
		claims, ok := r.Context().Value(UserContextKey).(*JWTClaims)
		if !ok || claims == nil {
			t.Error("Claims tidak ditemukan di context")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Authorized!"))
	})

	middleware := AuthMiddleware(nextHandler)

	t.Run("No Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/antrean", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Status Code = %v, want %v", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("Authorized Success", func(t *testing.T) {
		claims := JWTClaims{
			Sub: "user-123",
			Exp: time.Now().Add(1 * time.Hour).Unix(),
		}
		token, _ := generateTestJWT(secret, claims)

		req := httptest.NewRequest("GET", "/admin/antrean", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Status Code = %v, want %v", rr.Code, http.StatusOK)
		}
		if rr.Body.String() != "Authorized!" {
			t.Errorf("Body = %v, want %v", rr.Body.String(), "Authorized!")
		}
	})
}
