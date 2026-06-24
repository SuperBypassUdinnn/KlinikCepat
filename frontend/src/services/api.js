/**
 * services/api.js
 *
 * API Fetcher — menangani seluruh komunikasi HTTP ke Go Backend.
 * Semua request diarahkan ke backend Go, BUKAN langsung ke Supabase.
 *
 * Base URL menggunakan prefix /api/v1 sesuai routing di backend Go (go-chi).
 */
import { supabase } from "./supabaseClient";

const BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

/**
 * Helper fetch wrapper dengan JSON parsing dan error handling.
 */
async function request(endpoint, options = {}) {
  const url = `${BASE_URL}${endpoint}`;

  const defaultHeaders = {
    "Content-Type": "application/json",
  };

  // Ambil access token dari session Supabase yang aktif.
  const {
    data: { session },
    error: sessionError,
  } = await supabase.auth.getSession();

  if (sessionError) {
    throw new Error(`Gagal membaca sesi autentikasi: ${sessionError.message}`);
  }

  if (session?.access_token) {
    defaultHeaders.Authorization = `Bearer ${session.access_token}`;
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

// ─── Auth Profile ────────────────────────────────────────
export function getCurrentUser() {
  return request("/auth/me");
}

// ─── Klinik (Publik) ────────────────────────────────────
export function getClinics() {
  return request("/klinik");
}

export function getClinicById(id) {
  return request(`/klinik/${id}`);
}

// ─── Gejala (Publik) ────────────────────────────────────
export function getGejala() {
  return request("/gejala");
}

export function getGejalaById(id) {
  return request(`/gejala/${id}`);
}

// ─── Triage (Publik) ────────────────────────────────────
export function submitTriage(payload) {
  return request("/triage", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

// ─── Antrean Admin (Terproteksi JWT) ────────────────────
export function getQueue(status = "Menunggu") {
  const params = new URLSearchParams({
    status,
  });

  return request(`/admin/antrean?${params.toString()}`);
}

export function updateStatusAntrean(antreanId, status) {
  return request(`/admin/antrean/${antreanId}/status`, {
    method: "PUT",
    body: JSON.stringify({ status }),
  });
}

// ─── CRUD Klinik (Terproteksi JWT — Super Admin) ────────
export function createKlinik(payload) {
  return request("/klinik", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateKlinik(id, payload) {
  return request(`/klinik/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteKlinik(id) {
  return request(`/klinik/${id}`, {
    method: "DELETE",
  });
}

// ─── CRUD Gejala (Terproteksi JWT — Super Admin) ────────
export function createGejala(payload) {
  return request("/gejala", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateGejala(id, payload) {
  return request(`/gejala/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteGejala(id) {
  return request(`/gejala/${id}`, {
    method: "DELETE",
  });
}

// ─── Manajemen Admin Klinik — Superadmin ────────────────
export function createClinicAdmin(payload) {
  return request("/superadmin/admin-klinik", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

// ─── Tiket Publik ───────────────────────────────────────
export function getPublicTicket(publicToken) {
  return request(`/ticket/${encodeURIComponent(publicToken)}`);
}

export function checkPublicTicket(payload) {
  return request("/ticket/check", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export default {
  getCurrentUser,
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
  createClinicAdmin,
  getPublicTicket,
  checkPublicTicket,
};
