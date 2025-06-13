import { ThemedText } from "@/components/ThemedText";
import "@expo/match-media";
import React from "react";
import {
  GestureResponderEvent,
  NativeSyntheticEvent,
  StyleProp,
  TargetedEvent,
  TouchableOpacity,
  View,
} from "react-native";

interface ThemedButtonProps {
  styles?: {
    touchableOpacity?: StyleProp<any>;
    text?: StyleProp<any>;
  };
  text?: string;
  icon?: React.ReactNode;
  iconPosition?: "left" | "right";
  iconSize?: number;
  iconColor?: string;
  iconStyle?: StyleProp<any>;
  accessibilityLabel?: string;
  accessibilityHint?: string;

  onPress?: (event: GestureResponderEvent) => void;
  onLongPress?: (event: GestureResponderEvent) => void;
  onFocus?: (event: NativeSyntheticEvent<TargetedEvent>) => void;
  onPressIn?: (event: GestureResponderEvent) => void;
  onPressOut?: (event: GestureResponderEvent) => void;
}

export default function ThemedButton(props: ThemedButtonProps) {
  const renderText = (text: string, style: StyleProp<any>) => {
    return <ThemedText style={style}>{text}</ThemedText>;
  };

  const renderIcon = (
    icon: React.ReactNode,
    position: "left" | "right",
    size?: number,
    color?: string,
    style?: StyleProp<any>
  ) => {
    return (
      <View
        style={[
          props.styles?.touchableOpacity,
          position === "left" ? { marginRight: 8 } : { marginLeft: 8 },
        ]}
      >
        {React.cloneElement(icon as React.ReactElement, {
          size: size,
          color: color,
          style: style,
        })}
      </View>
    );
  };

  const renderLeftIcon = (
    icon: React.ReactNode,
    size?: number,
    color?: string,
    style?: StyleProp<any>
  ) => {
    return renderIcon(icon, "left", size, color, style);
  };

  const renderRightIcon = (
    icon: React.ReactNode,
    size?: number,
    color?: string,
    style?: StyleProp<any>
  ) => {
    return renderIcon(icon, "right", size, color, style);
  };

  return (
    <TouchableOpacity
      style={props.styles?.touchableOpacity}
      onPress={props.onPress}
      onLongPress={props.onLongPress}
      onFocus={props.onFocus}
      onPressIn={props.onPressIn}
      onPressOut={props.onPressOut}
      accessibilityLabel={props.accessibilityLabel ?? props.text}
      accessibilityHint={props.accessibilityHint ?? props.text}
      accessibilityRole="button"
    >
      {props.icon &&
        props.iconPosition === "left" &&
        renderLeftIcon(
          props.icon,
          props.iconSize,
          props.iconColor,
          props.iconStyle
        )}
      {props.text && renderText(props.text, props.styles?.text)}
      {props.icon &&
        props.iconPosition === "right" &&
        renderRightIcon(
          props.icon,
          props.iconSize,
          props.iconColor,
          props.iconStyle
        )}
    </TouchableOpacity>
  );
}
