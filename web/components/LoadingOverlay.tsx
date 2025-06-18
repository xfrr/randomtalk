import { animationConfig } from "@/utils/animationConfig";
import React, { createContext, useEffect, useRef, useState } from "react";
import { Animated, Easing, StyleSheet, View } from "react-native";

// Create a context to manage loading state
type LoadingContextType = {
  isLoading: boolean;
  setIsLoading: React.Dispatch<React.SetStateAction<boolean>>;
};

export const LoadingContext = createContext<LoadingContextType>({
  isLoading: false,
  setIsLoading: () => {},
});

// Provider that wraps your app and conditionally renders the loading overlay
export const LoadingOverlay = ({ children }: { children: React.ReactNode }) => {
  const [isLoading, setIsLoading] = useState(false);

  return (
    <LoadingContext.Provider value={{ isLoading, setIsLoading }}>
      {children}
      {isLoading && <LoadingSpinner />}
    </LoadingContext.Provider>
  );
};

// Animated loading spinner overlay component
const LoadingSpinner = () => {
  const rotation = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    // Continuous rotation animation
    Animated.loop(
      Animated.timing(rotation, {
        toValue: 1,
        duration: 1000,
        easing: Easing.linear,
        ...animationConfig,
      })
    ).start();
  }, [rotation]);

  // Interpolate rotation value to degrees
  const rotate = rotation.interpolate({
    inputRange: [0, 1],
    outputRange: ["0deg", "360deg"],
  });

  return (
    <View
      style={styles.loadingOverlay}
      accessible
      accessibilityLabel="Loading"
      accessibilityRole="alert"
    >
      <Animated.View style={[styles.spinner, { transform: [{ rotate }] }]} />
    </View>
  );
};

// Essential styling for a modern, intuitive UI
const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "#f5f5f5",
  },
  title: {
    fontSize: 24,
    marginBottom: 20,
    color: "#333",
  },
  button: {
    backgroundColor: "#6200ee",
    paddingVertical: 12,
    paddingHorizontal: 20,
    borderRadius: 8,
  },
  buttonText: {
    color: "#fff",
    fontSize: 16,
  },
  loadingOverlay: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: "rgba(0,0,0,0.3)",
    alignItems: "center",
    justifyContent: "center",
  },
  spinner: {
    width: 60,
    height: 60,
    borderWidth: 6,
    borderColor: "#fff",
    borderTopColor: "#6200ee",
    borderRadius: 30,
  },
});

export default LoadingOverlay;
