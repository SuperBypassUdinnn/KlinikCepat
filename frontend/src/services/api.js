/**
 * services/api.js
 *
 * API Fetcher — menangani seluruh komunikasi HTTP ke Go Backend.
 * Semua request diarahkan ke backend Go, BUKAN langsung ke Supabase.
 *
 * Base URL menggunakan prefix /api/v1 sesuai routing di backend Go (go-chi).
 */

const BASE_URL = '/api/v1';

/**
 * Helper fetch wrapper dengan JSON parsing dan error handling.
 */
async function request(endpoint, options = {}) {
  const url = `${BASE_URL}${endpoint}`;

  const defaultHeaders = {
    'Content-Type': 'application/json',
  };

  // Tambahkan token auth jika tersedia di localStorage
  const token = localStorage.getItem('access_token');
  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`;
  }

  const config = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  };

  const response = await fetch(url, config);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));
    throw new Error(errorBody.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// ─── Klinik (Publik) ────────────────────────────────────
export function getClinics() {
  return request('/klinik');
}

export function getClinicById(id) {
  return request(`/klinik/${id}`);
}

// ─── Gejala (Publik) ────────────────────────────────────
export function getGejala() {
  return request('/gejala');
}

export function getGejalaById(id) {
  return request(`/gejala/${id}`);
}

// ─── Triage (Publik) ────────────────────────────────────
export function submitTriage(payload) {
  return request('/triage', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// ─── Antrean Admin (Terproteksi JWT) ────────────────────
export function getQueue(klinikId, status = 'Menunggu') {
  return request(`/admin/antrean?klinik_id=${klinikId}&status=${status}`);
}

export function updateStatusAntrean(antreanId, status) {
  return request(`/admin/antrean/${antreanId}/status`, {
    method: 'PUT',
    body: JSON.stringify({ status }),
  });
}

// ─── CRUD Klinik (Terproteksi JWT — Super Admin) ────────
export function createKlinik(payload) {
  return request('/klinik', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function updateKlinik(id, payload) {
  return request(`/klinik/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export function deleteKlinik(id) {
  return request(`/klinik/${id}`, {
    method: 'DELETE',
  });
}

// ─── CRUD Gejala (Terproteksi JWT — Super Admin) ────────
export function createGejala(payload) {
  return request('/gejala', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function updateGejala(id, payload) {
  return request(`/gejala/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export function deleteGejala(id) {
  return request(`/gejala/${id}`, {
    method: 'DELETE',
  });
}

export default {
  getClinics,
  getClinicById,
  getGejala,
  getGejalaById,
  submitTriage,
  getQueue,
  updateStatusAntrean,
  createKlinik,
  updateKlinik,
  deleteKlinik,
  createGejala,
  updateGejala,
  deleteGejala,
};
