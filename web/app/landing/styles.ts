// styles.ts
import { StyleSheet } from "react-native";

/**
 * Base styles for your component, used when
 * no media query overrides are triggered.
 */
export const baseStyles = StyleSheet.create({
  gradientContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
  },
  safeArea: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
  },
  container: {
    flex: 1,
    paddingHorizontal: 20,
    justifyContent: "center",
    maxWidth: "100%",
  },
  heroContainer: {
    marginBottom: 50,
    alignItems: "center",
    width: "100%",
  },
  heroTitle: {
    fontSize: 28,
    fontWeight: "700",
    color: "#fff",
    marginBottom: 8,
    textAlign: "center",
    textShadowColor: "rgba(0, 0, 0, 0.2)",
    textShadowOffset: { width: 1, height: 2 },
    textShadowRadius: 4,
  },
  heroSubtitle: {
    fontSize: 16,
    color: "#bbb",
    textAlign: "center",
  },

  // Form
  formContainer: {
    marginBottom: 40,
    paddingHorizontal: 20,
    width: "100%",
  },
  input: {
    backgroundColor: "#3E3760",
    color: "#fff",
    paddingVertical: 12,
    paddingHorizontal: 15,
    borderRadius: 10,
    width: "100%",
  },
  pickerContainer: {
    marginBottom: 15,
    width: "100%",
  },
  startButton: {
    backgroundColor: "#543AB7",
    paddingVertical: 15,
    borderRadius: 10,
    alignItems: "center",
    // Subtle shadow for emphasis
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.3,
    shadowRadius: 5,
    elevation: 5,
  },
  startButtonText: {
    color: "#fff",
    fontSize: 18,
    fontWeight: "600",
  },
});

/**
 * Returns style overrides based on device size booleans.
 * For instance:
 * - `isSmallScreen` (max-width: 400)
 * - `isTablet` (min-width: 768)
 */
export function getResponsiveOverrides(
  isSmallScreen: boolean,
  isTablet: boolean
) {
  switch (true) {
    case isSmallScreen:
      // If small devices (e.g. max-width: 400)
      return {
        ...baseStyles,
        heroTitle: {
          fontSize: 24,
        },
        heroSubtitle: {
          fontSize: 14,
        },
        input: {
          paddingVertical: 10,
          paddingHorizontal: 10,
        },
        startButtonText: {
          fontSize: 16,
        },
      } as any;
    case isTablet:
      // If tablet devices (e.g. min-width: 768)
      return {
        ...baseStyles,
        container: {
          paddingHorizontal: 40,
        },
        heroTitle: {
          fontSize: 32,
        },
        heroSubtitle: {
          fontSize: 18,
        },
        input: {
          paddingVertical: 14,
          paddingHorizontal: 18,
        },
        startButton: {
          paddingVertical: 18,
          borderRadius: 12,
        },
        startButtonText: {
          fontSize: 20,
        },
      } as any;
    default:
      return baseStyles;
  }
}

export default baseStyles;
