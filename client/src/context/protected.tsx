import React from "react";
import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "./authContext";

const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, authChecked } = useAuth();

  if (!authChecked) {
    // Optional: show spinner or splash screen
    return <div>Loading...</div>;
  }

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
