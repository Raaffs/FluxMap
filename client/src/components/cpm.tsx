import React, { useState } from "react";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { Box, Button } from "@mui/material";
import { CpmApiResponse, CpmResult, tasks } from "../hooks/types";
import { AddCPMTaskModal } from "./modals/cpm";
import EditIcon from "@mui/icons-material/Edit";
import { usePostCpmData } from "../hooks/cpm";
import { useParams } from "react-router-dom";
const CpmTable: React.FC<{
  data: CpmResult | null;
  tasks: tasks[];
}> = ({ data, tasks }) => {
  const {id}=useParams()
  console.log("id of project",id)
  const [open, setOpen] = useState<boolean>(false);
  const [CPMData, setCPMData] = useState<CpmApiResponse[]>(data?.Result || []);
  const {addCPM,cpmloading,cpmerror,cpmsuccess}=usePostCpmData(Number(id))
  
  const handleAddCPMTask = (newCPMTask: CpmApiResponse) => {
    setCPMData([...CPMData, newCPMTask]);
    console.log("new cpm task in component: ",newCPMTask)
    addCPM(newCPMTask)
  };
  const columns: GridColDef[] = [
    { field: "taskId", headerName: "Task ID", flex: 1 },
    { field: "dependenciesName", headerName: "Dependencies", flex: 1 },
    { field: "duration", headerName: "Duration", flex: 1 },
    { field: "earliestStart", headerName: "Earliest Start", flex: 1 },
    { field: "earliestFinish", headerName: "Earliest Finish", flex: 1 },
    { field: "latestStart", headerName: "Latest Start", flex: 1 },
    { field: "latestFinish", headerName: "Latest Finish", flex: 1 },
    { field: "totalFloat", headerName: "Total Float", flex: 1 },
    { field: "freeFloat", headerName: "Free Float", flex: 1 },
    { field: "independentFloat", headerName: "Independent Float", flex: 1 },
    {
      field:"action",
      flex:1,
      headerName:"Actions",
      renderCell:(params)=>{
        return(
          <Box>
            <EditIcon/>
          </Box>
        )
      }
    }
  ];
  let CPMAddableTasks: CpmApiResponse[] = [];
  CPMAddableTasks = getCPMAddableTask(tasks, data?.Result || []);
  let formattedCPMData = setTaskNames(CPMData,tasks);
  getDependenciesName(formattedCPMData)
  console.log("formatted cpm",formattedCPMData)
  // Prepare rows for the DataGrid
  const rows = formattedCPMData.map((task) => ({
    id: task.taskId,
    taskId: task.taskId,
    dependenciesName: task.dependenciesName, 
    duration: task.duration,
    earliestStart: task.earliestStart,
    earliestFinish: task.earliestFinish,
    latestStart: task.latestStart,
    latestFinish: task.latestFinish,
    totalFloat: task.totalFloat,
    freeFloat: task.freeFloat,
    independentFloat: task.independentFloat,
  }));

  return (
    <Box sx={{ height: 400, width: "100%" }}>
      <Box sx={{ padding: 1, display: "flex", justifyContent: "flex-end" }}>
        <Button
          onClick={() => setOpen(true)} // Opens the modal
          variant="contained"
          color="primary"
          sx={{ marginBottom: 2, backgroundColor: "royalblue" }}
        >
          Add New Task
        </Button>
      </Box>
      <AddCPMTaskModal
        open={open}
        onClose={() => setOpen(false)}
        onAddTask={handleAddCPMTask}
        CPMAddableTasks={CPMAddableTasks || []}
        CPMTasks={formattedCPMData || []}
      />

      <DataGrid rows={rows} columns={columns} />
    </Box>
  );
};

