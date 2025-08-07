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
  Filler,
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
  Legend,
  Filler // <--- register this
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
  title: string;
  xLabel: string;
  yLabel: string;
  maintainRatio: boolean;
}

export const TaskStatusGraph: React.FC<GraphProps> = ({
  data,
  title,
  xLabel,
  yLabel,
  maintainRatio
}) => {
  const options: ChartOptions<"line"> = {
    responsive: true,
    maintainAspectRatio:maintainRatio,
    plugins: {
      legend: {
        position: "top",
      },
      title: {
        display: true,
        text: title,
      },
    },
    scales: {
      x: {
        title: {
          display: true,
          text: xLabel,
        },
      },
      y: {
        title: {
          display: true,
          text: yLabel,
        },
        beginAtZero: true,
      },
    },
  };

  return (
    <Box sx={{ width: "100%", height: "100%" }}>
      <Typography variant="h6" gutterBottom>
        Tasks Approved and Completed
      </Typography>
      <Box sx={{ width: "100%", height: "calc(100% - 32px)" }}>
        <Line data={data} options={options} />
      </Box>
    </Box>
  );
};

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
  const [fetch, fetchTrigger] = useState<boolean>(false);
  const [tasks, loading, error] = useFetchTaskData(id, fetch, fetchTrigger);
  console.log("tasks: ", tasks);
  const safeTasks = tasks || [];
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
        fill: true,
      },
      {
        label: "Tasks Completed",
        data: weeklyCompletedTasks.map((w) => w.count),
        borderColor: "rgba(153, 102, 255, 1)",
        backgroundColor: "rgba(153, 102, 255, 0.2)",
        tension: 0.3,
        fill: true,
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
      <TaskStatusGraph
        data={taskStatusData}
        title="Tasks Approved and Completed Over Time"
        xLabel="Date"
        yLabel="Number of Tasks"
        maintainRatio={true}
      />
    </Box>
  );
};

export default Graphs;
