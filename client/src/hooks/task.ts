import React from "react";
import { useState, useEffect } from "react";
import { tasks } from "./types";

const useFetchTaskData = (
  projectId: string | undefined,
  fetchTrigger: boolean,
  setFetchTrigger:React.Dispatch<React.SetStateAction<boolean>>
): [tasks[], boolean, string] => {
  const [tasks, setTasks] = useState<tasks[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(
          `http://localhost:4000/api/project/${projectId}/tasks`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
          }
        );
        const data = await response.json();
        if (!response.ok) {
          return;
        }
        setTasks(data);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch tasks");
        setLoading(false);
      }finally{
        setFetchTrigger(false)
      }
    };

    fetchTasks();
  }, [projectId, fetchTrigger]);

  return [tasks, loading, error];
};


export const useFetchOverdueTasks = () => {
  const [tasks, setTasks] = useState<tasks[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(
          `http://localhost:4000/api/graph/tasks/overdue`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
          }
        );
        const data = await response.json();
        if (!response.ok) {
          return;
        }
        setTasks(data.overdue);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch tasks");
        setLoading(false);
      }
    };

    fetchTasks(); 
  },[]);

  return {tasks, loading, error};
};

export const useFetchUpcomingTasks = () => {
  const [tasks, setTasks] = useState<tasks[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(
          `http://localhost:4000/api/graph/tasks/upcoming`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // Include credentials (cookies, authentication tokens, etc.)
          }
        );
        const data = await response.json();
        if (!response.ok) {
          return;
        }
        setTasks(data.upcoming);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch tasks");
        setLoading(false);
      }
    };

    fetchTasks();
  },[]);
}
export default useFetchTaskData;
