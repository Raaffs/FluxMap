import React, { useState } from "react";
import {
  Modal,
  Box,
  Typography,
  Button,
  TextField,
  MenuItem,
} from "@mui/material";
import { PertRows } from "../pert";

interface AddPertTaskModalProps {
  open: boolean;
  onClose: () => void;
  onAddTask: (task: PertRows) => void;
  pertAddableTasks: PertRows[];
  pertTasks: PertRows[];
  newTask:PertRows;
  setNewTask:React.Dispatch<React.SetStateAction<PertRows>>
}

const AddPertTaskModal: React.FC<AddPertTaskModalProps> = ({
  open,
  onClose,
  onAddTask,
  pertAddableTasks,
  pertTasks,
  newTask,
  setNewTask
}) => {
  
  const [errors, setErrors] = useState({
    optimistic: "",
    mostLikely: "",
    pessimistic: "",
  });

  const handleChange = (field: keyof PertRows, value: number | null) => {
    setNewTask((prev) => ({ ...prev, [field]: value }));

    // Validation logic
    if (
      field === "optimistic" ||
      field === "mostLikely" ||
      field === "pessimistic"
    ) {
      if (value === null || isNaN(value) || value < 0) {
        setErrors((prev) => ({ ...prev, [field]: "Must be a positive number" }));
      } else {
        setErrors((prev) => ({ ...prev, [field]: "" }));
      }
    }
  };

  const handleSubmit = () => {
    console.log('new pert task',newTask,newTask.parentTaskID)
    if (!newTask.parentTaskID) {
      alert("Please select a task.");
      return;
    }

    // Check for validation errors
    if (Object.values(errors).some((error) => error)) {
      alert("Please fix input errors before submitting.");
      return;
    }
    console.log("edited pert task:",newTask)
    onAddTask(newTask);
    onClose();
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Box
        sx={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          width: 400,
          bgcolor: "background.paper",
          boxShadow: 24,
          p: 4,
          borderRadius: 2,
        }}
      >
        <Typography variant="h6" gutterBottom>
          Add New PERT Task
        </Typography>

        <TextField
          select
          fullWidth
          label="Select Task"
          value={newTask.parentTaskID}
          onChange={(e) =>
            handleChange("parentTaskID", Number(e.target.value))
          }
          sx={{ mb: 2 }}
        >
          {pertAddableTasks.map((task) => (
            <MenuItem key={task.parentTaskID} value={task.parentTaskID}>
              Task {task.taskName}
            </MenuItem>
          ))}
        </TextField>

        <TextField
          select
          fullWidth
          label="Predecessor Task"
          value={newTask.predecessorTaskId ?? ""}
          onChange={(e) =>
            handleChange(
              "predecessorTaskId",
              e.target.value === "" ? null : Number(e.target.value)
            )
          }
          sx={{ mb: 2 }}
        >
          <MenuItem value="">None</MenuItem>
          {pertTasks.map((task) => (
            <MenuItem key={task.parentTaskID} value={task.parentTaskID}>
               {task.taskName}
            </MenuItem>
          ))}
        </TextField>

        {/* Optimistic Time */}
        <TextField
          fullWidth
          label="Optimistic Time"
          type="number"
          value={newTask.optimistic}
          onChange={(e) => handleChange("optimistic", Number(e.target.value))}
          error={!!errors.optimistic}
          helperText={errors.optimistic}
          sx={{ mb: 2 }}
        />

        {/* Most Likely Time */}
        <TextField
          fullWidth
          label="Most Likely Time"
          type="number"
          value={newTask.mostLikely}
          onChange={(e) => handleChange("mostLikely", Number(e.target.value))}
          error={!!errors.mostLikely}
          helperText={errors.mostLikely}
          sx={{ mb: 2 }}
        />

        {/* Pessimistic Time */}
        <TextField
          fullWidth
          label="Pessimistic Time"
          type="number"
          value={newTask.pessimistic}
          onChange={(e) => handleChange("pessimistic", Number(e.target.value))}
          error={!!errors.pessimistic}
          helperText={errors.pessimistic}
          sx={{ mb: 2 }}
        />

        {/* Actions */}
        <Box sx={{ display: "flex", justifyContent: "space-between", mt: 2 }}>
          <Button variant="outlined" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="contained" onClick={handleSubmit} disabled={!newTask.parentTaskID}>
            Add Task
          </Button>
        </Box>
      </Box>
    </Modal>
  );
};

export default AddPertTaskModal;
