import React, { useState } from "react";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { ChartOptions } from "chart.js";
import { Box, Typography } from "@mui/material";
import useFetchTaskData from "../../hooks/task";
import { useParams } from "react-router-dom";
import { tasks } from "../../hooks/types";
import { groupTasksWithMap, groupTasksByPeriod } from "../../helpers/group";
// Register Chart.js components
ChartJS.register(
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
  Title,
  Tooltip,
  Legend
);

// Define Props Type
interface GraphProps {
  data: {
    labels: string[];
    datasets: {
      label: string;
      data: any[];
      borderColor: string;
      backgroundColor: string;
      tension: number;
    }[];
  };
}

// TaskStatusGraph Component
const TaskStatusGraph: React.FC<GraphProps> = ({ data }) => {
  const options: ChartOptions<"line"> = {
    responsive: true,
    plugins: {
      legend: {
        position: "top", // Correctly typed position
      },
      title: {
        display: true,
        text: "Tasks Approved and Completed Over Time",
      },
    },
    scales: {
      x: {
        title: {
          display: true,
          text: "Date",
        },
      },
      y: {
        title: {
          display: true,
          text: "Number of Tasks",
        },
        beginAtZero: true,
      },
    },
  };

  return (
    <Box mb={4}>
      <Typography variant="h6" gutterBottom>
        Tasks Approved and Completed
      </Typography>
      <Line data={data} options={options} />
    </Box>
  );
};

// ContributorGraph Component
export const ContributorGraph: React.FC<GraphProps> = ({ data }) => {
  const options: ChartOptions<"line"> = {
    responsive: true,
    plugins: {
      legend: {
        position: "top",
      },
      title: {
        display: true,
        text: "Tasks Completed by Contributors Over Time",
      },
    },
    scales: {
      x: {
        title: {
          display: true,
          text: "Date",
        },
      },
      y: {
        title: {
          display: true,
          text: "Tasks Completed",
        },
        beginAtZero: true,
      },
    },
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Contributor Tasks
      </Typography>
      <Line data={data} options={options} />
    </Box>
  );
};

// Main App Component with Data
const Graphs: React.FC = () => {
  const { id } = useParams();
  const [fetch,fetchTrigger]=useState<boolean>(false)
  const [tasks, loading, error] = useFetchTaskData(id, fetch,fetchTrigger);
  console.log("tasks: ", tasks);
  const safeTasks = tasks || [];
  const completedDateData = safeTasks
    .map((task) => task.taskCompletedDate) // Extract the dates
    .filter((date): date is string => typeof date === "string" && date !== "0") // Narrow to valid strings
    .sort((a, b) => new Date(a).getTime() - new Date(b).getTime()) // Sort by date
    .map((date) => new Date(date).toISOString().split("T")[0]); // Format to 'YYYY-MM-DD'

  const approvedDateData = safeTasks
    .map((task) => task.taskApprovedDate) // Extract approval dates
    .filter((date): date is string => typeof date === "string" && date !== "0") // Ensure only valid strings
    .sort((a, b) => new Date(a).getTime() - new Date(b).getTime()) // Sort by timestamp
    .map((date) => new Date(date).toISOString().split("T")[0]); // Format to 'YYYY-MM-DD'

  const weeklyApprovedTasks = groupTasksByPeriod(
    safeTasks,
    "taskApprovedDate",
    "week"
  );
  const weeklyCompletedTasks = groupTasksByPeriod(
    safeTasks,
    "taskCompletedDate",
    "week"
  );

  const contributerTasks = groupTasksWithMap(
    safeTasks,
    "taskCompletedDate",
    "taskCompletedDate"
  );
  // Prepare graph data
  const taskStatusData = {
    labels: Array.from(
      new Set(
        [...weeklyApprovedTasks, ...weeklyCompletedTasks].map((w) => w.period)
      )
    ),
    datasets: [
      {
        label: "Tasks Approved",
        data: weeklyApprovedTasks.map((w) => w.count),
        borderColor: "rgba(75, 192, 192, 1)",
        backgroundColor: "rgba(75, 192, 192, 0.2)",
        tension: 0.3,
      },
      {
        label: "Tasks Completed",
        data: weeklyCompletedTasks.map((w) => w.count),
        borderColor: "rgba(153, 102, 255, 1)",
        backgroundColor: "rgba(153, 102, 255, 0.2)",
        tension: 0.3,
      },
    ],
  };
  // Sample data for ContributorGraph
  const contributorData = {
    labels: completedDateData,
    datasets: [
      {
        label: "Contributor A",
        data: [3, 5, 4, 6, 8, 10],
        borderColor: "rgba(255, 99, 132, 1)",
        backgroundColor: "rgba(255, 99, 132, 0.2)",
        tension: 0.3,
      },
      {
        label: "Contributor B",
        data: [4, 3, 5, 7, 9, 11],
        borderColor: "rgba(54, 162, 235, 1)",
        backgroundColor: "rgba(54, 162, 235, 0.2)",
        tension: 0.3,
      },
      {
        label: "Contributor C",
        data: [2, 4, 3, 5, 7, 8],
        borderColor: "rgba(255, 206, 86, 1)",
        backgroundColor: "rgba(255, 206, 86, 0.2)",
        tension: 0.3,
      },
    ],
  };

  return (
    <Box
      p={4}
      maxHeight="50%"
      width="99%"
      overflow="auto" // Ensures scrolling is enabled for overflowing content
    >
      <TaskStatusGraph data={taskStatusData} />
    </Box>
  );
};

export default Graphs;
