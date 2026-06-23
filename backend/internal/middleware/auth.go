package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// UserContextKey digunakan untuk menyimpan data klaim JWT ke konteks HTTP Request
const UserContextKey contextKey = "user"

// JWTClaims menampung klaim esensial dari Supabase JWT Token
type JWTClaims struct {
	Sub      string  `json:"sub"` // UUID Pengguna di Supabase Auth
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	Exp      int64   `json:"exp"`
	KlinikID *string `json:"-"`
}

// supabaseJWTClaims digunakan saat parsing JWT asli dari Supabase.
type supabaseJWTClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`

	jwt.RegisteredClaims
}

type TokenVerifier func(
	token string,
) (*JWTClaims, error)

var (
	supabaseJWKS   keyfunc.Keyfunc
	supabaseJWKSMu sync.Mutex
)

// getSupabaseJWKS mengambil dan menyimpan public signing keys Supabase.
func getSupabaseJWKS() (keyfunc.Keyfunc, error) {
	supabaseJWKSMu.Lock()
	defer supabaseJWKSMu.Unlock()

	if supabaseJWKS != nil {
		return supabaseJWKS, nil
	}

	supabaseURL := strings.TrimRight(
		os.Getenv("SUPABASE_URL"),
		"/",
	)

	if supabaseURL == "" {
		return nil, errors.New(
			"server error: SUPABASE_URL tidak disetel",
		)
	}

	jwksURL := supabaseURL +
		"/auth/v1/.well-known/jwks.json"

	jwks, err := keyfunc.NewDefault(
		[]string{jwksURL},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal memuat JWKS Supabase: %w",
			err,
		)
	}

	supabaseJWKS = jwks

	return supabaseJWKS, nil
}

// AuthMiddleware mengamankan rute admin dengan memverifikasi JWT dari Supabase Auth
func AuthMiddleware(
	next http.Handler,
) http.Handler {
	return authMiddlewareWithVerifier(
		VerifySupabaseJWT,
	)(next)
}

func authMiddlewareWithVerifier(
	verifier TokenVerifier,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get(
					"Authorization",
				)

				if authHeader == "" {
					http.Error(
						w,
						`{"error": "Unauthorized: Authorization header tidak ditemukan"}`,
						http.StatusUnauthorized,
					)
					return
				}

				parts := strings.Fields(authHeader)

				if len(parts) != 2 ||
					!strings.EqualFold(
						parts[0],
						"Bearer",
					) {
					http.Error(
						w,
						`{"error": "Unauthorized: Format header Authorization harus 'Bearer <token>'"}`,
						http.StatusUnauthorized,
					)
					return
				}

				claims, err := verifier(parts[1])
				if err != nil {
					http.Error(
						w,
						`{"error": "Unauthorized: `+
							err.Error()+`"}`,
						http.StatusUnauthorized,
					)
					return
				}

				ctx := WithClaims(
					r.Context(),
					claims,
				)

				next.ServeHTTP(
					w,
					r.WithContext(ctx),
				)
			},
		)
	}
}

// VerifySupabaseJWT memverifikasi JWT Supabase
// menggunakan public signing key dari endpoint JWKS.
func VerifySupabaseJWT(
	tokenString string,
) (*JWTClaims, error) {
	jwks, err := getSupabaseJWKS()
	if err != nil {
		return nil, err
	}

	supabaseURL := strings.TrimRight(
		os.Getenv("SUPABASE_URL"),
		"/",
	)

	expectedIssuer := supabaseURL + "/auth/v1"

	return verifySupabaseJWTWithKeyfunc(
		tokenString,
		jwks.Keyfunc,
		expectedIssuer,
	)
}

// WithClaims menyimpan JWT claims ke request context.
func WithClaims(
	ctx context.Context,
	claims *JWTClaims,
) context.Context {
	return context.WithValue(
		ctx,
		UserContextKey,
		claims,
	)
}

// GetClaimsFromContext mengambil JWT claims
// yang sudah disimpan dalam request context.
func GetClaimsFromContext(
	ctx context.Context,
) (*JWTClaims, bool) {
	claims, ok := ctx.Value(
		UserContextKey,
	).(*JWTClaims)

	return claims, ok && claims != nil
}

func verifySupabaseJWTWithKeyfunc(
	tokenString string,
	keyFunc jwt.Keyfunc,
	expectedIssuer string,
) (*JWTClaims, error) {
	claims := &supabaseJWTClaims{}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		keyFunc,
		jwt.WithValidMethods(
			[]string{
				jwt.SigningMethodES256.Alg(),
			},
		),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithAudience("authenticated"),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"token tidak valid: %w",
			err,
		)
	}

	if !parsedToken.Valid {
		return nil, errors.New("token tidak valid")
	}

	if claims.Subject == "" {
		return nil, errors.New(
			"token tidak memiliki subject pengguna",
		)
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New(
			"token tidak memiliki waktu kedaluwarsa",
		)
	}

	return &JWTClaims{
		Sub:   claims.Subject,
		Email: claims.Email,
		Role:  claims.Role,
		Exp:   claims.ExpiresAt.Unix(),
	}, nil
}
