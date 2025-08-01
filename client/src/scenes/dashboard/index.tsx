import { Box, useTheme } from "@mui/material";
import AdminPanelSettingsOutlinedIcon from "@mui/icons-material/AdminPanelSettingsOutlined";
import ManageAccountsOutlinedIcon from "@mui/icons-material/ManageAccountsOutlined";
import PermIdentityOutlinedIcon from "@mui/icons-material/PermIdentityOutlined";
import { tokens } from "../../theme";
import StatCard from "../../components/cards/statcard";
const Dashboard = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  return (
    <Box m="30px"
        sx={{
            backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[500] : "white",
                    maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
            padding: "25px",
            borderRadius: "16px",
        }}
    >
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
                : "rgba(0, 200, 150, 0.3)",
          }}
          justifyContent="center"
          alignItems="center"
        >
          <StatCard
            title="4"
            subtitle="Your Projects"
            fontColor={colors.greenAccent[200]}
            icon={
              <AdminPanelSettingsOutlinedIcon
                sx={{
                  color: colors.greenAccent[200],
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
                : "rgba(104, 112, 250, 0.4)",
          }}
        >
          <StatCard
            title="4"
            subtitle="Managed Projects"
            fontColor={colors.blueAccent[400]}
            icon={
              <ManageAccountsOutlinedIcon
                sx={{
                  color: colors.blueAccent[400],
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
                : "rgba(162, 89, 255, 0.3)",
          }}
        >
          <StatCard
            title="2"
            subtitle="Allocated Projects"
            fontColor="#a259ff"
            icon={
              <PermIdentityOutlinedIcon
                sx={{
                  color: "#a259ff"
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
