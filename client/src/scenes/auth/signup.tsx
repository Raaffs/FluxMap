import { useState } from "react";
import {  Box, TextField, Typography, Avatar, CardActions } from "@mui/material";
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { useNavigate } from "react-router-dom";
import Button from "@mui/material/Button";
import Link from "@mui/material/Link";
import Card from "@mui/material/Card";
import Grid2 from "@mui/material/Grid2";
// Type definition for the input state
interface InputState {
  username: string;
  email   : string;
  password: string;
}


const SignUpUser: React.FC = () => {
  const btnStyle = { margin: '50px 0', width: '200px' };
  const navigate = useNavigate();

  // State for input fields
  const [input, setInput] = useState<InputState>({
    username: "",
    password: "",
    email: ""
  });

  const [error, setError] = useState<string | null>(null); // State for error message

  const handleClick = () => {
    fetch('http://api:4000/api/register',{
        method:'POST',
        headers:{
          "Content-Type": "application/json",
        },
        body: JSON.stringify(input),
        credentials:'include'
      })
      .then(async (response) => {
        if (!response.ok) {
          const errorData = await response.json();
          console.log("Error : ",errorData.Message||errorData.Errors)
          throw new Error(errorData.Message||errorData.Errors.email||errorData.Errors.password);
        }
        const data = await response.json(); // Parse the response data
        console.log("Login successful:", data);
        navigate('/'); // Redirect to the home page
      })  
      .catch(err=>{
        console.error("Error fetching projects:", err);
        setError(err.message)
      })
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setInput({
      ...input,
      [event.target.name]: event.target.value
    });
  };

  return (
    <Box
      display="flex"
      padding={2}
      justifyContent="center"
      alignItems="center"
      sx={{
        backgroundColor: 'transparent',
        width: "100%",
        height: "86vh",
        borderRadius: "20px",
        margin: '5px'
      }}
    >
      <Card sx={{
        width: "50%",
        minHeight: '60%',
        borderRadius: "16px",
        backgroundColor: 'transparent',
      }} elevation={20}>
        <Grid2 >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main', top: '0px' }}>
            <LockOutlinedIcon />
          </Avatar>
          <h1>SignUp</h1>
        </Grid2>
        <CardActions sx={{ backgroundColor: 'transparent' }}>
          <Box display="flex" flexDirection="column" width="100%">
            {error && (
              <Typography color="error" align="center" style={{ marginBottom: '10px' }}>
                {error}
              </Typography>
            )}
            <TextField
              label='Username'
              placeholder='Enter username'
              variant="outlined"
              name="username"
              value={input.username}
              fullWidth
              required
              style={{ padding: "10px" }}
              onChange={handleChange}
            />
            <TextField
              label='Email'
              placeholder='Enter email'
              variant="outlined"
              name="email"
              value={input.email}
              fullWidth
              required
              style={{ padding: "10px" }}
              onChange={handleChange}
            />

            <TextField
              label='Password'
              placeholder='Enter password'
              name="password"
              value={input.password}
              type='password'
              onChange={handleChange}
              variant="outlined"
              fullWidth
              required
              style={{ padding: "10px" }}
            />
            <Box display="flex" justifyContent="center" alignItems="center">
              <Button
                type='submit'
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
              <Link href="/login">Log in</Link>
            </Box>
          </Box>
        </CardActions>
      </Card>
      <Card 
        sx={{
        width: "100%",
        minHeight: "63%",
        borderRadius: "16px",
        backgroundColor: 'transparent',
        marginLeft: "20px",
      }} 
        elevation={20}
      >
      </Card>
    </Box>
  );
}

export default SignUpUser;
