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
      console.log("📨 Message from server:", event.data);
      try {
        const data = JSON.parse(event.data);
        setUpdateCount(data.updateNotification);
        setInvitationtCount(data.invitationNotification);
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
