import React, { useState } from 'react';
import './App.css';
import { Route, Routes } from 'react-router-dom';
import {ColorModeContext, useMode} from './theme'
import { CssBaseline,ThemeProvider } from '@mui/material';
import Topbar from './scenes/global/Topbar';
import SidebarEx from './scenes/global/Sidebar';
import LoginUser from './scenes/auth/login';
import SignUpUser from './scenes/auth/signup';
import { ProjectTaskDetailPage } from './components/task';
import { ProjectOverview } from './scenes/projects/overview';
import Graphs from './components/graphs/LineGraphs';
import VoiceRecorder from './components/voiceRecorder';
import { AdminProjects,ManagerProjects,AllocatedProjects } from './scenes/projects/projects';
function App() {
  const [theme, colorMode] = useMode();
  const [isSidebar,setIsSidebar]=useState(true)
  return (
    <ColorModeContext.Provider value={colorMode}>
      <ThemeProvider theme={theme}>
        <CssBaseline/>
          <div className="App">
          {isSidebar && <SidebarEx />}
            <main className="content">
              <Topbar />
              <Routes>
                <Route path='/login' element={<LoginUser/>}/>
                <Route path='/register' element={<SignUpUser/>}/>
                <Route path='/projects/admin' element={<AdminProjects/>}/>
                <Route path='/projects/manager' element={<ManagerProjects/>}/>
                <Route path='/projects/allocated' element={<AllocatedProjects/>}/>
                <Route path='/project/:id' element={<ProjectOverview/>}/>
                <Route path='/dashboard/:id' element={<Graphs/>}/>
                <Route path='/voice' element={<VoiceRecorder/>}/>
              </Routes>
            </main>
          </div>
      </ThemeProvider>
    </ColorModeContext.Provider>
  );
}

export default App;
