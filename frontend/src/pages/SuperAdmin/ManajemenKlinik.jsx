import { useState, useEffect } from 'react';
import { getClinics, createKlinik, updateKlinik, deleteKlinik } from '../../services/api';
import Button from '../../components/Button';
import Modal from '../../components/Modal';
import LoadingSpinner from '../../components/LoadingSpinner';
import { FiPlus, FiEdit2, FiTrash2 } from 'react-icons/fi';
import './SuperAdmin.css';

const EMPTY_FORM = { nama_klinik: '', lat: '', lng: '', kapasitas_aktif: '' };

export default function ManajemenKlinik() {
  const [clinics, setClinics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  // Modal state
  const [showForm, setShowForm] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);

  const fetchData = async () => {
    try {
      const data = await getClinics();
      setClinics(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openCreate = () => {
    setForm(EMPTY_FORM);
    setEditingId(null);
    setShowForm(true);
  };

  const openEdit = (clinic) => {
    setForm({
      nama_klinik: clinic.nama_klinik,
      lat: String(clinic.lat),
      lng: String(clinic.lng),
      kapasitas_aktif: String(clinic.kapasitas_aktif),
    });
    setEditingId(clinic.id);
    setShowForm(true);
  };

  const openDelete = (clinic) => {
    setDeleteTarget(clinic);
    setShowDelete(true);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    const payload = {
      nama_klinik: form.nama_klinik,
      lat: parseFloat(form.lat),
      lng: parseFloat(form.lng),
      kapasitas_aktif: parseInt(form.kapasitas_aktif, 10),
    };

    try {
      if (editingId) {
        await updateKlinik(editingId, payload);
      } else {
        await createKlinik(payload);
      }
      setShowForm(false);
      await fetchData();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    setSubmitting(true);
    try {
      await deleteKlinik(deleteTarget.id);
      setShowDelete(false);
      setDeleteTarget(null);
      await fetchData();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <LoadingSpinner text="Memuat data klinik..." />;
  }

  return (
    <div className="container page-wrapper">
      {/* Header */}
      <div className="sa-header">
        <div>
          <h1>Manajemen Klinik</h1>
          <p>Kelola data fasilitas kesehatan yang terdaftar di platform</p>
        </div>
        <Button variant="primary" onClick={openCreate} id="btn-tambah-klinik">
          <FiPlus size={16} />
          Tambah Klinik
        </Button>
      </div>

      {/* Error */}
      {error && (
        <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {/* Table */}
      {clinics.length > 0 ? (
        <div className="table-wrapper">
          <table className="table" id="klinik-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Nama Klinik</th>
                <th>Latitude</th>
                <th>Longitude</th>
                <th>Kapasitas</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {clinics.map((clinic, index) => (
                <tr key={clinic.id}>
                  <td>{index + 1}</td>
                  <td style={{ fontWeight: 600 }}>{clinic.nama_klinik}</td>
                  <td>{clinic.lat}</td>
                  <td>{clinic.lng}</td>
                  <td>{clinic.kapasitas_aktif}</td>
                  <td>
                    <div className="sa-table-actions">
                      <Button variant="ghost" size="sm" onClick={() => openEdit(clinic)}>
                        <FiEdit2 size={14} />
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => openDelete(clinic)}>
                        <FiTrash2 size={14} />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="queue-empty">
          <div className="queue-empty-icon">🏥</div>
          <h3>Belum Ada Klinik</h3>
          <p style={{ color: 'var(--color-gray-500)' }}>
            Klik tombol "Tambah Klinik" untuk mendaftarkan faskes baru.
          </p>
        </div>
      )}

      {/* Form Modal */}
      <Modal
        isOpen={showForm}
        onClose={() => setShowForm(false)}
        title={editingId ? 'Edit Klinik' : 'Tambah Klinik Baru'}
      >
        <form className="sa-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="form-nama-klinik">Nama Klinik</label>
            <input
              id="form-nama-klinik"
              type="text"
              className="form-input"
              placeholder="Klinik Sehat Selalu"
              value={form.nama_klinik}
              onChange={(e) => setForm({ ...form, nama_klinik: e.target.value })}
              required
            />
          </div>

          <div className="sa-coord-row">
            <div className="form-group">
              <label className="form-label" htmlFor="form-lat">Latitude</label>
              <input
                id="form-lat"
                type="number"
                step="any"
                className="form-input"
                placeholder="-6.2000"
                value={form.lat}
                onChange={(e) => setForm({ ...form, lat: e.target.value })}
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="form-lng">Longitude</label>
              <input
                id="form-lng"
                type="number"
                step="any"
                className="form-input"
                placeholder="106.8166"
                value={form.lng}
                onChange={(e) => setForm({ ...form, lng: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="form-kapasitas">Kapasitas Aktif</label>
            <input
              id="form-kapasitas"
              type="number"
              className="form-input"
              placeholder="50"
              value={form.kapasitas_aktif}
              onChange={(e) => setForm({ ...form, kapasitas_aktif: e.target.value })}
              required
            />
          </div>

          <div className="sa-form-actions">
            <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>
              Batal
            </Button>
            <Button variant="primary" type="submit" loading={submitting}>
              {editingId ? 'Simpan Perubahan' : 'Tambah Klinik'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={showDelete}
        onClose={() => setShowDelete(false)}
        title="Konfirmasi Hapus"
      >
        <div className="sa-delete-warning">
          <div className="sa-delete-warning-icon">⚠️</div>
          <p>
            Apakah Anda yakin ingin menghapus klinik{' '}
            <span className="item-name">{deleteTarget?.nama_klinik}</span>?
            <br />
            Tindakan ini tidak dapat dibatalkan.
          </p>
          <div className="sa-form-actions" style={{ justifyContent: 'center' }}>
            <Button variant="secondary" onClick={() => setShowDelete(false)}>
              Batal
            </Button>
            <Button variant="danger" onClick={handleDelete} loading={submitting}>
              Hapus
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
