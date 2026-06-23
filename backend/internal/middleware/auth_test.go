package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestES256Token(
	t *testing.T,
	privateKey *ecdsa.PrivateKey,
	issuer string,
	expiresAt time.Time,
) string {
	t.Helper()

	claims := supabaseJWTClaims{
		Email: "admin@klinik.com",
		Role:  "authenticated",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "test-user-id",
			Issuer:  issuer,
			Audience: jwt.ClaimStrings{
				"authenticated",
			},
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
			ExpiresAt: jwt.NewNumericDate(
				expiresAt,
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		claims,
	)

	token.Header["kid"] = "test-key"

	signedToken, err := token.SignedString(
		privateKey,
	)
	if err != nil {
		t.Fatalf(
			"gagal membuat token test: %v",
			err,
		)
	}

	return signedToken
}

func testKeyFunc(
	publicKey *ecdsa.PublicKey,
) jwt.Keyfunc {
	return func(
		token *jwt.Token,
	) (interface{}, error) {
		return publicKey, nil
	}
}

func TestVerifySupabaseJWT(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(
		elliptic.P256(),
		rand.Reader,
	)
	if err != nil {
		t.Fatalf(
			"gagal membuat private key: %v",
			err,
		)
	}

	issuer := "https://test-project.supabase.co/auth/v1"

	t.Run("Valid Token", func(t *testing.T) {
		token := generateTestES256Token(
			t,
			privateKey,
			issuer,
			time.Now().Add(time.Hour),
		)

		claims, err :=
			verifySupabaseJWTWithKeyfunc(
				token,
				testKeyFunc(
					&privateKey.PublicKey,
				),
				issuer,
			)

		if err != nil {
			t.Fatalf(
				"gagal memverifikasi token valid: %v",
				err,
			)
		}

		if claims.Sub != "test-user-id" {
			t.Errorf(
				"Sub = %s, want test-user-id",
				claims.Sub,
			)
		}

		if claims.Email != "admin@klinik.com" {
			t.Errorf(
				"Email = %s, want admin@klinik.com",
				claims.Email,
			)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		token := generateTestES256Token(
			t,
			privateKey,
			issuer,
			time.Now().Add(-time.Hour),
		)

		_, err :=
			verifySupabaseJWTWithKeyfunc(
				token,
				testKeyFunc(
					&privateKey.PublicKey,
				),
				issuer,
			)

		if err == nil {
			t.Fatal(
				"token kedaluwarsa seharusnya ditolak",
			)
		}

		if !errors.Is(
			err,
			jwt.ErrTokenExpired,
		) {
			t.Errorf(
				"error = %v, want ErrTokenExpired",
				err,
			)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		otherKey, err := ecdsa.GenerateKey(
			elliptic.P256(),
			rand.Reader,
		)
		if err != nil {
			t.Fatalf(
				"gagal membuat key kedua: %v",
				err,
			)
		}

		token := generateTestES256Token(
			t,
			privateKey,
			issuer,
			time.Now().Add(time.Hour),
		)

		_, err =
			verifySupabaseJWTWithKeyfunc(
				token,
				testKeyFunc(
					&otherKey.PublicKey,
				),
				issuer,
			)

		if err == nil {
			t.Fatal(
				"signature tidak valid seharusnya ditolak",
			)
		}
	})

	t.Run("Invalid Format", func(t *testing.T) {
		_, err :=
			verifySupabaseJWTWithKeyfunc(
				"token-ngawur",
				testKeyFunc(
					&privateKey.PublicKey,
				),
				issuer,
			)

		if err == nil {
			t.Fatal(
				"token dengan format invalid harus ditolak",
			)
		}
	})
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("Authorized Success", func(t *testing.T) {
		verifier := func(
			token string,
		) (*JWTClaims, error) {
			if token != "valid-test-token" {
				return nil, errors.New(
					"token tidak valid",
				)
			}

			return &JWTClaims{
				Sub:   "test-user-id",
				Email: "admin@klinik.com",
				Role:  "authenticated",
				Exp: time.Now().
					Add(time.Hour).
					Unix(),
			}, nil
		}

		protectedHandler :=
			authMiddlewareWithVerifier(
				verifier,
			)(
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write(
							[]byte("Authorized!"),
						)
					},
				),
			)

		req := httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

		req.Header.Set(
			"Authorization",
			"Bearer valid-test-token",
		)

		rr := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf(
				"Status Code = %v, want %v",
				rr.Code,
				http.StatusOK,
			)
		}

		if rr.Body.String() != "Authorized!" {
			t.Errorf(
				"Body = %s, want Authorized!",
				rr.Body.String(),
			)
		}
	})
}
