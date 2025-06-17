import { Box, IconButton, useTheme, Badge } from "@mui/material";
import { useContext } from "react";
import { ColorModeContext, tokens } from "../../theme";
import InputBase from "@mui/material/InputBase";
import LightModeOutlinedIcon from "@mui/icons-material/LightModeOutlined";
import DarkModeOutlinedIcon from "@mui/icons-material/DarkModeOutlined";
import SettingsOutlinedIcon from "@mui/icons-material/SettingsOutlined";
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined";
import NotificationsOutlinedIcon from "@mui/icons-material/NotificationsOutlined";
import MarkunreadOutlinedIcon from "@mui/icons-material/MarkunreadOutlined";
import SearchIcon from "@mui/icons-material/Search";
import { useNavigate } from "react-router-dom";

const Topbar = (
  updateCount: any, 
  invitationCount: any
) => {
  const theme = useTheme();
  console.log("Update count; ",updateCount)
  const navigate=useNavigate()
  const colors = tokens(theme.palette.mode);
  const colorMode = useContext(ColorModeContext);
  console.log("inside top bar notification: ", updateCount);
  return (
    <Box display="flex" justifyContent="space-between" p={2} maxWidth="99%">
      <Box
        display="flex"
        borderRadius="3px"
        sx={{
          backgroundColor: colors.primary[400],
        }}
      >
        <InputBase sx={{ ml: 2, flex: 1 }} placeholder="Search" />
        <IconButton>
          <SearchIcon />
        </IconButton>
      </Box>

      <Box display="flex">
        <IconButton onClick={colorMode.toggleColorMode}>
          {theme.palette.mode === "dark" ? (
            <DarkModeOutlinedIcon />
          ) : (
            <LightModeOutlinedIcon />
          )}
        </IconButton>
        <IconButton
          onClick={()=>{navigate('/updates')}}
        >
          <Badge
            badgeContent={
              updateCount.updateCount !== 0 || updateCount.updateCount!==undefined ? updateCount.updateCount : null
            }
            color="error"
          >
            {" "}
            <NotificationsOutlinedIcon />
          </Badge>
        </IconButton>
        <IconButton
          onClick={()=>{navigate('/invitations')}}
        >
          <MarkunreadOutlinedIcon />
        </IconButton>
        <IconButton>
          <SettingsOutlinedIcon />
        </IconButton>
        <IconButton>
          <PersonOutlinedIcon />
        </IconButton>
      </Box>
    </Box>
  );
};

export default Topbar;
