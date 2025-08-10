import { useEffect, useState } from "react";
import { Graphs } from "./types";

export const useRetrieveTaskCompletedGraph = () => {
  const [completedGraph, setCompletedGraph] = useState<Graphs>();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchCompletedGraph = async () => {
      try {
        const response = await fetch('http://localhost:4000/api/graph/tasks/completed', {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json"
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        setCompletedGraph(data.graph);
      } catch (err: any) {
        console.error("Error fetching completed graph:", err);
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCompletedGraph(); 
  }, []);

  return { completedGraph, loading, error };
};


export const useRetrieveTaskApprovedGraph = () => {
  const [approvedGraph, setApprovedGraph] = useState<Graphs>();
  const [approvedLoading, setLoading] = useState<boolean>(true);
  const [approvedError, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchApprovedGraph = async () => {
      try {
        const response = await fetch('http://localhost:4000/api/graph/tasks/approved', {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json"
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        setApprovedGraph(data.graph);
      } catch (err: any) {
        console.error("Error fetching completed graph:", err);
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchApprovedGraph(); 
  }, []);

  return { approvedGraph, approvedLoading, approvedError };
};


export const useRetrieveTaskStatusBreakdown = () => {
  const [breakdownGraph, setbreakdownGraph] = useState<Graphs>();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchbreakdownGraph = async () => {
      try {
        const response = await fetch('http://localhost:4000/api/graph/tasks/breakdown', {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json"
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const data = await response.json();
        setbreakdownGraph(data.graph);
      } catch (err: any) {
        console.error("Error fetching completed graph:", err);
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchbreakdownGraph(); 
  }, []);

  return { breakdownGraph, loading, error };
};
