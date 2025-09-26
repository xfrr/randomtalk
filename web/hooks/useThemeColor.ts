/**
 * Hydration-safe theme color hook.
 * Learn more about light and dark modes:
 * https://docs.expo.dev/guides/color-schemes/
 */

import { appTheme } from "@/theme/app";
import { useHydratedTheme } from "./useHydratedTheme";

type ThemeProps = Partial<{ light: string; dark: string }>;

export function useThemeColor(
  props: ThemeProps,
  colorName: keyof typeof appTheme.light & keyof typeof appTheme.dark
) {
  const { theme } = useHydratedTheme();
  const colorFromProps = props[theme] ?? props.light ?? props.dark;
  return colorFromProps ?? appTheme[theme][colorName];
}
