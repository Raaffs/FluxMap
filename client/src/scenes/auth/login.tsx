import { useState } from "react";
import {
  Box,
  TextField,
  Typography,
  Avatar,
  CardActions,
  useTheme,
} from "@mui/material";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import { useNavigate } from "react-router-dom";
import Button from "@mui/material/Button";
import Link from "@mui/material/Link";
import Card from "@mui/material/Card";
import Grid2 from "@mui/material/Grid2";
import { useAuth } from "../../context/authContext";
import { normalizeAccessMap } from "../../helpers/filter";
import { UserRole } from "../../hooks/types";
import React from "react";
import { tokens } from "../../theme";
interface userAuth {
  username: string;
  password: string;
}
const LoginUser: React.FC<{
  startWebSocket: () => void;
  setUserProjectRoleMap: React.Dispatch<
    React.SetStateAction<Record<number, UserRole>>
  >;
}> = ({ startWebSocket, setUserProjectRoleMap }) => {
  const btnStyle = { margin: "50px 0", width: "200px" };
  const navigate = useNavigate();
  const [input, setInput] = useState<userAuth>({
    username: "",
    password: "",
  });

  const theme = useTheme();
  const colors = tokens(theme.palette.mode);

  const [error, setError] = useState<string | null>(null); // State for error message
  const { login } = useAuth();

  const handleClick = () => {
    fetch("http://localhost:4000/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(input),
      credentials: "include",
    })
      .then(async (response) => {
        if (!response.ok) {
          console.log(response);
          const errorData = await response.json();
          setError(errorData.error);
          throw new Error(errorData.error);
        }
        const data = await response.json(); // Parse the response data
        console.log("Login successful:", data.roles);
        const roles = data.roles as Record<UserRole, number[] | null>;
        const normal = normalizeAccessMap(roles);
        setUserProjectRoleMap(normal);
        login();
        startWebSocket();
        navigate("/"); // Redirect to the home page
      })
      .catch((err) => {
        console.error("Error fetching projects:", err);
        setError(err.error);
      });
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setInput({
      ...input,
      [event.target.name]: event.target.value,
    });
  };

  return (
    <Box
      display="flex"
      padding={2}
      justifyContent="center"
      alignItems="center"
      sx={{
        backgroundColor: "transparent",
        width: "100%",
        height: "86vh",
        borderRadius: "20px",
        margin: "5px",
      }}
    >
      <Card
        sx={{
          width: "50%",
          minHeight: "60%",
          borderRadius: "16px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[400] : "white",
          boxShadow:
            theme.palette.mode === "dark"
              ? "0px 2px 10px rgba(0, 0, 0, 0.4)"
              : "0px 4px 20px rgba(0, 0, 0, 0.4)",
          transition: "all 0.3s ease-in-out",
        }}
      >
        <Grid2>
          <Avatar sx={{ m: 1, bgcolor: "secondary.main", top: "0px" }}>
            <LockOutlinedIcon />
          </Avatar>
          <h1>LogIn</h1>
        </Grid2>
        <CardActions sx={{ backgroundColor: "transparent" }}>
          <Box display="flex" flexDirection="column" width="100%">
            {error && (
              <Typography
                color="error"
                align="center"
                style={{ marginBottom: "10px" }}
              >
                {error}
              </Typography>
            )}
            <TextField
              label="Username"
              placeholder="Enter username"
              variant="outlined"
              name="username"
              value={input.username}
              fullWidth
              required
              style={{ padding: "10px" }}
              onChange={handleChange}
            />
            <TextField
              label="Password"
              placeholder="Enter password"
              name="password"
              value={input.password}
              type="password"
              onChange={handleChange}
              variant="outlined"
              fullWidth
              required
              style={{ padding: "10px" }}
            />
            <Box display="flex" justifyContent="center" alignItems="center">
              <Button
                type="submit"
                size="medium"
                variant="contained"
                style={btnStyle}
                onClick={handleClick}
              >
                Log In
              </Button>
            </Box>
            <Box>
              <Typography variant="h5"> Don't have an account?</Typography>
              <Link href="Register">Link</Link>
            </Box>
          </Box>
        </CardActions>
      </Card>
      <Card
        sx={{
          width: "100%",
          minHeight: "60%",
          borderRadius: "16px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[400] : "white",
          boxShadow:
            theme.palette.mode === "dark"
              ? "0px 2px 10px rgba(0, 0, 0, 0.4)"
              : "0px 4px 20px rgba(0, 0, 0, 0.4)",
          transition: "all 0.3s ease-in-out",

          marginLeft: "20px",
        }}
        // elevation={2}
      ></Card>
    </Box>
  );
};

export default LoginUser;
