import { useRetrieveProjectsFrom } from "../hooks/projects";
import { Card, Box, Typography, useTheme } from "@mui/material";
import { tokens } from "../theme";
import LinearProgress from '@mui/material/LinearProgress';
import { useNavigate } from "react-router-dom";  // Import useNavigate

export const ProjectComponent = ({ URI }: { URI: string }) => {
  const { projects, loading, error } = useRetrieveProjectsFrom(URI);
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const navigate = useNavigate();  // Initialize navigate hook

  if (loading) {
    return (
      <Box
        sx={{
          display: "flex", 
          justifyContent: "center",
          alignItems: "center",
          minHeight: "50vh", 
        }}
      >
        <LinearProgress
          color="success"
          sx={{
            width: "50%", 
            height: "5px", 
          }}
        />
      </Box>
    );
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!projects.length) {
    return (
      <Box
        sx={{
          margin: '5px',
          overflowY: "auto",
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          textAlign: 'center',
          height: "50vh",
          padding: "16px",
          border: "5px #ddd",
          borderRadius: "8px",
          backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : 'white',
          boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        }}
      >
        <Typography variant="h4" color="text.secondary" sx={{ color: "orangered" }}>
          No projects available.
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        margin: '5px',
        minHeight: "100%",
        overflowY: "auto",
        padding: "16px",
        border: "5px #ddd",
        borderRadius: "8px",
        backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : 'white',
        boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        alignContent: 'left'
      }}
    >
      {projects.map((project, index) => (
        <Card
          key={index}
          sx={{
            maxWidth: '100%',
            minHeight: '20%',
            marginBottom: "16px",
            padding: "16px",
            borderRadius: "8px",
            backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : '#f',
            boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
            alignContent: 'left',
            alignItems: 'left',
            textAlign: 'left',
            cursor: 'pointer', // This makes the cursor a pointer on hover, indicating it's clickable
            '&:hover': {
            boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
          }

          }}
          
          onClick={() => navigate(`/project/${project.projectID}`)} // Navigate to the project detail page on click
        >
          <Typography variant="h5" sx={{ fontWeight: "bold" }}>
            {project.projectName}
          </Typography>
          <Typography
            variant="h6"
            color="text.secondary"
            sx={{ marginBottom: "8px" }}
          >
            {project.projectDescription || "No description available"}
          </Typography>
          <Typography
            variant="h6"
            color="text.secondary"
            sx={{
              color: new Date(project.projectDueDate) < new Date() ? colors.redAccent[500] : colors.greenAccent[400],
              fontWeight: 'bold'
            }}
          >
            <strong>Due:</strong> {new Date(project.projectDueDate).toLocaleDateString('en-GB', {
              day: '2-digit',
              month: 'short',
              year: 'numeric'
            })}
          </Typography>
        </Card>
      ))}
    </Box>
  );
};
