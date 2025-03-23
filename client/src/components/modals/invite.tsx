import { useState } from "react";
import {
  Modal,
  Box,
  TextField,
  Select,
  MenuItem,
  Button,
  InputLabel,
  FormControl,
  Alert,
  CircularProgress,
} from "@mui/material";
import { useInvite } from "../../hooks/invite";
import { useParams } from "react-router-dom";

interface InviteModalProps {
  open: boolean;
  onClose: () => void;
}

const InviteModal: React.FC<InviteModalProps> = ({ open, onClose }) => {
  const [username, setUsername] = useState("");
  const [role, setRole] = useState("user");
  const { id } = useParams();
  const {inviteUser, loading, error, success }= useInvite(id!);

  const handleInvite = async () => {
    console.log("username and role: ",username,role)
    if (username.trim()) {
      await inviteUser(username, role);
      if (success) {
        setUsername("");
        setRole("user");
        onClose();
      }else{
        
      }
    }
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Box
        sx={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          width: 300,
          bgcolor: "background.paper",
          boxShadow: 24,
          p: 3,
          borderRadius: 2,
          display: "flex",
          flexDirection: "column",
          gap: 2,
        }}
      >
        {error && <Alert severity="error">{error}</Alert>}
        {success && <Alert severity="success">Invite sent successfully!</Alert>}

        <TextField
          label="Username"
          variant="outlined"
          fullWidth
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          sx={{ marginBottom: 2 }}
        />
        <FormControl fullWidth>
          <InputLabel>Role</InputLabel>
          <Select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            sx={{ marginBottom: 2 }}
          >
            <MenuItem value="manager">Manager</MenuItem>
            <MenuItem value="user">User</MenuItem>
          </Select>
        </FormControl>
        <Button
          variant="contained"
          color="primary"
          onClick={handleInvite}
        >
          {loading ? <CircularProgress size={24} /> : "Send Invite"}
        </Button>
      </Box>
    </Modal>
  );
};

export default InviteModal;
