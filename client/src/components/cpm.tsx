import React from 'react';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { Box } from '@mui/material';
import { CpmResult } from '../hooks/types';


const CpmTable: React.FC<{ data: CpmResult |null}> = ({ data }) => {
  console.log('data cpm: ',data)
  // if(!data?.Result){
  //   return (
  //     <div>
  //       <h1>NO DATA</h1>
  //     </div>
  //   )
  // }
  // Define columns for the DataGrid
  const columns: GridColDef[] = [
    { field: 'taskId', headerName: 'Task ID', flex:1 },
    { field: 'dependencies', headerName: 'Dependencies', flex: 1 },
    { field: 'duration', headerName: 'Duration', flex: 1 },
    { field: 'earliestStart', headerName: 'Earliest Start', flex: 1 },
    { field: 'earliestFinish', headerName: 'Earliest Finish', flex: 1 },
    { field: 'latestStart', headerName: 'Latest Start', flex: 1 },
    { field: 'latestFinish', headerName: 'Latest Finish', flex: 1 },
    { field: 'totalFloat', headerName: 'Total Float', flex: 1 },
    { field: 'freeFloat', headerName: 'Free Float', flex: 1 },
    { field: 'independentFloat', headerName: 'Independent Float', flex: 1 },
  ];

  // Prepare rows for the DataGrid
  const rows = data?.Result.map((task) => ({
    id: task.taskId, // MUI DataGrid requires 'id' as a unique field
    taskId: task.taskId,
    dependencies: task.dependencies.join(', '), // Join dependencies into a string
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
    <Box sx={{ height: 400, width: '100%' }}>
      <DataGrid rows={rows} columns={columns}  />
    </Box>
  );
};

export default CpmTable;