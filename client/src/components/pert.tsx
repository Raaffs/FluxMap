import React, { useEffect, useRef, useState } from "react";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { Chart, ChartData, ChartOptions, ChartDataset } from "chart.js";
import {
  Chart as ChartJS,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
  Tooltip,
  Legend,
} from "chart.js";
import { TextField, Button, Box, SelectChangeEvent } from "@mui/material";
import { erf } from "mathjs";
import { ApiResponse, PertData, tasks } from "../hooks/types";
import AddPertTaskModal from "./modals/pert";
import { useParams } from "react-router-dom";
import { useAddPert } from "../hooks/pert";
ChartJS.register(
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
  Tooltip,
  Legend
);

interface PertRows {
  id: number; // Update this to number
  parentTaskID: number; // Also update this field to match
  predecessorTaskId?: number | null;
  optimistic: number;
  pessimistic: number;
  mostLikely: number;
  taskName?: string;
}

export const PertTable: React.FC<{
  pertTasks: PertData[] | undefined;
  tasks: tasks[];
}> = ({ pertTasks, tasks }) => {
  const {id}=useParams()
  console.log("USER PARAMS ID ", Number(id))
  const [open, setOpen] = useState(false);
  const [predecessorPertTasks, setpredecessorPertTasks] = useState<PertRows[]>([]);
  const [pertData, setPertData] = useState<PertData[]>(pertTasks || []);

  const {addPert,pertloading,perterror,pertsuccess}=useAddPert(String(id))

  const handleAddTask = (newTask: PertData) => {
    setPertData([...pertData, newTask]);
    console.log("nex perx tsx ",newTask)
    addPert(newTask)
  };


  console.log("pertTasks", pertTasks);

  if (pertTasks === undefined || pertTasks === null) {
    return <div>no pert data</div>;
  }

  const tasksMap = new Map<number, boolean>();
  const pertTaskMap = new Map<number, boolean>();
  tasks.forEach((task) => tasksMap.set(task.taskID, true));
  pertTasks.forEach((pertTask) => pertTaskMap.set(pertTask.parentTaskID, true));

  let PertAddableTask: PertRows[] = [];
  for (const task of tasks) {
    if (!pertTaskMap.has(task.taskID)) {
      PertAddableTask.push({
        id:0,
        parentTaskID: task.taskID,
        predecessorTaskId: 0,
        optimistic: 0,
        pessimistic: 0,
        mostLikely: 0,
        taskName: task.taskName,
      });
    }
  }

  // const PertAddableTask: PertData[] = (pertTasks || []).filter((task) => !tasksMap.has(task.parentTaskID))
  let rows: PertRows[] = [];
  //for some reason even though PertData has field ParentTaskID
  //here we've to use parentTaskId instead for tasks to render correctly
  //I've no clue why. It might be because of how backend is send data but not gonna
  //mess with it for now

  //20/30/25: solved, it was indeed the problem with wrong json format in backend
  //I don't understand why ts doesn't throw an error when it gets wrong json
  for (const pertTask of pertTasks) {
    rows.push({
      id: pertTask.parentTaskID,
      parentTaskID: pertTask.parentTaskID,
      predecessorTaskId: pertTask.predecessorTaskId,
      optimistic: pertTask.optimistic,
      pessimistic: pertTask.pessimistic,
      mostLikely: pertTask.mostLikely,
      taskName:
        tasks.find((task) => task.taskID === pertTask.parentTaskID)?.taskName ||
        "",
    });
  }

  let addAbleRows: PertRows[] = [];
  for (const pertTask of PertAddableTask) {
    addAbleRows.push({
      id: pertTask.parentTaskID,
      parentTaskID: pertTask.parentTaskID,
      predecessorTaskId: pertTask.predecessorTaskId,
      optimistic: pertTask.optimistic,
      pessimistic: pertTask.pessimistic,
      mostLikely: pertTask.mostLikely,
      taskName:
        tasks.find((task) => task.taskID === pertTask.parentTaskID)?.taskName ||
        "",
    });
  }

  console.log("rows: ", rows);
  const columns: GridColDef[] = [
    { field: "taskName", headerName: "Name", width: 150 },
    { field: "predecessorTaskId", headerName: "Predecessor", width: 150 },
    { field: "optimistic", headerName: "Optimistic", width: 150 },
    { field: "pessimistic", headerName: "Pessimistic", width: 150 },
    { field: "mostLikely", headerName: "Most Likely", width: 150 },
    { field: "edit", headerName: "Edit", width: 150 },
  ];
  return (
    <Box>
      <Button variant="outlined" onClick={() => setOpen(true)}>
        Add New Task
      </Button>
      <AddPertTaskModal
        open={open}
        onClose={() => setOpen(false)}
        onAddTask={handleAddTask}
        pertAddableTasks={PertAddableTask}
        pertTasks={rows}
      />
      <DataGrid rows={rows} columns={columns} />
      {/* Modal */}
    </Box>
  );
};

