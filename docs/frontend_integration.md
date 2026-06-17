# Panduan Integrasi Frontend - Autentikasi Supabase & API Backend

Dokumen ini berisi panduan teknis bagi tim Frontend untuk melakukan integrasi sistem autentikasi menggunakan Supabase dan cara berkomunikasi dengan API Backend kita.

## 1. Persiapan Kredensial Supabase

Tim Frontend akan membutuhkan `Supabase URL` dan `Anon Key` untuk berinteraksi dengan layanan Supabase dari sisi klien (browser).

> [!WARNING]
> **JANGAN PERNAH** menggunakan atau menyimpan kunci `service_role` (secret) di kode Frontend. Kunci tersebut memiliki akses admin penuh yang bisa mem-bypass semua aturan keamanan. Hanya gunakan `Anon Key` (public) di Frontend.

Dapatkan kredensial ini dari Admin Proyek (yang memiliki akses Dasbor Supabase) atau lihat di file `.env` Frontend jika sudah disediakan:
- `VITE_SUPABASE_URL`: (contoh: `https://xxxxxx.supabase.co`)
- `VITE_SUPABASE_ANON_KEY`: (contoh: `eyJhbGciOiJIUzI1NiIsInR5c...`)

## 2. Instalasi SDK Supabase

Pastikan modul `@supabase/supabase-js` sudah terinstal di proyek Frontend.

```bash
npm install @supabase/supabase-js
```

## 3. Inisialisasi Klien Supabase

Buat file konfigurasi (misal: `src/lib/supabase.js` atau `src/config/supabase.js`) untuk melakukan inisialisasi klien Supabase. Klien ini akan digunakan di seluruh aplikasi untuk urusan autentikasi.

```javascript
import { createClient } from '@supabase/supabase-js'

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

export const supabase = createClient(supabaseUrl, supabaseAnonKey)
```

## 4. Contoh Implementasi Login di Frontend

Gunakan fungsi `signInWithPassword` dari klien Supabase untuk melakukan login menggunakan Email dan Password.

```javascript
// Contoh pemanggilan login di komponen React/Vue
import { supabase } from '../lib/supabase'

const handleLogin = async (email, password) => {
  const { data, error } = await supabase.auth.signInWithPassword({
    email: email,
    password: password,
  })

  if (error) {
    console.error('Login gagal:', error.message)
    // Tampilkan pesan error ke user (misal: Email/Password salah)
    return null
  }
  
  console.log('Login sukses!', data.user)
  // Lakukan redirect ke halaman dashboard atau simpan state login
  return data.session
}
```

Fungsi pendukung lainnya yang mungkin dibutuhkan:
- Logout: `await supabase.auth.signOut()`
- Cek Sesi Aktif: `await supabase.auth.getSession()`
- Mendengarkan Perubahan Auth (Login/Logout): `supabase.auth.onAuthStateChange((event, session) => { ... })`

## 5. Mengakses Backend API (Sangat Penting!)

Backend Go kita menggunakan **custom middleware** untuk memverifikasi token JWT dari Supabase. Ini berarti **setiap kali Frontend memanggil endpoint Backend yang dilindungi**, Frontend WAJIB melampirkan Token Akses JWT (Access Token) ke dalam *header* `Authorization`.

### Langkah-langkah Memanggil Backend API:

1. Dapatkan token sesi saat ini dari Supabase.
2. Sisipkan token tersebut ke format `Bearer <token>`.
3. Kirimkan dalam objek *headers* pada fungsi `fetch` atau `axios`.

### Contoh menggunakan `fetch`:

```javascript
import { supabase } from '../lib/supabase'

const fetchUserProfile = async () => {
  // 1. Ambil sesi saat ini dari Supabase
  const { data: { session } } = await supabase.auth.getSession()
  
  // Pastikan user sudah login
  if (!session) {
    console.error('User belum login!')
    return
  }

  const token = session.access_token

  // 2. Lakukan request ke backend Go (misal: GET /api/users/me)
  try {
    const response = await fetch('http://localhost:8080/api/users/me', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        // PENTING: Masukkan token JWT ke header Authorization
        'Authorization': `Bearer ${token}` 
      }
    })

    if (!response.ok) {
      if (response.status === 401) {
        console.error('Token tidak valid atau sudah kadaluarsa (Unauthorized)')
        // Opsional: Lakukan proses refresh token atau paksa user re-login
      }
      throw new Error(`Error HTTP: ${response.status}`)
    }

    const result = await response.json()
    console.log('Data profil:', result)
    return result

  } catch (error) {
    console.error('Gagal mengambil data profil:', error)
  }
}
```

## 6. Tips Tambahan untuk Frontend

- **Interceptor Axios**: Jika Anda menggunakan Axios, sangat disarankan untuk membuat *interceptor* yang akan secara otomatis menyisipkan *header* `Authorization` ke setiap request ke backend, sehingga Anda tidak perlu menulisnya berulang-ulang di setiap pemanggilan fungsi API.
- **Auto Refresh Token**: Supabase JS Client secara otomatis menangani *refresh token* di belakang layar. Namun, pastikan Anda selalu mengambil token menggunakan `supabase.auth.getSession()` sebelum memanggil backend untuk memastikan token yang dikirim adalah yang terbaru (tidak kadaluarsa).
- **Backend URL**: Gunakan *environment variable* untuk Base URL Backend (contoh: `VITE_BACKEND_URL=http://localhost:8080`), jangan di-*hardcode*.
