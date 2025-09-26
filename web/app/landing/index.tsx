import "@expo/match-media";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Dimensions, Platform } from "react-native";
import { v4 as uuidv4 } from "uuid";

import { LinearGradient } from "expo-linear-gradient";
import {
  KeyboardAvoidingView,
  SafeAreaView,
  StatusBar,
  View,
} from "react-native";

import { LoadingContext } from "@/components/LoadingOverlay";
import { ThemedText } from "@/components/ThemedText";
import { Wizard, WizardFooter, WizardStep } from "@/components/wizard";
import { useCommandDispatcher } from "@/internal/command";
import { Gender } from "@/utils/genders";
import { useTranslation } from "react-i18next";
import ExtendedForm from "./ExtendedForm";
import { FormButtons } from "./FormButtons";
import MainForm from "./MainForm";
import { baseStyles, getResponsiveOverrides } from "./styles";

// Debounce helper
function debounce(fn: () => void, delay: number) {
  let timer: ReturnType<typeof setTimeout>;
  return () => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      fn();
    }, delay);
  };
}

export default function LandingView() {
  const { setIsLoading } = useContext(LoadingContext);
  const { dispatchCommand } = useCommandDispatcher();
  const { t } = useTranslation();

  const [formState, setFormState] = useState({
    username: "",
    interests: "",
    country: "ES",
    birthdate: (() => {
      const date = new Date();
      date.setFullYear(date.getFullYear() - 18);
      return date;
    })(),
    gender: Gender.Unspecified,
  });

  // SSR-safe initial width
  const getInitialWidth = () => {
    if (Platform.OS !== "web") {
      return Dimensions.get("window").width;
    }

    // For web, try to get width safely
    if (typeof window !== "undefined") {
      return window.innerWidth;
    }

    // Fallback
    return 768;
  };

  const [screenWidth, setScreenWidth] = useState(getInitialWidth);
  const [isReady, setIsReady] = useState(false);

  const updateScreenWidth = useCallback(
    debounce(() => {
      const width =
        Platform.OS === "web"
          ? window.innerWidth
          : Dimensions.get("window").width;
      setScreenWidth(width);
    }, 150),
    []
  );

  // Set up resize listener
  useEffect(() => {
    setIsReady(true);
    if (Platform.OS === "web") {
      window.addEventListener("resize", updateScreenWidth);
      return () => window.removeEventListener("resize", updateScreenWidth);
    } else {
      const subscription = Dimensions.addEventListener("change", ({ window }) =>
        setScreenWidth(window.width)
      );
      return () => subscription.remove();
    }
  }, [updateScreenWidth]);

  const isSmallScreen = screenWidth <= 400;
  const isTablet = screenWidth >= 768;

  const responsiveStyle = useMemo(
    () => getResponsiveOverrides(isSmallScreen, isTablet),
    [isSmallScreen, isTablet]
  );

  const styles = useMemo(() => {
    return {
      ...baseStyles,
      container: [baseStyles.container, responsiveStyle.container],
      heroTitle: [baseStyles.heroTitle, responsiveStyle.heroTitle],
      heroSubtitle: [baseStyles.heroSubtitle, responsiveStyle.heroSubtitle],
      input: [baseStyles.input, responsiveStyle.input],
      startButton: [baseStyles.startButton, responsiveStyle.startButton],
      startButtonText: [
        baseStyles.startButtonText,
        responsiveStyle.startButtonText,
      ],
      safeArea: responsiveStyle.safeArea || baseStyles.safeArea,
      gradientContainer:
        responsiveStyle.gradientContainer || baseStyles.gradientContainer,
      heroContainer: responsiveStyle.heroContainer || baseStyles.heroContainer,
    };
  }, [responsiveStyle]);

  if (!isReady) {
    return null;
  }

  const handleComplete = () => {
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
    }, 1000);

    const { username, interests, gender, birthdate } = formState;

    const calculateAge = (birthdate: Date) => {
      const today = new Date();
      let age = today.getFullYear() - birthdate.getFullYear();
      const monthDiff = today.getMonth() - birthdate.getMonth();
      if (
        monthDiff < 0 ||
        (monthDiff === 0 && today.getDate() < birthdate.getDate())
      ) {
        age--;
      }
      return age;
    };

    dispatchCommand({
      type: "randomtalk.chat.create_chat_session",
      payload: {
        user_id: uuidv4(),
        user_nickname: username,
        user_interests: interests,
        user_age: calculateAge(birthdate),
        user_gender: gender,
        user_match_preference_min_age: 18,
        user_match_preference_max_age: 35,
        user_match_preference_interests: Array.from(interests.split(",")),
      },
      timestamp: new Date(),
    });
  };

  return (
    <LinearGradient
      colors={["#2E2157", "#1E1B29"]}
      style={styles.gradientContainer}
    >
      <SafeAreaView style={styles.safeArea}>
        <KeyboardAvoidingView
          style={styles.container}
          behavior={Platform.OS === "ios" ? "padding" : "height"}
        >
          <StatusBar barStyle="light-content" />

          <View style={styles.heroContainer}>
            <ThemedText style={styles.heroTitle}>
              {t("landing.title")}
            </ThemedText>
            <ThemedText style={styles.heroSubtitle}>
              {t("landing.subtitle")}
            </ThemedText>
          </View>

          <Wizard>
            <WizardStep>
              <MainForm
                state={formState}
                actions={{
                  setUsername: (username: string) =>
                    setFormState((prev) => ({ ...prev, username })),
                  setInterests: (interests: string) =>
                    setFormState((prev) => ({ ...prev, interests })),
                  setCountry: (country: string) =>
                    setFormState((prev) => ({ ...prev, country })),
                }}
                styles={styles}
              />
            </WizardStep>
            <WizardStep>
              <ExtendedForm
                state={formState}
                actions={{
                  setBirthdate: (birthdate: Date) =>
                    setFormState((prev) => ({ ...prev, birthdate })),
                  setGender: (gender: Gender) =>
                    setFormState((prev) => ({ ...prev, gender })),
                }}
                styles={styles}
              />
            </WizardStep>
            <WizardFooter>
              <FormButtons styles={styles} onComplete={handleComplete} />
            </WizardFooter>
          </Wizard>
        </KeyboardAvoidingView>
      </SafeAreaView>
    </LinearGradient>
  );
}
