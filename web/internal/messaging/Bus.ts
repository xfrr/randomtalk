// Message Interface defines the contract for a message object.
// It has a type property that should be a string and a payload property that can be any type.
export interface Message {
  type: string;
  payload?: any;
  timestamp?: Date;
}

// MessageBus Interface defines the contract for a message bus implementation.
// It allows sending messages and registering handlers for specific message types.
export interface MessageBus {
  send(message: Message): Promise<void>;
  registerHandler(
    messageType: string,
    handler: (message: Message) => Promise<void>
  ): void;
}

// InternalMessageBus is a simple in-memory implementation of the IMessageBus interface.
// It allows registering handlers for specific message types and sending messages.
export class InternalMessageBus implements MessageBus {
  private handlers: Map<string, (message: Message) => Promise<void>> =
    new Map();

  constructor() {}

  async send(message: Message): Promise<void> {
    const handler = this.handlers.get(message.type);
    if (handler) {
      return await handler(message);
    }

    console.error(`Message handler not found for message: ${message.type}`);
  }

  registerHandler(
    messageType: string,
    handler: (message: Message) => Promise<void>
  ): void {
    if (this.handlers.has(messageType)) {
      console.warn(
        `Handler for message type ${messageType} already exists. It will be overwritten.`
      );
    }
    this.handlers.set(messageType, handler);
  }
}

export default InternalMessageBus;
