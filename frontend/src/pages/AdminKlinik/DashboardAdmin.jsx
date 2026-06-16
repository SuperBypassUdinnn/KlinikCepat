import { useState, useEffect, useCallback } from 'react';
import { getClinics, getQueue, updateStatusAntrean } from '../../services/api';
import Card from '../../components/Card';
import Badge from '../../components/Badge';
import Button from '../../components/Button';
import LoadingSpinner from '../../components/LoadingSpinner';
import { FiRefreshCw, FiCheck, FiSkipForward, FiAlertCircle, FiAlertTriangle, FiHeart, FiUsers } from 'react-icons/fi';
import './DashboardAdmin.css';

export default function DashboardAdmin() {
  const [clinics, setClinics] = useState([]);
  const [selectedClinic, setSelectedClinic] = useState('');
  const [queue, setQueue] = useState([]);
  const [statusFilter, setStatusFilter] = useState('Menunggu');
  const [loading, setLoading] = useState(true);
  const [queueLoading, setQueueLoading] = useState(false);
  const [actionLoading, setActionLoading] = useState(null);
  const [error, setError] = useState(null);

  // Ambil daftar klinik untuk dropdown
  useEffect(() => {
    getClinics()
      .then((data) => {
        setClinics(data || []);
        if (data?.length > 0) {
          setSelectedClinic(data[0].id);
        }
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  // Fetch queue data
  const fetchQueue = useCallback(async () => {
    if (!selectedClinic) return;
    setQueueLoading(true);
    try {
      const data = await getQueue(selectedClinic, statusFilter);
      setQueue(data || []);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setQueueLoading(false);
    }
  }, [selectedClinic, statusFilter]);

  // Fetch saat klinik atau filter berubah
  useEffect(() => {
    fetchQueue();
  }, [fetchQueue]);

  // Auto-refresh setiap 10 detik
  useEffect(() => {
    if (!selectedClinic) return;
    const interval = setInterval(fetchQueue, 10000);
    return () => clearInterval(interval);
  }, [selectedClinic, statusFilter, fetchQueue]);

  // Handle aksi status
  const handleUpdateStatus = async (antreanId, newStatus) => {
    setActionLoading(antreanId);
    try {
      await updateStatusAntrean(antreanId, newStatus);
      // Refresh queue
      await fetchQueue();
    } catch (err) {
      setError(err.message);
    } finally {
      setActionLoading(null);
    }
  };

  // Hitung statistik
  const stats = {
    total: queue.length,
    merah: queue.filter((q) => q.status_triage === 'Merah').length,
    kuning: queue.filter((q) => q.status_triage === 'Kuning').length,
    hijau: queue.filter((q) => q.status_triage === 'Hijau').length,
  };

  const formatTime = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleTimeString('id-ID', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return <LoadingSpinner text="Memuat dashboard..." />;
  }

  return (
    <div className="container page-wrapper">
      {/* Header */}
      <div className="dashboard-header">
        <div className="dashboard-header-info">
          <h1>Dashboard Antrean</h1>
          <p>Pantau dan kelola antrean pasien secara real-time</p>
        </div>
        <div className="dashboard-header-actions">
          <Button
            variant="outline"
            size="sm"
            onClick={fetchQueue}
            disabled={queueLoading}
          >
            <FiRefreshCw size={14} className={queueLoading ? 'spinning' : ''} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Klinik Selector */}
      <div className="klinik-selector">
        <label htmlFor="klinik-select">Klinik:</label>
        <select
          id="klinik-select"
          className="form-select"
          value={selectedClinic}
          onChange={(e) => setSelectedClinic(e.target.value)}
        >
          {clinics.map((c) => (
            <option key={c.id} value={c.id}>
              {c.nama_klinik}
            </option>
          ))}
        </select>
      </div>

      {/* Stats Cards */}
      {statusFilter === 'Menunggu' && (
        <div className="dashboard-stats">
          <Card className="stat-card">
            <div className="stat-card-header">
              <span className="stat-card-label">Total Menunggu</span>
              <div className="stat-card-icon total">
                <FiUsers size={18} />
              </div>
            </div>
            <div className="stat-card-value">{stats.total}</div>
          </Card>
          <Card className="stat-card" accent="merah">
            <div className="stat-card-header">
              <span className="stat-card-label">Merah</span>
              <div className="stat-card-icon merah">
                <FiAlertCircle size={18} />
              </div>
            </div>
            <div className="stat-card-value">{stats.merah}</div>
          </Card>
          <Card className="stat-card" accent="kuning">
            <div className="stat-card-header">
              <span className="stat-card-label">Kuning</span>
              <div className="stat-card-icon kuning">
                <FiAlertTriangle size={18} />
              </div>
            </div>
            <div className="stat-card-value">{stats.kuning}</div>
          </Card>
          <Card className="stat-card" accent="hijau">
            <div className="stat-card-header">
              <span className="stat-card-label">Hijau</span>
              <div className="stat-card-icon hijau">
                <FiHeart size={18} />
              </div>
            </div>
            <div className="stat-card-value">{stats.hijau}</div>
          </Card>
        </div>
      )}

      {/* Error */}
      {error && (
        <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {/* Queue Section */}
      <div className="queue-section">
        <div className="queue-section-header">
          <h2>Daftar Antrean</h2>
          <div className="auto-refresh-badge">
            <span className="auto-refresh-dot" />
            Auto-refresh 10 detik
          </div>
        </div>

        {/* Status Filter Tabs */}
        <div className="status-tabs">
          {['Menunggu', 'Selesai', 'Dilewati'].map((status) => (
            <button
              key={status}
              className={`status-tab ${statusFilter === status ? 'active' : ''}`}
              onClick={() => setStatusFilter(status)}
              id={`status-tab-${status.toLowerCase()}`}
            >
              {status}
            </button>
          ))}
        </div>

        {/* Queue Table */}
        {queueLoading && queue.length === 0 ? (
          <LoadingSpinner text="Memuat antrean..." />
        ) : queue.length > 0 ? (
          <div className="table-wrapper">
            <table className="table" id="queue-table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Nama Pasien</th>
                  <th>Status Triage</th>
                  <th>Skor</th>
                  <th>Waktu Daftar</th>
                  {statusFilter === 'Menunggu' && <th>Aksi</th>}
                </tr>
              </thead>
              <tbody>
                {queue.map((item, index) => (
                  <tr key={item.id}>
                    <td>{index + 1}</td>
                    <td style={{ fontWeight: 600 }}>{item.nama_pasien}</td>
                    <td>
                      <Badge status={item.status_triage?.toLowerCase()}>
                        {item.status_triage}
                      </Badge>
                    </td>
                    <td>{item.total_skor}</td>
                    <td>
                      <span className="queue-timestamp">
                        {formatTime(item.created_at)}
                      </span>
                    </td>
                    {statusFilter === 'Menunggu' && (
                      <td>
                        <div className="queue-actions">
                          <Button
                            variant="success"
                            size="sm"
                            loading={actionLoading === item.id}
                            onClick={() => handleUpdateStatus(item.id, 'Selesai')}
                            id={`btn-selesai-${item.id}`}
                          >
                            <FiCheck size={14} />
                            Selesai
                          </Button>
                          <Button
                            variant="secondary"
                            size="sm"
                            loading={actionLoading === item.id}
                            onClick={() => handleUpdateStatus(item.id, 'Dilewati')}
                            id={`btn-lewati-${item.id}`}
                          >
                            <FiSkipForward size={14} />
                            Lewati
                          </Button>
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="queue-empty">
            <div className="queue-empty-icon">📋</div>
            <h3>Tidak Ada Antrean</h3>
            <p style={{ color: 'var(--color-gray-500)' }}>
              Belum ada pasien dengan status "{statusFilter}" saat ini.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
