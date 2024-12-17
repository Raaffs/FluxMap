import { useState } from "react";
import { Sidebar, Menu, MenuItem,SubMenu } from "react-pro-sidebar";
import { Box, IconButton, Typography, useTheme } from "@mui/material";
import { Link } from "react-router-dom";
import { tokens } from "../../theme";
import HomeOutlinedIcon from "@mui/icons-material/HomeOutlined";
import PeopleOutlinedIcon from "@mui/icons-material/PeopleOutlined";
import ContactsOutlinedIcon from "@mui/icons-material/ContactsOutlined";
import ReceiptOutlinedIcon from "@mui/icons-material/ReceiptOutlined";
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined";
import CalendarTodayOutlinedIcon from "@mui/icons-material/CalendarTodayOutlined";
import HelpOutlineOutlinedIcon from "@mui/icons-material/HelpOutlineOutlined";
import BarChartOutlinedIcon from "@mui/icons-material/BarChartOutlined";
import PieChartOutlineOutlinedIcon from "@mui/icons-material/PieChartOutlineOutlined";
import TimelineOutlinedIcon from "@mui/icons-material/TimelineOutlined";
import MenuOutlinedIcon from "@mui/icons-material/MenuOutlined";
import SettingsOutlinedIcon from '@mui/icons-material/SettingsOutlined';
import ManageAccountsOutlinedIcon from '@mui/icons-material/ManageAccountsOutlined';
import AdminPanelSettingsSharpIcon from '@mui/icons-material/AdminPanelSettingsSharp';
import TaskAltSharpIcon from '@mui/icons-material/TaskAltSharp';
import PendingActionsOutlinedIcon from '@mui/icons-material/PendingActionsOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import Inventory2Outlined from '@mui/icons-material/Inventory2Outlined';
import LoginOutlinedIcon from '@mui/icons-material/LoginOutlined';
import PersonAddOutlinedIcon from '@mui/icons-material/PersonAddOutlined';
import NotificationsNoneOutlinedIcon from '@mui/icons-material/NotificationsNoneOutlined';

const SidebarEx = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [selected, setSelected] = useState("Dashboard");

  return (
   
      <Sidebar collapsed={isCollapsed}
        backgroundColor={ theme.palette.mode=='dark'? colors.primary[400]: 'white'}
        rootStyles={{
          ".ps-icon-wrapper": {
            backgroundColor: "transparent !important",
          },
          // ".ps-item": {
          //   padding: "5px 35px 5px 20px !important",
          // },
          ".ps-menu-button:hover": {
            backgroundColor:` ${colors.blueAccent[700]} !important`,
            color: ` ${colors.blueAccent[500]} !important`,
          },
          ".ps-item.active": {
            color: "#blue !important",
          },
          alignContent:'space-between'
        }}
        style={{
          // borderRadius: '9px',        // Optional: keep rounded corners
          // margin: '5px',              // Add margin to ensure spacing
        }}
      >
  <Menu
       menuItemStyles={{
        button: {
          textAlign: 'left',
        },
      }}
  >
    <MenuItem icon={<HomeOutlinedIcon/>}>Dashboard</MenuItem>
    <SubMenu
      icon={<TimelineOutlinedIcon/>}
      label="Projects"
 
    >
      <SubMenu
        label="Ownership"
        style={{       
          backgroundColor: theme.palette.mode==='dark'? colors.primary[500]:'white' 
        }}
      >
        <MenuItem 
          icon={<AdminPanelSettingsSharpIcon/>}
          style={{
            color:'#4CAF50',
            backgroundColor: theme.palette.mode==='dark'? colors.primary[600]:'white' 
          }} 
          >
            Admin
          </MenuItem>
        <MenuItem  
          icon={<ManageAccountsOutlinedIcon/>}
          style={{
              color:'#8E44AD',
              backgroundColor: theme.palette.mode==='dark'? colors.primary[600]:'white' 
            }}>Manager</MenuItem>
        <MenuItem 
         icon={<PeopleOutlinedIcon/>}
         style={{
            color:'#3498DB',
            backgroundColor: theme.palette.mode==='dark'? colors.primary[600]:'white' 
          }} >Assigned</MenuItem>
      </SubMenu>
      <SubMenu
        label="Status"
        style={{ backgroundColor: theme.palette.mode==='dark'? colors.primary[500]:'white' }}
      >
        <MenuItem  
          icon={<TaskAltSharpIcon/>}
          style={{
            color: colors.greenAccent[400],
            backgroundColor:  theme.palette.mode=='dark'? colors.primary[600]:'white'
          }}> Completed</MenuItem>
        <MenuItem  
          icon={<PendingActionsOutlinedIcon/>}
          style={{
            color:'#FFEB3B',
            backgroundColor:theme.palette.mode=='dark'? colors.primary[600]:'white'
          }}>Pending</MenuItem>
      </SubMenu>
    </SubMenu>
    <MenuItem icon={<NotificationsNoneOutlinedIcon/>} >Notifications</MenuItem>
    <MenuItem icon={<SettingsOutlinedIcon/>}>Settings</MenuItem>
    <MenuItem
      icon={<LogoutOutlinedIcon/>}
    >Logout</MenuItem>
    <MenuItem
      icon={<LoginOutlinedIcon/>}
    >Login</MenuItem>
    <MenuItem
      icon={<PersonAddOutlinedIcon/>}
    >Sign Up</MenuItem>
  </Menu>
      </Sidebar>
  );
};

export default SidebarEx;