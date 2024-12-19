import { useEffect, useState } from "react";
import { ProjectTaskDetailPage } from "../../components/task";
import { useParams } from "react-router-dom";
import { Projects } from "../../hooks/types";
import { Tab, Tabs, Box, Card, useTheme, Typography } from "@mui/material";
import { tokens } from "../../theme";

export const ProjectOverview = () => {
    const { id } = useParams();
    console.log("projectid: ", id);
    
    const [project, setProject] = useState<Projects | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<any | null>(null);
    const [selectedTab, setSelectedTab] = useState<number>(0);
    
    const theme = useTheme();
    const colors = tokens(theme.palette.mode);
    
    useEffect(() => {
        const fetchProjectOverviewByID = async () => {
            try {
                const response = await fetch(`http://localhost:4000/api/project/${id}`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
                });
                const data = await response.json();
                console.log("data : ", data);
                setProject(data);
                setLoading(false);
            } catch (err) {
                setError(err);
                setLoading(false);
            }
        };
        fetchProjectOverviewByID();
    }, [id]);

    const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedTab(newValue);
    };

    return (
        <Box
            sx={{
                margin: '5px',
                minWidth: '100%',
                minHeight: "100%",
                overflowY: "auto",
                padding: "16px",
                border: "5px #ddd",
                borderRadius: "8px",
                backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : 'white',
                boxShadow: "0 4px 8px rgba(0,0,0,0.5)",
                alignContent: 'left',
            }}
        >
            <Card
                sx={{
                    maxWidth: '100%',
                    minHeight: '20%',
                    marginBottom: "16px",
                    padding: "16px",
                    borderRadius: "8px",
                    backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : '#f',
                    boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
                    alignContent: 'left',
                    alignItems: 'left',
                    textAlign: 'left',
                    cursor: 'pointer', // This makes the cursor a pointer on hover, indicating it's clickable
                    '&:hover': {
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
                        color: project?.projectDueDate && new Date(project.projectDueDate) < new Date() ? colors.redAccent[500] : colors.greenAccent[400],
                        fontWeight: 'bold',
                    }}
                >
                    <strong>Due:</strong> {project?.projectDueDate ? new Date(project.projectDueDate).toLocaleDateString('en-GB', {
                        day: '2-digit',
                        month: 'short',
                        year: 'numeric',
                    }) : 'No due date'}
                </Typography>
            </Card>

            {/* Tabs Component */}
            <Tabs value={selectedTab} onChange={handleTabChange} aria-label="project tabs">
                <Tab label="Task Details" />
                <Tab label="Other Tab" />
            </Tabs>

            {/* Tab Panels */}
            <Box sx={{ padding: '16px' }}>
                {selectedTab === 0 && (
                    <Card
                        sx={{
                            maxWidth: '100%',
                            minHeight: '20%',
                            marginBottom: "16px",
                            padding: "16px",
                            borderRadius: "8px",
                            backgroundColor: theme.palette.mode === 'dark' ? colors.primary[400] : '#f',
                            boxShadow: "0 2px 4px rgba(0,0,0,0.5)",
                            alignContent: 'left',
                            alignItems: 'left',
                            textAlign: 'left',
                            cursor: 'pointer', // This makes the cursor a pointer on hover, indicating it's clickable
                            '&:hover': {
                                boxShadow: "0 4px 8px rgba(0,0,0,0.8)", // Optional: add a hover effect to emphasize the card
                            },
                        }}
                    >
                        <ProjectTaskDetailPage />
                    </Card>
                )}
                {selectedTab === 1 && (
                    <Typography variant="body1">This is the content for the other tab.</Typography>
                )}
            </Box>
        </Box>
    );
};
