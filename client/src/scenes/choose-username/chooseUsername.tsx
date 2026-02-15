import React, { useState } from "react";
import {
  Box,
  Button,
  Card,
  TextField,
  Typography,
  useTheme,
  Stack,
  Fade,
  InputAdornment,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import AlternateEmailIcon from "@mui/icons-material/AlternateEmail";
import { tokens } from "../../theme";

const ChooseUsername = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [username, setUsername] = useState("");
  const [backendError, setBackendError] = useState("");
  const navigate = useNavigate();

  // Real-time validation check
  const isTooShort = username.length > 0 && username.length < 3;

  const handleFinalize = async () => {
    if (username.length < 3) return;

    try {
      const res = await fetch("http://localhost:4000/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: username.trim() }),
        credentials: "include",
      });

      if (res.ok) {
        window.location.href = "/"; // Force a full refresh      } else {
        const data = await res.json();
        setBackendError(data.error || "Username taken.");
      }
    } catch (err) {
      setBackendError("Server connection failed.");
    }
  };

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      sx={{ height: "80vh", width: "100%" }}
    >
      <Fade in={true} timeout={800}>
        <Card
          elevation={0}
          sx={{
            p: 6,
            width: "440px",
            textAlign: "center",
            borderRadius: "24px",
            border: `1px solid ${colors.grey[800]}`,
            bgcolor: theme.palette.mode === "dark" ? "#0A0A0A" : "#fff",
            boxShadow: "0 40px 80px -12px rgba(0,0,0,0.7)",
          }}
        >
          <Stack spacing={1} mb={5}>
            <Typography
              variant="h2"
              fontWeight="900"
              sx={{ letterSpacing: "-2px", textTransform: "uppercase" }}
            >
              USERNAME
            </Typography>
            <Typography
              variant="body2"
              sx={{
                color: colors.grey[400],
                letterSpacing: "2px",
                fontWeight: 600,
              }}
            >
              Register your unique identity
            </Typography>
          </Stack>

          <TextField
            fullWidth
            autoFocus
            value={username}
            onChange={(e) => {
              setUsername(e.target.value.toLowerCase().replace(/\s/g, ""));
              setBackendError("");
            }}
            placeholder="username"
            // Displays error if too short OR if backend rejects it
            error={isTooShort || !!backendError}
            helperText={
              isTooShort
                ? "Minimum 3 characters required"
                : backendError
                  ? backendError
                  : " " // Keeps constant height so UI doesn't jump
            }
            variant="outlined"
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <AlternateEmailIcon sx={{ color: "#1E88E5" }} />
                </InputAdornment>
              ),
            }}
            sx={{
              mb: 2,
              "& .MuiOutlinedInput-root": {
                borderRadius: "12px",
                "& fieldset": { borderColor: colors.grey[700] },
                "&:hover fieldset": { borderColor: "#1E88E5" },
                "&.Mui-focused fieldset": { borderColor: "#1E88E5" },
              },
              "& .MuiFormHelperText-root": {
                fontWeight: 600,
                marginTop: "8px",
                height: "20px",
              },
            }}
          />

          <Button
            fullWidth
            onClick={handleFinalize}
            disabled={username.length < 3}
            variant="contained"
            sx={{
              py: 2,
              borderRadius: "12px",
              fontSize: "1rem",
              fontWeight: 900,
              textTransform: "uppercase",
              letterSpacing: "1.5px",

              // Your working blue logic
              backgroundColor: "#1E88E5 !important",
              color: "#fff !important",
              boxShadow: "0px 4px 20px rgba(30, 136, 229, 0.4)",
              transition: "all 0.2s ease-in-out",

              "&:hover": {
                backgroundColor: "#1565C0 !important",
                transform: "translateY(-2px)",
                boxShadow: "0px 8px 25px rgba(30, 136, 229, 0.6)",
              },

              "&:active": {
                transform: "translateY(0px)",
              },

              "&.Mui-disabled": {
                backgroundColor: "#1E88E5 !important",
                color: "#ffffff80 !important",
                opacity: 0.6,
                boxShadow: "none",
              },
            }}
          >
            Finalize Connection
          </Button>

          <Box sx={{ mt: 5, opacity: 0.4 }}>
            <Typography
              variant="caption"
              sx={{
                color: colors.grey[500],
                fontFamily: "monospace",
                letterSpacing: "1px",
              }}
            >
              This username will be your unique identifier across the FluxMap.
            </Typography>
          </Box>
        </Card>
      </Fade>
    </Box>
  );
};

export default ChooseUsername;
