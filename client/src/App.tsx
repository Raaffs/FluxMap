import React, { useState } from "react";
import "./App.css";
import { Route, Routes } from "react-router-dom";
import { ColorModeContext, useMode } from "./theme";
import { CssBaseline, ThemeProvider } from "@mui/material";
import Topbar from "./scenes/global/Topbar";
import SidebarEx from "./scenes/global/Sidebar";
import LoginUser from "./scenes/auth/login";
import SignUpUser from "./scenes/auth/signup";
import { ProjectOverview } from "./scenes/projects/overview";
import Graphs from "./components/graphs/LineGraphs";
import VoiceRecorder from "./components/voiceRecorder";
import { Invitation } from "./scenes/invitation";
import Updates from "./scenes/notification";
import {
  AdminProjects,
  ManagerProjects,
  AllocatedProjects,
} from "./scenes/projects/projects";
import { AuthProvider } from "./context/authContext";
import ProtectedRoute from "./context/protected"; 
import { useEffect } from "react";

function App() {
  const [theme, colorMode] = useMode();
  const [isSidebar, setIsSidebar] = useState(true);
  useEffect(() => {
    const socket = new WebSocket('http://localhost:4000/api/ws');

    socket.onopen = () => {
      console.log('🟢 WebSocket connected');
    };

    socket.onmessage = (event) => {
      console.log('📨 Message from server:', event.data);
    };

    socket.onclose = () => {
      console.log('🔴 WebSocket disconnected');
    };

    socket.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    return () => {
      socket.close();
    };
  }, []);


  return (
    <AuthProvider>
      <ColorModeContext.Provider value={colorMode}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <div className="App">
            {isSidebar && <SidebarEx />}
            <main className="content">
              <Topbar />
              <Routes>
                {/* Public Routes */}
                <Route path="/login" element={<LoginUser />} />
                <Route path="/register" element={<SignUpUser />} />

                {/* Protected Routes */}
                <Route element={<ProtectedRoute />}>
                  <Route path="/projects/admin" element={<AdminProjects />} />
                  <Route path="/projects/manager" element={<ManagerProjects />} />
                  <Route path="/projects/allocated" element={<AllocatedProjects />} />
                  <Route path="/project/:id" element={<ProjectOverview />} />
                  <Route path="/" element={<Graphs />} />
                  <Route path="/invitations" element={<Invitation />} />
                  <Route path="/voice" element={<VoiceRecorder />} />
                  <Route path="/updates" element={<Updates/>}/>
                </Route>
              </Routes>
            </main>
          </div>
        </ThemeProvider>
      </ColorModeContext.Provider>
    </AuthProvider>
  );
}


export default App;
