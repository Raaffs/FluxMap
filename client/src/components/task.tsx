import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Box, Typography, MenuItem, Select, FormControl, InputLabel, Checkbox, Modal, TextField, Button, useTheme } from "@mui/material";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { tokens } from "../theme";
export const ProjectTaskDetailPage = () => {
  const { id } = useParams(); // Get the project id from the URL
  const [tasks, setTasks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [openTaskNameModal, setOpenTaskNameModal] = useState(false);
  const [openDescriptionModal, setOpenDescriptionModal] = useState(false);
  const [currentTaskName, setCurrentTaskName] = useState("");
  const [currentDescription, setCurrentDescription] = useState("");
  const [selectedTaskID, setSelectedTaskID] = useState<number | null>(null);
  const theme=useTheme()
  const colors=tokens(theme.palette.mode)
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(`http://localhost:4000/api/project/${id}/tasks`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
        });
        const data = await response.json();
        setTasks(data);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch tasks");
        setLoading(false);
      }
    };

    fetchTasks();
  }, [id]);

  const toggleApproval = async (taskID: number, approved: boolean) => {
    try {
      const response = await fetch(`http://localhost:4000/api/tasks/${taskID}/approve`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ approved: !approved }), // Toggle approval status
      });
      if (response.ok) {
        setTasks(prevTasks => 
          prevTasks.map(task => 
            task.taskID === taskID ? { ...task, approved: !approved } : task
          )
        );
      } else {
        setError("Failed to update approval status");
      }
    } catch (err) {
      setError("Failed to update approval status");
    }
  };

  const handleStatusChange = async (taskID: number, status: string) => {
    try {
      const response = await fetch(`http://localhost:4000/api/tasks/${taskID}/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ taskStatus: status }), // Update task status
      });
      if (response.ok) {
        setTasks(prevTasks => 
          prevTasks.map(task => 
            task.taskID === taskID ? { ...task, taskStatus: status } : task
          )
        );
      } else {
        setError("Failed to update task status");
      }
    } catch (err) {
      setError("Failed to update task status");
    }
  };

  const handleOpenTaskNameModal = (taskID: number, taskName: string) => {
    setSelectedTaskID(taskID);
    setCurrentTaskName(taskName);
    setOpenTaskNameModal(true);
  };

  const handleOpenDescriptionModal = (taskID: number, taskDescription: string) => {
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

  const handleDescriptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setCurrentDescription(event.target.value);
  };

  const handleSaveTaskName = async () => {
    if (selectedTaskID !== null) {
      try {
        const response = await fetch(`http://localhost:4000/api/tasks/${selectedTaskID}/name`, {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ taskName: currentTaskName }),
        });
        if (response.ok) {
          setTasks(prevTasks => 
            prevTasks.map(task => 
              task.taskID === selectedTaskID ? { ...task, taskName: currentTaskName } : task
            )
          );
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
        const response = await fetch(`http://localhost:4000/api/tasks/${selectedTaskID}/description`, {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ taskDescription: currentDescription }),
        });
        if (response.ok) {
          setTasks(prevTasks => 
            prevTasks.map(task => 
              task.taskID === selectedTaskID ? { ...task, taskDescription: currentDescription } : task
            )
          );
          handleCloseModal();
        } else {
          setError("Failed to update task description");
        }
      } catch (err) {
        setError("Failed to update task description");
      }
    }
  };

  if (loading) {
    return <Typography>Loading tasks...</Typography>;
  }

  if (error) {
    return <Typography>{error}</Typography>;
  }

  const columns: GridColDef[] = [
    { field: 'taskID', headerName: 'Task ID', width: 150 },
    {
      field: 'taskName',
      headerName: 'Name',
      width: 200,
      renderCell: (params) => (
        <Button onClick={() => handleOpenTaskNameModal(params.row.taskID, params.row.taskName)}>{params.value}</Button>
      ),
    },
    {
      field: 'taskDescription',
      headerName: 'Description',
      width: 300,
      renderCell: (params) => (
        <Button onClick={() => handleOpenDescriptionModal(params.row.taskID, params.row.taskDescription)}>
          {params.value.length > 50 ? `${params.value.slice(0, 50)}...` : params.value}
        </Button>
      ),
    },
    { field: 'taskStartDate', headerName: 'Start Date', width: 180, valueFormatter: (params) => params ? new Date(params).toLocaleDateString() : 'N/A' },
    { field: 'taskDueDate', headerName: 'Due Date', width: 180, valueFormatter: (params) => params ? new Date(params).toLocaleDateString() : 'N/A' },
    { field: 'assignedUsername', headerName: 'Contributor', width: 200 },
    {
      field: 'taskStatus',
      headerName: 'Status',
      width: 150,
      renderCell: (params) => (
        <FormControl variant="standard" sx={{ width: '100%' }}>
          <InputLabel>Status</InputLabel>
          <Select
            value={params.value || ""}
            onChange={(e) => handleStatusChange(params.row.taskID, e.target.value)}
            label="Status"
          >
            <MenuItem value="completed">Completed</MenuItem>
            <MenuItem value="pending">Pending</MenuItem>
          </Select>
        </FormControl>
      ),
    },
    {
      field: 'approved',
      headerName: 'Approved',
      width: 150,
      renderCell: (params) => (
        <Checkbox
          checked={params.value}
          onChange={() => toggleApproval(params.row.taskID, params.value)}
        />
      ),
    },
  ];

  const rows = tasks.map(task => ({
    id: task.taskID,
    taskID: task.taskID,
    taskName: task.taskName,
    taskDescription: task.taskDescription || 'No description available',
    taskStartDate: task.taskStartDate,
    taskDueDate: task.taskDueDate,
    parentProjectID: task.parentProjectID,
    assignedUsername: task.assignedUsername,
    taskStatus: task.taskStatus || 'Not assigned',
    approved: task.approved,
  }));

  return (
    <Box sx={{ height:"100%", width: '90%', border:'5px', borderRadius:'10px' }}>
      <Typography variant="h4" gutterBottom>
        Tasks for Project {id}
      </Typography>
      <DataGrid
        rows={rows}
        columns={columns}
      />

      {/* Task Name Modal */}
      <Modal open={openTaskNameModal} onClose={handleCloseModal}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', padding: 2, backgroundColor: theme.palette.mode==='dark'?colors.primary[400]:'grey', borderRadius: 2, maxWidth: 400, margin: 'auto', marginTop: 10 }}>
          <Typography variant="h6">Edit Task Name</Typography>
          <TextField
            label="Task Name"
            value={currentTaskName}
            onChange={handleTaskNameChange}
            fullWidth
            variant="outlined"
            margin="normal"
          />
          <Button onClick={handleSaveTaskName} color="primary" variant="contained">Save</Button>
        </Box>
      </Modal>

      {/* Task Description Modal */}
      <Modal open={openDescriptionModal} onClose={handleCloseModal}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', padding: 2, backgroundColor: 'white', borderRadius: 2, maxWidth: 400, margin: 'auto', marginTop: 10 }}>
          <Typography variant="h6">Edit Task Description</Typography>
          <TextField
            label="Description"
            value={currentDescription}
            onChange={handleDescriptionChange}
            fullWidth
            variant="outlined"
            margin="normal"
            multiline
            rows={4}
          />
          <Button onClick={handleSaveDescription} color="primary" variant="contained">Save</Button>
        </Box>
      </Modal>
    </Box>
  );
};
