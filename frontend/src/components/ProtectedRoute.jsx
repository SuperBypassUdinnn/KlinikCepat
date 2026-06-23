import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import LoadingSpinner from "./LoadingSpinner";

/**
 * ProtectedRoute — memeriksa autentikasi dan role pengguna.
 */
export default function ProtectedRoute({ children, allowedRoles = [] }) {
  const { user, role, loading, authError } = useAuth();

  const location = useLocation();

  if (loading) {
    return <LoadingSpinner fullPage text="Memverifikasi sesi..." />;
  }

  if (!user) {
    return <Navigate to="/admin/login" replace state={{ from: location }} />;
  }

  if (authError) {
    return (
      <div className="container page-wrapper">
        <div className="alert alert-danger">
          Gagal memuat hak akses akun: {authError}
        </div>
      </div>
    );
  }

  const isRoleAllowed =
    allowedRoles.length === 0 || allowedRoles.includes(role);

  if (!isRoleAllowed) {
    if (role === "superadmin") {
      return <Navigate to="/superadmin/klinik" replace />;
    }

    if (role === "klinik_admin") {
      return <Navigate to="/admin/dashboard" replace />;
    }

    return <Navigate to="/admin/login" replace />;
  }

  return children;
}
