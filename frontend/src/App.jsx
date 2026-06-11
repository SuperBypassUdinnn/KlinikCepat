import { Routes, Route } from 'react-router-dom';

function App() {
  return (
    <Routes>
      {/* Patient Routes */}
      <Route path="/" element={<div>Halaman Utama Pasien</div>} />

      {/* Admin Klinik Routes */}
      <Route path="/admin/*" element={<div>Dashboard Admin Klinik</div>} />

      {/* Super Admin Routes */}
      <Route path="/superadmin/*" element={<div>Dashboard Super Admin</div>} />
    </Routes>
  );
}

export default App;
