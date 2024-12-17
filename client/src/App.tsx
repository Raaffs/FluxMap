import React, { useState } from 'react';
import './App.css';
import { Router } from 'react-router-dom';
import {ColorModeContext, useMode} from './theme'
import { CssBaseline,ThemeProvider } from '@mui/material';
import Topbar from './scenes/global/Topbar';
import SidebarEx from './scenes/global/Sidebar';
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
              <Topbar setIsSidebar={setIsSidebar} />
            </main>
          </div>
      </ThemeProvider>
    </ColorModeContext.Provider>
  );
}

export default App;
