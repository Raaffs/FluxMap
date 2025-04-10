import React, { useState } from "react";
import { useRetrieveProjectsFrom } from "../hooks/projects";
import { usePostProject } from "../hooks/projects";
import {
  Card,
  Box,
  Typography,
  Button,
  Modal,
  TextField,
  useTheme,
} from "@mui/material";
import { tokens } from "../theme";
import LinearProgress from "@mui/material/LinearProgress";
import { useNavigate } from "react-router-dom";
import { Projects } from "../hooks/types";
// export interface Projects {
//   projectID?: number; // Matches `omitempty`
//   projectName: string; // Required field
//   projectDescription?: string | null; // Matches `null.String`
//   projectStartDate?: string | null; // Matches `null.Time`
//   projectDueDate?: string | null; // Matches `null.Time`
//   ownername: string; // Required field
// }
export const ProjectComponent = ({ URI }: { URI: string }) => {
  const { projects, loading, error } = useRetrieveProjectsFrom(URI);

  const {
    postProject,
    loading: postLoading,
    error: postError,
    success,
  } = usePostProject();
  const [modalOpen, setModalOpen] = useState(false);
  const [newProject, setNewProject] = useState<Projects>({
    projectID: 0,
    projectName: "",
    projectDescription: "",
    projectStartDate: null,
    projectDueDate: null,
    ownername: "",
  });
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const navigate = useNavigate();

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setNewProject((prev) => ({ ...prev, [name]: value }));
  };

  const handleOpenModal = () => {
    console.log("opened");
    setModalOpen(true);
    console.log(modalOpen);
  };
  const handleCloseModal = () => setModalOpen(false);

  const handleSubmit = () => {
    const formattedProject = {
      ...newProject,
      projectDueDate: newProject.projectDueDate
        ? new Date(newProject.projectDueDate).toISOString()
        : null,
      projectStartDate: newProject.projectStartDate
        ? new Date(newProject.projectStartDate).toISOString()
        : null,
    };

    postProject(formattedProject, URI);

    if (!postError) {
      handleCloseModal(); // Close modal on successful creation
    }
  };

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
    return (
      <Box
        sx={{
          margin: "20px",
          padding: "24px",
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          textAlign: "center",
          height: "60vh",
          borderRadius: "12px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.grey[800] : "white",
          boxShadow: "0 6px 12px rgba(0,0,0,0.1)",
          border: `2px solid ${theme.palette.divider}`,
        }}
      >
        <Typography
          variant="h4"
          color="error"
          sx={{
            marginBottom: "16px",
            fontWeight: "bold",
          }}
        >
          Error Loading Projects
        </Typography>
        <Typography
          variant="body1"
          color="text.secondary"
          sx={{
            marginBottom: "24px",
            maxWidth: "400px",
            lineHeight: "1.5",
            fontSize: "1.2rem",
          }}
        >
          {error}
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => window.location.reload()}
        >
          Try Again
        </Button>
      </Box>
    );
  }

  if (!projects.length) {
    return (
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          textAlign: "center",
          height: "60vh", // Increased height for more space
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[500] : "white",
        }}
      >
        <Typography
          variant="h4"
          color="text.secondary"
          sx={{
            marginBottom: "16px",
            color: theme.palette.mode === "dark" ? "orangered" : "crimson",
            fontWeight: "bold",
          }}
        >
          No projects available.
        </Typography>
        <Typography
          variant="body1"
          color="text.secondary"
          sx={{
            marginBottom: "24px",
            maxWidth: "400px",
            lineHeight: "1.5",
            fontSize: "1.2rem",
            fontStyle: "italic",
          }}
        >
          Looks like you don't have any projects yet. Create one to get started
          and bring your ideas to life!
        </Typography>
        <Button
          variant="contained"
          color="secondary"
          size="large"
          onClick={handleOpenModal}
          sx={{
            backgroundColor: "royalblue",
            padding: "12px 32px",
            fontSize: "1.1rem",
            borderRadius: "50px",
            boxShadow: "0 6px 12px rgba(0,0,0,0.2)",
            "&:hover": {
              backgroundColor: "darkorange",
              boxShadow: "0 8px 16px rgba(0,0,0,0.3)",
            },
          }}
        >
          Create a New Project
        </Button>
        {/* I should really make a separate component for this modal */}
        <Modal open={modalOpen} onClose={handleCloseModal}>
          <Box
            sx={{
              width: "400px",
              padding: "16px",
              backgroundColor: "white",
              borderRadius: "8px",
              boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
              margin: "auto",
              marginTop: "10%",
            }}
          >
            <Typography variant="h6" sx={{ marginBottom: "16px" }}>
              Create New Project
            </Typography>
            <TextField
              fullWidth
              label="Project Name"
              name="projectName"
              value={newProject.projectName || ""}
              onChange={handleInputChange}
              sx={{ marginBottom: "16px" }}
            />
            <TextField
              fullWidth
              label="Description"
              name="projectDescription"
              value={newProject.projectDescription || ""}
              onChange={handleInputChange}
              sx={{ marginBottom: "16px" }}
            />
            <TextField
              fullWidth
              label="Due Date"
              name="projectDueDate"
              type="date"
              variant="outlined"
              value={newProject.projectDueDate || ""}
              onChange={handleInputChange}
              sx={{ marginBottom: "16px" }}
            />
            <Button
              fullWidth
              variant="contained"
              color="primary"
              onClick={handleSubmit}
              disabled={postLoading}
            >
              Create
            </Button>
            {postError && (
              <Typography color="error" sx={{ marginTop: "16px" }}>
                {postError}
              </Typography>
            )}
          </Box>
        </Modal>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        margin: "5px",
        maxHeight: "100%",
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
      {projects.map((project, index) => (
        <Card
          key={index}
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
            alignItems: "left",
            textAlign: "left",
            cursor: "pointer",
            "&:hover": {
              boxShadow: "0 4px 8px rgba(0,0,0,0.8)",
            },
          }}
          onClick={() => navigate(`/project/${project.projectID}`)}
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
              color:
                project.projectDueDate &&
                new Date(project.projectDueDate) < new Date()
                  ? colors.redAccent[500]
                  : colors.greenAccent[400],
              fontWeight: "bold",
            }}
          >
            <strong>Due:</strong>{" "}
            {project.projectDueDate
              ? new Date(project.projectDueDate).toLocaleDateString("en-GB", {
                  day: "2-digit",
                  month: "short",
                  year: "numeric",
                })
              : "No due date"}
          </Typography>
        </Card>
      ))}
      <Button
        variant="contained"
        color="primary"
        onClick={handleOpenModal}
        sx={{
          marginTop: "16px",
          display: "block",
        }}
      >
        Create New Project
      </Button>

      <Modal open={modalOpen} onClose={handleCloseModal}>
        <Box
          sx={{
            width: "400px",
            padding: "16px",
            backgroundColor: "white",
            borderRadius: "8px",
            boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
            margin: "auto",
            marginTop: "10%",
          }}
        >
          <Typography variant="h6" sx={{ marginBottom: "16px" }}>
            Create New Project
          </Typography>
          <TextField
            fullWidth
            label="Project Name"
            name="projectName"
            value={newProject.projectName || ""}
            onChange={handleInputChange}
            sx={{ marginBottom: "16px" }}
          />
          <TextField
            fullWidth
            label="Description"
            name="projectDescription"
            value={newProject.projectDescription || ""}
            onChange={handleInputChange}
            sx={{ marginBottom: "16px" }}
          />
          <TextField
            fullWidth
            label="Due Date"
            name="projectDueDate"
            type="date"
            variant="outlined"
            value={newProject.projectDueDate || ""}
            onChange={handleInputChange}
            sx={{ marginBottom: "16px" }}
          />
          <Button
            fullWidth
            variant="contained"
            color="primary"
            onClick={handleSubmit}
            disabled={postLoading}
          >
            Create
          </Button>
          {postError && (
            <Typography color="error" sx={{ marginTop: "16px" }}>
              {postError}
            </Typography>
          )}
        </Box>
      </Modal>
    </Box>
  );
};
