import React, {
  createContext,
  useContext,
  useRef,
  useEffect,
  ReactNode,
  useState,
} from "react";
import {
  WebSocketService,
  WebSocketConfig,
} from "@/internal/websocket/service";
import { Command } from "@/internal/command";

interface WebSocketContextType {
  sendCommand: (command: Command) => void;
  latestNotification: any;
}

// Create the context
const WebSocketContext = createContext<WebSocketContextType>({
  sendCommand: () => {
    throw new Error("sendCommand must be used within a WebSocketProvider");
  },
  latestNotification: null,
});

// Provider component
interface WebSocketProviderProps {
  config: WebSocketConfig;
  children: ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  config,
  children,
}) => {
  // Keep the service in a ref so it's not recreated on every render
  const serviceRef = useRef<WebSocketService | null>(null);

  // Track the most recent message (optional, for convenience)
  const [latestNotification, setLatestNotification] = useState<any>(null);

  // Initialize the service once
  if (!serviceRef.current) {
    serviceRef.current = new WebSocketService({
      ...config,
      // We capture messages globally here, then also call user-supplied onNotification if provided
      onNotification: (msg: any) => {
        setLatestNotification(msg);
        config.onNotification?.(msg);
      },
    });
  }

  // Connect on mount only
  useEffect(() => {
    serviceRef.current?.connect();
    // Cleanup by closing on unmount if desired
    return () => {
      serviceRef.current?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // The context value includes the service and latest notification
  const contextValue: WebSocketContextType = {
    sendCommand: serviceRef.current.sendCommand.bind(serviceRef.current),
    latestNotification,
  };

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  );
};

export function useWebSocketContext() {
  const context = useContext(WebSocketContext);
  if (!context.sendCommand) {
    throw new Error(
      "useWebSocketContext must be used within a WebSocketProvider"
    );
  }
  return context;
}
