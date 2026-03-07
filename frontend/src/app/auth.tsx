import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren
} from "react";
import { useNavigate } from "react-router-dom";

import { api, unwrap } from "../api/client";
import type { User } from "../types/api";

type AuthContextValue = {
  user: User | null;
  ready: boolean;
  login: (payload: { identifier: string; password: string; rememberMe: boolean }) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: PropsWithChildren) {
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);
  const navigate = useNavigate();

  const refresh = async () => {
    try {
      const result = await unwrap<{ user: User }>(api.get("/auth/me"));
      setUser(result.user);
    } catch {
      setUser(null);
    } finally {
      setReady(true);
    }
  };

  useEffect(() => {
    void refresh();
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      ready,
      refresh,
      async login(payload) {
        await unwrap<{ user: User }>(api.post("/auth/login", payload));
        await refresh();
        navigate("/");
      },
      async logout() {
        await unwrap(api.post("/auth/logout", {}));
        setUser(null);
        navigate("/login");
      }
    }),
    [navigate, ready, user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
