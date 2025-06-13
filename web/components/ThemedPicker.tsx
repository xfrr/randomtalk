import React, { useEffect, useState } from "react";
import { StyleSheet, View, StyleProp, ViewStyle } from "react-native";
import { Picker, PickerProps } from "@react-native-picker/picker";
import { useThemeColor } from "@/hooks/useThemeColor";
import { useHydratedTheme } from "@/hooks/useHydratedTheme";
import { ActivityIndicator } from "react-native-paper";

type PickerStyleType = "default" | "underlined" | "rounded";

export interface ThemedPickerProps<T> extends Omit<PickerProps<T>, "style"> {
  lightColor?: string;
  darkColor?: string;
  type?: PickerStyleType;
  containerStyle?: StyleProp<ViewStyle>;
}

export function ThemedPicker<T>({
  lightColor,
  darkColor,
  type = "default",
  containerStyle,
  ...props
}: ThemedPickerProps<T>) {
  const { isHydrated, theme } = useHydratedTheme();

  const textColor = useThemeColor(
    { light: lightColor, dark: darkColor },
    "text"
  );

  const backgroundColor = useThemeColor({}, "background");

  if (!isHydrated) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <View style={[styles.container, containerStyle]}>
      <Picker
        {...props}
        style={[
          { color: textColor, backgroundColor },
          type === "default" && styles.default,
          type === "underlined" && styles.underlined,
          type === "rounded" && styles.rounded,
        ]}
        dropdownIconColor={textColor} // Android-specific icon color
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginVertical: 8,
  },
  default: {
    fontSize: 16,
    lineHeight: 24,
    paddingHorizontal: 6,
    paddingVertical: 8,
    borderRadius: 4,
  },
  underlined: {
    fontSize: 16,
    borderBottomWidth: 1,
    borderBottomColor: "#999",
    paddingHorizontal: 12,
    paddingVertical: 6,
  },
  rounded: {
    fontSize: 16,
    lineHeight: 24,
    paddingHorizontal: 6,
    paddingVertical: 8,
    borderRadius: 10,
  },
});
