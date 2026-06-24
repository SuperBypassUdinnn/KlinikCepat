import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { checkPublicTicket } from "../../services/api";
import Card from "../../components/Card";
import Button from "../../components/Button";
import { FiSearch, FiArrowLeft } from "react-icons/fi";
import "./CekTiket.css";

export default function CekTiket() {
  const navigate = useNavigate();

  const [kodeTiket, setKodeTiket] = useState("");
  const [email, setEmail] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (event) => {
    event.preventDefault();

    setSubmitting(true);
    setError(null);

    try {
      const ticket = await checkPublicTicket({
        kode_tiket: kodeTiket.trim().toUpperCase(),
        email: email.trim().toLowerCase(),
      });

      navigate(`/ticket/${ticket.public_token}`);
    } catch (err) {
      setError(err.message || "Kode tiket atau email tidak sesuai.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="container page-wrapper check-ticket-page">
      <Card className="check-ticket-card">
        <div className="card-body">
          <div className="check-ticket-icon">
            <FiSearch />
          </div>

          <h1>Cek Tiket Antrean</h1>

          <p className="check-ticket-subtitle">
            Masukkan kode tiket dan email yang digunakan saat pendaftaran.
          </p>

          {error && (
            <div
              className="alert alert-danger"
              style={{
                marginBottom: "1rem",
              }}
            >
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label" htmlFor="kode-tiket">
                Kode Tiket
              </label>

              <input
                id="kode-tiket"
                type="text"
                className="form-input"
                placeholder="KC-ABC123"
                value={kodeTiket}
                onChange={(event) =>
                  setKodeTiket(event.target.value.toUpperCase())
                }
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="email-tiket">
                Email
              </label>

              <input
                id="email-tiket"
                type="email"
                className="form-input"
                placeholder="nama@email.com"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                required
              />
            </div>

            <div className="check-ticket-actions">
              <Button
                type="button"
                variant="secondary"
                onClick={() => navigate("/")}
              >
                <FiArrowLeft size={16} />
                Kembali
              </Button>

              <Button
                type="submit"
                variant="primary"
                loading={submitting}
                block
              >
                <FiSearch size={16} />
                Tampilkan Tiket
              </Button>
            </div>
          </form>
        </div>
      </Card>
    </div>
  );
}
