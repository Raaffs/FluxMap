import React from "react";
import { Doughnut } from "react-chartjs-2";
import { Box, Typography } from "@mui/material";
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  Title,
  ChartOptions
} from "chart.js";
import { ChartData } from "chart.js/auto";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

interface DonutProps {
  data: ChartData<"doughnut", number[], string>;
  title?: string;
  maintainRatio?: boolean;
}

export const TaskStatusDonutChart: React.FC<DonutProps> = ({
  data,
  title = "",
  maintainRatio = true
}) => {
  const options: ChartOptions<"doughnut"> = {
    responsive: true,
    maintainAspectRatio: maintainRatio,
    plugins: {
      legend: {
        position: "bottom"
      },
    },
    cutout: "60%" 
  };

  return (
    <Box sx={{ width: "100%", height: "100%" }}>
      <Typography variant="h6" gutterBottom>
        {title}
      </Typography>
      <Box sx={{ width: "100%", height: "calc(100% - 32px)" }}>
        <Doughnut data={data} options={options} />
      </Box>
    </Box>
  );
};
