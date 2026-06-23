# Integrasi Frontend KlinikCepat

Dokumen ini menjelaskan integrasi frontend React dengan Supabase Auth dan backend Go KlinikCepat.

**Pembaruan terakhir:** 24 Juni 2026

---

## 1. Arsitektur Integrasi

Frontend KlinikCepat menggunakan:

* React
* Vite
* React Router
* Supabase JavaScript Client
* Backend REST API berbasis Go

Alur komunikasinya:

```text
React Frontend
    │
    ├── Supabase Auth
    │       └── login, session, access token
    │
    └── Backend API
            ├── data klinik
            ├── katalog gejala
            ├── triage
            ├── antrean
            └── CRUD superadmin
```

Supabase Auth hanya menangani identitas dan session pengguna.

Role aplikasi seperti `superadmin` dan `klinik_admin` berasal dari tabel `user_roles` melalui backend.

---

## 2. Environment Variable

Frontend menggunakan environment variable berikut:

```env
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
VITE_API_URL=/api/v1
```

Keterangan:

* `VITE_SUPABASE_URL`
  URL project Supabase.

* `VITE_SUPABASE_ANON_KEY`
  Public anon key Supabase.

* `VITE_API_URL`
  Base URL backend API.

Untuk development lokal:

```env
VITE_API_URL=/api/v1
```

Vite proxy meneruskan request `/api` ke backend lokal.

Untuk deployment dengan domain backend terpisah:

```env
VITE_API_URL=https://api.example.com/api/v1
```

File `.env` yang berisi kredensial tidak boleh dimasukkan ke Git.

---

## 3. Supabase Client

Supabase client dibuat pada:

```text
frontend/src/services/supabaseClient.js
```

Contoh:

```js
import { createClient } from '@supabase/supabase-js';

const supabaseUrl =
  import.meta.env.VITE_SUPABASE_URL;

const supabaseAnonKey =
  import.meta.env.VITE_SUPABASE_ANON_KEY;

export const supabase = createClient(
  supabaseUrl,
  supabaseAnonKey,
);
```

Frontend hanya menggunakan anon key.

Service role key tidak boleh diletakkan di frontend.

---

## 4. API Service

Semua komunikasi dengan backend dipusatkan pada:

```text
frontend/src/services/api.js
```

Base URL:

```js
const BASE_URL =
  import.meta.env.VITE_API_URL || '/api/v1';
```

### 4.1 Session Token

Sebelum mengirim request, frontend mengambil session aktif dari Supabase:

```js
const {
  data: { session },
  error: sessionError,
} = await supabase.auth.getSession();
```

Jika access token tersedia, frontend menambahkan header:

```http
Authorization: Bearer <access_token>
```

Contoh:

```js
if (session?.access_token) {
  defaultHeaders.Authorization =
    `Bearer ${session.access_token}`;
}
```

Frontend tidak menyimpan salinan access token secara manual melalui:

```js
localStorage.setItem('access_token', token);
```

Supabase menjadi satu-satunya sumber session.

---

## 5. Request Helper

Contoh struktur helper request:

```js
async function request(endpoint, options = {}) {
  const url = `${BASE_URL}${endpoint}`;

  const headers = {
    'Content-Type': 'application/json',
  };

  const {
    data: { session },
    error: sessionError,
  } = await supabase.auth.getSession();

  if (sessionError) {
    throw new Error(
      `Gagal membaca sesi: ${sessionError.message}`,
    );
  }

  if (session?.access_token) {
    headers.Authorization =
      `Bearer ${session.access_token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...headers,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorBody = await response
      .json()
      .catch(() => ({}));

    throw new Error(
      errorBody.error ||
        `Request gagal dengan status ${response.status}`,
    );
  }

  return response.json();
}
```

---

## 6. Endpoint Publik

Endpoint berikut tidak membutuhkan login:

```http
GET /api/v1/klinik
GET /api/v1/klinik/{id}
GET /api/v1/gejala
GET /api/v1/gejala/{id}
POST /api/v1/triage
```

Contoh fungsi:

```js
export function getClinics() {
  return request('/klinik');
}

