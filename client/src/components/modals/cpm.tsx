import {
  Box,
  Modal,
  TextField,
  Typography,
  MenuItem,
  useTheme,
  Checkbox,
  SelectChangeEvent,
  Autocomplete,
  FormControlLabel,
  Button,
  FormControl,
  Select,
  InputLabel,
} from "@mui/material";
import CheckBoxOutlineBlankIcon from "@mui/icons-material/CheckBoxOutlineBlank";
import CheckBoxIcon from "@mui/icons-material/CheckBox";
import React, { useEffect, useState } from "react";
import { CpmApiResponse } from "../../hooks/types";
import { tokens } from "../../theme";

interface AddCPMTaskModalProps {
  open: boolean;
  onClose: () => void;
  onAddTask: (newCPMTask: CpmApiResponse) => void;
  CPMAddableTasks: CpmApiResponse[];
  CPMTasks: CpmApiResponse[];
  newCPMTask: CpmApiResponse;
  setNewCPMTask: React.Dispatch<React.SetStateAction<CpmApiResponse>>;
}
export const AddCPMTaskModal: React.FC<AddCPMTaskModalProps> = ({
  open,
  onClose,
  onAddTask,
  CPMAddableTasks,
  CPMTasks,
  newCPMTask,
  setNewCPMTask,
}) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const icon = <CheckBoxOutlineBlankIcon fontSize="small" />;
  const checkedIcon = <CheckBoxIcon fontSize="small" />;
  const [CPMAddableTasksState, setCPMAddableTasksState] = useState<CpmApiResponse[]>(CPMTasks);
  const [errors, setErrors] = useState({
    duration: "",
  });
  useEffect(() => {
  if (
    newCPMTask.taskId &&
    !CPMAddableTasksState.some((task) => task.taskId === newCPMTask.taskId)
  ) {
    setCPMAddableTasksState((prev) => [...prev, newCPMTask]);
  }
}, [newCPMTask]);

    console.log("CPMAddableTasksState", CPMAddableTasksState, "new task ",newCPMTask);    
  
  const handleChange = (field: keyof CpmApiResponse, value: number | null) => {
    setNewCPMTask((prev) => ({ ...prev, [field]: value }));
    if (field === "duration") {
      if (
        value === null ||
        isNaN(value) ||
        value < 0 ||
        typeof value !== "number"
      ) {
        setErrors((prev) => ({
          ...prev,
          [field]: "Must be a positive number",
        }));
      } else {
        setErrors((prev) => ({ ...prev, [field]: "" }));
      }
    }
  };

  const handleSubmit = () => {
    console.log("new cpm task", newCPMTask, newCPMTask?.taskId);
    if (!newCPMTask?.taskId) {
      alert("Please select a task.");
      return;
    }
    if (Object.values(errors).some((error) => error)) {
      alert("Please fix input errors before submitting.");
      return;
    }

    onAddTask(newCPMTask);
    onClose();
  };
  return (
    <Modal open={open} onClose={onClose}>
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
          minHeight: "50%",
          margin: "auto",
          marginTop: 10,
        }}
      >
        <Typography variant="h6" gutterBottom>
          Add New CPM Task
        </Typography>
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel id="task-select-label">Select Task</InputLabel>
          <Select
            labelId="task-select-label"
            value={newCPMTask.taskId ?? ""}
            label="Select Task"
            onChange={(e) => handleChange("taskId", Number(e.target.value))}
          >
            {CPMAddableTasksState.map((CPMTask) => (
              <MenuItem key={CPMTask.taskId} value={CPMTask.taskId}>
                {CPMTask.taskName}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <TextField
          fullWidth
          label="Duration"
          type="number"
          value={newCPMTask.duration}
          onChange={(e) => handleChange("duration", Number(e.target.value))}
          error={!!errors.duration}
          helperText={errors.duration}
          sx={{ mb: 2 }}
        />
        <Autocomplete
          multiple
          id="checkboxes-tags-demo"
          options={CPMTasks}
          disableCloseOnSelect
          getOptionLabel={(option: CpmApiResponse) => option.taskName}
          isOptionEqualToValue={(
            option: CpmApiResponse,
            value: CpmApiResponse
          ) => option.taskId === value.taskId}
          renderOption={(props, option, { selected }) => {
            return (
              <li {...props}>
                <Checkbox
                  icon={icon}
                  checkedIcon={checkedIcon}
                  style={{ marginRight: 8 }}
                  checked={selected}
                />
                {option.taskName}
              </li>
            );
          }}
          sx={{
            alignItems: "center",
            mb: 2,
            width: 370,
          }}
          renderInput={(params) => (
            <TextField {...params} label="Checkboxes" placeholder="Favorites" />
          )}
          onChange={(
            event: React.SyntheticEvent<Element, Event>,
            newValue: CpmApiResponse[]
          ) => {
            const selectedIds = newValue.map((task) => task.taskId);
            setNewCPMTask((prev) => ({ ...prev, dependencies: selectedIds }));          }}
        />
        <FormControlLabel
          control={
            <Checkbox
              onChange={() => {
                setNewCPMTask((prev) => ({
                  ...prev,
                  isCriticalPath: !newCPMTask.isCriticalPath,
                }));
              }}
            />
          }
          label="is critical?"
          sx={{
            marginLeft: 0,
            paddingLeft: 0,
            justifyContent: "flex-start", // push everything to the left
            width: "100%", // full width so alignment works
            "& .MuiCheckbox-root": {
              marginLeft: 0, // remove default margin on checkbox
            },
          }}
        />
        <Button onClick={handleSubmit}>Add CPM task</Button>
      </Box>
    </Modal>
  );
};
