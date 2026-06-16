import { useLocation, useNavigate, Navigate } from 'react-router-dom';
import Card from '../../components/Card';
import Badge from '../../components/Badge';
import Button from '../../components/Button';
import { FiHome } from 'react-icons/fi';
import './TicketAntrean.css';

const STATUS_EMOJI = {
  Merah: '🚨',
  Kuning: '⚠️',
  Hijau: '✅',
};

export default function TicketAntrean() {
  const location = useLocation();
  const navigate = useNavigate();
  const data = location.state;

  // Redirect jika tidak ada data (akses langsung URL)
  if (!data) {
    return <Navigate to="/" replace />;
  }

  const statusLower = data.status_triage?.toLowerCase() || 'hijau';
  const shortId = data.antrean_id?.slice(-8)?.toUpperCase() || '-';

  return (
    <div className="page-wrapper ticket-page">
      <Card className="ticket-card">
        <div className="card-body">
          {/* Status Icon */}
          <div className={`ticket-status-icon ${statusLower}`}>
            {STATUS_EMOJI[data.status_triage] || '📋'}
          </div>

          <h2 className="ticket-title">Tiket Antrean Digital</h2>
          <p className="ticket-subtitle">
            Pendaftaran berhasil! Simpan tiket ini.
          </p>

          {/* Badge Status */}
          <div style={{ marginBottom: 'var(--space-lg)' }}>
            <Badge status={statusLower} size="lg">
              Status {data.status_triage}
            </Badge>
          </div>

          {/* Antrean ID */}
          <div className="ticket-antrean-label">Nomor Antrean</div>
          <div className="ticket-antrean-id">#{shortId}</div>

          <div className="ticket-divider" />

          {/* Detail */}
          <div className="ticket-details">
            <div className="ticket-detail-row">
              <span className="ticket-detail-label">Nama Pasien</span>
              <span className="ticket-detail-value">{data.nama_pasien}</span>
            </div>
            <div className="ticket-detail-row">
              <span className="ticket-detail-label">Klinik</span>
              <span className="ticket-detail-value">{data.nama_klinik}</span>
            </div>
            <div className="ticket-detail-row">
              <span className="ticket-detail-label">Skor Urgensi</span>
              <span className="ticket-detail-value">{data.total_skor}</span>
            </div>
          </div>

          {/* Pesan */}
          <div className={`ticket-message ${statusLower}`}>
            {data.pesan}
          </div>

          {/* Kembali */}
          <Button
            variant="secondary"
            block
            onClick={() => navigate('/')}
          >
            <FiHome size={16} />
            Kembali ke Beranda
          </Button>
        </div>
      </Card>
    </div>
  );
}