export function getGejala() {
  return request('/gejala');
}

export function submitTriage(payload) {
  return request('/triage', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
```

Walaupun request helper dapat menyertakan token jika session tersedia, endpoint publik tidak bergantung pada token.

---

## 7. Endpoint Profil Auth

Setelah session Supabase ditemukan, frontend memanggil:

```http
GET /api/v1/auth/me
```

Fungsi API:

```js
export function getCurrentUser() {
  return request('/auth/me');
}
```

Contoh respons Admin Klinik:

```json
{
  "id": "uuid-user",
  "email": "admin@klinik.com",
  "role": "klinik_admin",
  "klinik_id": "uuid-klinik"
}
```

Contoh respons Superadmin:

```json
{
  "id": "uuid-user",
  "email": "superadmin@klinik.com",
  "role": "superadmin",
  "klinik_id": null
}
```

Role aplikasi tidak diambil langsung dari metadata JWT frontend.

Backend mengambil role dari tabel `user_roles`.

---

## 8. AuthContext

State autentikasi dikelola pada:

```text
frontend/src/context/AuthContext.jsx
```

Context menyimpan:

```text
user
profile
role
clinicId
loading
authError
```

Keterangan:

* `user`
  User dari Supabase Auth.

* `profile`
  Respons lengkap dari `/auth/me`.

* `role`
  `superadmin` atau `klinik_admin`.

* `clinicId`
  ID klinik yang terkait dengan Admin Klinik.

* `loading`
  Status pemeriksaan session.

* `authError`
  Error saat sinkronisasi session dan role.

### 8.1 Sinkronisasi Session

Alurnya:

```text
Aplikasi dibuka
→ Supabase getSession()
→ session ditemukan
→ GET /api/v1/auth/me
→ simpan profile, role, dan clinicId
```

Jika session tidak tersedia:

```text
user = null
profile = null
role = null
clinicId = null
```

### 8.2 Login

Login dilakukan melalui:

```js
supabase.auth.signInWithPassword({
  email,
  password,
});
```

Setelah login berhasil, frontend menunggu respons `/auth/me`.

Fungsi `signIn` mengembalikan profile aplikasi:

```js
const currentUser = await signIn(
  email,
  password,
);
```

Redirect dilakukan berdasarkan role:

```js
if (currentUser.role === 'superadmin') {
  navigate('/superadmin/klinik');
}

if (currentUser.role === 'klinik_admin') {
  navigate('/admin/dashboard');
}
```

### 8.3 Logout

Logout dilakukan melalui:

```js
await supabase.auth.signOut();
```

Setelah logout, state autentikasi dibersihkan.

---

## 9. Role-Based Routing

Route dilindungi menggunakan komponen:

```text
frontend/src/components/ProtectedRoute.jsx
```

Contoh Admin Klinik:

```jsx
<ProtectedRoute
  allowedRoles={['klinik_admin']}
>
  <DashboardAdmin />
</ProtectedRoute>
```

Contoh Superadmin:

```jsx
<ProtectedRoute
  allowedRoles={['superadmin']}
>
  <ManajemenKlinik />
</ProtectedRoute>
```

### Perilaku Route

Pengguna tanpa session:

```text
→ /admin/login
```

Admin Klinik membuka route Superadmin:

```text
→ /admin/dashboard
```

Superadmin membuka dashboard Admin Klinik:

```text
→ /superadmin/klinik
```

Frontend route protection merupakan perlindungan UX.

Backend tetap menjadi sumber otorisasi utama.

---

## 10. Navigasi Berdasarkan Role

Navbar membaca `user` dan `role` dari AuthContext.

### Pengguna publik

Menu:

```text
Cari Klinik
Login Admin
```

### Admin Klinik

Menu:

```text
Dashboard
Halaman Pasien
Logout
```

### Superadmin

Menu:

```text
Kelola Klinik
Kelola Gejala
Halaman Pasien
Logout
```

Navbar tidak menentukan akses berdasarkan URL saja.

Menu ditampilkan berdasarkan role pengguna.

---

## 11. Integrasi Dashboard Admin Klinik

Admin Klinik mengambil antrean melalui:

```http
GET /api/v1/admin/antrean?status=Menunggu
```

Frontend tidak mengirimkan `klinik_id`.

Contoh fungsi:

```js
export function getQueue(
  status = 'Menunggu',
) {
  const params = new URLSearchParams({
    status,
  });

  return request(
    `/admin/antrean?${params.toString()}`,
  );
}
```

Backend menentukan `klinik_id` berdasarkan user yang login.

Frontend tidak menampilkan dropdown untuk memilih klinik.

### Update Status

```js
export function updateStatusAntrean(
  antreanId,
  status,
) {
  return request(
    `/admin/antrean/${antreanId}/status`,
    {
      method: 'PUT',
      body: JSON.stringify({ status }),
    },
  );
}
```

Backend memastikan antrean tersebut merupakan milik klinik Admin yang login.

---

## 12. Integrasi Superadmin

Superadmin dapat melakukan CRUD klinik:

```http
POST /api/v1/klinik
PUT /api/v1/klinik/{id}
DELETE /api/v1/klinik/{id}
```

Superadmin juga dapat melakukan CRUD gejala:

```http
POST /api/v1/gejala
PUT /api/v1/gejala/{id}
DELETE /api/v1/gejala/{id}
```

Request tersebut menggunakan token Supabase yang sama.

Backend memeriksa bahwa role pengguna adalah `superadmin`.

---

## 13. Penanganan Error

API service membaca error JSON dari backend:

```json
{
  "error": "Pesan error"
}
```

Status yang umum:

```text
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
500 Internal Server Error
```

Contoh arti:

* `401 Unauthorized`
  Token tidak tersedia, tidak valid, atau expired.

* `403 Forbidden`
  User valid tetapi role atau kliniknya tidak diizinkan.

* `404 Not Found`
  Data tidak ditemukan atau berada di luar scope tenant.

Frontend menampilkan pesan error melalui state komponen atau `authError`.

---

## 14. Keamanan

Frontend tidak boleh dianggap sebagai pengaman utama.

Pengguna masih dapat:

* mengubah JavaScript melalui DevTools;
* memodifikasi request;
* mencoba URL secara langsung;
* mengubah query parameter;
* memanggil backend tanpa UI.

Karena itu backend tetap memvalidasi:

* JWT;
* role;
* `klinik_id`;
* kepemilikan antrean;
* endpoint yang diperbolehkan.

Frontend hanya meningkatkan UX dengan menyembunyikan menu dan route yang tidak relevan.

---

## 15. Testing Frontend

Validasi build:

```bash
cd frontend
npm run build
```

Manual test:

```text
Pasien membuka halaman publik
Admin Klinik login
Admin Klinik diarahkan ke dashboard
Admin Klinik tidak dapat membuka route Superadmin
Superadmin login
Superadmin diarahkan ke halaman kelola klinik
Superadmin tidak dapat membuka dashboard Admin Klinik
Session bertahan setelah refresh
Logout membersihkan session
Navbar berubah berdasarkan role
```

---

## 16. Fitur yang Belum Diimplementasikan

Integrasi berikut masih berstatus Planned:

* login pasien;
* registrasi pasien;
* tiket permanen;
* polling status tiket pasien;
* live tracking posisi antrean;
* estimasi waktu tunggu;
* WebSocket;
* notifikasi real-time;
* dashboard analitik global;
* manajemen akun admin dari frontend;
* deployment production.