// utils/responsive.ts
import { Dimensions, PixelRatio } from "react-native";

const { width: SCREEN_WIDTH } = Dimensions.get("window");

// iPhone 11 Pro reference width for scaling. Adjust as desired.
const REFERENCE_WIDTH = 375;

/**
 * Scale a given size (e.g., font size, spacing) relative to the reference width.
 */
export const scaleSize = (size: number) => {
  const scale = SCREEN_WIDTH / REFERENCE_WIDTH;
  const newSize = size * scale;
  return Math.round(PixelRatio.roundToNearestPixel(newSize));
}