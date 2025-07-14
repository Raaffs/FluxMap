import { useCallback } from "react";

export function useWebSocket(
  url: string,
  setUpdateCount: React.Dispatch<React.SetStateAction<number>>,
  setInvitationtCount: React.Dispatch<React.SetStateAction<number>>
) {
  const startWebSocket = useCallback(() => {
    const socket = new WebSocket(url);

    socket.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setUpdateCount(() => {
          console.log("Updated count to: ",data.updateNotification)
          return data.updateNotification
        });
        setInvitationtCount(() => data.invitateNotification);
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
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
  }, [url, setUpdateCount, setInvitationtCount]);

  return [startWebSocket];
}
