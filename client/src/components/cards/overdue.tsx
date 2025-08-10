import React from "react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  CardHeader,
  CardContent,
  Typography,
  Box,
  useTheme,
  Divider,
} from "@mui/material";
import { CheckCircleOutline } from "@mui/icons-material";
import {
  Timeline,
  TimelineItem,
  TimelineSeparator,
  TimelineConnector,
  TimelineContent,
  TimelineDot,
} from "@mui/lab";

interface Task {
  taskID: number;
  taskName: string;
  taskDueDate?: string | null;
  parentProjectID: number;
}

interface OverdueTasksProps {
  tasks: Task[];
}

const OverdueTasks: React.FC<OverdueTasksProps> = ({ tasks }) => {
  const navigate = useNavigate();
  const theme = useTheme();
  console.log("Overdue Tasks:", tasks);
  const formatDate = (dateString?: string | null) => {
    if (!dateString) return "";
    return new Date(dateString).toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  return (
    <Card
      sx={{
        height: "100%",
        width: "100%",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <CardHeader
        title={
          <Box
            display="flex"
            alignItems="center"
            justifyContent="space-between"
            px={2}
            pb={1}
          >
            <Typography variant="h6" fontWeight="bold" color="text.primary">
              Overdue Tasks
            </Typography>
            <Box
              sx={{
                width: 16,
                height: 16,
                bgcolor: "error.main",
                borderRadius: "50%",
              }}
              aria-label="Overdue indicator"
            />
          </Box>
        }
        sx={{ pb: 0 }}
      />
      <Divider />

      <CardContent sx={{ pt: 0, px: 2, flexGrow: 1 }}>
        {tasks.length === 0 ? (
                    <Box
            flex={1}
            display="flex"
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            py={5}
            px={2}
            color="text.secondary"
          >
            <CheckCircleOutline sx={{ fontSize: 48, mb: 1, color: "success.main" }} />
            <Typography variant="body2" fontStyle="italic">
              No overdue tasks — nice work!
            </Typography>
          </Box>
        ) : (
          <Timeline sx={{ p: 0, mt: 1 }} position="alternate">
            {tasks.map((task, index) => (
              <TimelineItem
                key={task.taskID}
                sx={{
                  cursor: "pointer",
                  "&:hover": {
                    backgroundColor: theme.palette.action.hover,
                    borderRadius: 1,
                  },
                  px: 1,
                  mb: 1.5,
                }}
                onClick={() => navigate(`/project/${task.parentProjectID}`)}
              >
                <TimelineSeparator>
                  <TimelineDot
                    color="error"
                    sx={{
                      width: 14,
                      height: 14,
                      boxShadow: `0 0 8px ${theme.palette.error.main}`,
                    }}
                  />
                  {index !== tasks.length - 1 && (
                    <TimelineConnector sx={{ bgcolor: "error.main" }} />
                  )}
                </TimelineSeparator>
                <TimelineContent sx={{ pb: 1 }}>
                  <Typography
                    variant="subtitle1"
                    fontWeight="600"
                    sx={{
                      "&:hover": { color: "error.main" },
                      transition: "color 0.3s ease",
                      userSelect: "none",
                    }}
                  >
                    {task.taskName}
                  </Typography>
                  <Typography
                    variant="body2"
                    color="error.light"
                    fontWeight="500"
                    mt={0.3}
                  >
                    {formatDate(task.taskDueDate)}
                  </Typography>
                </TimelineContent>
              </TimelineItem>
            ))}
          </Timeline>
        )}
      </CardContent>
    </Card>
  );
};

export default OverdueTasks;
