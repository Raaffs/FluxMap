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
    login: () => void;
    logout: () => void;
  }
  const AuthContext = createContext<AuthContextType>({
    isAuthenticated: false,
    login: () => {},
    logout: () => {},
  });
  
  interface AuthProviderProps {
    children: ReactNode;
  }
  
  export const AuthProvider: React.FC<AuthProviderProps> =  ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const navigate=useNavigate()
    useEffect(() => {
      const checkAuth = async () => {
        try {
          const res = await fetch("http://localhost:4000/api/auth/session", {
            method:'GET',
            credentials: "include", 
          });
          const data = await res.json();
          setIsAuthenticated(data.isAuthenticated)
          if(isAuthenticated){
            navigate("/")
          }
          console.log("setisauth :",isAuthenticated,data.isAuthenticated)
        } catch (err) {
          console.error("Failed to check auth", err);
        }
      };
       checkAuth();
    }, []);
    const login = () => {
      // optional: do a redirect to login page
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
      <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
        {children}
      </AuthContext.Provider>
    );
  };
  
  export const useAuth = () => useContext(AuthContext);
  