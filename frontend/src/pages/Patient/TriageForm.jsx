import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getGejala, getClinicById, submitTriage } from '../../services/api';
import Card from '../../components/Card';
import Button from '../../components/Button';
import LoadingSpinner from '../../components/LoadingSpinner';
import { FiCheck, FiArrowLeft, FiSend } from 'react-icons/fi';
import './TriageForm.css';

export default function TriageForm() {
  const { klinikId } = useParams();
  const navigate = useNavigate();

  const [gejalaList, setGejalaList] = useState([]);
  const [clinicName, setClinicName] = useState('');
  const [namaPasien, setNamaPasien] = useState('');
  const [selectedGejala, setSelectedGejala] = useState({}); // { gejala_id: skala_keparahan }
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);
  const [step, setStep] = useState(1); // 1: nama, 2: gejala

  // Ambil data gejala & nama klinik
  useEffect(() => {
    Promise.all([getGejala(), getClinicById(klinikId)])
      .then(([gejala, clinic]) => {
        setGejalaList(gejala || []);
        setClinicName(clinic?.nama_klinik || '');
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [klinikId]);

  const toggleGejala = (id) => {
    setSelectedGejala((prev) => {
      if (prev[id] !== undefined) {
        const next = { ...prev };
        delete next[id];
        return next;
      }
      return { ...prev, [id]: 1 }; // default skala 1
    });
  };

  const setSeverity = (id, skala) => {
    setSelectedGejala((prev) => ({ ...prev, [id]: skala }));
  };

  const handleSubmit = async () => {
    const selectedIds = Object.keys(selectedGejala);
    if (!namaPasien.trim()) {
      setError('Nama pasien wajib diisi.');
      return;
    }
    if (selectedIds.length === 0) {
      setError('Pilih minimal satu gejala.');
      return;
    }

    setSubmitting(true);
    setError(null);

    try {
      const payload = {
        klinik_id: klinikId,
        nama_pasien: namaPasien.trim(),
        gejala: selectedIds.map((id) => ({
          gejala_id: id,
          skala_keparahan: selectedGejala[id],
        })),
      };

      const result = await submitTriage(payload);

      // Navigasi ke halaman tiket dengan data hasil triage
      navigate('/ticket', {
        state: {
          ...result,
          nama_pasien: namaPasien.trim(),
          nama_klinik: clinicName,
        },
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <LoadingSpinner text="Memuat kuesioner gejala..." />;
  }

  const selectedCount = Object.keys(selectedGejala).length;

  return (
    <div className="container page-wrapper triage-page">
      {/* Header */}
      <div className="triage-header">
        <h1>Kuesioner Triage</h1>
        <p style={{ color: 'var(--color-gray-500)', marginBottom: 'var(--space-sm)' }}>
          Isi data diri dan pilih gejala yang Anda alami saat ini
        </p>
        {clinicName && (
          <span className="clinic-name-badge">🏥 {clinicName}</span>
        )}
      </div>

      {/* Step Indicator */}
      <div className="triage-step-indicator">
        <div className={`step-dot ${step >= 1 ? (step > 1 ? 'completed' : 'active') : 'inactive'}`}>
          {step > 1 ? <FiCheck size={16} /> : '1'}
        </div>
        <div className={`step-line ${step > 1 ? 'completed' : ''}`} />
        <div className={`step-dot ${step >= 2 ? 'active' : 'inactive'}`}>
          2
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {/* Step 1: Nama Pasien */}
      {step === 1 && (
        <Card className="triage-form-card">
          <div className="card-body">
            <h3 style={{ marginBottom: 'var(--space-lg)' }}>Data Diri Pasien</h3>

            <div className="form-group">
              <label className="form-label" htmlFor="nama-pasien">
                Nama Lengkap
              </label>
              <input
                id="nama-pasien"
                type="text"
                className="form-input"
                placeholder="Masukkan nama lengkap Anda"
                value={namaPasien}
                onChange={(e) => setNamaPasien(e.target.value)}
                autoFocus
              />
            </div>

            <div className="triage-form-actions">
              <Button
                variant="secondary"
                onClick={() => navigate('/')}
              >
                <FiArrowLeft size={16} />
                Kembali
              </Button>
              <Button
                variant="primary"
                block
                disabled={!namaPasien.trim()}
                onClick={() => {
                  setError(null);
                  setStep(2);
                }}
              >
                Lanjut — Pilih Gejala
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Step 2: Pilih Gejala */}
      {step === 2 && (
        <Card className="triage-form-card">
          <div className="card-body">
            <h3 style={{ marginBottom: 'var(--space-xs)' }}>Pilih Gejala yang Dialami</h3>
            <p style={{ fontSize: '0.875rem', marginBottom: 'var(--space-lg)' }}>
              Pilih satu atau lebih gejala, lalu tentukan tingkat keparahannya
            </p>

            <div className="gejala-list">
              {gejalaList.map((gejala) => {
                const isSelected = selectedGejala[gejala.id] !== undefined;
                const currentSeverity = selectedGejala[gejala.id] || 1;

                return (
                  <div
                    key={gejala.id}
                    className={`gejala-item ${isSelected ? 'selected' : ''}`}
                    id={`gejala-item-${gejala.id}`}
                  >
                    <div
                      className="gejala-checkbox"
                      onClick={() => toggleGejala(gejala.id)}
                    >
                      {isSelected && <FiCheck size={14} />}
                    </div>

                    <div className="gejala-info" onClick={() => toggleGejala(gejala.id)}>
                      <div className="gejala-nama">{gejala.nama_gejala}</div>
                      <div className={`gejala-bobot ${gejala.bobot_urgensi >= 8 ? 'high' : ''}`}>
                        Bobot urgensi: {gejala.bobot_urgensi}
                        {gejala.bobot_urgensi === 10 && ' ⚠ Kritis'}
                      </div>

                      {/* Severity selector (muncul jika dipilih) */}
                      {isSelected && (
                        <div className="severity-section" onClick={(e) => e.stopPropagation()}>
                          <div className="severity-label">Tingkat Keparahan:</div>
                          <div className="severity-options">
                            <button
                              className={`severity-btn ${currentSeverity === 1 ? 'active-1' : ''}`}
                              onClick={() => setSeverity(gejala.id, 1)}
                            >
                              1 — Ringan
                            </button>
                            <button
                              className={`severity-btn ${currentSeverity === 2 ? 'active-2' : ''}`}
                              onClick={() => setSeverity(gejala.id, 2)}
                            >
                              2 — Sedang
                            </button>
                            <button
                              className={`severity-btn ${currentSeverity === 3 ? 'active-3' : ''}`}
                              onClick={() => setSeverity(gejala.id, 3)}
                            >
                              3 — Berat
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>

            {selectedCount > 0 && (
              <p className="selected-count">
                <strong>{selectedCount}</strong> gejala dipilih
              </p>
            )}

            <div className="triage-form-actions">
              <Button variant="secondary" onClick={() => setStep(1)}>
                <FiArrowLeft size={16} />
                Kembali
              </Button>
              <Button
                variant="primary"
                block
                loading={submitting}
                disabled={selectedCount === 0}
                onClick={handleSubmit}
              >
                <FiSend size={16} />
                Daftar Antrean
              </Button>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
}
