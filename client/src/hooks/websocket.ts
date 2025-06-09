import React, { useState, useEffect } from "react";

export function useWebSocket(
  url:string, 
  setUpdateCount: React.Dispatch<React.SetStateAction<string>> 
) {

  const startWebSocket =()=>{
      const socket = new WebSocket(url);
  
      socket.onopen = () => {
        console.log("WebSocket connected");
      };
  
      socket.onmessage = (event) => {
        console.log("📨 Message from server:", event.data);
        // Parse data & update count, for example:
        try {
          const data = JSON.parse(event.data);
            setUpdateCount(data);
          
        } catch(error) {
          console.error(error)
        }
      };
  
      socket.onclose = () => {
        console.log("WebSocket disconnected");
      };
  
      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
  
      return () => {
        socket.close();
      };
  }
  return [startWebSocket]
}
