import { ThemedPicker, ThemedPickerProps } from "@/components/ThemedPicker";
import { Picker } from "@react-native-picker/picker";
import React from "react";
import { ThemedText } from "./ThemedText";
import { View } from "react-native";

/**
 * Minimal data structure for genders, including an emoji/icon.
 */
const GENDERS_WITH_EMOJIS = [
  { label: "Male", value: "male", icon: "👨" },
  { label: "Female", value: "female", icon: "👩" },
  { label: "Prefer not to say", value: "prefer_not_to_say", icon: "🤫" },
];

export type GenderPickerProps = Omit<ThemedPickerProps<string>, "children"> & {
  /**
   * Optional custom data for genders, if you want to override the default
   * with your own list of label/value/icon.
   */
  genders?: Array<{ label: string; value: string; icon: string }>;

  /**
   * Optional label for the picker.
   */
  label?: string;
};

export function GenderPicker({
  genders = GENDERS_WITH_EMOJIS,
  label = "Your Gender",
  ...props
}: GenderPickerProps) {
  return (
    <View style={{ marginTop: 16 }}>
      {label && (
        <ThemedText style={{ marginBottom: 4, fontWeight: "500" }}>
          {label}
        </ThemedText>
      )}
      <ThemedPicker<string> {...props}>
        {genders.map(({ label, value, icon }) => (
          <Picker.Item
            key={value}
            value={value}
            // Combine the icon + label for the user-friendly text
            label={`${icon}  ${label}`}
          />
        ))}
      </ThemedPicker>
    </View>
  );
}
