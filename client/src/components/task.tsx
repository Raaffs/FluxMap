import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  Box,
  Typography,
  MenuItem,
  Select,
  FormControl,
  Checkbox,
  Modal,
  TextField,
  Button,
  useTheme,
} from "@mui/material";
import EditIcon from '@mui/icons-material/Edit';
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { tokens } from "../theme";
import { tasks } from "../hooks/types";
import CreateTaskModal from "./modals/createTask";
import PopUp from "../scenes/global/Popup";
export const ProjectTaskDetailPage = ({
  tasks,
  //to re-renders the component on adding/modifying a task
  setFetchTrigger,
}: {
  tasks: tasks[];
  setFetchTrigger: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const { id } = useParams(); // Get the project id from the URL
  const [openTaskNameModal, setOpenTaskNameModal] = useState(false);
  const [openDescriptionModal, setOpenDescriptionModal] = useState(false);
  const [currentTaskName, setCurrentTaskName] = useState("");
  const [currentDescription, setCurrentDescription] = useState("");
  const [selectedTaskID, setSelectedTaskID] = useState<number | null>(null);
  const [openNewTaskModal, setOpenNewTaskModal] = useState(false);
  const [isEditTask,setIsEditTask]=useState(false)
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
    createdBy: "",
  });
  const [error, setError] = useState<string | null>(null); // Define error state
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

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
    } finally {
      handleCloseNewTaskModal();
    }
  };

  const updateTask = async () => {
    try {
      console.log("new updaated task : ",newTask)
      const response = await fetch(
        `http://localhost:4000/api/project/${id}/task/${newTask.taskID}/manager`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            taskID: newTask.taskID,
            taskName: newTask.taskName,
            taskDescription: newTask.taskDescription,
            taskStatus: newTask.taskStatus,
            taskStartDate: newTask.taskStartDate,
            taskDueDate: newTask.taskDueDate,
            parentProjectId: newTask.parentProjectID,
            assignedUsername: newTask.assignedUsername,
            taskCompletedDate: newTask.taskCompletedDate,
            approved: newTask.approved
          }),
        }
      );
      const data=await response.json()
      if (!response.ok) {
        const errors = [
          data.Errors.description,
          data.Errors.name,
          data.Errors.email,
          data.Errors.phone
        ].filter(Boolean).join('\n'); // Or use '\n' if you want line breaks
        
        setError(errors);
      }
    } catch (err) {
      console.error("failed to udate task" ,err);
      setError("Something went wrong while updating the task");
    } finally {
      setFetchTrigger(true); // re-fetch like a boss
    }
  };

  const toggleApproval = async (taskID: number, approved: boolean) => {
    let task = tasks.find((task) => task.taskID === Number(taskID));
    if (task) {
      console.log("updated task: ", task);
      task.approved = !approved;
    }

    console.log("updated task: ", task);

    try {
      const response = await fetch(
        `http://localhost:4000/api/project/${id}/task/${taskID}/approve`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            taskID: task?.taskID,
            approved: !approved, // toggling approval
            taskName: task?.taskName,
            taskDescription: task?.taskDescription,
            taskStatus: task?.taskStatus,
            taskStartDate: task?.taskStartDate,
            taskDueDate: task?.taskDueDate,
            parentProjectId: task?.parentProjectID,
            assignedUsername: task?.assignedUsername,
            taskCompletedDate: task?.taskCompletedDate,
            taskApprovedDate: task?.taskApprovedDate,
            createdBy: task?.createdBy,
            targetUsername: task?.assignedUsername,
            targetType: "user",
            msg: `task ${task?.taskName} ${
              !approved ? "approved" : "disapproved"
            } by`,
          }),
        }
      );
      if (response.ok) {

      } else {
        setError("Failed to update approval status");
      }
    } catch (err) {
      setError("Failed to update approval status");
    } finally {
      setFetchTrigger(true);
    }
  };

  const handleStatusChange = async (
    taskID: number,
    taskName: string,
    status: string,
    username: string
  ) => {
    console.log("status: ", status);
    try {
      const response = await fetch(
        `http://localhost:4000/api/project/${id}/task/${taskID}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            taskStatus: status,
            targetUsername: username,
            targetType: "user",
            msg: `task ${taskName} status set to ${status} by`,
          }), // Update task status
        }
      );
      if (response.ok) {
        // Update tasks state after status change
      } else {
        setError("Failed to update task status");
      }
    } catch (err) {
      setError("Failed to update task status");
    } finally {
      setFetchTrigger(true);
    }
  };

  const handleOpenTaskNameModal = (taskID: number, taskName: string) => {
    setSelectedTaskID(taskID);
    setCurrentTaskName(taskName);
    setOpenTaskNameModal(true);
  };

  const handleOpenDescriptionModal = (
    taskID: number,
    taskDescription: string
  ) => {
    setSelectedTaskID(taskID);
    setCurrentDescription(taskDescription);
    setOpenDescriptionModal(true);
  };

  const handleCloseModal = () => {
    setOpenTaskNameModal(false);
    setOpenDescriptionModal(false);
    setSelectedTaskID(null);
  };

  const handleTaskNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setCurrentTaskName(event.target.value);
  };

  const handleDescriptionChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setCurrentDescription(event.target.value);
  };

  const handleSaveTaskName = async () => {
    if (selectedTaskID !== null) {
      try {
        const response = await fetch(
          `http://localhost:4000/api/tasks/${selectedTaskID}/name`,
          {
            method: "PATCH",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ taskName: currentTaskName }),
          }
        );
        if (response.ok) {
          // Update tasks state after task name change
          handleCloseModal();
        } else {
          setError("Failed to update task name");
        }
      } catch (err) {
        setError("Failed to update task name");
      }
    }
  };

  const handleSaveDescription = async () => {
    if (selectedTaskID !== null) {
      try {
        const response = await fetch(
          `http://localhost:4000/api/tasks/${selectedTaskID}/description`,
          {
            method: "PATCH",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ taskDescription: currentDescription }),
          }
        );
        if (response.ok) {
          // Update tasks state after description change
          handleCloseModal();
        } else {
          setError("Failed to update task description");
        }
      } catch (err) {
        setError("Failed to update task description");
      }
    }
  };

  const columns: GridColDef[] = [
    { field: "taskID", headerName: "ID", width: 10 },
    {
      field: "taskName",
      headerName: "Name",
      flex: 2,
      renderCell: (params) => (
        <Button
          onClick={() =>
            handleOpenTaskNameModal(params.row.taskID, params.row.taskName)
          }
        >
          {params.value}
        </Button>
      ),
    },
    {
      field: "taskDescription",
      headerName: "Description",
      flex: 2,
      renderCell: (params) => (
        <Button
          onClick={() =>
            handleOpenDescriptionModal(
              params.row.taskID,
              params.row.taskDescription
            )
          }
        >
          {params.value.length > 50
            ? `${params.value.slice(0, 50)}...`
            : params.value}
        </Button>
      ),
    },
    {
      field: "taskStartDate",
      headerName: "Start Date",
      flex: 2,
      valueFormatter: (params) =>
        params ? new Date(params).toLocaleDateString() : "N/A",
    },
    {
      field: "taskDueDate",
      headerName: "Due Date",
      flex: 2,
      valueFormatter: (params) =>
        params ? new Date(params).toLocaleDateString() : "N/A",
    },
    { field: "assignedUsername", headerName: "Contributor", flex: 2 },
    { field: "createdBy", headerName: "Assigned By", flex: 2 },
    {
      field: "taskStatus",
      headerName: "Status",
      flex: 2,
      renderCell: (params) => {
        const statusColor = params.value === "pending" ? "#f44336" : "#4caf50"; // red for pending, green for completed

        return (
          <FormControl
            variant="outlined"
            sx={{
              width: "100%",
              maxHeight: "100%",
              display: "flex",
              justifyContent: "center",
            }}
          >
            <Select
              value={params.value || ""}
              onChange={(e) =>
                handleStatusChange(
                  params.row.taskID,
                  params.row.taskName,
                  e.target.value,
                  params.row.assignedUsername
                )
              }
              label="Status"
              sx={{
                backgroundColor: `${statusColor}22`, // subtle background color
                borderRadius: "16px", // rounded corners
                height: "40px", // Ensure it doesn't exceed its container's height
                marginTop: "4px",
              }}
            >
              <MenuItem value="completed">Completed</MenuItem>
              <MenuItem value="pending">Pending</MenuItem>
            </Select>
          </FormControl>
        );
      },
    },
    {
      field: "approved",
      headerName: "Approved",
      flex: 2,
      renderCell: (params) => (
        <Checkbox
          color="info"
          checked={params.value}
          onChange={() => toggleApproval(params.row.taskID, params.value)}
        />
      ),
    },
    {
      field:"edit",
      headerName:"edit",
      renderCell: (params)=>(
        <EditIcon
          onClick={()=>{
            setIsEditTask(true)
            setOpenNewTaskModal(true);
            setNewTask(params.row); 
          }}
        />
      )
    }
  ];

  const rows = tasks.map((task) => ({
    id: task.taskID,
    taskID: task.taskID,
    taskName: task.taskName,
    taskDescription: task.taskDescription || "No description available",
    taskStartDate: task.taskStartDate,
    taskDueDate: task.taskDueDate,
    parentProjectID: task.parentProjectID,
    assignedUsername: task.assignedUsername,
    taskStatus: task.taskStatus || "Not assigned",
    approved: task.approved,
  }));

  return (
    <Box
      sx={{ height: "100%", width: "99%", border: "5px", borderRadius: "10px" }}
    >
      {error && <PopUp Error={error} Message="" onClose={()=> setError(null)}/>}
      <Box sx={{ padding: 1, display: "flex", justifyContent: "flex-end" }}>
        <Button
          onClick={handleOpenNewTaskModal} // Opens the modal
          variant="contained"
          color="primary"
          sx={{ marginBottom: 2, backgroundColor: "royalblue" }}
        >
          Create New Task
        </Button>
      </Box>
      <DataGrid rows={rows} columns={columns} />
      {/*Create new task modal*/}
      <CreateTaskModal
        open={openNewTaskModal}
        newTask={newTask}
        onClose={handleCloseNewTaskModal}
        handleCreateNewTask={ isEditTask? updateTask : handleCreateNewTask}
        setNewTask={setNewTask}
      />
      {/* Task Name Modal */}
      <Modal open={openTaskNameModal} onClose={handleCloseModal}>
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            padding: 2,
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[400] : "white",
            borderRadius: 2,
            maxWidth: 400,
            margin: "auto",
            marginTop: 10,
          }}
        >
          <Typography variant="h6">Edit Task Name</Typography>
          <TextField
            label="Task Name"
            value={currentTaskName}
            onChange={handleTaskNameChange}
            fullWidth
            variant="outlined"
            margin="normal"
          />
          <Button
            onClick={handleSaveTaskName}
            color="primary"
            variant="contained"
          >
            Save
          </Button>
        </Box>
      </Modal>

      {/* Task Description Modal */}
      <Modal open={openDescriptionModal} onClose={handleCloseModal}>
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            padding: 2,
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[400] : "white",
            borderRadius: 2,
            maxWidth: 400,
            margin: "auto",
            marginTop: 10,
          }}
        >
          <Typography variant="h6">Edit Task Description</Typography>
          <TextField
            label="Task Description"
            value={currentDescription}
            onChange={handleDescriptionChange}
            fullWidth
            variant="outlined"
            margin="normal"
            multiline
            rows={6} // You can adjust this to fit your desired height
          />
          <Button
            onClick={handleSaveDescription}
            color="primary"
            variant="contained"
          >
            Save
          </Button>
        </Box>
      </Modal>
    </Box>
  );
};

export default ProjectTaskDetailPage;
