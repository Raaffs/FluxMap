// src/scenes/NotFound.tsx
import React from "react";
import { Box, Typography, Button, useTheme } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { tokens } from "../../theme";

const NotFound: React.FC = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  // Dynamic text color based on theme mode
  const textColor = theme.palette.mode === "dark" ? "#fff" : "#000";

  return (
    <Box
      sx={{
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        textAlign: "center",
        backgroundColor: colors.primary[400], // optional background from your theme
        color: textColor,
        px: 2,
      }}
    >
      <Typography
        variant="h1"
        sx={{ fontSize: "8rem", fontWeight: "bold", mb: 2, color: textColor }}
      >
        404
      </Typography>

      <Typography variant="h4" sx={{ mb: 2, color: textColor }}>
        Oops! Page Not Found
      </Typography>

      <Typography variant="body1" sx={{ mb: 4, maxWidth: "400px", color: textColor }}>
        The page you are looking for might have been removed, had its name changed,
        or is temporarily unavailable.
      </Typography>

      <Button
        variant="contained"
        color="secondary"
        sx={{
          px: 4,
          py: 1,
          fontSize: "1rem",
          borderRadius: "8px",
          textTransform: "none",
          backgroundColor: "#ff6b6b",
          "&:hover": { backgroundColor: "#ff4c4c" },
        }}
        onClick={() => navigate("/")}
      >
        Go Back Home
      </Button>
    </Box>
  );
};

export default NotFound;
