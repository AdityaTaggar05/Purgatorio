import { useState, useEffect, useCallback, type ReactNode } from 'react';
import { AuthContext } from './AuthContext';
import type { User } from '../../types/user';
import * as api from '../../features/auth/authApi'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<User | null>(() => {
    const savedUser = localStorage.getItem('purgatorio_user');
    return savedUser ? JSON.parse(savedUser) : null;
  });

  const login = (token: string, userData: User) => {
    setAccessToken(token);
    setUser(userData);
    localStorage.setItem('purgatorio_user', JSON.stringify(userData));
  };

  const logout = useCallback(() => {
    setAccessToken(null);
    setUser(null);
    localStorage.removeItem('purgatorio_user');
  }, []);

  const getFreshToken = useCallback(async (): Promise<string | null> => {
    try {
      const response = await api.refresh()
      if (!response.success) throw new Error('Session expired');
      setAccessToken(response.data.access_token);
      return response.data.access_token;
    } catch {
      logout();
      return null;
    }
  }, [logout]);

  useEffect(() => {
    let isMounted = true;
    const checkSession = async () => {
      if (user) {
        await getFreshToken();
      }

      if (isMounted) {
        setIsLoading(false)
      }
    };
    checkSession();
    return () => { isMounted = false; };
  }, [user, getFreshToken]);

  return (
    <AuthContext.Provider value={{ user, accessToken, isLoading, login, logout, getFreshToken }}>
      {children}
    </AuthContext.Provider>
  );
}
