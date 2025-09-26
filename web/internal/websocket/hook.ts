// useWebSocket.ts
// -----------------------------------------------------------------------------
// This custom React hook builds on the WebSocketService, initiating the connection
// and managing state for the latest notification.
// -----------------------------------------------------------------------------

import { useState, useEffect, useRef } from "react";
import {
  WebSocketService,
  WebSocketConfig,
} from "@/internal/websocket/service";

interface UseWebSocketResult {
  latestNotification: any;
  sendCommand: (cmd: unknown) => void;
  closeConnection: () => void;
}

export function useWebSocket(config: WebSocketConfig): UseWebSocketResult {
  const [latestNotification, setLatestNotification] = useState<any>(null);

  // We use a ref to store the service instance so it persists across renders
  const wsServiceRef = useRef<WebSocketService | null>(null);

  // Create the WebSocketService on mount
  useEffect(() => {
    wsServiceRef.current = new WebSocketService({
      ...config,
      // Provide an onNotification to hook into React state
      onNotification: (msg) => {
        setLatestNotification(msg);
        if (config.onNotification) {
          config.onNotification(msg);
        }
      },
    });

    // Initiate the connection
    wsServiceRef.current.connect();

    // Cleanup on unmount
    return () => {
      wsServiceRef.current?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const sendCommand = (cmd: unknown) => {
    wsServiceRef.current?.sendCommand(cmd);
  };

  const closeConnection = () => {
    wsServiceRef.current?.close();
  };

  return {
    latestNotification,
    sendCommand,
    closeConnection,
  };
}
