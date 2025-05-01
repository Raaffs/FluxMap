import React from "react";
import { Box, Typography, Paper, Card } from "@mui/material";
import { Update } from "../../hooks/types";
import { useTheme } from "@mui/material";
import { tokens } from "../../theme";
export const UpdateCard = ({ update }: { update: Update }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  return (
    <Card
      sx={{
        marginBottom: "20px",
        padding: "20px",
        borderRadius: "8px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "#fafafa",
        boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
        "&:hover": { boxShadow: "0 4px 8px rgba(0,0,0,0.8)" },
        textAlign: "left",
      }}
    >
      <Typography variant="h5" gutterBottom>
        Update created by{" "}
        <Typography
          component="span"
          fontWeight="bold"
          color={colors.blueAccent[500]}
          variant="h5"
        >
          {update.createdBy}
        </Typography>{" "}
        {/* for Project{" "}
        <Typography
          component="span"
          fontWeight="bold"
          color={colors.blueAccent[500]}
          variant="h5"
        >
          {update.projectName}
        </Typography>{" "} */}

      </Typography>
      <Typography
        variant="body1"
        color="text.primary"
        sx={{ marginTop: "10px", lineHeight: 1.6 }}
      >
        {update.msg}
      </Typography>

      <Typography
        variant="h6"
        sx={{
          marginTop: "16px",
        }}
      >
        Created at:{" "}
        <Typography
          component="span"
          fontWeight="bold"
          color={
            theme.palette.mode === "dark"
              ? colors.blueAccent[400]
              : colors.blueAccent[400]
          }
        >
          {update.createdAt}
        </Typography>
      </Typography>
    </Card>
  );
};
