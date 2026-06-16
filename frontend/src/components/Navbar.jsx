import { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { HiOutlineMenu, HiOutlineX } from 'react-icons/hi';
import { FiActivity } from 'react-icons/fi';
import './Navbar.css';

export default function Navbar() {
  const { user, signOut } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  // Detect scroll for navbar shadow
  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Close mobile menu on route change
  useEffect(() => {
    setIsOpen(false);
  }, [location.pathname]);

  const isActive = (path) => location.pathname === path;

  const handleLogout = async () => {
    await signOut();
    navigate('/');
  };

  // Determine which nav to show based on current route
  const isAdminRoute = location.pathname.startsWith('/admin');
  const isSuperAdminRoute = location.pathname.startsWith('/superadmin');

  return (
    <nav className={`navbar ${scrolled ? 'scrolled' : ''}`}>
      <div className="navbar-inner">
        <Link to="/" className="navbar-brand">
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

        <ul className={`navbar-nav ${isOpen ? 'open' : ''}`}>
          {!isAdminRoute && !isSuperAdminRoute && (
            <>
              <li>
                <Link
                  to="/"
                  className={`navbar-nav-link ${isActive('/') ? 'active' : ''}`}
                  id="nav-link-home"
                >
                  Cari Klinik
                </Link>
              </li>
            </>
          )}

          {isAdminRoute && (
            <>
              <li>
                <Link
                  to="/admin/dashboard"
                  className={`navbar-nav-link ${isActive('/admin/dashboard') ? 'active' : ''}`}
                  id="nav-link-admin-dashboard"
                >
                  Dashboard
                </Link>
              </li>
            </>
          )}

          {isSuperAdminRoute && (
            <>
              <li>
                <Link
                  to="/superadmin/klinik"
                  className={`navbar-nav-link ${isActive('/superadmin/klinik') ? 'active' : ''}`}
                  id="nav-link-sa-klinik"
                >
                  Kelola Klinik
                </Link>
              </li>
              <li>
                <Link
                  to="/superadmin/gejala"
                  className={`navbar-nav-link ${isActive('/superadmin/gejala') ? 'active' : ''}`}
                  id="nav-link-sa-gejala"
                >
                  Kelola Gejala
                </Link>
              </li>
            </>
          )}

          <li>
            {user ? (
              <button
                className="navbar-auth-btn logout"
                onClick={handleLogout}
                id="navbar-logout-btn"
              >
                Logout
              </button>
            ) : (
              <Link to="/admin/login">
                <button className="navbar-auth-btn login" id="navbar-login-btn">
                  Login Admin
                </button>
              </Link>
            )}
          </li>
        </ul>
      </div>
    </nav>
  );
}
