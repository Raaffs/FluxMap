import { Box, Button, Typography } from "@mui/material";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { LinearProgress } from "@mui/material";
import { Bar } from "react-chartjs-2";
import InviteModal from "./modals/invite";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { tasks } from "../hooks/types";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { useMemo } from "react";
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

// Function to render the Task Table 
export const ContributerTaskTable = ({ tasksData }: { tasksData: tasks[] }) => {
  const [isInviteModal, setInviteModal] = useState(false);
  const contributerMap: {
    [username: string]: {
      completed: number;
      approved: number;
      pending: number;
      progress: number;
    };
  } = {};

  for (const task of tasksData) {
    const { assignedUsername, taskStatus, approved } = task;

    if (!contributerMap[assignedUsername]) {
      contributerMap[assignedUsername] = {
        completed: 0,
        approved: 0,
        pending: 0,
        progress: 0,
      };
    }

    const userStats = contributerMap[assignedUsername];

    if (taskStatus === "completed") {
      userStats.completed++;
    } else {
      userStats.pending++;
    }

    if (approved) {
      userStats.approved++;
    }

    const totalTasks = userStats.completed + userStats.pending;
    userStats.progress =
      totalTasks > 0 ? (userStats.completed / totalTasks) * 100 : 0;
  }

  const rows = Object.entries(contributerMap).map(
    ([username, stats], index) => ({
      id: index + 1,
      assignedUsername: username,
      taskCompleted: stats.completed,
      taskPending: stats.pending,
      approved: stats.approved,
      progressPercent: stats.progress,
    })
  );

  const columns: GridColDef[] = [
    { field: "assignedUsername", headerName: "Contributor", flex: 2 },
    { field: "taskCompleted", headerName: "Task Completed", flex: 2 },
    { field: "taskPending", headerName: "Pending", flex: 2 },
    { field: "approved", headerName: "Approved", flex: 2 },
    {
      field: "progressPercent",
      headerName: "Progress (%)",
      renderCell: (params) => (
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            width: "100%",
          }}
        >
          <LinearProgress
            variant="determinate"
            value={params.value}
            color="success"
            sx={{ width: "100%", marginTop: "10px" }}
          />
          <Typography variant="body2" align="center" sx={{ marginTop: "4px" }}>
            {params.value.toFixed(2)}%
          </Typography>
        </Box>
      ),
    },
  ];

  return (
    <Box sx={{ height: "400px", width: "100%" }}>
      <Box sx={{ padding: 1, display: "flex", justifyContent: "flex-end" }}>
        <Button
          variant="contained"
          color="primary"
          sx={{ backgroundColor: "royalblue", marginBottom: 2 }}
          onClick={() => {
            setInviteModal(true);
          }}
        >
          Add Contributer
        </Button>
        <InviteModal
          open={isInviteModal}
          onClose={() => setInviteModal(false)}
        />
      </Box>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={2}
      >
        <DataGrid rows={rows} columns={columns} />
      </Box>
    </Box>
  );
};

// Function to render the Task Progress Chart
export const ContributerTaskChart = ({ tasksData }: { tasksData: tasks[] }) => {
  const contributerMap = useMemo(() => {
    const map: {
      [username: string]: {
        completed: number;
        approved: number;
        pending: number;
      };
    } = {};

    for (const task of tasksData) {
      const { assignedUsername, taskStatus, approved } = task;

      if (!map[assignedUsername]) {
        map[assignedUsername] = {
          completed: 0,
          approved: 0,
          pending: 0,
        };
      }

      if (taskStatus === "completed") {
        map[assignedUsername].completed++;
      } else {
        map[assignedUsername].pending++;
      }

      if (approved) {
        map[assignedUsername].approved++;
      }
    }

    return map;
  }, [tasksData]);

  const chartData = useMemo(
    () => ({
      labels: Object.keys(contributerMap),
      datasets: [
        {
          label: "Completed Tasks",
          data: Object.values(contributerMap).map((stats) => stats.completed),
          backgroundColor: "#4caf50",
        },
        {
          label: "Pending Tasks",
          data: Object.values(contributerMap).map((stats) => stats.pending),
          backgroundColor: "#fbc02d",
        },
        {
          label: "Approved Tasks",
          data: Object.values(contributerMap).map((stats) => stats.approved),
          backgroundColor: "#1976d2",
        },
      ],
    }),
    [contributerMap]
  );

  const chartOptions = {
    responsive: true,
    scales: {
      x: {
        type: "category" as const,
        categoryPercentage: 1,
        barPercentage: 0.9,
      },
    },
    plugins: {
      legend: {
        position: "top" as const,
      },
    },
  };


  return (
    <Box p={4} maxHeight="50%" width="70%">
      <Typography variant="h6" gutterBottom>
        Task Progress Chart
      </Typography>
      <Bar data={chartData} options={chartOptions} />
    </Box>
  );
};
