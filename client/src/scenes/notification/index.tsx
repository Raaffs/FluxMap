import { useFetchUpdates } from "../../hooks/updates";
import { UpdateCard } from "../../components/cards/updates";
import {
  Box,
  Card,
  Typography,
  Button,
  Collapse,
  CardContent,
  useTheme,
  LinearProgress,
} from "@mui/material";
import { Update } from "../../hooks/types";
import { useState } from "react";
import { tokens } from "../../theme";
import { NoUpdates } from "../../components/cards/noUpdates";
const Updates = () => {
  const { updates, loading, error } = useFetchUpdates();
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  console.log("updates: ", updates);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(
    null
  );
  let updatesPerProject = new Map<number, Update[]>();
  const handleCardClick = (projectId: number) => {
    setSelectedProjectId((prev) => (prev === projectId ? null : projectId));
  };

  if (updates === null || updates === undefined ||updates.length===0) {
    console.log('nothing update')
    return <NoUpdates
        title="You're All Caught Up"
        description="There are currently no updates requiring your attention."
    />;
  }
    if (loading) {
      return (
        <Box
          display="flex"
          justifyContent="center"
          alignItems="center"
          height="90vh"
          width="80vw"
        >
          <Box sx={{ width: "50%" }}>
            <LinearProgress color="info" />
          </Box>
        </Box>
      );
    }
      
  for (const update of updates) {
    updatesPerProject.set(
      update.projectID,
      (updatesPerProject.get(update.projectID) ?? []).concat(update)
    );
  }
  console.log("mapppp", updatesPerProject);

  return (
    <Box
      sx={{
        margin: "20px",
        maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
        padding: "25px",
        borderRadius: "16px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white", // Soft light background for light mode
        textAlign: "left",
        border: "1px solid",
        borderColor: theme.palette.mode === "dark" ? "grey.800" : "#e0e7ff", // Subtle border in light mode
        boxShadow: theme.palette.mode === "dark" ? "0 4px 6px rgba(0, 0, 0, 0.1)" : "0 4px 8px rgba(0, 0, 0, 0.05)", // Soft shadow for depth
      }}
    >
      {Array.from(updatesPerProject.entries()).map(([projectId, updates], index) => {
        const { projectName, projectDescription } = updates[0];
        return (
          <Box
            key={projectId}
            sx={{
              paddingY: 3,
              borderBottom:
                index !== updatesPerProject.size - 1 ? "1px solid" : "none",
              borderColor: theme.palette.mode === "dark" ? "grey.800" : "#e0e7ff", // Soft divider line
            }}
          >
            <Typography
              variant="h5"
              sx={{
                fontWeight: "600",
                color: theme.palette.mode === "dark" ? "primary.main" : "#2c3e50", // Bold, vibrant heading for dark mode
                marginBottom: 1.5,
                textTransform: "capitalize",
                letterSpacing: "0.5px",
              }}
            >
              {projectName}
            </Typography>
            <Typography
              variant="body1"
              sx={{
                marginBottom: 2,
                color: "text.secondary",
                lineHeight: 1.8,
                fontSize: "1rem",
                fontStyle: "italic", // Adds a subtle flair
              }}
            >
              {projectDescription}
            </Typography>
            <Button
              onClick={() => handleCardClick(projectId)}
              variant="outlined"
              fullWidth
              sx={{
                borderRadius: 30, // More rounded button
                padding: "12px 24px",
                textTransform: "none",
                fontWeight: 600,
                backgroundColor: "white",
                color: "#1d72b8", // Subtle modern blue
                border: "1px solid #1d72b8",
                transition: "all 0.3s ease-in-out",
                "&:hover": {
                  backgroundColor: "#e6f7ff", // Lighter blue on hover
                  borderColor: "#1d72b8",
                },
                "&:focus": {
                  outline: "none",
                },
              }}
            >
              {selectedProjectId === projectId ? "Hide Updates" : "Show Updates"}
            </Button>
  
            <Collapse
              in={selectedProjectId === projectId}
              timeout="auto"
              unmountOnExit
            >
              <Box sx={{ marginTop: 2 }}>
                {updates.map((update) => (
                  <UpdateCard key={update.id} update={update} />
                ))}
              </Box>
            </Collapse>
          </Box>
        );
      })}
    </Box>
  );
  
  
};

export default Updates;
