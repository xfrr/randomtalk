// WebSocketService.ts
// -----------------------------------------------------------------------------
// This class handles bi-directional communication with a WebSocket server.
// It provides methods to connect, send data, manage reconnection, and emit
// or dispatch received messages. Following SOLID's SRP, it focuses only
// on WebSocket communication logic.
// -----------------------------------------------------------------------------

import { Command } from "../command";
import { v4 as uuidv4 } from "uuid";

export interface WebSocketConfig {
  url: string; // The WebSocket server endpoint
  maxRetries?: number; // Maximum number of reconnection attempts
  reconnectDelay?: number; // Delay between attempts in ms
  protocols?: string | string[]; // WebSocket subprotocols, if needed
  onNotification?: (message: any) => void; // Callback for incoming messages
  onOpen?: () => void; // Callback triggered on successful connection
  onClose?: (event: CloseEvent) => void; // Callback triggered on close event
  onError?: (event: Event) => void; // Callback triggered on error
}

export class WebSocketService {
  private ws: WebSocket | null = null;
  private config: WebSocketConfig;
  private retryCount = 0;
  private forcedClose = false;

  constructor(config: WebSocketConfig) {
    this.config = config;
  }

  /**
   * Initiates the WebSocket connection.
   */
  public connect(): void {
    if (!("WebSocket" in global)) {
      console.error("WebSocket is not supported in this environment.");
      return;
    }

    this.forcedClose = false;

    this.ws = new WebSocket(this.config.url, this.config.protocols || []);

    this.ws.onopen = () => {
      this.retryCount = 0;
      this.config.onOpen && this.config.onOpen();
      console.info("WebSocket connected.");
    };

    this.ws.onmessage = (event) => {
      try {
        // Attempt to parse the incoming data as JSON
        const parsedData = JSON.parse(event.data);
        if (this.config.onNotification) {
          this.config.onNotification(parsedData);
        }
      } catch (e) {
        // Fallback: pass raw string if not JSON
        this.config.onNotification && this.config.onNotification(event.data);
      }
    };

    this.ws.onerror = (event) => {
      this.config.onError && this.config.onError(event);
      console.error("WebSocket error observed:", event);
    };

    this.ws.onclose = (event) => {
      this.config.onClose && this.config.onClose(event);
      // Attempt reconnection unless closed intentionally
      if (!this.forcedClose && this.shouldReconnect()) {
        this.reconnect();
      }
    };
  }

  /**
   * Sends a command (string or object). If object, it's stringified as JSON.
   */
  public sendCommand(command: Command): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn("WebSocket is not open. Cannot send command:", command);
      return;
    }

    const payload = JSON.stringify({
      kind: "command",
      data: {
        id: uuidv4(),
        type: command.type,
        payload: command.payload,
        timestamp: command.timestamp || new Date(),
      },
    });

    this.ws.send(payload);
  }

  /**
   * Closes the WebSocket connection gracefully.
   */
  public close(): void {
    this.forcedClose = true;
    this.ws?.close();
  }

  /**
   * Decide whether to attempt reconnection based on retryCount vs maxRetries.
   */
  private shouldReconnect(): boolean {
    const { maxRetries = 5 } = this.config;
    return this.retryCount < maxRetries;
  }

  /**
   * Attempts to reconnect after a specified delay.
   * Uses exponential backoff or a simple linear delay, as preferred.
   */
  private reconnect(): void {
    const { reconnectDelay = 1000 } = this.config;
    this.retryCount += 1;
    const delay = this.retryCount * reconnectDelay; // simple linear backoff
    console.info(`Reconnecting in ${delay} ms...`);
    setTimeout(() => {
      this.connect();
    }, delay);
  }
}
