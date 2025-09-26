import { Command } from "./Command";

/**
 * Command dispatcher
 * @export
 * @abstract
 * @class Dispatcher
 * @description Dispatcher dispatches commands to the appropriate handler.
 **/
export abstract class Dispatcher {
  /**
   * Dispatches a command to the appropriate handler.
   * @param {Command} command The command to dispatch
   * @returns {Promise<void>}
   */
  abstract dispatchCommand(command: Command): Promise<void>;
}
