import { useEffect, useMemo, useState } from "react";
import {
  createClinicAdmin,
  getClinics,
} from "../../services/api";
import Button from "../../components/Button";
import LoadingSpinner from "../../components/LoadingSpinner";
import {
  FiCheckCircle,
  FiClipboard,
  FiUserPlus,
} from "react-icons/fi";
import "./ManajemenAdminKlinik.css";

const EMPTY_FORM = {
  email: "",
  klinik_id: "",
};

export default function ManajemenAdminKlinik() {
  const [clinics, setClinics] = useState([]);
  const [form, setForm] = useState(EMPTY_FORM);

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [error, setError] = useState(null);
  const [credentials, setCredentials] = useState(null);
  const [copiedField, setCopiedField] = useState(null);

  useEffect(() => {
    async function fetchClinics() {
      try {
        setError(null);

        const data = await getClinics();
        const clinicList = data || [];

        setClinics(clinicList);

        if (clinicList.length > 0) {
          setForm((currentForm) => ({
            ...currentForm,
            klinik_id:
              currentForm.klinik_id ||
              clinicList[0].id,
          }));
        }
      } catch (err) {
        setError(
          err.message ||
            "Gagal mengambil daftar klinik.",
        );
      } finally {
        setLoading(false);
      }
    }

    fetchClinics();
  }, []);

  const selectedClinic = useMemo(
    () =>
      clinics.find(
        (clinic) => clinic.id === form.klinik_id,
      ) || null,
    [clinics, form.klinik_id],
  );

  const replyTemplate = useMemo(() => {
    if (!credentials) {
      return "";
    }

    return `Halo,

Akun Admin Klinik Anda telah berhasil dibuat dengan rincian berikut:

Email: ${credentials.email}
Password sementara: ${credentials.temporary_password}
Klinik: ${credentials.nama_klinik}

Silakan login melalui halaman Admin Klinik menggunakan kredensial tersebut.

Untuk keamanan akun, segera ubah password setelah berhasil login. Jangan membagikan kredensial ini kepada pihak lain.

Terima kasih,
Superadmin KlinikCepat`;
  }, [credentials]);

  const handleChange = (event) => {
    const { name, value } = event.target;

    setForm((currentForm) => ({
      ...currentForm,
      [name]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!form.klinik_id) {
      setError("Klinik wajib dipilih.");
      return;
    }

    setSubmitting(true);
    setError(null);
    setCredentials(null);

    try {
      const result = await createClinicAdmin({
        email: form.email.trim(),
        klinik_id: form.klinik_id,
      });

      setCredentials({
        ...result,
        nama_klinik:
          selectedClinic?.nama_klinik ||
          "Klinik tidak diketahui",
      });

      setForm((currentForm) => ({
        ...currentForm,
        email: "",
      }));
    } catch (err) {
      setError(
        err.message ||
          "Gagal membuat akun Admin Klinik.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  const copyText = async (value, field) => {
    try {
      await navigator.clipboard.writeText(value);
      setCopiedField(field);

      window.setTimeout(() => {
        setCopiedField(null);
      }, 2000);
    } catch {
      setError(
        "Gagal menyalin teks. Salin secara manual.",
      );
    }
  };

  const clearCredentials = () => {
    setCredentials(null);
    setCopiedField(null);
  };

  if (loading) {
    return (
      <LoadingSpinner text="Memuat daftar klinik..." />
    );
  }

  return (
    <div className="container page-wrapper">
      <div className="sa-header">
        <div>
          <h1>Manajemen Admin Klinik</h1>
          <p>
            Buat akun awal dan kaitkan Admin Klinik
            dengan fasilitas kesehatan yang dikelolanya
          </p>
        </div>
      </div>

      {error && (
        <div
          className="alert alert-danger animate-fade-in"
          style={{ marginBottom: "1rem" }}
        >
          {error}
        </div>
      )}

      <div className="admin-management-grid">
        <section className="admin-management-card">
          <div className="admin-management-card-header">
            <FiUserPlus size={22} />

            <div>
              <h2>Buat Akun Admin</h2>
              <p>
                Password sementara akan dibuat
                otomatis oleh sistem.
              </p>
            </div>
          </div>

          {clinics.length === 0 ? (
            <div className="alert alert-warning">
              Belum ada klinik. Tambahkan klinik
              terlebih dahulu.
            </div>
          ) : (
            <form
              className="sa-form"
              onSubmit={handleSubmit}
            >
              <div className="form-group">
                <label
                  className="form-label"
                  htmlFor="admin-email"
                >
                  Email Admin
                </label>

                <input
                  id="admin-email"
                  name="email"
                  type="email"
                  className="form-input"
                  placeholder="admin@klinik.com"
                  value={form.email}
                  onChange={handleChange}
                  required
                />

                <span className="form-hint">
                  Email ini akan dipakai sebagai
                  username saat login.
                </span>
              </div>

              <div className="form-group">
                <label
                  className="form-label"
                  htmlFor="admin-klinik"
                >
                  Klinik yang Dikelola
                </label>

                <select
                  id="admin-klinik"
                  name="klinik_id"
                  className="form-input"
                  value={form.klinik_id}
                  onChange={handleChange}
                  required
                >
                  <option value="">
                    Pilih klinik
                  </option>

                  {clinics.map((clinic) => (
                    <option
                      key={clinic.id}
                      value={clinic.id}
                    >
                      {clinic.nama_klinik}
                    </option>
                  ))}
                </select>
              </div>

              <Button
                type="submit"
                variant="primary"
                loading={submitting}
                disabled={clinics.length === 0}
                block
              >
                <FiUserPlus size={16} />
                Buat Akun Admin Klinik
              </Button>
            </form>
          )}
        </section>

        <section className="admin-management-card">
          <div className="admin-management-card-header">
            <FiCheckCircle size={22} />

            <div>
              <h2>Kredensial Awal</h2>
              <p>
                Informasi ini hanya tampil setelah
                akun berhasil dibuat.
              </p>
            </div>
          </div>

          {!credentials ? (
            <div className="admin-credential-empty">
              <p>
                Belum ada akun baru yang dibuat.
              </p>

              <span>
                Email, password sementara, dan nama
                klinik akan muncul di sini.
              </span>
            </div>
          ) : (
            <div className="admin-credential-result">
              <div className="alert alert-success">
                {credentials.message}
              </div>

              <div className="credential-row">
                <span>Email</span>

                <div>
                  <strong>{credentials.email}</strong>

                  <button
                    type="button"
                    className="credential-copy-button"
                    onClick={() =>
                      copyText(
                        credentials.email,
                        "email",
                      )
                    }
                  >
                    <FiClipboard />
                    {copiedField === "email"
                      ? "Tersalin"
                      : "Salin"}
                  </button>
                </div>
              </div>

              <div className="credential-row">
                <span>Password sementara</span>

                <div>
                  <code>
                    {credentials.temporary_password}
                  </code>

                  <button
                    type="button"
                    className="credential-copy-button"
                    onClick={() =>
                      copyText(
                        credentials.temporary_password,
                        "password",
                      )
                    }
                  >
                    <FiClipboard />
                    {copiedField === "password"
                      ? "Tersalin"
                      : "Salin"}
                  </button>
                </div>
              </div>

              <div className="credential-row">
                <span>Klinik</span>
                <strong>
                  {credentials.nama_klinik}
                </strong>
              </div>

              <div className="form-group">
                <label
                  className="form-label"
                  htmlFor="reply-template"
                >
                  Template Balasan Email
                </label>

                <textarea
                  id="reply-template"
                  className="form-input admin-reply-template"
                  value={replyTemplate}
                  readOnly
                />
              </div>

              <div className="sa-form-actions">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={clearCredentials}
                >
                  Selesai
                </Button>

                <Button
                  type="button"
                  variant="primary"
                  onClick={() =>
                    copyText(
                      replyTemplate,
                      "template",
                    )
                  }
                >
                  <FiClipboard size={16} />
                  {copiedField === "template"
                    ? "Template Tersalin"
                    : "Salin Template Email"}
                </Button>
              </div>

              <p className="credential-security-note">
                Password tidak disimpan pada halaman ini.
                Setelah halaman ditutup atau di-refresh,
                password tidak dapat ditampilkan kembali.
              </p>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}