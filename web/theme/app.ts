const tintColorLight = "#543AB7";
const tintColorDark = "#543AB7";

export const appTheme = {
  light: {
    text: "#1E1B29",
    background: "#FFFFFF",
    tint: tintColorLight,
    icon: "#3E3760",
    tabIconDefault: "#332D45",
    tabIconSelected: tintColorLight,

    // Extended Colors
    heading: "#2E2157",
    subheading: "#543AB7",
    secondaryBackground: "#F2F1F7",
    border: "#DDDDDD",

    // Typography
    fontSizeBase: 16,
    fontSizeLarge: 20,
    fontSizeSmall: 14,
    fontWeightRegular: "400",
    fontWeightBold: "700",

    // Buttons
    buttonBackground: tintColorLight,
    buttonText: "#FFFFFF",

    // Etc.
    placeholderText: "#AAA",
  },
  dark: {
    text: "#FFFFFF",
    background: "#2E2157",
    tint: tintColorDark,
    icon: "#BFBFBF",
    tabIconDefault: "#7A7A7A",
    tabIconSelected: tintColorDark,

    // Extended Colors
    heading: "#FFFFFF",
    subheading: "#BFBFBF",
    secondaryBackground: "#2B2540",
    border: "#3E3760",

    // Typography
    fontSizeBase: 16,
    fontSizeLarge: 20,
    fontSizeSmall: 14,
    fontWeightRegular: "400",
    fontWeightBold: "700",

    // Buttons
    buttonBackground: tintColorDark,
    buttonText: "#FFFFFF",

    // Etc.
    placeholderText: "#AAA",
  },
};
