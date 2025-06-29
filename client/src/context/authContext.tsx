import React, {
    createContext,
    useContext,
    useEffect,
    useState,
    ReactNode,
  } from "react";
import { Navigate, useNavigate } from "react-router-dom";
  
  interface AuthContextType {
    isAuthenticated: boolean;
    authChecked: boolean; 
    login: () => void;
    logout: () => void;
  }
  const AuthContext = createContext<AuthContextType>({
    isAuthenticated: false,
    authChecked: false,
    login: () => {},
    logout: () => {},
  });
  
  interface AuthProviderProps {
    children: ReactNode;
  }
  
  export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authChecked, setAuthChecked] = useState(false); // ✅ NEW
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch("http://localhost:4000/api/auth/session", {
          method: 'GET',
          credentials: "include",
        });

        const data = await res.json();
        setIsAuthenticated(data.isAuthenticated);

        // ✅ Only redirect if they’re stuck on login while already authenticated
        if (window.location.pathname === "/login" && data.isAuthenticated) {
          navigate("/");
        }

      } catch (err) {
        console.error("Failed to check auth", err);
      } finally {
        setAuthChecked(true); // ✅ Done checking!
      }
    };

    checkAuth();
  }, []);

  const login = () => {
    setIsAuthenticated(true);
  };

  const logout = async () => {
    try {
      await fetch("http://localhost:4000/api/logout", {
        method: "POST",
        credentials: "include",
      });
    } catch (err) {
      console.error("Logout failed", err);
    }
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, authChecked, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
export const useAuth = () => useContext(AuthContext);