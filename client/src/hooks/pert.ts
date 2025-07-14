import { useState, useEffect } from "react";
import { ApiResponse, PertData } from "./types";

export const useFetchPertData = (
  id: any,
  fetchPertTrigger: boolean,
  setFetchPertTrigger: React.Dispatch<React.SetStateAction<boolean>>
): [ApiResponse | null, boolean, any] => {
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
        setFetchPertTrigger(false)
      }
    };

    fetchData();
  }, [id,fetchPertTrigger]); // Re-run the effect when `id` changes

  return [apiResponse, loading, error];
};


export const useAddPert=(id:string)=>{
  const [pertloading, setLoading] = useState<boolean>(false);
  const [perterror, setError] = useState<string | null>(null);
  const [pertsuccess, setSuccess] = useState<boolean>(false);
  const addPert= async (pert:PertData)=>{
    setLoading(true);
    setError(null);
    setSuccess(false);
    console.log("add pert: ",pert)
    try{

      const res=await fetch(`http://localhost:4000/api/project/${id}/pert`,{
        method:"POST",
        credentials:'include',
        headers:{
          'Content-Type':'application/json'
        },
        //backend takes this as an array.
        //since endpoint is created so that you can technically add
        //multiple pert tasks at once.
        //but UI is designed to add one at a time.
        //so we wrap it in an array.
        body:JSON.stringify([{
          parentTaskID: pert.parentTaskID,
          predecessorTaskId: pert.predecessorTaskId===""? null : pert.predecessorTaskId,
          optimistic: pert.optimistic,
          pessimistic: pert.pessimistic,
          mostLikely: pert.mostLikely,
        }])
      })
    }catch(err:any){
      console.error(err)
      setError(err.error)
    }finally{
      setLoading(false)
    }
  }
  return {addPert,pertloading,perterror,pertsuccess}
}