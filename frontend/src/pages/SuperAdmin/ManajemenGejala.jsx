import { useState, useEffect } from 'react';
import { getGejala, createGejala, updateGejala, deleteGejala } from '../../services/api';
import Button from '../../components/Button';
import Badge from '../../components/Badge';
import Modal from '../../components/Modal';
import LoadingSpinner from '../../components/LoadingSpinner';
import { FiPlus, FiEdit2, FiTrash2 } from 'react-icons/fi';
import './SuperAdmin.css';

const EMPTY_FORM = { nama_gejala: '', bobot_urgensi: '' };

export default function ManajemenGejala() {
  const [gejalaList, setGejalaList] = useState([]);
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
      const data = await getGejala();
      setGejalaList(data || []);
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

  const openEdit = (gejala) => {
    setForm({
      nama_gejala: gejala.nama_gejala,
      bobot_urgensi: String(gejala.bobot_urgensi),
    });
    setEditingId(gejala.id);
    setShowForm(true);
  };

  const openDelete = (gejala) => {
    setDeleteTarget(gejala);
    setShowDelete(true);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    const payload = {
      nama_gejala: form.nama_gejala,
      bobot_urgensi: parseInt(form.bobot_urgensi, 10),
    };

    try {
      if (editingId) {
        await updateGejala(editingId, payload);
      } else {
        await createGejala(payload);
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
      await deleteGejala(deleteTarget.id);
      setShowDelete(false);
      setDeleteTarget(null);
      await fetchData();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  const getBobotBadge = (bobot) => {
    if (bobot >= 8) return 'merah';
    if (bobot >= 4) return 'kuning';
    return 'hijau';
  };

  if (loading) {
    return <LoadingSpinner text="Memuat katalog gejala..." />;
  }

  return (
    <div className="container page-wrapper">
      {/* Header */}
      <div className="sa-header">
        <div>
          <h1>Manajemen Gejala</h1>
          <p>Kelola katalog gejala dan kalibrasi bobot urgensi triage</p>
        </div>
        <Button variant="primary" onClick={openCreate} id="btn-tambah-gejala">
          <FiPlus size={16} />
          Tambah Gejala
        </Button>
      </div>

      {/* Error */}
      {error && (
        <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {/* Table */}
      {gejalaList.length > 0 ? (
        <div className="table-wrapper">
          <table className="table" id="gejala-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Nama Gejala</th>
                <th>Bobot Urgensi</th>
                <th>Level</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {gejalaList.map((gejala, index) => (
                <tr key={gejala.id}>
                  <td>{index + 1}</td>
                  <td style={{ fontWeight: 600 }}>{gejala.nama_gejala}</td>
                  <td>{gejala.bobot_urgensi}</td>
                  <td>
                    <Badge status={getBobotBadge(gejala.bobot_urgensi)}>
                      {gejala.bobot_urgensi >= 8
                        ? 'Kritis'
                        : gejala.bobot_urgensi >= 4
                        ? 'Sedang'
                        : 'Ringan'}
                    </Badge>
                  </td>
                  <td>
                    <div className="sa-table-actions">
                      <Button variant="ghost" size="sm" onClick={() => openEdit(gejala)}>
                        <FiEdit2 size={14} />
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => openDelete(gejala)}>
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
          <div className="queue-empty-icon">🩺</div>
          <h3>Belum Ada Gejala</h3>
          <p style={{ color: 'var(--color-gray-500)' }}>
            Klik tombol "Tambah Gejala" untuk menambahkan ke katalog.
          </p>
        </div>
      )}

      {/* Form Modal */}
      <Modal
        isOpen={showForm}
        onClose={() => setShowForm(false)}
        title={editingId ? 'Edit Gejala' : 'Tambah Gejala Baru'}
      >
        <form className="sa-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="form-nama-gejala">Nama Gejala</label>
            <input
              id="form-nama-gejala"
              type="text"
              className="form-input"
              placeholder="Contoh: Pendarahan Hebat"
              value={form.nama_gejala}
              onChange={(e) => setForm({ ...form, nama_gejala: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="form-bobot">Bobot Urgensi (1-10)</label>
            <input
              id="form-bobot"
              type="number"
              min="1"
              max="10"
              className="form-input"
              placeholder="10"
              value={form.bobot_urgensi}
              onChange={(e) => setForm({ ...form, bobot_urgensi: e.target.value })}
              required
            />
            <span className="form-hint">
              1 = Sangat Ringan, 10 = Kritis (otomatis status Merah)
            </span>
          </div>

          <div className="sa-form-actions">
            <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>
              Batal
            </Button>
            <Button variant="primary" type="submit" loading={submitting}>
              {editingId ? 'Simpan Perubahan' : 'Tambah Gejala'}
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
            Apakah Anda yakin ingin menghapus gejala{' '}
            <span className="item-name">{deleteTarget?.nama_gejala}</span>?
            <br />
            Tindakan ini tidak dapat dibatalkan dan akan mempengaruhi kalkulasi triage.
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
