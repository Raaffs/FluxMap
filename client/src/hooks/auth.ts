import { error } from "console"

export const uselogout=()=>{
    fetch('http://localhost:4000/api/logout',{
        method:'POST',
        credentials: 'include'
    })
    .then(respone=>{
        if(!respone.ok){
            throw new Error("Error logging out")
        }
    })
}