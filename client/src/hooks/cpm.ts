import { useState, useEffect } from "react";
import { CpmResult } from "./types";

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