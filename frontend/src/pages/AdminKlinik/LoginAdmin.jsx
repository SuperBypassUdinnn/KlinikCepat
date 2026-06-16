import { useState } from 'react';
import { useNavigate, Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import Card from '../../components/Card';
import Button from '../../components/Button';
import { FiShield } from 'react-icons/fi';
import './LoginAdmin.css';

export default function LoginAdmin() {
  const { signIn, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Redirect jika sudah login
  if (isAuthenticated) {
    return <Navigate to="/admin/dashboard" replace />;
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await signIn(email, password);
      navigate('/admin/dashboard');
    } catch (err) {
      setError(err.message || 'Login gagal. Periksa email dan password Anda.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-wrapper login-page">
      <Card className="login-card">
        <div className="card-body">
          <div className="login-icon">
            <FiShield />
          </div>
          <h2 className="login-title">Login Admin Klinik</h2>
          <p className="login-subtitle">
            Masuk untuk mengelola antrean pasien di faskes Anda
          </p>

          {error && (
            <div className="alert alert-danger animate-fade-in" style={{ marginBottom: '1rem' }}>
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label" htmlFor="admin-email">
                Email
              </label>
              <input
                id="admin-email"
                type="email"
                className="form-input"
                placeholder="admin@klinik.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoFocus
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="admin-password">
                Password
              </label>
              <input
                id="admin-password"
                type="password"
                className="form-input"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <Button
              type="submit"
              variant="primary"
              block
              size="lg"
              loading={loading}
              style={{ marginTop: 'var(--space-md)' }}
            >
              Masuk
            </Button>
          </form>

          <div className="login-divider">atau</div>

          <Button
            variant="ghost"
            block
            onClick={() => navigate('/')}
          >
            Kembali ke Halaman Pasien
          </Button>
        </div>
      </Card>
    </div>
  );
}