const PertNormalDistributionChart: React.FC<{
  apiResponse: ApiResponse | null;
  pertTasks: PertData[] | undefined;
  tasks: tasks[];
}> = ({ apiResponse, pertTasks, tasks }) => {
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);
  const [zValue, setZValue] = useState<number | null>(null);
  const [probability, setProbability] = useState<number | null>(null);
  const [xInput, setXInput] = useState<string>("");
  const selectTasks = tasks?.map((task) => task.taskID);
  const standardNormalCDF = (z: number): number => {
    return 0.5 * (1 + erf(z / Math.sqrt(2)));
  };

  useEffect(() => {
    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
    }

    if (!apiResponse) {
      return; // Skip rendering the chart if there is no data
    }

    const labels = Array.from({ length: 101 }, (_, i) => (i / 10).toFixed(2)); // x-axis range: 0 to +10

    const datasets: ChartDataset<"line">[] = apiResponse.result.taskResults.map(
      (task) => {
        const { mean, stddev, taskId } = task;

        const normalDistribution = (x: number): number => {
          const exponent = -((x - mean) ** 2) / (2 * stddev ** 2);
          return (1 / (stddev * Math.sqrt(2 * Math.PI))) * Math.exp(exponent);
        };

        const data = labels.map((x) =>
          parseFloat(normalDistribution(parseFloat(x)).toFixed(2))
        );

        return {
          label: `Task ${taskId}`,
          data,
          borderColor: `hsl(${Math.random() * 360}, 70%, 50%)`,
          fill: false,
        };
      }
    );

    const chartData: ChartData<"line"> = {
      labels,
      datasets,
    };

    const chartOptions: ChartOptions<"line"> = {
      responsive: true,
      plugins: {
        legend: {
          display: true,
        },
        tooltip: {
          mode: "index",
          intersect: false,
        },
      },
      scales: {
        x: {
          title: {
            display: true,
            text: "Standard Deviation from Mean",
          },
        },
        y: {
          title: {
            display: true,
            text: "Probability Density",
          },
        },
      },
    };

    if (chartRef.current) {
      chartInstanceRef.current = new Chart(chartRef.current, {
        type: "line",
        data: chartData,
        options: chartOptions,
      });
    }

    return () => {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.destroy();
      }
    };
  }, [apiResponse]);

  const calculateZValue = () => {
    if (!apiResponse || !xInput) return;

    const criticalTasks = apiResponse.result.taskResults.filter((task) =>
      apiResponse.result.criticalPath.includes(task.taskId)
    );

    const mean = criticalTasks.reduce((acc, task) => acc + task.mean, 0);
    const stddev = criticalTasks.reduce((acc, task) => acc + task.stddev, 0);

    const x = parseFloat(xInput);
    if (!isNaN(x)) {
      const z = (x - mean) / stddev;
      const probability = standardNormalCDF(z); // Calculate probability using CDF
      setZValue(z);
      setProbability(probability);
    }
  };

  if (!apiResponse) {
    return <p>No data available</p>;
  }

  const criticalTasks = apiResponse.result.taskResults.filter((task) =>
    apiResponse.result.criticalPath.includes(task.taskId)
  );
  const mean = criticalTasks.reduce((acc, task) => acc + task.mean, 0);
  const stddev = criticalTasks.reduce((acc, task) => acc + task.stddev, 0);

  return (
    <div>
      <div
        style={{
          marginBottom: "20px",
          padding: "10px",
          border: "1px solid #ccc",
          borderRadius: "5px",
        }}
      >
        <h3>Critical Path Information</h3>
        <p>
          <strong>Critical Path:</strong>{" "}
          {apiResponse.result.criticalPath.map((taskId, index) => (
            <span key={taskId}>
              {taskId}
              {index < apiResponse.result.criticalPath.length - 1 ? " → " : ""}
            </span>
          ))}
        </p>
        <p>
          <strong>Mean:</strong> {mean.toFixed(2)}
        </p>
        <p>
          <strong>Standard Deviation:</strong> {stddev.toFixed(2)}
        </p>
        <div style={{ marginTop: "10px" }}>
          <TextField
            label="Input X"
            type="number"
            value={xInput}
            onChange={(e) => setXInput(e.target.value)}
            variant="outlined"
            size="small"
            style={{ marginRight: "10px" }}
          />
          <Button variant="contained" color="primary" onClick={calculateZValue}>
            Calculate Z
          </Button>
        </div>
        {zValue !== null && (
          <p>
            <strong>Z-Value:</strong> {zValue.toFixed(2)}
          </p>
        )}
        {probability !== null && (
          <p>
            <strong>Probability (P(Z ≤ z)):</strong> {probability.toFixed(4)}
          </p>
        )}
      </div>
      <canvas ref={chartRef}></canvas>
    </div>
  );
};

export default PertNormalDistributionChart;
