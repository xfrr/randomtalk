import { createContext, useContext, useMemo } from "react";
import { Command } from "./Command";
import { Dispatcher } from "./Dispatcher";
import { InternalMessageBus } from "../messaging";

/**
 * CommandDispatcherContext
 * @type {React.Context<Dispatcher<Command>>}
 * @description React Context for Command Dispatcher
 * @export
 * @constant
 * @default
 * @memberof CommandDispatcherContext
 * @example
 * const { dispatchCommand } = useContext(CommandDispatcherContext);
 * dispatchCommand(MyCommand);
 * @see CommandDispatcherProvider
 */
export const CommandDispatcherContext: React.Context<Dispatcher | null> =
  createContext<Dispatcher | null>(null);

/**
 * useCommandDispatcher
 * @function
 * @description React Hook for using the Command Dispatcher
 * @returns {Dispatcher<Command>}
 * @export
 * @function
 * @memberof CommandDispatcherContext
 * @example
 * const { dispatchCommand } = useCommandDispatcher();
 * dispatchCommand(MyCommand);
 * @see CommandDispatcherProvider
 */
export const useCommandDispatcher = (): Dispatcher => {
  const dispatcher = useContext(CommandDispatcherContext);

  if (!dispatcher) {
    throw new Error(
      "useCommandDispatcher must be used within a CommandDispatcherProvider"
    );
  }

  return dispatcher;
};

/**
 * CommandDispatcherProviderProps
 * @interface
 * @description Command Dispatcher Provider Props
 * @export
 * @memberof CommandDispatcherProvider
 * @example
 * <CommandDispatcherProvider>
 *  <App />
 * </CommandDispatcherProvider>
 * @see CommandDispatcherProvider
 * @see CommandDispatcherContext
 * @see CommandDispatcher
 * @see Command
 */
export interface CommandDispatcherProviderProps {
  /**
   * The children
   * @type {React.ReactNode}
   */
  children: React.ReactNode;
  /**
   * An optional list of command handlers to register
   * @type {{ messageName: string; handlerFn: (command: Command) => Promise<void>; }[]}
   */
  handlers?: {
    CommandName: string;
    Do: (command: Command) => Promise<void>;
  }[];
}

/**
 * CommandDispatcherProvider
 * @param {CommandDispatcherProviderProps} props The props
 * @returns {React.ReactElement}
 * @description Command Dispatcher Provider
 * @export
 * @function
 * @memberof CommandDispatcherProvider
 * @example
 * <CommandDispatcherProvider>
 *  <App />
 * </CommandDispatcherProvider>
 * @see CommandDispatcherProviderProps
 * @see CommandDispatcherContext
 * @see CommandDispatcher
 * @see Command
 */
export const CommandDispatcherProvider: React.FC<
  CommandDispatcherProviderProps
> = ({
  children,
  handlers,
}: CommandDispatcherProviderProps): React.ReactElement => {
  // memoize the instance to avoid re-creating on every render.
  const bus = useMemo(() => new InternalMessageBus(), []);

  if (handlers && handlers.length > 0) {
    handlers.forEach((handler) => {
      bus.registerHandler(handler.CommandName, (msg) =>
        handler.Do(msg as Command)
      );
    });
  }

  const dispatchCommand = (command: Command): Promise<void> => {
    return bus.send(command);
  };

  return (
    <CommandDispatcherContext.Provider
      value={{
        dispatchCommand,
      }}
    >
      {children}
    </CommandDispatcherContext.Provider>
  );
};

export default CommandDispatcherProvider;
