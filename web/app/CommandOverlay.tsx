import { Command, CommandDispatcherProvider } from "@/internal/command";
import { useWebSocketContext } from "@/internal/websocket";

export default function CommandOverlay({
  children,
}: {
  children: React.ReactNode;
}) {
  const { sendCommand } = useWebSocketContext();

  const commandHandlers = [
    {
      CommandName: "randomtalk.chat.create_chat_session",
      Do: async (cmd: Command) => {
        sendCommand(cmd);
      },
    },
  ];

  return (
    <CommandDispatcherProvider handlers={commandHandlers}>
      {children}
    </CommandDispatcherProvider>
  );
}
