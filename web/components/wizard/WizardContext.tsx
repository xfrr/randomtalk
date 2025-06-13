import React from "react";

import "@expo/match-media";

// WizardContextType defines the shape of the context that will be provided
// to the Wizard component and its children. It includes functions for
// navigating between steps, the current step index, total steps, and
// booleans indicating if there are next or previous steps available.
export interface WizardContextType {
  nextStep: () => void;
  previousStep: () => void;
  currentStep: number;
  totalSteps: number;
  hasNextStep: boolean;
  hasPreviousStep: boolean;
}

// WizardContext provides the current step and navigation functions to its children.
// It uses React's Context API to allow any component in the tree to access
// the current step and navigation functions without having to pass them down
// through props.
export const WizardContext = React.createContext<WizardContextType>({
  nextStep: () => {
    throw new Error("nextStep must be used within a WizardProvider");
  },
  previousStep: () => {
    throw new Error("previousStep must be used within a WizardProvider");
  },
  currentStep: 0,
  totalSteps: 0,
  hasNextStep: false,
  hasPreviousStep: false,
});
