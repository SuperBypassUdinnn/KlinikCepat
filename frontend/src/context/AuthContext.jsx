import { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext(null);

/**
 * AuthProvider — menyimpan session user dari Supabase Auth.
 * Bungkus komponen root (App) dengan provider ini.
 */
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // TODO: Integrasikan dengan Supabase Auth untuk mendengarkan perubahan session
    setLoading(false);
  }, []);

  const value = {
    user,
    setUser,
    loading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

/**
 * useAuth — custom hook untuk mengakses AuthContext.
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth harus digunakan di dalam AuthProvider');
  }
  return context;
}
