import React, { ReactNode, useState, useMemo, useContext } from "react";
import { WizardContext } from "./WizardContext";
import WizardFooter from "./WizardFooter";
import WizardStep from "./WizardStep";

// --- Types ---
interface WizardProps {
  children: ReactNode;
  initialStep?: number;
  onStepChange?: (step: number) => void;
}

// --- Custom Hook with proper typing ---
export const useWizard = () => {
  const context = useContext(WizardContext);
  if (!context) {
    throw new Error("useWizard must be used within a WizardProvider");
  }
  return context;
};

// --- Main Component ---
export const Wizard: React.FC<WizardProps> = ({
  children,
  initialStep = 0,
  onStepChange,
}) => {
  const [currentStep, setCurrentStep] = useState(initialStep);

  // Split children into steps and footer
  const childArray = React.Children.toArray(children);

  const steps = childArray.filter(
    (child: any) => React.isValidElement(child) && child.type === WizardStep
  );

  const footer = childArray.find(
    (child: any) => React.isValidElement(child) && child.type === WizardFooter
  );

  const totalSteps = steps.length;

  const goToStep = (step: number) => {
    if (step >= 0 && step < totalSteps) {
      setCurrentStep(step);
      onStepChange?.(step);
    }
  };

  const nextStep = () => goToStep(currentStep + 1);
  const previousStep = () => goToStep(currentStep - 1);

  const contextValue = useMemo(
    () => ({
      currentStep,
      totalSteps,
      hasNextStep: currentStep < totalSteps - 1,
      hasPreviousStep: currentStep > 0,
      nextStep,
      previousStep,
      goToStep,
    }),
    [currentStep, totalSteps]
  );

  return (
    <WizardContext.Provider value={contextValue}>
      {/* Render current step */}
      {steps[currentStep]}

      {/* Render footer (always) */}
      {footer}
    </WizardContext.Provider>
  );
};
