import React, { useState } from "react";
import { useRetrieveProjectsFrom } from "../hooks/projects";
import { usePostProject } from "../hooks/projects";
import {
  Box,
  Typography,
  Button,
  TextField,
  useTheme,
  FormControl,
  Select,
  MenuItem,
} from "@mui/material";
import { tokens } from "../theme";
import LinearProgress from "@mui/material/LinearProgress";
import { useNavigate } from "react-router-dom";
import { Projects, UserRole } from "../hooks/types";
import CreateProjectModal from "./modals/createProject";
import PopUp from "../scenes/global/Popup";
type ProjectListProps = {
  URI: string;
  setUserProjectRoleMap: React.Dispatch<React.SetStateAction<Record<number, UserRole>>>;
};

export const ProjectListComponent = ({ URI, setUserProjectRoleMap }: ProjectListProps) => {
  const [fetchTrigger, setFetchTrigger] = useState<boolean>(false);
  const { projects, loading, error } = useRetrieveProjectsFrom(URI, fetchTrigger, setFetchTrigger);
  const [searchQuery, setSearchQuery] = useState("");
  const [sortOption, setSortOption] = useState("dueDate");
  const {
    postProject,
    loading: postLoading,
    error: postError,
    success,
  } = usePostProject(setUserProjectRoleMap);
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

  const filteredProjects = projects.filter((project) =>
    project.projectName.toLowerCase().includes(searchQuery.toLowerCase())
  );
  const handleSortChange = (event: any) => {
    setSortOption(event.target.value as string); // Cast the value to a string
  };
  const sortedAndFilteredProjects = filteredProjects.sort((a, b) => {
    if (sortOption === "dueDate") {
      // Ensure projectDueDate is a valid date string or Date object
      const dateA = a.projectDueDate ? new Date(a.projectDueDate) : new Date(0); // Default to an "epoch" date if not available
      const dateB = b.projectDueDate ? new Date(b.projectDueDate) : new Date(0); // Same here
      return dateA.getTime() - dateB.getTime(); // Compare the dates
    } else if (sortOption === "startDate") {
      // Ensure projectStartDate is a valid date string or Date object
      const startDateA = a.projectStartDate
        ? new Date(a.projectStartDate)
        : new Date(0); // Default to an "epoch" date if not available
      const startDateB = b.projectStartDate
        ? new Date(b.projectStartDate)
        : new Date(0); // Same here
      return startDateA.getTime() - startDateB.getTime(); // Compare the start dates
    } else if (sortOption === "projectName") {
      // Sort by project name
      return a.projectName.localeCompare(b.projectName);
    }
    return 0; // Default return if no sort option is matched
  });
  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setNewProject((prev) => ({ ...prev, [name]: value }));
  };

  const handleOpenModal = () => {
    setModalOpen(true);
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

    postProject(formattedProject, URI, setFetchTrigger);
    console.log("error here 2 ",postError)
    if (postError?.length === 0) {
      setModalOpen(false);
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
      {error && <PopUp Error={postError} Message="" onClose={() => {}} />}

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
      {error && <PopUp Error={postError} Message="" onClose={() => {}} />}

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
      {error && <PopUp Error={postError} Message="" onClose={() => {}} />}

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
          Looks like you don't have any projects here yet. Create one to get
          started and bring your ideas to life!
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
      {error && <PopUp Error={postError} Message="" onClose={() => {}} />}

        <CreateProjectModal
          open={modalOpen}
          onClose={handleCloseModal}
          newProject={newProject}
          handleInputChange={handleInputChange}
          handleSubmit={handleSubmit}
        />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        m: 3,
        p: 4,
        borderRadius: 2,
        bgcolor: theme.palette.mode === "dark" ? colors.primary[400] : "#fff",
        border: "1px solid",
        borderColor: theme.palette.mode === "dark" ? "#333" : "#e1e4e8",
        overflowY: "auto",
        boxShadow: 1,
      }}
    >
      {error && <PopUp Error={postError} Message="" onClose={() => {}} />}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <TextField
          variant="outlined"
          placeholder="Search projects"
          size="small"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          sx={{
            width: "60%",
            "& .MuiOutlinedInput-root": {
              borderRadius: "12px",
              "& fieldset": {
                borderColor: theme.palette.mode === "dark" ? "#666" : "#e1e4e8",
              },
            },
          }}
        />

        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          <FormControl sx={{ minWidth: 120, width: "auto" }}>
            <Select
              value={sortOption}
              onChange={handleSortChange}
              displayEmpty
              sx={{
                fontSize: "0.75rem",
                fontWeight: 500,
                borderRadius: "11px",
                width: "130px", // Make the width smaller, in line with the "New" button
                height: "46px",
                backgroundColor:
                  theme.palette.mode === "dark" ? "#388e3c" : "#4caf50", // Darker, richer green
                color: "#fff", // Keep text white for contrast
                "& .MuiSelect-icon": {
                  fontSize: "1.25rem",
                  color: "#fff", // White icon for contrast
                },
                "& .MuiOutlinedInput-notchedOutline": {
                  borderColor:
                    theme.palette.mode === "dark" ? "#444" : "#e1e4e8",
                },
                "&:hover": {
                  backgroundColor:
                    theme.palette.mode === "dark" ? "#2c6e2d" : "#388e3c", // Darker green on hover
                  borderColor: theme.palette.primary.main,
                },
                "&.Mui-focused": {
                  backgroundColor:
                    theme.palette.mode === "dark" ? "#2c6e2d" : "#388e3c", // Keep it darker when focused
                  borderColor: theme.palette.primary.main,
                },
                transition: "all 0.3s ease",
                boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                "& .MuiSelect-root": {
                  paddingRight: "2rem", // Room for the icon
                },
              }}
            >
              <MenuItem value="" disabled>
                Select Sort Option
              </MenuItem>{" "}
              {/* Placeholder */}
              <MenuItem value="dueDate">Due Date</MenuItem>
              <MenuItem value="projectName">Project Name</MenuItem>
              <MenuItem value="startDate">Created Date</MenuItem>
            </Select>
          </FormControl>

          <Button
            variant="contained"
            color="primary"
            onClick={handleOpenModal}
            sx={{
              fontWeight: 600,
              textTransform: "none",
              whiteSpace: "nowrap",
              width: "120px",
              borderRadius: "12px",
              backgroundColor: "royalblue",
              px: 3,
              py: 1.25,
              boxShadow: 2,
            }}
          >
            + New
          </Button>
        </Box>
      </Box>

      {sortedAndFilteredProjects.map((project, index) => (
        <Box
          key={index}
          onClick={() => navigate(`/project/${project.projectID}`)}
          sx={{
            display: "flex",
            // borderColor: theme.palette.mode==="dark"?"grey.800":"#d0d7de",
            flexDirection: "column",
            textAlign: "left",
            alignItems: "flex-start",
            justifyContent: "flex-start",
            borderBottom:
              theme.palette.mode === "dark"
                ? "1px solid #2E3C57"
                : "1px solid #e1e4e8",
            py: 2,
            px: 2,
            cursor: "pointer",
            transition: "background-color 0.2s ease",
            borderRadius: "10px",
            "&:hover": {
              backgroundColor:
                theme.palette.mode === "dark"
                  ? colors.blueAccent[700]
                  : "#f6f8fa",
            },
          }}
        >
          <Typography
            variant="h6"
            fontWeight={600}
            sx={{ color: "text.primary" }}
          >
            {project.projectName}
          </Typography>

          <Typography
            variant="body2"
            color="text.secondary"
            sx={{ mt: 0.5, mb: 1 }}
          >
            {project.projectDescription || "No description available"}
          </Typography>

          <Typography
            variant="caption"
            sx={{
              color:
                project.projectDueDate &&
                new Date(project.projectDueDate) < new Date()
                  ? colors.redAccent[500]
                  : colors.greenAccent[400],
              fontWeight: 500,
              fontStyle: "italic",
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
        </Box>
      ))}

      <Button
        variant="contained"
        color="primary"
        onClick={handleOpenModal}
        sx={{
          mt: 4,
          width: "100%",
          py: 1.5,
          fontWeight: 600,
          boxShadow: 3,
          textTransform: "none",
          borderRadius: "12px",
          backgroundColor: "royalblue",
        }}
      >
        + Create New Project
      </Button>

      <CreateProjectModal
        open={modalOpen}
        onClose={handleCloseModal}
        newProject={newProject}
        handleInputChange={handleInputChange}
        handleSubmit={handleSubmit}
      />
    </Box>
  );
};
