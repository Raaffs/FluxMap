import { useState, useEffect } from "react";
import { CpmApiResponse, CpmResult } from "./types";

export const useFetchCpmData =(id: any):[CpmResult|null,boolean,any]=>{
    const [data, setData] = useState<CpmResult|null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(()=>{
        setLoading(true)
        setError(null)
        setData(null)

        const fetchData = async () => {
            try {
              const response = await fetch(`http://localhost:4000/api/project/${id}/cpm`, {
                method: "GET",
                credentials: "include",
              });
    
              if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
              }
              const data = await response.json();
              setData(data?.result);
            } catch (err: any) {
              console.error(err);
              setError(err.message);
            } finally {
              setLoading(false);
            }
          };
      
          fetchData();
    },[id])
    return [data,loading,error]
}     

export const usePostCpmData=(id: Number)=>{
  const [cpmloading, setLoading] = useState<boolean>(false);
  const [cpmerror, setError] = useState<string | null>(null);
  const [cpmsuccess, setSuccess] = useState<boolean>(false);
const addCPM= async (CPM:CpmApiResponse)=>{
    setLoading(true);
    setError(null);
    setSuccess(false);
    try{
      const res=await fetch(`http://localhost:4000/api/project/${id}/cpm`,{
        method:"POST",
        credentials:'include',
        headers:{
          'Content-Type':'application/json'
        },
        //backend takes this as an array for some reason. I don't know why. 
        body:JSON.stringify([{
            taskId:CPM.taskId,
            taskName:CPM.taskName,
            parentProjectID:id,
            dependencies: CPM.dependencies,
            duration:CPM.duration,
            earliestStart:0,
            earliestFinish:0,
            latestStart:0,
            latestFinish:0,
            totalFloat:0,
            freeFloat:0,
            independentFloat:0,
            isCriticalPath:0,
        }])
      })
      if(!res.ok){
      }
    }catch(err:any){
      console.error(err)
      setError(err.error)
    }finally{
      setLoading(false)
    }
  }
  return {addCPM,cpmloading,cpmerror,cpmsuccess}
}