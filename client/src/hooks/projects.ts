import { useEffect, useState } from "react";
import { Projects, UserRole } from "./types";
import { normalizeAccessMap } from "../helpers/filter";
export const useRetrieveProjectsFrom = (
  link: string,
  fetchTrigger: boolean,
  setFetchTrigger:React.Dispatch<React.SetStateAction<boolean>>
  
) => {
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
      })
      .finally(() => {
        setFetchTrigger(false); 
      });
  },[link,fetchTrigger]);

  return { projects, loading, error };
};

export const usePostProject = (  
  setUserProjectRoleMap: React.Dispatch<React.SetStateAction<Record<number,UserRole>>>

) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const postProject = async (
    newProject: Projects, 
    link: string,
    setFetchTrigger: React.Dispatch<React.SetStateAction<boolean>>
  ) => {
    setLoading(true);
    setError(null);
    setSuccess(false);
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
        
        const errors=[
          errorData.error,
          errorData.Errors.description,
          errorData.Errors.name,
        ]
        .filter(Boolean)
        .join("\n")
        setError(errors)
        return
      }
      const data=await response.json()
      let noramlizedRoleMap=normalizeAccessMap(data.roles as Record<UserRole, number[] | null>)
      setUserProjectRoleMap(noramlizedRoleMap)
      localStorage.setItem("userProjectRoleMap", JSON.stringify(noramlizedRoleMap));
      setSuccess(true);
      setFetchTrigger(true);
    } catch (err) {
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