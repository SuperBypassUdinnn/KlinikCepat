import { useState, useEffect } from 'react';

/**
 * useHaversine — menghitung jarak (km) antara posisi pengguna dan koordinat target
 * menggunakan rumus Haversine.
 *
 * @param {number} targetLat - Latitude target (klinik)
 * @param {number} targetLng - Longitude target (klinik)
 * @returns {{ distance: number|null, error: string|null, loading: boolean }}
 */
export function useHaversine(targetLat, targetLng) {
  const [distance, setDistance] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!navigator.geolocation) {
      setError('Geolocation tidak didukung oleh browser ini.');
      setLoading(false);
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        const d = haversine(latitude, longitude, targetLat, targetLng);
        setDistance(d);
        setLoading(false);
      },
      (err) => {
        setError(err.message);
        setLoading(false);
      },
    );
  }, [targetLat, targetLng]);

  return { distance, error, loading };
}

/**
 * Rumus Haversine — menghitung jarak antara dua titik koordinat di permukaan bumi.
 * @returns {number} Jarak dalam kilometer.
 */
function haversine(lat1, lon1, lat2, lon2) {
  const R = 6371; // Radius bumi dalam km
  const dLat = toRad(lat2 - lat1);
  const dLon = toRad(lon2 - lon1);
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2;
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

function toRad(deg) {
  return (deg * Math.PI) / 180;
}
