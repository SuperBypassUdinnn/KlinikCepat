import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getPublicTicket } from "../../services/api";
import Card from "../../components/Card";
import Badge from "../../components/Badge";
import Button from "../../components/Button";
import LoadingSpinner from "../../components/LoadingSpinner";
import { FiHome, FiSearch } from "react-icons/fi";
import "./TicketAntrean.css";

const STATUS_EMOJI = {
  Merah: "🚨",
  Kuning: "⚠️",
  Hijau: "✅",
};

const STATUS_MESSAGE = {
  Merah:
    "Kondisi darurat medis. Segera menuju fasilitas kesehatan untuk mendapatkan prioritas penanganan.",
  Kuning:
    "Kondisi membutuhkan perhatian medis. Silakan datang dan menunggu sesuai prioritas antrean.",
  Hijau:
    "Pendaftaran berhasil. Silakan datang dan menunggu sesuai urutan antrean.",
};

export default function TicketAntrean() {
  const { publicToken } = useParams();
  const navigate = useNavigate();

  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    let active = true;

    async function fetchTicket() {
      try {
        setLoading(true);
        setError(null);

        const result = await getPublicTicket(publicToken);

        if (active) {
          setData(result);
        }
      } catch (err) {
        if (active) {
          setError(err.message || "Tiket tidak ditemukan.");
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }

    fetchTicket();

    return () => {
      active = false;
    };
  }, [publicToken]);

  if (loading) {
    return <LoadingSpinner text="Memuat tiket antrean..." />;
  }

  if (error || !data) {
    return (
      <div className="page-wrapper ticket-page">
        <Card className="ticket-card">
          <div className="card-body">
            <h2>Tiket Tidak Ditemukan</h2>

            <p className="ticket-subtitle">
              {error || "Data tiket tidak tersedia."}
            </p>

            <div className="ticket-actions">
              <Button
                variant="primary"
                block
                onClick={() => navigate("/cek-tiket")}
              >
                <FiSearch size={16} />
                Cek Tiket
              </Button>

              <Button variant="secondary" block onClick={() => navigate("/")}>
                <FiHome size={16} />
                Beranda
              </Button>
            </div>
          </div>
        </Card>
      </div>
    );
  }

  const statusLower = data.status_triage?.toLowerCase() || "hijau";

  const formattedDate = new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(data.created_at));

  return (
    <div className="page-wrapper ticket-page">
      <Card className="ticket-card">
        <div className="card-body">
          <div className={`ticket-status-icon ${statusLower}`}>
            {STATUS_EMOJI[data.status_triage] || "📋"}
          </div>

          <h2 className="ticket-title">Tiket Antrean Digital</h2>

          <p className="ticket-subtitle">
            Simpan kode tiket dan tautan halaman ini.
          </p>

          <div
            style={{
              marginBottom: "var(--space-lg)",
            }}
          >
            <Badge status={statusLower} size="lg">
              Status {data.status_triage}
            </Badge>
          </div>

          <div className="ticket-antrean-label">Kode Tiket</div>

          <div className="ticket-antrean-id">{data.kode_tiket || "-"}</div>

          <div className="ticket-divider" />

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
              <span className="ticket-detail-label">Status Antrean</span>

              <span className="ticket-detail-value">{data.status_antrean}</span>
            </div>

            <div className="ticket-detail-row">
              <span className="ticket-detail-label">Waktu Daftar</span>

              <span className="ticket-detail-value">{formattedDate}</span>
            </div>

            <div className="ticket-detail-row">
              <span className="ticket-detail-label">Skor Urgensi</span>

              <span className="ticket-detail-value">{data.total_skor}</span>
            </div>
          </div>

          <div className={`ticket-message ${statusLower}`}>
            {STATUS_MESSAGE[data.status_triage]}
          </div>

          <div className="ticket-actions">
            <Button
              variant="primary"
              block
              onClick={() => navigate("/cek-tiket")}
            >
              <FiSearch size={16} />
              Cek Tiket Lain
            </Button>

            <Button variant="secondary" block onClick={() => navigate("/")}>
              <FiHome size={16} />
              Kembali ke Beranda
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}
