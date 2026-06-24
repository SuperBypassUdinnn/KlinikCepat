package models

import (
	"time"
)

// Klinik merepresentasikan tabel klinik di database
type Klinik struct {
	ID             string    `json:"id" db:"id"`
	NamaKlinik     string    `json:"nama_klinik" db:"nama_klinik"`
	Lat            float64   `json:"lat" db:"lat"`
	Lng            float64   `json:"lng" db:"lng"`
	KapasitasAktif int       `json:"kapasitas_aktif" db:"kapasitas_aktif"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Gejala merepresentasikan tabel katalog_gejala di database
type Gejala struct {
	ID           string    `json:"id" db:"id"`
	NamaGejala   string    `json:"nama_gejala" db:"nama_gejala"`
	BobotUrgensi int       `json:"bobot_urgensi" db:"bobot_urgensi"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Antrean merepresentasikan baris antrean di database
type Antrean struct {
	ID            string    `json:"id" db:"id"`
	KlinikID      string    `json:"klinik_id" db:"klinik_id"`
	NamaPasien    string    `json:"nama_pasien" db:"nama_pasien"`
	EmailPasien   string    `json:"email_pasien" db:"email_pasien"`
	PublicToken   string    `json:"public_token" db:"public_token"`
	KodeTiket     string    `json:"kode_tiket" db:"kode_tiket"`
	TotalSkor     int       `json:"total_skor" db:"total_skor"`
	StatusTriage  string    `json:"status_triage" db:"status_triage"`
	StatusAntrean string    `json:"status_antrean" db:"status_antrean"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// UserAccess menyimpan role aplikasi dan klinik
// yang diasosiasikan dengan user.
type UserAccess struct {
	UserID   string  `json:"user_id"`
	Role     string  `json:"role"`
	KlinikID *string `json:"klinik_id"`
}

// AuthMeResponse adalah respons endpoint GET /api/v1/auth/me.
type AuthMeResponse struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	KlinikID   *string `json:"klinik_id"`
	NamaKlinik *string `json:"nama_klinik"`
}

// CreateClinicAdminRequest adalah payload Superadmin
// untuk membuat akun awal Admin Klinik.
type CreateClinicAdminRequest struct {
	Email    string `json:"email"`
	KlinikID string `json:"klinik_id"`
}

// CreatedAuthUser merepresentasikan akun yang
// berhasil dibuat pada Supabase Auth.
type CreatedAuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// CreateClinicAdminResponse dikembalikan setelah
// akun Auth dan role Admin Klinik berhasil dibuat.
type CreateClinicAdminResponse struct {
	Message           string `json:"message"`
	UserID            string `json:"user_id"`
	Email             string `json:"email"`
	TemporaryPassword string `json:"temporary_password"`
	Role              string `json:"role"`
	KlinikID          string `json:"klinik_id"`
}

// Request & Response Payloads

// GejalaInput adalah sub-struct untuk menampung input skala dari pasien
type GejalaInput struct {
	GejalaID       string `json:"gejala_id"`
	SkalaKeparahan int    `json:"skala_keparahan"` // Skala 1-3
}

// TriageRequest adalah payload JSON yang dikirim Frontend saat pasien mendaftar
type TriageRequest struct {
	KlinikID    string        `json:"klinik_id"`
	NamaPasien  string        `json:"nama_pasien"`
	EmailPasien string        `json:"email_pasien"`
	Gejala      []GejalaInput `json:"gejala"`
}

// TriageResponse dikembalikan setelah perhitungan sukses
type TriageResponse struct {
	AntreanID    string `json:"antrean_id"`
	StatusTriage string `json:"status_triage"`
	TotalSkor    int    `json:"total_skor"`
	KodeTiket    string `json:"kode_tiket"`
	PublicToken  string `json:"public_token"`
	Pesan        string `json:"pesan"`
}

// PublicTicket adalah data tiket yang aman ditampilkan
// kepada pasien tanpa login.
type PublicTicket struct {
	PublicToken   string    `json:"public_token"`
	KodeTiket     string    `json:"kode_tiket"`
	NamaPasien    string    `json:"nama_pasien"`
	NamaKlinik    string    `json:"nama_klinik"`
	TotalSkor     int       `json:"total_skor"`
	StatusTriage  string    `json:"status_triage"`
	StatusAntrean string    `json:"status_antrean"`
	CreatedAt     time.Time `json:"created_at"`
}

// CheckTicketRequest adalah payload halaman cek tiket.
type CheckTicketRequest struct {
	KodeTiket string `json:"kode_tiket"`
	Email     string `json:"email"`
}
