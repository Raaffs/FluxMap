import { Box, useTheme } from "@mui/material";
import AdminPanelSettingsOutlinedIcon from "@mui/icons-material/AdminPanelSettingsOutlined";
import ManageAccountsOutlinedIcon from "@mui/icons-material/ManageAccountsOutlined";
import TimerOutlinedIcon from "@mui/icons-material/TimerOutlined";
import PermIdentityOutlinedIcon from "@mui/icons-material/PermIdentityOutlined";
import { tokens } from "../../theme";
import StatCard from "../../components/cards/statcard";
const Dashboard = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  return (
    <Box m="30px">
      <Box
        display="grid"
        gridTemplateColumns="repeat(12, 1fr)"
        gridAutoRows="140px"
        gap="20px"
      >
        <Box
          gridColumn="span 3"
          display="flex"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : colors.grey[900],
          }}
          justifyContent="center"
          alignItems="center"
        >
          <StatCard
            title="4"
            subtitle="Your Projects"
            fontColor={colors.greenAccent[500]}
            icon={
              <AdminPanelSettingsOutlinedIcon
                sx={{
                  color: colors.greenAccent[500],
                  fontSize: "26px",
                  margin: "5px",
                }}
              />
            }
            path="/projects/admin"
          />
        </Box>
        <Box
          gridColumn="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : colors.grey[900],
          }}
        >
          <StatCard
            title="4"
            subtitle="Managed Projects"
            fontColor={colors.blueAccent[500]}
            icon={
              <ManageAccountsOutlinedIcon
                sx={{
                  color: colors.blueAccent[500],
                  fontSize: "26px",
                  margin: "5px",
                }}
              />
            }
            path="/projects/manager"
          />
        </Box>
        <Box
          gridColumn="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : colors.grey[900],
          }}
        >
          <StatCard
            title="2"
            subtitle="Allocated Projects"
            fontColor={colors.redAccent[500]}
            icon={
              <PermIdentityOutlinedIcon
                sx={{
                  color: colors.redAccent[500],
                }}
              />
            }
            path="/projects/allocated"
          />
        </Box>
      </Box>
    </Box>
  );
};

export default Dashboard;
