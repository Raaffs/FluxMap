import React, { useState } from 'react';
import './App.css';
import { Route, Router, Routes } from 'react-router-dom';
import {ColorModeContext, useMode} from './theme'
import { CssBaseline,ThemeProvider } from '@mui/material';
import Topbar from './scenes/global/Topbar';
import SidebarEx from './scenes/global/Sidebar';
import LoginUser from './scenes/auth/login';
import SignUpUser from './scenes/auth/signup';
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
              <Routes>
                <Route path='/login' element={<LoginUser/>}/>
                <Route path='/register' element={<SignUpUser/>}/>
              </Routes>
            </main>
          </div>
      </ThemeProvider>
    </ColorModeContext.Provider>
  );
}

export default App;
