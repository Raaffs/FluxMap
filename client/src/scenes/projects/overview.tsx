import { useEffect, useState } from "react";
import { ProjectTaskDetailPage } from "../../components/task";
import { useParams } from "react-router-dom";
import { Projects } from "../../hooks/types";
import { Tab, Tabs, Box, Card, useTheme, Typography } from "@mui/material";
import { tokens } from "../../theme";
import Graphs from "../../components/graphs/LineGraphs";
import PertNormalDistributionChart, { PertTable } from "../../components/pert";
import {
  ContributerTaskChart,
  ContributerTaskTable,
} from "../../components/contributer";
import useFetchTaskData from "../../hooks/task";
import CpmNormalDistributionChart from "../../components/cpm";
import { useFetchPertData } from "../../hooks/pert";
import { useFetchCpmData } from "../../hooks/cpm";
export const ProjectOverview = () => {
  const { id } = useParams();
  console.log("projectid: ", id);

  const [project, setProject] = useState<Projects | null>(null);
  const [apiResponse, pertLoading, pertError] = useFetchPertData(id);
  const [cpmApiResponse, cpmLoading, cpmError] = useFetchCpmData(id);
  const [fetchTrigger, setFetchTrigger] = useState(false);
  console.log("PERT API: ", apiResponse);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any | null>(null);
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [tasks, fetchLoading, fetchError] = useFetchTaskData(id, fetchTrigger,setFetchTrigger);
  useEffect(() => {
    const fetchProjectOverviewByID = async () => {
      try {
        const response = await fetch(
          `http://localhost:4000/api/project/${id}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
          }
        );
        const data = await response.json();
        console.log("data : ", data);
        setProject(data);
        setLoading(false);
      } catch (err) {
        console.log(err)
        setError(err);
        setLoading(false);
      }
    };
    fetchProjectOverviewByID();
  }, [id]);
  if (tasks === null) {
    return <div>This Project is empty D:</div>;
  }

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  return (
    <Box
      sx={{
        overflowY: "auto",
        margin: "5px",
        width: "100%",
        minHeight: "100%",
        padding: "16px",
        border: "5px #ddd",
        borderRadius: "8px",
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[400] : "white",
        boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
        alignContent: "left",
        overflowX: "auto",
      }}
    >
      <Card
        sx={{
          width: "99%",
          minHeight: "20%",
          marginBottom: "16px",
          padding: "16px",
          borderRadius: "8px",
          backgroundColor:
            theme.palette.mode === "dark" ? colors.primary[400] : "#f",
          boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
          alignContent: "left",
          alignItems: "left",
          textAlign: "left",
          cursor: "pointer", // This makes the cursor a pointer on hover, indicating it's clickable
          "&:hover": {
            boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
          },
        }}
      >
        <Typography variant="h4" sx={{ fontWeight: "bold" }}>
          {project?.projectName}
        </Typography>
        <Typography variant="h5" color="textSecondary">
          {project?.projectDescription}
        </Typography>
        <Typography
          variant="h6"
          color="text.secondary"
          sx={{
            color:
              project?.projectDueDate &&
              new Date(project.projectDueDate) < new Date()
                ? colors.redAccent[500]
                : colors.greenAccent[400],
            fontWeight: "bold",
          }}
        >
          <strong>Due:</strong>{" "}
          {project?.projectDueDate
            ? new Date(project.projectDueDate).toLocaleDateString("en-GB", {
                day: "2-digit",
                month: "short",
                year: "numeric",
              })
            : "No due date"}
        </Typography>
      </Card>

      {/* Tabs Component */}
      <Tabs
        value={selectedTab}
        onChange={handleTabChange}
        aria-label="project tabs"
      >
        <Tab label="Task Details" />
        <Tab label="Contributer" />
        <Tab label="PERT" />
        <Tab label="CPM" />
      </Tabs>

      {/* Tab Panels */}
      <Box sx={{ padding: "16px" }}>
        <Card
          sx={{
            maxWidth: "99%",
            minHeight: "20%",
            marginBottom: "16px",
            padding: "16px",
            borderRadius: "8px",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[400] : "#f",
            boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
            alignContent: "left",
            alignItems: "left",
            textAlign: "left",
            cursor: "pointer", // This makes the cursor a pointer on hover, indicating it's clickable
            "&:hover": {
              boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
            },
          }}
        >
          {selectedTab === 0 && (
            <Box>
              <ProjectTaskDetailPage
                tasks={tasks}
                setFetchTrigger={setFetchTrigger}
              />
              <Graphs />
            </Box>
          )}
          {selectedTab === 1 && (
            <Box>
              <ContributerTaskTable tasksData={tasks} />
              <ContributerTaskChart tasksData={tasks} />
            </Box>
          )}
          {selectedTab === 2 && (
            <Box>
              <PertTable pertTasks={apiResponse?.data} tasks={tasks} />
              <PertNormalDistributionChart
                apiResponse={apiResponse}
                pertTasks={apiResponse?.data}
                tasks={tasks}
              />
            </Box>
          )}
          {selectedTab === 3 && (
            <Box>
              <CpmNormalDistributionChart data={cpmApiResponse} />
            </Box>
          )}
        </Card>
      </Box>
    </Box>
  );
};
