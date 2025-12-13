import React from "react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  CardHeader,
  CardContent,
  Typography,
  Box,
  Divider,
} from "@mui/material";
import { CheckCircleOutline } from "@mui/icons-material";
import {
  VerticalTimeline,
  VerticalTimelineElement,
} from "react-vertical-timeline-component";
import "react-vertical-timeline-component/style.min.css";

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

  const formatDate = (dateString?: string | null) => {
    if (!dateString) return "";
    return new Date(dateString).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <Card sx={{ height: "100%", width: "100%", display: "flex", flexDirection: "column" }}>
      <CardHeader
        title={
          <Box display="flex" alignItems="center" justifyContent="space-between" px={2} pb={1}>
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

      <CardContent sx={{ flexGrow: 1, overflowY: "auto", pt: 1 }}>
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
          <VerticalTimeline layout="2-columns" lineColor="#f44336">
            {tasks.map((task, index) => (
              <Box
                key={task.taskID}
                sx={{
                  width: "100%",
                  cursor: "pointer",
                  "&:hover": { backgroundColor: "rgba(244, 67, 54, 0.08)" },
                  transition: "background-color 0.2s ease",
                }}
                onClick={() => navigate(`/project/${task.parentProjectID}`)}
              >
                <VerticalTimelineElement
                  position={index % 2 === 0 ? "left" : "right"}
                  contentStyle={{
                    background: "transparent",
                    color: "#000",
                    boxShadow: "none",
                    padding: "0.8rem 1rem",
                    marginBottom: "0.75rem",
                    borderRadius: "6px",
                  }}
                  contentArrowStyle={{ display: "none" }}
                  iconStyle={{
                    background: "#f44336",
                    boxShadow: "0 0 6px #f44336",
                    width: "14px",
                    height: "14px",
                    marginLeft: "-7px",
                  }}
                  style={{ padding: "0", minHeight: "48px" }}
                >
                  <Typography
                    variant="subtitle1"
                    fontWeight="600"
                    sx={{ userSelect: "none" }}
                  >
                    {task.taskName}
                  </Typography>
                  <Typography
                    variant="subtitle2"
                    color="error.main"
                    fontWeight="500"
                  >
                    {formatDate(task.taskDueDate)}
                  </Typography>
                </VerticalTimelineElement>
              </Box>
            ))}
          </VerticalTimeline>
        )}
      </CardContent>
    </Card>
  );
};

export default OverdueTasks;
