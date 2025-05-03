import { useEffect, useState } from "react";
import { ProjectTaskDetailPage } from "../../components/task";
import { useParams } from "react-router-dom";
import { Projects, tasks } from "../../hooks/types";
import { Tab, Tabs, Box, Card, useTheme, Typography, Button, TextField, Modal } from "@mui/material";
import { tokens } from "../../theme";
import Graphs from "../../components/graphs/LineGraphs";
import PertNormalDistributionChart, { PertTable } from "../../components/pert";
import {
  ContributerTaskChart,
  ContributerTaskTable,
} from "../../components/contributer";
import useFetchTaskData from "../../hooks/task";
import CpmNormalDistributionChart from "../../components/cpm";
import { useFetchPertData } from "../../hooks/pert";
import { useFetchCpmData } from "../../hooks/cpm";
import CreateTaskModal from "../../components/modals/createTask";
export const ProjectOverview = () => {
  const { id } = useParams();
  const [openNewTaskModal, setOpenNewTaskModal] = useState(false);
    const [newTask, setNewTask] = useState<tasks>({
      taskID: 0, // Set to 0 or another default value if necessary
      taskName: "",
      taskDescription: "",
      taskStatus: "",
      taskStartDate: null,
      taskDueDate: null,
      parentProjectID: Number(id), // Assuming the project ID is available
      assignedUsername: "",
      approved: false,
      taskCompletedDate: null,
      taskApprovedDate: null,
      createdBy:""
    });
  
  const [project, setProject] = useState<Projects | null>(null);
  const [apiResponse, pertLoading, pertError] = useFetchPertData(id);
  const [cpmApiResponse, cpmLoading, cpmError] = useFetchCpmData(id);
  const [fetchTrigger, setFetchTrigger] = useState(false);
  console.log("PERT API: ", apiResponse);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const [selectedTab, setSelectedTab] = useState<number>(0);
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [tasks, fetchLoading, fetchError] = useFetchTaskData(id, fetchTrigger,setFetchTrigger);
  useEffect(() => {
    const fetchProjectOverviewByID = async () => {
      try {
        const response = await fetch(
          `http://localhost:4000/api/project/${id}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
          }
        );
        const data = await response.json();
        console.log("data : ", data);
        setProject(data);
        setLoading(false);
      } catch (err) {
        console.log(err)
        setError(err);
        setLoading(false);
      }
    };
    fetchProjectOverviewByID();
  }, [id]);
  const handleOpenModal = () => {
    console.log("opened");
    setModalOpen(true);
    console.log(modalOpen);
  };
  const handleOpenNewTaskModal = () => setOpenNewTaskModal(true);  
  const handleCloseNewTaskModal = () => {
    setOpenNewTaskModal(false);
    setError(null); // Reset errors when closing modal
  };



  const handleCreateNewTask = async () => {
    try {
      if (!newTask.taskName || !newTask.assignedUsername) {
        setError("Task Name and Assigned User are required.");
        return;
      }
      const formattedTask = {
        ...newTask,
        taskDueDate: newTask.taskDueDate
          ? new Date(newTask.taskDueDate).toISOString()
          : null,
        taskStartDate: newTask.taskStartDate
          ? new Date(newTask.taskStartDate).toISOString()
          : null,
      };
      const response = await fetch(
        `http://localhost:4000/api/project/${id}/task`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify(formattedTask),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        console.error(errorData);
        setError(
          errorData.error || "Failed to save the task. Please try again."
        );
      } else {
        setFetchTrigger(true);
        console.log("created task new task");
      }
    } catch (err) {
      console.error("Failed to create task", err);
    }finally{
      handleCloseNewTaskModal();
    }
  };


  if (tasks === null) {
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
          No Task available.
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
          Looks like you don't have any task yet. Create one to get started
          and bring your ideas to life!
        </Typography>
        <Button
          variant="contained"
          color="secondary"
          size="large"
          onClick={handleOpenNewTaskModal}
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
          Create a New Task
        </Button>
        <CreateTaskModal
          open={openNewTaskModal}
          onClose={handleCloseNewTaskModal}
          newTask={newTask}
          handleCreateNewTask={handleCreateNewTask}
          setNewTask={setNewTask}
        />
      </Box>
    );

  }

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  return (
    <Box
      sx={{
        overflowY: "auto",
        margin: "5px",
        width: "100%",
        minHeight: "100%",
        padding: "16px",
        border: "5px #ddd",
        borderRadius: "8px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white",
        boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        alignContent: "left",
        overflowX: "auto",
      }}
    >
      <Card
        sx={{
          width: "99%",
          minHeight: "20%",
          marginBottom: "16px",
          padding: "16px",
          borderRadius: "8px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[400] : "#f",
          boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
          alignContent: "left",
          alignItems: "left",
          textAlign: "left",
          cursor: "pointer", // This makes the cursor a pointer on hover, indicating it's clickable
          "&:hover": {
            boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
          },
        }}
      >
        <Typography variant="h4" sx={{ fontWeight: "bold" }}>
          {project?.projectName}
        </Typography>
        <Typography variant="h5" color="textSecondary">
          {project?.projectDescription}
        </Typography>
        <Typography
          variant="h6"
          color="text.secondary"
          sx={{
            color:
              project?.projectDueDate &&
              new Date(project.projectDueDate) < new Date()
                ? colors.redAccent[500]
                : colors.greenAccent[400],
            fontWeight: "bold",
          }}
        >
          <strong>Due:</strong>{" "}
          {project?.projectDueDate
            ? new Date(project.projectDueDate).toLocaleDateString("en-GB", {
                day: "2-digit",
                month: "short",
                year: "numeric",
              })
            : "No due date"}
        </Typography>
      </Card>

      {/* Tabs Component */}
      <Tabs
        value={selectedTab}
        onChange={handleTabChange}
        aria-label="project tabs"
      >
        <Tab label="Task Details" />
        <Tab label="Contributer" />
        <Tab label="PERT" />
        <Tab label="CPM" />
      </Tabs>

      {/* Tab Panels */}
      <Box sx={{ padding: "16px" }}>
        <Card
          sx={{
            maxWidth: "99%",
            minHeight: "20%",
            marginBottom: "16px",
            padding: "16px",
            borderRadius: "8px",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[400] : "#f",
            boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
            alignContent: "left",
            alignItems: "left",
            textAlign: "left",
            cursor: "pointer", // This makes the cursor a pointer on hover, indicating it's clickable
            "&:hover": {
              boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
            },
          }}
        >
          {selectedTab === 0 && (
            <Box>
              <ProjectTaskDetailPage
                tasks={tasks}
                setFetchTrigger={setFetchTrigger}
              />
              <Graphs />
            </Box>
          )}
          {selectedTab === 1 && (
            <Box>
              <ContributerTaskTable tasksData={tasks} />
              <ContributerTaskChart key={tasks.length} tasksData={tasks} />
            </Box>
          )}
          {selectedTab === 2 && (
            <Box>
              <PertTable pertTasks={apiResponse?.data} tasks={tasks} />
              <PertNormalDistributionChart
                apiResponse={apiResponse}
                pertTasks={apiResponse?.data}
                tasks={tasks}
              />
            </Box>
          )}
          {selectedTab === 3 && (
            <Box>
              <CpmNormalDistributionChart data={cpmApiResponse} />
            </Box>
          )}
        </Card>
      </Box>
    </Box>
  );
};
