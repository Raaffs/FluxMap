import { Box, useTheme } from "@mui/material";
import AdminPanelSettingsOutlinedIcon from "@mui/icons-material/AdminPanelSettingsOutlined";
import ManageAccountsOutlinedIcon from "@mui/icons-material/ManageAccountsOutlined";
import PermIdentityOutlinedIcon from "@mui/icons-material/PermIdentityOutlined";
import { tokens } from "../../theme";
import StatCard from "../../components/cards/statcard";
import { TaskStatusGraph } from "../../components/graphs/LineGraphs";
import {
  useRetrieveTaskApprovedGraph,
  useRetrieveTaskCompletedGraph,
} from "../../hooks/graphs";
import { number } from "mathjs";
import { useFetchRecentUpdates, useFetchUpdates } from "../../hooks/updates";
import { UpdateCard } from "../../components/cards/updates";
import { TaskStatusDonutChart } from "../../components/graphs/Donut";
import { useFetchOverdueTasks } from "../../hooks/task";
import OverdueTasks from "../../components/cards/overdue";
const Dashboard = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const {
    completedGraph,loading: completedGraphLoading, error: completedGraphError} = useRetrieveTaskCompletedGraph();
  const { approvedGraph, approvedLoading, approvedError } = useRetrieveTaskApprovedGraph();
  const { updates, loading, error } = useFetchRecentUpdates();
  const {tasks: overdueTasks, loading: overdueLoading, error: overdueError} = useFetchOverdueTasks();
  const taskStatusData = {
    labels:
      completedGraph?.XAxis?.map((x) => new Date(x).toLocaleDateString()) || [],
    datasets: [
      {
        label: "Tasks Completed",
        data: completedGraph?.YAxis?.map((x) => Number(x)) || [],
        borderColor: "rgba(0, 123, 255, 1)", // Bright blue for white mode
        backgroundColor: "rgba(0, 123, 255, 0.15)", // Soft blue fill
        tension: 0.3,
        fill: true,
      },
      {
        label: "Tasks Approved",
        data: approvedGraph?.YAxis?.map((x) => Number(x)) || [],
        borderColor: "rgba(40, 167, 69, 1)", // Clean green for white mode
        backgroundColor: "rgba(40, 167, 69, 0.15)", // Soft green fill
        tension: 0.3,
        fill: true,
      },
    ],
  };
  const donutData = {
    labels: ["Completed", "Active Pending", "Overdue"],
    datasets: [
      {
        label: "Tasks",
        data: [10, 5, 2],
        backgroundColor: ["#7CD1B8", "#FFE182", "#FF9AA2"],

        borderColor: ["#4FB59C", "#FFD447", "#FF6B81"],
        borderWidth: 1,
      },
    ],
  };

  return (
    <Box
      m="30px"
      sx={{
        backgroundColor:
          theme.palette.mode === "dark" ? colors.primary[500] : "white",
        maxHeight: "100%",
        height: "100%",
        overflowY: "auto",
        padding: "25px",
        borderRadius: "16px",
      }}
    >
      <Box
        display="grid"
        gridTemplateColumns="repeat(12, 1fr)"
        gridAutoRows="140px"
        gap="20px"
      >
        <Box
          gridColumn="span 3"
          display="flex"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : "rgba(0, 200, 150, 0.3)",
            borderRadius: "16px",
            boxShadow: "0 6px 24px rgba(0, 0, 0, 0.15)",
          }}
          justifyContent="center"
          alignItems="center"
        >
          <StatCard
            title="4"
            subtitle="Your Projects"
            fontColor={colors.greenAccent[200]}
            icon={
              <AdminPanelSettingsOutlinedIcon
                sx={{
                  color: colors.greenAccent[200],
                  fontSize: "26px",
                  margin: "5px",
                }}
              />
            }
            path="/projects/admin"
          />
        </Box>
        <Box
          gridColumn="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : "rgba(104, 112, 250, 0.4)",
            borderRadius: "16px",
            boxShadow: "0 6px 24px rgba(0, 0, 0, 0.15)",
          }}
        >
          <StatCard
            title="4"
            subtitle="Managed Projects"
            fontColor={colors.blueAccent[400]}
            icon={
              <ManageAccountsOutlinedIcon
                sx={{
                  color: colors.blueAccent[400],
                  fontSize: "26px",
                  margin: "5px",
                }}
              />
            }
            path="/projects/manager"
          />
        </Box>
        <Box
          gridColumn="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : "rgba(162, 89, 255, 0.3)",
            borderRadius: "16px",
            boxShadow: "0 6px 24px rgba(0, 0, 0, 0.15)",
          }}
        >
          <StatCard
            title="2"
            subtitle="Allocated Projects"
            fontColor="#a259ff"
            icon={
              <PermIdentityOutlinedIcon
                sx={{
                  color: "#a259ff",
                }}
              />
            }
            path="/projects/allocated"
          />
        </Box>
        <Box
          gridColumn="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            backgroundColor:
              theme.palette.mode === "dark"
                ? colors.primary[600]
                : "rgba(162, 89, 255, 0.3)",
            borderRadius: "16px",
            boxShadow: "0 6px 24px rgba(0, 0, 0, 0.15)",
          }}
        >
          <StatCard
            title="2"
            subtitle="Allocated Projects"
            fontColor="#a259ff"
            icon={
              <PermIdentityOutlinedIcon
                sx={{
                  color: "#a259ff",
                }}
              />
            }
            path="/projects/allocated"
          />
        </Box>
        <Box
          gridColumn="span 8"
          gridRow="span 3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          sx={{
            background:
              theme.palette.mode === "dark"
                ? "#0d141f"
                : "linear-gradient(135deg, #f0fffc, #f3fff0)",
            width: "100%",
            height: "100%",
            borderRadius: "16px",
            boxShadow: "0 6px 24px rgba(0, 0, 0, 0.15)",
          }}
        >
          <Box flex={1} width="100%" height="100%">
            <TaskStatusGraph
              data={taskStatusData}
              title="Tasks Approved and Completed Over Time"
              xLabel="last 7 days"
              yLabel="Number of Tasks"
              maintainRatio={false}
            />
          </Box>
        </Box>
        <Box
          gridColumn="span 4"
          gridRow="span 3"
          display="flex"
          flexDirection="column"
          alignItems="flex-start"
          justifyContent="flex-start"
          width="100%"
          gap={2}
          sx={{
            boxShadow: "0 8px 24px rgba(0, 0, 0, 0.12)",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[600] : "#ffffff",
            borderRadius: "16px",
            padding: 2,
            overflow: "hidden",
          }}
        >
          <Box
            display="flex"
            flexDirection="column"
            gap={2}
            width="100%"
            sx={{
              overflowY: "auto",
              paddingRight: 1,
              maxHeight: "100%", // or a fixed height like "500px" if needed
              "&::-webkit-scrollbar": {
                width: "6px",
              },
              "&::-webkit-scrollbar-thumb": {
                backgroundColor: "rgba(0, 0, 0, 0.15)",
                borderRadius: "4px",
              },
            }}
          >
            {updates?.map((update, index) => (
              <Box
                key={update.id || index}
                width="100%"
                borderRadius="8px"
                sx={{
                  backgroundColor:
                    theme.palette.mode === "dark"
                      ? colors.primary[500]
                      : "white",
                  transition: "background-color 0.2s ease",
                  "&:hover": {
                    backgroundColor:
                      theme.palette.mode === "dark"
                        ? colors.primary[400]
                        : "#eeeeee",
                  },
                }}
              >
                <UpdateCard update={update} />
              </Box>
            ))}
          </Box>
        </Box>
        <Box
          gridColumn="span 3"
          gridRow="span 3"
          display="flex"
          flexDirection="column"
          alignItems="flex-start"
          justifyContent="flex-start"
          width="100%"
          gap={2}
          sx={{
            boxShadow: "0 8px 24px rgba(0, 0, 0, 0.12)",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[600] : "#ffffff",
            borderRadius: "16px",
            padding: 2,
            overflow: "hidden",
          }}
        >
          <OverdueTasks
            tasks={overdueTasks}
          />
        </Box>

        <Box
          gridColumn="span 3"
          gridRow="span 3"
          display="flex"
          flexDirection="column"
          alignItems="flex-start"
          justifyContent="flex-start"
          width="100%"
          gap={2}
          sx={{
            boxShadow: "0 8px 24px rgba(0, 0, 0, 0.12)",
            backgroundColor:
              theme.palette.mode === "dark" ? colors.primary[600] : "#ffffff",
            borderRadius: "16px",
            padding: 2,
            overflow: "hidden",
          }}
        >
          <TaskStatusDonutChart data={donutData} title="Task Status Overview" />
        </Box>
      </Box>
    </Box>
  );
};

export default Dashboard;
