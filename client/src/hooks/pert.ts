import { useState, useEffect } from "react";
import { ApiResponse } from "./types";

export const useFetchPertData = (id: any): [ApiResponse | null, boolean, any] => {
  const [apiResponse, setApiResponse] = useState<ApiResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    setLoading(true); // Reset loading state on id change
    setError(null); // Clear any previous errors
    setApiResponse(null); // Clear previous data

    const fetchData = async () => {
      try {
        const response = await fetch(`http://localhost:4000/api/project/${id}/pert`, {
          method: "GET",
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        setApiResponse(data);
      } catch (err: any) {
        console.error(err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]); // Re-run the effect when `id` changes

  return [apiResponse, loading, error];
};
