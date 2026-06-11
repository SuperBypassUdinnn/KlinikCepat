/**
 * services/api.js
 *
 * API Fetcher — menangani seluruh komunikasi HTTP ke Go Backend.
 * Semua request diarahkan ke backend Go, BUKAN langsung ke Supabase.
 */

const BASE_URL = '/api';

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

// ─── Triage ─────────────────────────────────────────────
export function submitTriage(payload) {
  return request('/triage', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// ─── Queue / Antrean ────────────────────────────────────
export function getQueue(clinicId) {
  return request(`/clinics/${clinicId}/queue`);
}

// ─── Clinics / Klinik ───────────────────────────────────
export function getClinics() {
  return request('/clinics');
}

export function getClinicById(id) {
  return request(`/clinics/${id}`);
}

export default {
  submitTriage,
  getQueue,
  getClinics,
  getClinicById,
};
