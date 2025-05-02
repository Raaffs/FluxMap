import { useState,useEffect } from "react";
import { DisplayInvitations } from "./types";

export const useFetchInvitations = (): [DisplayInvitations[]|null,boolean,any]=>{
    const [data, setData] = useState<DisplayInvitations[]|null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    useEffect(()=>{
        setLoading(true)
        setError(null)
        setData(null)

        const fetchInvitations=async()=>{
            try {
                const response = await fetch(`http://localhost:4000/api/invitation`,{
                    method: 'GET',
                    credentials:'include'
                })
                const data = await response.json();
                setData(data);
  
            }catch(err: any){
                console.error(err)
                setError(err)
            }finally{
                setLoading(false)
            }
        }

        fetchInvitations()
    },[])
    return [data,loading,error]
}

export const useSendInvitation = (projectID: string):[DisplayInvitations|null,boolean,any]=>{
    const [data, setData] = useState<DisplayInvitations|null>(null);
   const [loading, setLoading] = useState(true);
   const [error, setError] = useState(null);
   useEffect(()=>{
       setLoading(true)
       setError(null)
       setData(null)

       const fetchInvitations=async()=>{
           try {
               const response = await fetch(`http://localhost:3001/api/project/${projectID}/invite`,{
                   method: 'GET',
                   credentials:'include'
               })
               const data = await response.json();
               setData(data);
 
           }catch(err: any){
               console.error(err)
               setError(err)
           }finally{
               setLoading(false)
           }
       }

       fetchInvitations()
   })
   return [data,loading,error]
}

export const useInvite = (id: string) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const inviteUser = async (username: string, role: string) => {
    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const res = await fetch(`http://localhost:4000/api/project/${id}/invite`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, role }),
      });

      if (!res.ok) {
        const data = await res.json();
        console.log("error: ",data)
        throw new Error(data.error || "Failed to send invite");
      }

      setSuccess(true);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return { inviteUser, loading, error, success };
};


export const useConfirmInvitation = () => {
  const [inviteloading, setLoading] = useState<boolean>(false);
  const [inviteerror, setError] = useState<string | null>(null);
  const [invitesuccess, setSuccess] = useState<boolean>(false);

  const confirmInvitation = async (invitationId: number, status: string) => {
    setLoading(true);
    setError(null);
    setSuccess(false);
    console.log("id invite",invitationId,status)
    try {
      const res = await fetch(`http://localhost:4000/api/invitation/${invitationId}`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });

      setSuccess(true);
    } catch (err: any) {
      console.log('error invite',err)
      setError(err.error);
    } finally {
      setLoading(false);
    }
  };

  return { confirmInvitation, inviteloading, inviteerror, invitesuccess };
};
