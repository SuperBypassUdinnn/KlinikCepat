package middleware

import (
	"net/http"
	"strings"

	"KlinikCepat/internal/repository"
)

// RequireRole memeriksa apakah user memiliki salah satu dari role yang diizinkan
func RequireRole(repo repository.RepositoryInterface, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil klaim user dari konteks
			claims, ok := r.Context().Value(UserContextKey).(*JWTClaims)
			if !ok || claims == nil {
				http.Error(w, `{"error": "Forbidden: Informasi user tidak ditemukan di request"}`, http.StatusForbidden)
				return
			}

			role := claims.Role
			klinikID := claims.KlinikID

			if role == "" || role == "authenticated" || (role == "klinik_admin" && klinikID == nil) {
				access, err := repo.GetUserAccess(
					r.Context(),
					claims.Sub,
				)
				if err != nil {
					http.Error(
						w,
						`{"error": "Internal Server Error: Gagal memeriksa hak akses"}`,
						http.StatusInternalServerError,
					)
					return
				}

				if access == nil || strings.TrimSpace(access.Role) == "" {
					http.Error(
						w,
						`{"error": "Forbidden: Anda belum memiliki role yang ditetapkan"}`,
						http.StatusForbidden,
					)
					return
				}

				role = access.Role
				klinikID = access.KlinikID
			}

			if role == "klinik_admin" &&
				(klinikID == nil ||
					strings.TrimSpace(*klinikID) == "") {
				http.Error(
					w,
					`{"error": "Forbidden: Akun admin belum terhubung ke klinik"}`,
					http.StatusForbidden,
				)
				return
			}

			// 3. Periksa apakah role ada dalam daftar allowedRoles
			isAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, `{"error": "Forbidden: Anda tidak memiliki akses ke rute ini"}`, http.StatusForbidden)
				return
			}

			// Simpan role dan klinik_id yang tervalidasi ke claims agar bisa digunakan oleh handler
			claims.Role = role
			claims.KlinikID = klinikID
			ctx := WithClaims(
				r.Context(),
				claims,
			)

			// 4. Lanjutkan ke handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
