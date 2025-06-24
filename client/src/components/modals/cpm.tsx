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
} from "@mui/material";
import CheckBoxOutlineBlankIcon from "@mui/icons-material/CheckBoxOutlineBlank";
import CheckBoxIcon from "@mui/icons-material/CheckBox";
import React, { useState } from "react";
import { CpmApiResponse } from "../../hooks/types";
import { tokens } from "../../theme";

interface AddCPMTaskModalProps {
  open: boolean;
  onClose: () => void;
  onAddTask: (newCPMTask: CpmApiResponse) => void;
  CPMAddableTasks: CpmApiResponse[];
  CPMTask: CpmApiResponse[];
}

const names = [
  "Oliver Hansen",
  "Van Henry",
  "April Tucker",
  "Ralph Hubbard",
  "Omar Alexander",
  "Carlos Abbott",
  "Miriam Wagner",
  "Bradley Wilkerson",
  "Virginia Andrews",
  "Kelly Snyder",
];

export const AddCPMTaskModal: React.FC<AddCPMTaskModalProps> = ({
  open,
  onClose,
  onAddTask,
  CPMAddableTasks,
  CPMTask,
}) => {
  const [personName, setPersonName] = React.useState<string[]>([]);
  const [newCPMTask, setNewCPMTask] = useState<CpmApiResponse>(
    {} as CpmApiResponse
  );
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const icon = <CheckBoxOutlineBlankIcon fontSize="small" />;
  const checkedIcon = <CheckBoxIcon fontSize="small" />;
  const [errors, setErrors] = useState({
    duration: "",
  });
  const handleChanger = (event: SelectChangeEvent<typeof personName>) => {
    const {
      target: { value },
    } = event;
    setPersonName(
      // On autofill we get a stringified value.
      typeof value === "string" ? value.split(",") : value
    );
  };

  const handleChange = (field: keyof CpmApiResponse, value: number | null) => {
    setNewCPMTask((prev) => ({ ...prev, [field]: value }));

    // Validation logic
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
    console.log("new pert task", newCPMTask, newCPMTask?.taskId);
    if (!newCPMTask?.taskId) {
      alert("Please select a task.");
      return;
    }
    // Check for validation errors
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
        <TextField
          select
          fullWidth
          label="Select Task"
          value={newCPMTask?.taskId}
          onChange={
            (e) => console.log(CPMAddableTasks[0].taskName)
            // handleChange("task?.taskId", Number(e.target.value))
          }
          sx={{ mb: 2 }}
        >
          {CPMAddableTasks.map((CPMTask) => (
            <MenuItem >
              Task {CPMTask.taskName}
            </MenuItem>
          ))}
        </TextField>
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
          options={names}
          disableCloseOnSelect
          getOptionLabel={(option) => option}
          renderOption={(props, option, { selected }) => {
            const { key, ...optionProps } = props;
            return (
              <li key={key} {...optionProps}>
                <Checkbox
                  icon={icon}
                  checkedIcon={checkedIcon}
                  style={{ marginRight: 8 }}
                  checked={selected}
                />
                {option}
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
        />
      </Box>
    </Modal>
  );
};
