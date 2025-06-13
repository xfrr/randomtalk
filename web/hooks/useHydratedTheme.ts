import { useEffect, useState } from "react";
import { useColorScheme as useNativeColorScheme } from "react-native";

export function useHydratedTheme(initialGuess: "light" | "dark" = "light") {
  const colorScheme = useNativeColorScheme();
  const [isHydrated, setIsHydrated] = useState(!!colorScheme);
  const [theme, setTheme] = useState(colorScheme ?? initialGuess);

  useEffect(() => {
    if (colorScheme) {
      setTheme(colorScheme);
      setIsHydrated(true);
    }
  }, [colorScheme]);

  return { isHydrated, theme };
}
