import { useState,useEffect } from "react";
import { Update } from "./types";

export const useFetchUpdates=(): {updates: Update[]|null,loading:boolean,error:any}=>{
    const [updates, setUpdates] = useState<Update[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    useEffect(()=>{
        setLoading(true)
        setError(null)
        setUpdates([])
        const fetchUpdates=async()=>{
            try{
                const res=await fetch(`http://localhost:4000/api/updates`,{
                    method: 'GET',
                    credentials:'include',
                    headers: {
                        "Content-Type": "application/json",
                      },
                })
                if(!res.ok){
                    const data=await res.json()
                    throw new Error(data.error || "Failed to fetch updates")
                }
                const data=await res.json()
                console.log("updates data",data)
                setUpdates(data)
            }catch (err:any){
                setError(err.message)
            }finally{
                setLoading(false)
            }
        }
        fetchUpdates()
    },[])
    return {updates,loading,error}
}

export const useFetchRecentUpdates=(): {updates: Update[]|null,loading:boolean,error:any}=>{
    const [updates, setUpdates] = useState<Update[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    useEffect(()=>{
        setLoading(true)
        setError(null)
        setUpdates([])
        const fetchUpdates=async()=>{
            try{
                const res=await fetch(`http://localhost:4000/api/updates/recent`,{
                    method: 'GET',
                    credentials:'include',
                    headers: {
                        "Content-Type": "application/json",
                      },
                })
                if(!res.ok){
                    const data=await res.json()
                    throw new Error(data.error || "Failed to fetch updates")
                }
                const data=await res.json()
                console.log("updates data",data)
                setUpdates(data)
            }catch (err:any){
                setError(err.message)
            }finally{
                setLoading(false)
            }
        }
        fetchUpdates()
    },[])
    return {updates,loading,error}
}

export const useFetchProjectUpdates=(id:number)=>{
    const [updates, setUpdates] = useState<Update[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    useEffect(()=>{
        setLoading(true)
        setError(null)
        setUpdates([])
        const fetchUpdates=async()=>{
            try{
                const res=await fetch(`http://localhost:4000/api/project/${id}/update`,{
                    method: 'GET',
                    credentials:'include'
                })
                if(!res.ok){
                    const data=await res.json()
                    throw new Error(data.error || "Failed to fetch updates")
                }
            }catch (err:any){
                setError(err.message)
            }finally{
                setLoading(false)
            }
        }
        fetchUpdates()
    })
    return {updates,loading,error}
}