import { Platform } from "react-native";

export const animationConfig = {
  useNativeDriver: Platform.OS !== "web",
};
