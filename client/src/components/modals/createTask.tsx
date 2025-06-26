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
  handleCreateNewTask: () => void;
  setNewTask: React.Dispatch<React.SetStateAction<tasks>>;
}

const CreateTaskModal: React.FC<CreateTaskProps> = ({
  open,
  newTask,
  onClose,
  handleCreateNewTask,
  setNewTask,
}) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const { id } = useParams();
  const { users, usersLoading, usersError } = useGetConfirmedUsers(Number(id));
  const isValidTaskname=newTask.taskName.trim().length>=3 && newTask.taskName.length <=40
  const isValidTaskDescription=newTask.taskDescription!.trim().length>=10 && newTask.taskDescription!.length <=200
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
          onChange={(e) => setNewTask({ ...newTask, taskName: e.target.value })}
          fullWidth
          error={newTask.taskName!=='' && !isValidTaskname}
          helperText={!isValidTaskname?'must be more than 3 character':''}
          variant="outlined"
          margin="normal"
        />
        <Autocomplete
          disablePortal
          value={newTask.assignedUsername}
          onChange={(e, newVal) => {
            setNewTask({ ...newTask, assignedUsername: newVal?newVal:"" });
          }}
          options={users}
          sx={{ width: '100%' }}
          renderInput={(params) => <TextField {...params} label="Assigned User" />}
        />
        <TextField
          label="Start Date"
          type="date"
          value={newTask.taskStartDate || ""}
          onChange={(e) =>
            setNewTask({ ...newTask, taskStartDate: e.target.value })
          }
          fullWidth
          variant="outlined"
          margin="normal"
          InputLabelProps={{
            shrink: true, // Ensures the label doesn't overlap with the date value
          }}
        />
        <TextField
          label="Due Date"
          type="date"
          value={newTask.taskDueDate || ""}
          onChange={(e) =>
            setNewTask({ ...newTask, taskDueDate: e.target.value })
          }
          fullWidth
          variant="outlined"
          margin="normal"
          InputLabelProps={{
            shrink: true,
          }}
        />

        <TextField
          label="Description"
          value={newTask.taskDescription || ""}
          onChange={(e) =>
            setNewTask({ ...newTask, taskDescription: e.target.value })
          }
          fullWidth
          multiline
          error={newTask.taskDescription!=='' && !isValidTaskDescription}
          helperText={!isValidTaskname?'must be more than 10 character':''}
          variant="outlined"
          margin="normal"
          rows={4}
        />

        <Button
          onClick={handleCreateNewTask}
          color="primary"
          variant="contained"
          sx={{ mt: 2 }}
        >
          Save Task
        </Button>
      </Box>
    </Modal>
  );
};

export default CreateTaskModal
