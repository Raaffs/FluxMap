import { useEffect, useState } from "react";
import { Projects } from "./types";
export const useRetrieveProjectsFrom = (link: string) => {
  const [projects, setProjects] = useState<Projects[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(link, {
      method: "GET",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => response.json()) // Proceed without checking `response.ok`
      .then((data) => {
        if(!data){
          setProjects([])
          setLoading(false)
          return
        }
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
        setError(error); // Handle all errors here
        setLoading(false);
      });
  },[link]);

  return { projects, loading, error };
};

export const usePostProject = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const postProject = async (newProject: Projects, link: string) => {
    setLoading(true);
    setError(null);
    setSuccess(false);
    console.log("heelo",JSON.stringify(newProject), newProject)
    try {
      const response = await fetch("http://localhost:4000/api/projects", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newProject),
        credentials: "include",
      });

      if (!response.ok) {
        const errorData = await response.json();
        console.log("project post err",errorData)
        throw new Error(errorData.error || "Failed to post the project");
      }

      await response.json(); // Assuming success response has no additional data
      setSuccess(true);
    } catch (err) {
      console.error("Error posting project:", err);
      setError(err instanceof Error ? err.message : "An unknown error occurred");
    } finally {
      setLoading(false);
    }
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