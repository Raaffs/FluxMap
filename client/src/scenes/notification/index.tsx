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
        margin: "10px",
        maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
        padding: "20px",
        borderRadius: "12px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white",
        boxShadow: "0 8px 20px rgba(0,0,0,0.1)",
        textAlign: "left",
      }}
    >
      {Array.from(updatesPerProject.entries()).map(([projectId, updates]) => {
        const { projectName, projectDescription } = updates[0];
        return (
          <Card
            key={projectId}
            sx={{
              marginBottom: 4,
              boxShadow: 5,
              borderRadius: 4,
              transition: "all 0.3s ease-in-out",
              "&:hover": { boxShadow: 8 },
            }}
          >
            <CardContent sx={{ padding: 4 }}>
              <Typography
                variant="h5"
                sx={{
                  fontWeight: "600",
                  color: "primary.main",
                  marginBottom: 1,
                  textTransform: "uppercase",
                  letterSpacing: "0.5px",
                }}
              >
                {projectName}
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  marginBottom: 3,
                  color: "text.secondary",
                  lineHeight: 1.6,
                  fontSize: "1rem",
                }}
              >
                {projectDescription}
              </Typography>
              <Button
                onClick={() => handleCardClick(projectId)}
                variant="outlined"
                color="primary"
                fullWidth
                sx={{
                  marginBottom: 3,
                  borderRadius: 4,
                  padding: "12px 24px",
                  textTransform: "none",
                  fontWeight: 600,
                  border: "2px solid #4169E1", // Royal Blue border color for default
                  backgroundColor: "#4169E1", // Solid Royal Blue background
                  color: "white", // White text color for default
                  transition: "all 0.3s ease",
                  "&:hover": {
                    backgroundColor: "#20B2AA", // Soft teal for hover
                    borderColor: "#20B2AA", // Matching border color
                    color: "white", // Text color stays white on hover
                    boxShadow: "0 4px 20px rgba(32, 178, 170, 0.5)", // Soft glowing effect
                    transform: "scale(1.005)", // Slightly enlarge the button on hover
                  },
                  "&:focus": {
                    outline: "none",
                  },
                }}
              >
                {selectedProjectId === projectId
                  ? "Hide Updates"
                  : "Show Updates"}
              </Button>

              <Collapse
                in={selectedProjectId === projectId}
                timeout="auto"
                unmountOnExit
              >
                <Box sx={{ marginTop: 3 }}>
                  {updates.map((update) => (
                    <UpdateCard key={update.id} update={update} />
                  ))}
                </Box>
              </Collapse>
            </CardContent>
          </Card>
        );
      })}
    </Box>
  );
};

export default Updates;
