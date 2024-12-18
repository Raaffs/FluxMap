import { useEffect, useState } from "react";
import { Projects } from "./types";

export const useRetrieveProjectsFrom = (link: string) => {
  const [projects, setProjects] = useState<Projects[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(link, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch projects");
        }
        return response.json();
      })
      .then((data) => {
        const mappedProjects: Projects[] = data.map((item: any) => ({
          projectID: item.projectID,
          projectName: item.projectName,
          projectDescription: item.projectDescription || null,
          projectStartDate: item.projectStartDate,
          projectDueDate: item.projectDueDate,
        }));
        setProjects(mappedProjects);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching projects:", err);
        setError("Failed to load projects.");
        setLoading(false);
      });
  }, [link]);

  return { projects, loading, error };
};

export const usePostProject = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const postProject = (newProject: Projects, link: string) => {
    setLoading(true);
    setError(null);
    setSuccess(false);

    fetch(link, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newProject),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to post the project");
        }
        return response.json();
      })
      .then(() => {
        setSuccess(true);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error posting project:", err);
        setError("Failed to post the project.");
        setLoading(false);
      });
  };

  return { postProject, loading, error, success };
};

export const usePutProject = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const putProject = (updatedProject: Projects, link: string) => {
    setLoading(true);
    setError(null);
    setSuccess(false);

    fetch(link, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(updatedProject),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to update the project");
        }
        return response.json();
      })
      .then(() => {
        setSuccess(true);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error updating project:", err);
        setError("Failed to update the project.");
        setLoading(false);
      });
  };

  return { putProject, loading, error, success };
};
