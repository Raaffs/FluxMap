import React from 'react';
import { Modal, Box, Typography, TextField, Button } from '@mui/material';
import { Projects } from '../../hooks/types';

interface ModalProps {
  open: boolean;
  onClose: () => void;
  newProject: Projects
  handleInputChange: (event: React.ChangeEvent<HTMLInputElement>)=>void
  handleSubmit: ()=>void
}

const CreateProjectModal: React.FC<ModalProps> = ({ open, onClose, newProject, handleInputChange, handleSubmit }) => {
    const isValidProjectName=newProject.projectName.trim().length>4
    const isValidProjectDescription=newProject.projectDescription!.trim().length>10 
  
  return (
    <Modal open={open} onClose={onClose}>
                  <Box
            sx={{
              width: "400px",
              padding: "16px",
              backgroundColor: "white",
              borderRadius: "8px",
              boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
              margin: "auto",
              marginTop: "10%",
            }}
          >
            <Typography variant="h6" sx={{ marginBottom: "16px" }}>
              Create New Project
            </Typography>
            <TextField
              fullWidth
              error={newProject.projectName!==''&& !isValidProjectName}
              value={newProject.projectName || ""}
              name="projectName"
              label="Project Name"
              onChange={handleInputChange}
              helperText={!isValidProjectName?'must be more than 4 character':''}
              sx={{ marginBottom: "16px" }}
            />
            <TextField
              fullWidth
              error={newProject.projectDescription!==''&& !isValidProjectDescription}
              helperText={!isValidProjectDescription?'must be more than 10 character':''}
              label="Description"
              name="projectDescription"
              value={newProject.projectDescription || ""}
              onChange={handleInputChange}
              sx={{ marginBottom: "16px" }}
            />
            <TextField
              fullWidth
              label="ETA"
              name="projectDueDate"
              type="date"
              variant="outlined"
              value={newProject.projectDueDate || ""}
              onChange={handleInputChange}
              InputLabelProps={{ shrink: true }}
              sx={{ marginBottom: "16px" }}
            />
            <Button
              fullWidth
              variant="contained"
              color="primary"
              onClick={handleSubmit}
            //   disabled={postLoading}
            >
              Create
            </Button>
          </Box>

    </Modal>
  );
};

export default CreateProjectModal;
