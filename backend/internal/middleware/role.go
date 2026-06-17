package middleware

import (
	"context"
	"net/http"

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

			// 2. Ambil role dari database
			role, err := repo.GetUserRole(r.Context(), claims.Sub)
			if err != nil {
				http.Error(w, `{"error": "Internal Server Error: Gagal memeriksa hak akses"}`, http.StatusInternalServerError)
				return
			}

			if role == "" {
				http.Error(w, `{"error": "Forbidden: Anda belum memiliki role yang ditetapkan"}`, http.StatusForbidden)
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

			// Jika role sesuai, simpan role yang didapat dari DB ke context untuk digunakan lebih lanjut jika diperlukan
			claims.Role = role
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			
			// 4. Lanjutkan ke handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
