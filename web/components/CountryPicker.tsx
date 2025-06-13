import React from "react";
import { Picker } from "@react-native-picker/picker";
import { ThemedPicker, ThemedPickerProps } from "@/components/ThemedPicker";

/**
 * Minimal data structure for countries, including an emoji flag.
 */
const COUNTRIES_WITH_EMOJIS = [
  { label: "United States", value: "US", flag: "🇺🇸" },
  { label: "Canada", value: "CA", flag: "🇨🇦" },
  { label: "Mexico", value: "MX", flag: "🇲🇽" },
  { label: "United Kingdom", value: "GB", flag: "🇬🇧" },
  { label: "Germany", value: "DE", flag: "🇩🇪" },
  { label: "France", value: "FR", flag: "🇫🇷" },
  { label: "Spain", value: "ES", flag: "🇪🇸" },
  { label: "Italy", value: "IT", flag: "🇮🇹" },
  { label: "Japan", value: "JP", flag: "🇯🇵" },
  { label: "Australia", value: "AU", flag: "🇦🇺" },
];

export type CountryPickerProps = Omit<ThemedPickerProps<string>, "children"> & {
  /**
   * Optional custom data for countries, if you want to override the default
   * with your own list of label/value/flag.
   */
  countries?: Array<{ label: string; value: string; flag: string }>;
};

export function CountryPicker({
  countries = COUNTRIES_WITH_EMOJIS,
  ...props
}: CountryPickerProps) {
  return (
    <ThemedPicker<string> {...props}>
      {countries.map(({ label, value, flag }) => (
        <Picker.Item
          key={value}
          value={value}
          // Combine the flag + label for the user-friendly text
          label={`${flag}  ${label}`}
        />
      ))}
    </ThemedPicker>
  );
}
