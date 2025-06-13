import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import * as Localization from "react-native-localize";
import en from "./locales/en.json";
import es from "./locales/es.json";
import { Platform } from "react-native";

// Detect user language
const getDeviceLanguage = () => {
  if (Platform.OS === "web") {
    if (typeof navigator !== "undefined" && navigator.language) {
      return navigator.language;
    }
    // Fallback if still not available
    return "en";
  } else {
    const locales = Localization.getLocales();
    if (Array.isArray(locales) && locales.length > 0) {
      return locales[0].languageTag;
    }
    return "en";
  }
};

i18n.use(initReactI18next).init({
  lng: getDeviceLanguage(),
  fallbackLng: "en", // if user locale is not available
  debug: true, // disable in production
  resources: {
    en: { translation: en },
    es: { translation: es },
  },
  interpolation: {
    escapeValue: false, // react already safes from xss
  },
});

export default i18n;
