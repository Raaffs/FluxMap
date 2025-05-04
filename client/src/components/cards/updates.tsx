import React from "react";
import { Box, Typography, Paper, Card } from "@mui/material";
import { Update } from "../../hooks/types";
import { useTheme } from "@mui/material";
import { tokens } from "../../theme";
export const UpdateCard = ({ update }: { update: Update }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  return (
    <Box
      sx={{
        padding: "16px 0",
        borderBottom: "1px solid",
        borderColor: theme.palette.mode === "dark" ? "grey.800" : "#d0d7de",
        textAlign: "left",
      }}
    >
      <Typography variant="body1" fontSize='1rem' gutterBottom>
        Update created by{" "}
        <Typography
          component="span"
          fontWeight="bold"
          color={colors.blueAccent[500]}
          variant="body1"
          fontSize='1rem'
        >
          {update.createdBy}
        </Typography>
      </Typography>

      <Typography
        variant="body2"
        color="text.primary"
        sx={{ marginTop: "4px", lineHeight: 1.6, fontSize: "1rem" }}
      >
        {update.msg}
      </Typography>

      <Typography
        variant="caption"
        sx={{
          marginTop: "12px",
          display: "block",
          color: "text.secondary",
          fontSize:'0.9rem'
        }}
      >
        Created at:{" "}
        <Typography
          component="span"
          fontWeight="bold"
          color={colors.blueAccent[400]}
        >
          {new Date(update.createdAt).toLocaleDateString("en-US", {
            year: "numeric",
            month: "long",
            day: "numeric",
          })}
        </Typography>
      </Typography>
    </Box>
  );
};
