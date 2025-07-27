import React, { useEffect, useRef } from "react";
import Chart from "chart.js/auto";

interface IndividualTaskPert {
  mean: number;
  predecessorId: number | null;
  stddev: number;
  taskId: number;
  variance: number;
}

interface PertResult {
  criticalPath: number[];
  taskPert: IndividualTaskPert[];
}

interface NormalDistributionChartProps {
  pertResult: PertResult;
}

const NormalDistributionChart: React.FC<NormalDistributionChartProps> = ({ pertResult }) => {
  const chartRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    // Helper to calculate normal distribution
    const normalDistribution = (x: number, mean: number, stddev: number): number => {
      const exponent = -Math.pow(x - mean, 2) / (2 * Math.pow(stddev, 2));
      return (1 / (stddev * Math.sqrt(2 * Math.PI))) * Math.exp(exponent);
    };

    // Prepare data for each task
    const datasets = pertResult.taskPert.map((task) => {
      const { mean, stddev, taskId } = task;
      const xValues = Array.from({ length: 100 }, (_, i) => mean - 3 * stddev + (i * 6 * stddev) / 100);
      const yValues = xValues.map((x) => normalDistribution(x, mean, stddev));

      return {
        label: `Task ${taskId}`,
        data: xValues.map((x, index) => ({ x, y: yValues[index] })),
        borderColor: `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.7)`,
        backgroundColor: "rgba(226, 11, 11, 0)",
        borderWidth: 2,
        tension: 0.3,
        fill:true
      };
    });

    // Create the chart
    const chart = new Chart(chartRef.current, {
      type: "line",
      data: {
        datasets,
      },
      options: {
        responsive: true,
        scales: {
          x: {
            type: "linear",
            title: {
              display: true,
              text: "Time (days)",
            },
          },
          y: {
            title: {
              display: true,
              text: "Probability Density",
            },
          },
        },
        plugins: {
          legend: {
            position: "top",
          },
        },
      },
    });

    // Cleanup
    return () => chart.destroy();
  }, [pertResult]);

  return <canvas ref={chartRef} />;
};

export default NormalDistributionChart;
