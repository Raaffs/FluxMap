import {
  Box,
  Button,
  Card,
  Divider,
  LinearProgress,
  Paper,
  Stack,
  Typography,
  useTheme,
} from "@mui/material";
import { tokens } from "../../theme";
import { useConfirmInvitation, useFetchInvitations } from "../../hooks/invite";
import { DisplayInvitations } from "../../hooks/types";
import { useState } from "react";
import { NoUpdates } from "../../components/cards/noUpdates";
import NotificationsOffOutlinedIcon from '@mui/icons-material/NotificationsOffOutlined';
export const Invitation = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [trigger,setTrigger]=useState<boolean>(false)
  const [invitations, loading, error] = useFetchInvitations();
  const {confirmInvitation,inviteloading,inviteerror,invitesuccess}=useConfirmInvitation()

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        height="90vh"
        width="80vw"
      >
        <Box sx={{ width: "50%" }}>
          <LinearProgress color="info" />
        </Box>
      </Box>
    );
  }

  if (invitations === null || invitations === undefined) {
    return <NoUpdates
        title="You're All Caught Up"
        description="There are currently no updates requiring your attention."
    />;
  }

  return (
    <Box
      sx={{
        margin: "5px",
        maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
        padding: "16px",
        borderRadius: "8px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white",
        boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        textAlign: "left",
      }}
    >
      {invitations.map((invite: DisplayInvitations) => (
        <Card
          key={invite.invitationID}
          sx={{
            marginBottom: "20px",
            padding: "20px",
            borderRadius: "8px",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[400] : "#fafafa",
            boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
            "&:hover": { boxShadow: "0 4px 8px rgba(0,0,0,0.8)" },
            textAlign: "left",
          }}
        >
          <Typography variant="h5" gutterBottom>
            You've been invited to{" "}
            <Typography
              component="span"
              fontWeight="bold"
              color={colors.blueAccent[500]}
              variant="h5"
            >
              {invite.projectName}
            </Typography>{" "}
            by{" "}
            <Typography
              component="span"
              fontWeight="bold"
              color={colors.blueAccent[500]}
              variant="h5"
            >
              {invite.ownerName}
            </Typography>
          </Typography>

          <Typography
            variant="body1"
            color="text.primary"
            sx={{ marginTop: "10px", lineHeight: 1.6 }}
          >
            {invite.projectDescription}
          </Typography>

          <Typography
            variant="h6"
            sx={{
              marginTop: "16px",
            }}
          >
            Role offered:{" "}
            <Typography
              component="span"
              fontWeight="bold"
              color={
                theme.palette.mode === "dark"
                  ? colors.blueAccent[400]
                  : colors.blueAccent[400]
              }
            >
              {invite.role}
            </Typography>
          </Typography>

          <Box
            sx={{
              display: "flex",
              justifyContent: "flex-start",
              gap: 1,
              marginTop: "18px",
            }}
          >
            <Button
              variant="contained"
              sx={{ backgroundColor: colors.greenAccent[500] }}
              onClick={() =>{
                console.log('invt id',invite)
                 confirmInvitation((invite.invitationID),"accepted")
                if(error){
                  console.log(error)
                }
              }
            }
              
            >
              Accept
            </Button>
            <Button
              variant="contained"
              sx={{ backgroundColor: colors.redAccent[500] }}
              onClick={() =>{
                  confirmInvitation((invite.invitationID),"accepted")
                  if(error){
                    console.log(error)
                  }
                }
              }
            >
              Decline
            </Button>
          </Box>
        </Card>
      ))}
    </Box>
  );
};
