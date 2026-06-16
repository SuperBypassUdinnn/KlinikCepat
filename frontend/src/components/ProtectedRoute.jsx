import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import LoadingSpinner from './LoadingSpinner';

/**
 * ProtectedRoute — wrapper route yang memeriksa autentikasi.
 * Redirect ke /admin/login jika user belum login.
 */
export default function ProtectedRoute({ children }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner fullPage text="Memverifikasi sesi..." />;
  }

  if (!user) {
    return <Navigate to="/admin/login" replace />;
  }

  return children;
}
