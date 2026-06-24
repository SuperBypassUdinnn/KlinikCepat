import { Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import ProtectedRoute from "./components/ProtectedRoute";
import { Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import ProtectedRoute from "./components/ProtectedRoute";

// Patient Pages
import CariKlinik from "./pages/Patient/CariKlinik";
import TriageForm from "./pages/Patient/TriageForm";
import TicketAntrean from "./pages/Patient/TicketAntrean";
import CariKlinik from "./pages/Patient/CariKlinik";
import TriageForm from "./pages/Patient/TriageForm";
import TicketAntrean from "./pages/Patient/TicketAntrean";

// Admin Klinik Pages
import LoginAdmin from "./pages/AdminKlinik/LoginAdmin";
import DashboardAdmin from "./pages/AdminKlinik/DashboardAdmin";
import LoginAdmin from "./pages/AdminKlinik/LoginAdmin";
import DashboardAdmin from "./pages/AdminKlinik/DashboardAdmin";

// Super Admin Pages
import ManajemenKlinik from "./pages/SuperAdmin/ManajemenKlinik";
import ManajemenGejala from "./pages/SuperAdmin/ManajemenGejala";
import ManajemenAdminKlinik from "./pages/SuperAdmin/ManajemenAdminKlinik";

function App() {
  return (
    <>
      <Navbar />


      <Routes>
        {/* Patient Routes */}
        <Route path="/" element={<CariKlinik />} />


        <Route path="/triage/:klinikId" element={<TriageForm />} />


        <Route path="/ticket" element={<TicketAntrean />} />

        {/* Login */}
        {/* Login */}
        <Route path="/admin/login" element={<LoginAdmin />} />

        {/* Admin Klinik Routes */}

        {/* Admin Klinik Routes */}
        <Route
          path="/admin/dashboard"
          element={
            <ProtectedRoute allowedRoles={["klinik_admin"]}>
            <ProtectedRoute allowedRoles={["klinik_admin"]}>
              <DashboardAdmin />
            </ProtectedRoute>
          }
        />

        {/* Super Admin Routes */}
        <Route
          path="/superadmin/klinik"
          element={
            <ProtectedRoute allowedRoles={["superadmin"]}>
            <ProtectedRoute allowedRoles={["superadmin"]}>
              <ManajemenKlinik />
            </ProtectedRoute>
          }
        />


        <Route
          path="/superadmin/gejala"
          element={
            <ProtectedRoute allowedRoles={["superadmin"]}>
            <ProtectedRoute allowedRoles={["superadmin"]}>
              <ManajemenGejala />
            </ProtectedRoute>
          }
        />
        
        <Route
          path="/superadmin/admin-klinik"
          element={
            <ProtectedRoute allowedRoles={["superadmin"]}>
              <ManajemenAdminKlinik />
            </ProtectedRoute>
          }
        />
      </Routes>
    </>
  );
}

export default App;
