import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { supabase } from "../services/supabaseClient";
import { getCurrentUser } from "../services/api";

const AuthContext = createContext(null);

/**
 * AuthProvider — mengelola session pengguna via Supabase Auth.
 */
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState(null);
  const [role, setRole] = useState(null);
  const [clinicId, setClinicId] = useState(null);
  const [authError, setAuthError] = useState(null);

  const clearAuthState = useCallback(() => {
    setUser(null);
    setProfile(null);
    setRole(null);
    setClinicId(null);
    setAuthError(null);
  }, []);

  const syncSession = useCallback(
    async (session) => {
      setLoading(true);
      setAuthError(null);

      if (!session?.user) {
        clearAuthState();
        setLoading(false);
        return null;
      }

      setUser(session.user);

      try {
        const currentUser = await getCurrentUser();

        setProfile(currentUser);
        setRole(currentUser.role);
        setClinicId(currentUser.klinik_id);

        return currentUser;
      } catch (error) {
        setProfile(null);
        setRole(null);
        setClinicId(null);
        setAuthError(error.message);

        return null;
      } finally {
        setLoading(false);
      }
    },
    [clearAuthState],
  );

  useEffect(() => {
    let mounted = true;

    const initializeAuth = async () => {
      const {
        data: { session },
        error,
      } = await supabase.auth.getSession();

      if (!mounted) return;

      if (error) {
        setAuthError(error.message);
        setLoading(false);
        return;
      }

      await syncSession(session);
    };

    initializeAuth();

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setTimeout(() => {
        if (mounted) {
          void syncSession(session);
        }
      }, 0);
    });

    return () => {
      mounted = false;
      subscription.unsubscribe();
    };
  }, [syncSession]);

  const signIn = useCallback(
    async (email, password) => {
      const { data, error } = await supabase.auth.signInWithPassword({
        email,
        password,
      });

      if (error) {
        throw error;
      }

      const currentUser = await syncSession(data.session);

      if (!currentUser) {
        throw new Error("Login berhasil, tetapi role akun tidak dapat dimuat.");
      }

      return currentUser;
    },
    [syncSession],
  );

  const signUp = useCallback(async (email, password) => {
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
    });
    if (error) throw error;
    return data;
  }, []);

  const signOut = useCallback(async () => {
    const { error } = await supabase.auth.signOut();

    if (error) {
      throw error;
    }

    clearAuthState();
  }, [clearAuthState]);

  const value = useMemo(
    () => ({
      user,
      loading,
      profile,
      role,
      clinicId,
      authError,
      signIn,
      signUp,
      signOut,
    }),
    [
      user,
      loading,
      profile,
      role,
      clinicId,
      authError,
      signIn,
      signUp,
      signOut,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * useAuth — custom hook untuk mengakses AuthContext.
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth harus digunakan di dalam AuthProvider");
  }
  return context;
}
