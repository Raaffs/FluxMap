import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Sidebar, Menu, MenuItem, SubMenu } from "react-pro-sidebar";
import { Box, useTheme, Typography } from "@mui/material";
import { tokens } from "../../theme";
import HomeOutlinedIcon from "@mui/icons-material/HomeOutlined";
import PeopleOutlinedIcon from "@mui/icons-material/PeopleOutlined";
import TimelineOutlinedIcon from "@mui/icons-material/TimelineOutlined";
import SettingsOutlinedIcon from '@mui/icons-material/SettingsOutlined';
import ManageAccountsOutlinedIcon from '@mui/icons-material/ManageAccountsOutlined';
import AdminPanelSettingsSharpIcon from '@mui/icons-material/AdminPanelSettingsSharp';
import TaskAltSharpIcon from '@mui/icons-material/TaskAltSharp';
import PendingActionsOutlinedIcon from '@mui/icons-material/PendingActionsOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import LoginOutlinedIcon from '@mui/icons-material/LoginOutlined';
import PersonAddOutlinedIcon from '@mui/icons-material/PersonAddOutlined';
import NotificationsNoneOutlinedIcon from '@mui/icons-material/NotificationsNoneOutlined';
import AutoModeOutlinedIcon from '@mui/icons-material/AutoModeOutlined';
import EngineeringOutlinedIcon from '@mui/icons-material/EngineeringOutlined';
import MarkunreadOutlinedIcon from '@mui/icons-material/MarkunreadOutlined';
import PublishedWithChangesIcon from '@mui/icons-material/PublishedWithChanges';

// 👇 Import your custom auth hook
import { useAuth } from "../../context/authContext";

const SidebarEx = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const navigate = useNavigate();

  const { isAuthenticated } = useAuth(); // 👈 Get authentication status
  const {logout}=useAuth()
  return (
    <Sidebar
      collapsed={isCollapsed}
      backgroundColor={theme.palette.mode === 'dark' ? colors.primary[400] : 'white'}
      rootStyles={{
        ".ps-icon-wrapper": {
          backgroundColor: "transparent !important",
        },
        ".ps-menu-button:hover": {
          backgroundColor: `${colors.blueAccent[700]} !important`,
          color: `${colors.blueAccent[500]} !important`,
        },
        ".ps-item.active": {
          color: "#blue !important",
        },
        alignContent: 'space-between',
      }}
    >
      <Box
        textAlign="center"
        sx={{
          backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : 'white',
          padding: "15px",
          borderRadius: "8px",
        }}
      >
        <Typography
          variant="h3"
          color={colors.grey[100]}
          fontWeight="600"
          sx={{ mb: "8px" }}
        >
          FluxMap
        </Typography>
        <Typography variant="h6" color={colors.greenAccent[500]}>
          A project Management System
        </Typography>
      </Box>

      <Menu
        menuItemStyles={{
          button: {
            textAlign: 'left',
          },
        }}
      >
        {isAuthenticated ? (
          <>
            {/* Authenticated Menu Items */}
            <MenuItem icon={<HomeOutlinedIcon />}
              onClick={()=>{navigate('/')}}
            >Dashboard</MenuItem>

            <SubMenu icon={<TimelineOutlinedIcon />} label="Projects">
              <SubMenu
                label="Roles"
                icon={<EngineeringOutlinedIcon />}
                style={{
                  backgroundColor: theme.palette.mode === 'dark' ? colors.primary[500] : 'white',
                }}
              >
                <MenuItem
                  onClick={() => navigate('/projects/admin')}
                  icon={<AdminPanelSettingsSharpIcon />}
                  style={{
                    color: '#4CAF50',
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[600] : 'white',
                  }}
                >
                  Admin
                </MenuItem>
                <MenuItem
                  onClick={() => navigate('/projects/manager')}
                  icon={<ManageAccountsOutlinedIcon />}
                  style={{
                    color: '#8E44AD',
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[600] : 'white',
                  }}
                >
                  Manager
                </MenuItem>
                <MenuItem
                  onClick={() => navigate('/projects/allocated')}
                  icon={<PeopleOutlinedIcon />}
                  style={{
                    color: '#3498DB',
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[600] : 'white',
                  }}
                >
                  Allocated
                </MenuItem>
              </SubMenu>

              <SubMenu
                label="Status"
                icon={<AutoModeOutlinedIcon />}
                style={{
                  backgroundColor: theme.palette.mode === 'dark' ? colors.primary[500] : 'white',
                }}
              >
                <MenuItem
                  icon={<TaskAltSharpIcon />}
                  style={{
                    color: colors.greenAccent[400],
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[600] : 'white',
                  }}
                >
                  Completed
                </MenuItem>
                <MenuItem
                  icon={<PendingActionsOutlinedIcon />}
                  style={{
                    color: '#FFEB3B',
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[600] : 'white',
                  }}
                >
                  Pending
                </MenuItem>
              </SubMenu>
            </SubMenu>

            <SubMenu icon={<NotificationsNoneOutlinedIcon />} label="Notifications">
              <MenuItem icon={<PublishedWithChangesIcon />}>Updates</MenuItem>
              <MenuItem icon={<MarkunreadOutlinedIcon />} onClick={() => navigate('/invitations')}>
                Invitations
              </MenuItem>
            </SubMenu>

            <MenuItem icon={<SettingsOutlinedIcon />}>Settings</MenuItem>
            <MenuItem icon={<LogoutOutlinedIcon />}
            onClick={async()=>{
              await logout()
              navigate('/login')
            }}
            >Logout</MenuItem>
          </>
        ) : null}

        {/* Public Items - Always visible */}
        {!isAuthenticated && (
          <>
            <MenuItem onClick={() => navigate('/login')} icon={<LoginOutlinedIcon />}>
              Login
            </MenuItem>
            <MenuItem onClick={() => navigate('/register')} icon={<PersonAddOutlinedIcon />}>
              Sign Up
            </MenuItem>
          </>
        )}
      </Menu>
    </Sidebar>
  );
};

export default SidebarEx;
