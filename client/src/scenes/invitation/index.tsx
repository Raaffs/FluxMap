import { Box, Button, Card, Typography, useTheme } from "@mui/material";
import { tokens } from "../../theme";
export const Invitation = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  return (
    <Box
      sx={{
        margin: "5px",
        maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
        padding: "16px",
        border: "5px #ddd",
        borderRadius: "8px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white",
        boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        alignContent: "left",
      }}
    >
      <Card
        sx={{
          minHeight: "15%",
          maxWidth: "100%",
          marginBottom: "16px",
          padding: "16px",
          borderRadius: "8px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[400] : "#f",
          boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
          alignContent: "left",
          textAlign: "left",
          margin: 5,
          display: "flex",
          flexDirection: "column", // Makes sure elements inside stack vertically
          "&:hover": {
            boxShadow: "0 4px 8px rgba(0,0,0,0.8)",
          },
        }}
      >
        {/* Invite Message */}
        <Typography variant="h5">
          You've been invited to{" "}
          <Typography component="span" fontWeight="bold" variant="h5">
            ProjectName
          </Typography>{" "}
          by{" "}
          <Typography component="span" fontWeight="bold" variant="h5">
            OwnerName
          </Typography>
        </Typography>
        <Typography variant="h5" color="text.secondary">
            Description of the project is here jfjefjwejfiwejfiewjfijewfiojewfiojweifojweoifjwefioewjfoiewjfioejwfioewjfoijewfo
        </Typography>

        {/* Buttons Container */}
        <Box
          sx={{
            display: "flex",
            justifyContent: "left", // Align buttons horizontally in center
            gap: 1, // Adds space between buttons
            marginTop: "auto", // Pushes buttons to the bottom
          }}
        >
          <Button
            variant="contained"
            color="primary"
            sx={{
              backgroundColor: colors.greenAccent[500],
            }}
          >
            Accept
          </Button>
          <Button
            variant="contained"
            color="primary"
            sx={{
              backgroundColor: colors.redAccent[500],
            }}
          >
            Decline
          </Button>
        </Box>
      </Card>
    </Box>
  );
};
