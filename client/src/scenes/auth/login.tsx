import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Avatar,
  useTheme,
  Button,
  Card,
  Divider,
  CircularProgress,
  Stack,
} from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";
import ExploreIcon from "@mui/icons-material/Explore";
import { useNavigate } from "react-router-dom";
import { tokens } from "../../theme";
import { useAuth } from "../../context/authContext";
import { normalizeAccessMap } from "../../helpers/filter";
import { UserRole } from "../../hooks/types";

interface LoginProps {
  startWebSocket: () => void;
  setUserProjectRoleMap: React.Dispatch<React.SetStateAction<Record<number, UserRole>>>;
}

const LoginUser: React.FC<LoginProps> = ({ startWebSocket, setUserProjectRoleMap }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const navigate = useNavigate();
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  // 1. HANDLER: Redirect to Go Backend OAuth initiator
  const handleGoogleLogin = () => {
    window.location.href = "http://localhost:4000/api/auth/google";
  };

  // 2. HANDSHAKE: Logic for returning from Go after successful OAuth
  useEffect(() => {
    if (window.location.pathname === "/") {
      setLoading(true);
      fetch("http://localhost:4000/api/auth/session", { credentials: "include" })
        .then(async (res) => {
          if (!res.ok) throw new Error("Session invalid");
          return res.json();
        })
        .then((data) => {
          // data.roles is the map from your Go session
          const roles = data.roles as Record<UserRole, number[] | null>;
          const normal = normalizeAccessMap(roles);
          
          setUserProjectRoleMap(normal);
          localStorage.setItem("userProjectRoleMap", JSON.stringify(normal));
          
          login();
          startWebSocket();
          navigate("/");
        })
        .catch((err) => {
          console.error("Auth error:", err);
          setLoading(false);
          navigate("/login");
        });
    }
  }, [navigate, login, startWebSocket, setUserProjectRoleMap]);

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      sx={{
        width: "100%",
        height: "80vh", // Does NOT occupy full screen; keeps your top bar visible
        backgroundColor: "transparent",
      }}
    >
      <Card
        sx={{
          display: "flex",
          width: "900px",
          height: "550px",
          borderRadius: "24px",
          overflow: "hidden",
          boxShadow: "0px 20px 50px rgba(0, 0, 0, 0.3)",
          backgroundColor: theme.palette.mode === "dark" ? colors.primary[400] : "#fff",
        }}
      >
        {/* LEFT SIDE: ACTION AREA */}
        <Box
          sx={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            padding: "40px",
            backgroundColor: theme.palette.mode === "dark" ? "rgba(0,0,0,0.2)" : "#fafafa",
          }}
        >
          {loading ? (
            <Stack spacing={2} alignItems="center">
              <CircularProgress size={40} sx={{ color: colors.blueAccent[500] }} />
              <Typography variant="h6">Syncing Roles...</Typography>
            </Stack>
          ) : (
            <>
              <Avatar
                sx={{
                  width: 60,
                  height: 60,
                  bgcolor: colors.blueAccent[500],
                  mb: 2,
                  boxShadow: `0 0 20px ${colors.blueAccent[500]}66`,
                }}
              >
                <GoogleIcon fontSize="large" />
              </Avatar>
              
              <Typography variant="h3" fontWeight="bold" mb={1}>
                Welcome Back
              </Typography>
              <Typography variant="body1" color="textSecondary" mb={5} textAlign="center">
                Authenticate securely using your Google account.
              </Typography>

              <Button
                variant="contained"
                fullWidth
                startIcon={<GoogleIcon />}
                onClick={handleGoogleLogin}
                sx={{
                  py: 1.5,
                  borderRadius: "12px",
                  fontSize: "1rem",
                  fontWeight: "600",
                  textTransform: "none",
                  backgroundColor: theme.palette.mode === "dark" ? "#fff" : "#1a73e8",
                  color: theme.palette.mode === "dark" ? "#000" : "#fff",
                  "&:hover": {
                    backgroundColor: theme.palette.mode === "dark" ? "#e0e0e0" : "#1557b0",
                    transform: "translateY(-2px)",
                    boxShadow: "0 5px 15px rgba(0,0,0,0.3)",
                  },
                  transition: "all 0.2s ease-in-out",
                }}
              >
                Sign in with Google
              </Button>

              <Divider sx={{ width: "100%", my: 4 }}>
                <Typography variant="body2" color="textSecondary">
                  SECURE OAUTH 2.0
                </Typography>
              </Divider>
              
              <Typography variant="caption" color="textSecondary">
                Internal SSO Access Only
              </Typography>
            </>
          )}
        </Box>

        {/* RIGHT SIDE: BRANDING AREA (No image, solid technical aesthetic) */}
        <Box
          sx={{
            flex: 1.2,
            position: "relative",
            backgroundColor: colors.blueAccent[700], // Solid color instead of image
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            color: "#fff",
            p: 4,
          }}
        >
          <ExploreIcon sx={{ fontSize: 80, mb: 2, filter: "drop-shadow(0 0 10px rgba(255,255,255,0.3))" }} />
          <Typography variant="h1" fontWeight="800" letterSpacing={2} gutterBottom>
            FLUXMAP
          </Typography>
          <Typography
            variant="h5"
            textAlign="center"
            sx={{ opacity: 0.9, maxWidth: "300px", lineHeight: 1.6 }}
          >
            Real-time data visualization and project synchronization.
          </Typography>
          
          <Box
            sx={{
              position: "absolute",
              bottom: 20,
              right: 20,
              display: "flex",
              gap: 1
            }}
          >
            <Box sx={{ width: 8, height: 8, borderRadius: "50%", bgcolor: colors.greenAccent[500] }} />
            <Typography variant="caption" sx={{ color: "rgba(255,255,255,0.7)" }}>
              System Operational
            </Typography>
          </Box>
        </Box>
      </Card>
    </Box>
  );
};

export default LoginUser;