import { Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';

// Patient Pages
import CariKlinik from './pages/Patient/CariKlinik';
import TriageForm from './pages/Patient/TriageForm';
import TicketAntrean from './pages/Patient/TicketAntrean';

// Admin Klinik Pages
import LoginAdmin from './pages/AdminKlinik/LoginAdmin';
import DashboardAdmin from './pages/AdminKlinik/DashboardAdmin';

// Super Admin Pages
import ManajemenKlinik from './pages/SuperAdmin/ManajemenKlinik';
import ManajemenGejala from './pages/SuperAdmin/ManajemenGejala';

function App() {
  return (
    <>
      <Navbar />
      <Routes>
        {/* Patient Routes */}
        <Route path="/" element={<CariKlinik />} />
        <Route path="/triage/:klinikId" element={<TriageForm />} />
        <Route path="/ticket" element={<TicketAntrean />} />

        {/* Admin Klinik Routes */}
        <Route path="/admin/login" element={<LoginAdmin />} />
        <Route
          path="/admin/dashboard"
          element={
            <ProtectedRoute>
              <DashboardAdmin />
            </ProtectedRoute>
          }
        />

        {/* Super Admin Routes */}
        <Route
          path="/superadmin/klinik"
          element={
            <ProtectedRoute>
              <ManajemenKlinik />
            </ProtectedRoute>
          }
        />
        <Route
          path="/superadmin/gejala"
          element={
            <ProtectedRoute>
              <ManajemenGejala />
            </ProtectedRoute>
          }
        />
      </Routes>
    </>
  );
}

export default App;
