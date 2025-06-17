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
import { useWebSocket } from "./hooks/websocket";

function App() {
  const [theme, colorMode] = useMode();
  const [isSidebar, setIsSidebar] = useState(true);
  const [updateCount, setUpdateCount]=useState(0)
  const [invitationCount, setInvitationcount]=useState(0)
  const [startWebSocket]=useWebSocket("ws://localhost:4000/api/ws",setUpdateCount,setInvitationcount)
  return (
    <AuthProvider>
      <ColorModeContext.Provider value={colorMode}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <div className="App">
            {isSidebar && <SidebarEx />}
            <main className="content">
              <Topbar updateCount={updateCount} invitationCount={invitationCount}/>
              <Routes>
                {/* Public Routes */}
                <Route path="/login" element={<LoginUser startWebSocket={startWebSocket} />} />
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
