import { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { getClinics } from '../../services/api';
import Card from '../../components/Card';
import Button from '../../components/Button';
import LoadingSpinner from '../../components/LoadingSpinner';
import { FiMapPin, FiNavigation, FiUsers } from 'react-icons/fi';
import './CariKlinik.css';

/**
 * Rumus Haversine — menghitung jarak (km) antara dua titik koordinat.
 */
function haversine(lat1, lon1, lat2, lon2) {
  const R = 6371;
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) ** 2;
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

export default function CariKlinik() {
  const navigate = useNavigate();
  const [clinics, setClinics] = useState([]);
  const [userPos, setUserPos] = useState(null);
  const [gpsError, setGpsError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Ambil data klinik dari backend
  useEffect(() => {
    getClinics()
      .then((data) => {
        setClinics(data || []);
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  // Minta lokasi GPS pengguna
  useEffect(() => {
    if (!navigator.geolocation) {
      setGpsError('Geolocation tidak didukung oleh browser ini.');
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setUserPos({
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        });
      },
      (err) => {
        setGpsError(err.message);
      }
    );
  }, []);

  // Hitung jarak & urutkan klinik dari terdekat
  const sortedClinics = useMemo(() => {
    if (!userPos) return clinics;

    return [...clinics]
      .map((clinic) => ({
        ...clinic,
        distance: haversine(userPos.lat, userPos.lng, clinic.lat, clinic.lng),
      }))
      .sort((a, b) => a.distance - b.distance);
  }, [clinics, userPos]);

  const handleSelectClinic = (clinicId) => {
    navigate(`/triage/${clinicId}`);
  };

  if (loading) {
    return <LoadingSpinner text="Memuat daftar klinik..." />;
  }

  return (
    <div className="container page-wrapper">
      {/* Hero */}
      <div className="cari-klinik-hero">
        <h1>Cari Klinik Terdekat</h1>
        <p>
          Pilih fasilitas kesehatan tujuan Anda untuk memulai pendaftaran
          antrean cerdas berbasis triage.
        </p>

        {/* GPS Status */}
        {userPos ? (
          <div className="gps-status active">
            <span className="gps-dot" />
            Lokasi GPS aktif
          </div>
        ) : (
          <div className="gps-status inactive">
            <span className="gps-dot" />
            {gpsError || 'Menunggu izin lokasi...'}
          </div>
        )}
      </div>

      {/* Error */}
      {error && (
        <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1.5rem' }}>
          Gagal memuat data: {error}
        </div>
      )}

      {/* Klinik Grid */}
      {sortedClinics.length > 0 ? (
        <div className="klinik-grid stagger-list">
          {sortedClinics.map((clinic) => (
            <Card
              key={clinic.id}
              accent="primary"
              className="klinik-card"
              onClick={() => handleSelectClinic(clinic.id)}
              id={`clinic-card-${clinic.id}`}
            >
              <div className="klinik-card-body">
                <div className="klinik-card-top">
                  <span className="klinik-name">{clinic.nama_klinik}</span>
                  {clinic.distance !== undefined && (
                    <span className="klinik-distance">
                      <FiNavigation size={12} />
                      {clinic.distance.toFixed(1)} km
                    </span>
                  )}
                </div>

                <div className="klinik-info-row">
                  <FiMapPin size={14} />
                  <span>
                    {clinic.lat.toFixed(4)}, {clinic.lng.toFixed(4)}
                  </span>
                </div>
              </div>

              <div className="klinik-card-footer">
                <div className="klinik-capacity">
                  <FiUsers size={14} />
                  <span>
                    Kapasitas:{' '}
                    <span className="capacity-value">
                      {clinic.kapasitas_aktif}
                    </span>
                  </span>
                </div>
                <Button variant="primary" size="sm">
                  Pilih Klinik
                </Button>
              </div>
            </Card>
          ))}
        </div>
      ) : (
        !error && (
          <div className="empty-state">
            <div className="empty-state-icon">🏥</div>
            <h3>Tidak Ada Klinik</h3>
            <p>Belum ada klinik yang terdaftar dalam sistem saat ini.</p>
          </div>
        )
      )}
    </div>
  );
}
