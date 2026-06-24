import { useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { HiOutlineMenu, HiOutlineX } from "react-icons/hi";
import { FiActivity } from "react-icons/fi";
import "./Navbar.css";

function getDashboardPath(role) {
  switch (role) {
    case "superadmin":
      return "/superadmin/klinik";

    case "klinik_admin":
      return "/admin/dashboard";

    default:
      return "/";
  }
}

export default function Navbar() {
  const { user, role, loading: authLoading, signOut, clinicName } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const [logoutLoading, setLogoutLoading] = useState(false);
  const [logoutError, setLogoutError] = useState(null);

  // Detect scroll for navbar shadow
  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  // Close mobile menu on route change
  useEffect(() => {
    setIsOpen(false);
  }, [location.pathname]);

  const isActive = (path) => location.pathname === path;

  const handleLogout = async () => {
    setLogoutLoading(true);
    setLogoutError(null);

    try {
      await signOut();
      navigate("/", { replace: true });
    } catch (error) {
      setLogoutError(error.message || "Gagal logout.");
    } finally {
      setLogoutLoading(false);
    }
  };

  const dashboardPath = getDashboardPath(role);

  const brandPath = user && role ? dashboardPath : "/";

  return (
    <>
      <nav className={`navbar ${scrolled ? "scrolled" : ""}`}>
        <div className="navbar-inner">
          <Link to={brandPath} className="navbar-brand">
            <div className="navbar-logo-icon">
              <FiActivity />
            </div>
            <span className="navbar-brand-text">
              Klinik<span>Cepat</span>
            </span>
          </Link>
          <button
            className="navbar-toggle"
            onClick={() => setIsOpen(!isOpen)}
            aria-label="Toggle navigation"
            id="navbar-toggle-btn"
          >
            {isOpen ? <HiOutlineX /> : <HiOutlineMenu />}
          </button>
          <ul className={`navbar-nav ${isOpen ? "open" : ""}`}>
            {!user && !authLoading && (
              <li>
                <Link
                  to="/"
                  className={`navbar-nav-link ${isActive("/") ? "active" : ""}`}
                  id="nav-link-home"
                >
                  Cari Klinik
                </Link>
              </li>
            )}

            {user && role === "klinik_admin" && (
              <>
                <li>
                  <Link
                    to="/admin/dashboard"
                    className={`navbar-nav-link ${
                      isActive("/admin/dashboard") ? "active" : ""
                    }`}
                    id="nav-link-admin-dashboard"
                  >
                    Dashboard
                  </Link>
                </li>
                <li>
                  <span
                    className="navbar-clinic-context"
                    title={clinicName || "Klinik tidak diketahui"}
                  >
                    <small>Login di</small>
                    <strong>{clinicName || "Klinik tidak diketahui"}</strong>
                  </span>
                </li>

                <li>
                  <Link
                    to="/"
                    className={`navbar-nav-link ${isActive("/") ? "active" : ""}`}
                  >
                    Halaman Pasien
                  </Link>
                </li>
              </>
            )}

            {user && role === "superadmin" && (
              <>
                <li>
                  <Link
                    to="/superadmin/klinik"
                    className={`navbar-nav-link ${
                      isActive("/superadmin/klinik") ? "active" : ""
                    }`}
                    id="nav-link-sa-klinik"
                  >
                    Kelola Klinik
                  </Link>
                </li>

                <li>
                  <Link
                    to="/superadmin/gejala"
                    className={`navbar-nav-link ${
                      isActive("/superadmin/gejala") ? "active" : ""
                    }`}
                    id="nav-link-sa-gejala"
                  >
                    Kelola Gejala
                  </Link>
                </li>

                <li>
                  <Link
                    to="/superadmin/admin-klinik"
                    className={`navbar-nav-link ${
                      isActive("/superadmin/admin-klinik") ? "active" : ""
                    }`}
                    id="nav-link-sa-admin-klinik"
                  >
                    Kelola Admin
                  </Link>
                </li>

                <li>
                  <Link
                    to="/"
                    className={`navbar-nav-link ${isActive("/") ? "active" : ""}`}
                  >
                    Halaman Pasien
                  </Link>
                </li>
              </>
            )}

            <li>
              {authLoading ? (
                <span className="navbar-auth-loading">Memuat...</span>
              ) : user ? (
                <button
                  type="button"
                  className="navbar-auth-btn logout"
                  onClick={handleLogout}
                  disabled={logoutLoading}
                  id="navbar-logout-btn"
                >
                  {logoutLoading ? "Keluar..." : "Logout"}
                </button>
              ) : (
                <Link
                  to="/admin/login"
                  className="navbar-auth-btn login"
                  id="navbar-login-btn"
                >
                  Login Admin
                </Link>
              )}
            </li>
          </ul>{" "}
        </div>
      </nav>
      {logoutError && <div className="navbar-error">{logoutError}</div>}
    </>
  );
}
