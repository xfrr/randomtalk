import { Message } from "@/internal/messaging";
/**
 * Command object
 * @class Command
 * @description Command is a type of message to perform an action in the application.
 * @export
 * @abstract
 */
export abstract class Command implements Message {
  /**
   * Creates an instance of Command.
   * @param {CommandType} type The command type
   * @param {CommandName} name The command name
   * @param {CommandPayload} payload The command payload
   * @param {Date} [timestamp] The command timestamp
   * @memberof Command
   * @description Command constructor
   * @constructor
   * @abstract
   * @export
   */
  constructor(
    public type: CommandType,
    public payload: CommandPayload,
    public timestamp: Date
  ) {}
}

/**
 * Command type
 * @type {string}
 */
export type CommandType = string;
/**
 * Command name
 * @type {string}
 */
export type CommandName = string;
/**
 * Command payload
 * @type {any}
 */
export type CommandPayload = any;

/**
 * Command handler
 * @interface
 * @export
 * @abstract
 * @class CommandHandler
 * @description CommandHandler handles a command.
 * @template T The type of Command
 */
export interface CommandHandler<T extends Command> {
  /**
   * Handles a command.
   * @param {T} command The command to handle
   * @returns {Promise<void>}
   */
  handle(command: T): Promise<void>;
}