//cause the backend/database only store foreign key reference, we need to manually
//join the task names on frontend for UX
function getDependenciesName(CPMTasks: CpmApiResponse[]) {
  const CPMTaskMap = new Map<number, string>();
  for (const CPMTask of CPMTasks) {
    CPMTaskMap.set(CPMTask.taskId, CPMTask.taskName);
  }
  for (let i = 0; i < CPMTasks.length; i++) {
    console.log("cpmtasks1",CPMTasks[i])
    if (
      CPMTasks[i].dependencies === undefined ||
      CPMTasks[i].dependencies === null ||
      CPMTasks[i].dependencies.length === 0
    ) {
      continue;
    }
    if (CPMTasks[i].dependenciesName === undefined) {
      CPMTasks[i].dependenciesName = [];
    }
    for (let id of CPMTasks[i].dependencies) {
      CPMTasks[i].dependenciesName.push(CPMTaskMap.get(id) || "");
    }
  }
}

function getCPMAddableTask(
  task: tasks[],
  CPMTasks: CpmApiResponse[]
): CpmApiResponse[] {
  const CPMTaskMap = new Map<number, boolean>();
  let CPMAddableTasks: CpmApiResponse[] = [];
  if (task === null || task === undefined) {
    return CPMAddableTasks;
  }
  if (CPMTasks === null || CPMTasks.length === 0) {
    for (const t of task) {
      CPMAddableTasks.push({
        taskId: t.taskID,
        taskName: t.taskName,
        parentProjectID: 0,
        dependencies: [],
        dependenciesName: [],
        duration: 0,
        earliestStart: 0,
        earliestFinish: 0,
        latestStart: 0,
        latestFinish: 0,
        totalFloat: 0,
        freeFloat: 0,
        independentFloat: 0,
        isCriticalPath: false,
      });
    }
  }
  CPMTasks.forEach((t) => CPMTaskMap.set(t.taskId, true));
  for (const t of task) {
    if (!CPMTaskMap.has(t.taskID)) {
      CPMAddableTasks.push({
        taskId: t.taskID,
        taskName: t.taskName,
        parentProjectID: 0,
        dependencies: [],
        dependenciesName: [],
        duration: 0,
        earliestStart: 0,
        earliestFinish: 0,
        latestStart: 0,
        latestFinish: 0,
        totalFloat: 0,
        freeFloat: 0,
        independentFloat: 0,
        isCriticalPath: false,
      });
    }
  }
  return CPMAddableTasks;
}

function setTaskNames(CPMTasks: CpmApiResponse[], task: tasks[]):CpmApiResponse[] {
  const CPMTaskMap = new Map<number, CpmApiResponse>();
  let formattedCPMData: CpmApiResponse[] = [];
  CPMTasks.forEach((ct) => CPMTaskMap.set(ct.taskId, ct));
  for (const t of task) {
    if (CPMTaskMap.has(t.taskID)) {
      formattedCPMData.push({
        taskId: t.taskID,
        taskName: t.taskName,
        parentProjectID: CPMTaskMap.get(t.taskID)?.parentProjectID||0,
        dependencies: CPMTaskMap.get(t.taskID)?.dependencies||[],
        dependenciesName: [],
        duration: CPMTaskMap.get(t.taskID)?.duration||0,
        earliestStart: CPMTaskMap.get(t.taskID)?.earliestStart||0,
        earliestFinish: CPMTaskMap.get(t.taskID)?.earliestFinish||0,
        latestStart: CPMTaskMap.get(t.taskID)?.latestStart||0,
        latestFinish: CPMTaskMap.get(t.taskID)?.latestFinish||0,
        totalFloat: CPMTaskMap.get(t.taskID)?.totalFloat||0,
        freeFloat: CPMTaskMap.get(t.taskID)?.freeFloat||0,
        independentFloat: CPMTaskMap.get(t.taskID)?.independentFloat||0,
        isCriticalPath: CPMTaskMap.get(t.taskID)?.isCriticalPath||false,
      });
    }
  }
  return formattedCPMData
}

export default CpmTable;
