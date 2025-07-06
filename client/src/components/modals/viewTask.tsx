import { tasks } from "../../hooks/types";
import {
  Box,
  Modal,
  Typography,
  TextField,
  Button,
  useTheme,
  Autocomplete,
} from "@mui/material";
import { tokens } from "../../theme";
import { useParams } from "react-router-dom";
import { useGetConfirmedUsers } from "../../hooks/invite";
interface CreateTaskProps {
  open: boolean;
  newTask: tasks;
  onClose: () => void;
}

const ViewTaskModal: React.FC<CreateTaskProps> = ({
  open,
  newTask,
  onClose,
}) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const { id } = useParams();
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
          minHeight:'50%',
          margin: "auto",
          marginTop: 10,
        }}
      >
        <Typography variant="h6">Create New Task</Typography>

        <TextField
          label="Task Name"
          value={newTask.taskName || ""}
          fullWidth
          variant="outlined"
          margin="normal"
          aria-readonly
        />
        <TextField
        label="assigned user"
        value={newTask.assignedUsername||""}
        fullWidth
       margin="normal"
        aria-readonly
        />
        <TextField
          label="Start Date"
          value={newTask.taskStartDate || "no start date"}
          fullWidth
          variant="outlined"
          margin="normal"
          InputLabelProps={{
            shrink: true, // Ensures the label doesn't overlap with the date value
          }}
          aria-readonly
        />
        <TextField
          label="Due Date"
          value={newTask.taskDueDate || "no due date"}
          fullWidth
          variant="outlined"
          margin="normal"
          InputLabelProps={{
            shrink: true,
          }}
          aria-readonly
        />

        <TextField
          label="Description"
          value={newTask.taskDescription || ""}
          fullWidth
          multiline
          variant="outlined"
          margin="normal"
          rows={4}
          aria-readonly
        />

        <Button
          color="primary"
          variant="contained"
          sx={{ mt: 2 }}
          onClick={onClose}
        >
          Close
        </Button>
      </Box>
    </Modal>
  );
};

export default ViewTaskModal
