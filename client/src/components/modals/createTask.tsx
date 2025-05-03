import React, { useState } from "react";
import { tasks } from "../../hooks/types";
import { Box, Modal,Typography, TextField, Button, useTheme } from "@mui/material";
import { tokens } from "../../theme";
interface CreateTaskProps {
    open: boolean
    newTask: tasks
    onClose: () => void
    handleCreateNewTask: () => void
    setNewTask: React.Dispatch<React.SetStateAction<tasks>>
}

const CreateTaskModal: React.FC<CreateTaskProps> = ({
        open,
        newTask,
        onClose,
        handleCreateNewTask,
        setNewTask
    })=>{
        const theme=useTheme()
        const colors=tokens(theme.palette.mode)
        return(
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
                margin: "auto",
                marginTop: 10,
              }}
            >
              <Typography variant="h6">Create New Task</Typography>
    
              <TextField
                label="Task Name"
                value={newTask.taskName}
                onChange={(e) =>
                  setNewTask({ ...newTask, taskName: e.target.value })
                }
                fullWidth
                variant="outlined"
                margin="normal"
              />
    
              <TextField
                label="Assigned User"
                value={newTask.assignedUsername}
                onChange={(e) =>{
                    setNewTask({ ...newTask, assignedUsername: e.target.value })
                  }
                }
                fullWidth
                variant="outlined"
                margin="normal"
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
                variant="outlined"
                margin="normal"
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
        )
}

export default CreateTaskModal