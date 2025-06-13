import CommandOverlay from "@/app/CommandOverlay";
import LoadingOverlay from "@/components/LoadingOverlay";
import { WebSocketProvider } from "@/internal/websocket";
import { useFonts } from "expo-font";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import { useEffect } from "react";
import StackLayout from "./stack";

import "react-native-reanimated";
import "@/locale/i18n";

// prevent the splash screen from auto-hiding before asset loading is complete.
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const [loaded] = useFonts({
    SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
  });

  useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  if (!loaded) {
    return null;
  }

  return (
    <>
      <WebSocketProvider
        config={{
          url: "ws://localhost:51000",
          maxRetries: 3,
          reconnectDelay: 1500,
          onOpen: () => console.log("Global WS opened"),
          onError: (err: any) => console.log("Global WS error:", err),
          onClose: () => console.log("Global WS closed"),
        }}
      >
        <LoadingOverlay>
          <CommandOverlay>
            <StackLayout />
          </CommandOverlay>
        </LoadingOverlay>
      </WebSocketProvider>
      <StatusBar style="auto" />
    </>
  );
}
