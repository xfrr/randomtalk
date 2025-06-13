import { appTheme } from "@/theme/app";
import { LinearGradient } from "expo-linear-gradient";
import { Stack } from "expo-router";
import React, { useEffect, useState } from "react";
import {
  FlatList,
  KeyboardAvoidingView,
  Platform,
  StatusBar,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";
import styles from "../styles";

export default function Index() {
  const [messages, setMessages] = useState<
    Array<{
      id: string;
      text: string;
      sender: string;
      timestamp: string;
    }>
  >([]);
  const [input, setInput] = useState("");

  useEffect(() => {
    // const ws = new WebSocket('wss://your-realtime-endpoint');
    // ws.onmessage = (event) => {
    //   const newMessage = JSON.parse(event.data);
    //   setMessages((prev) => [...prev, newMessage]);
    // };
    // return () => {
    //   ws.close();
    // };
  }, []);

  const handleSend = () => {
    if (input.trim()) {
      const newMessage = {
        id: Date.now().toString(),
        text: input.trim(),
        sender: "me", // or user ID/username
        timestamp: new Date().toLocaleTimeString(),
      };
      // In real app, send via WebSocket:
      // ws.send(JSON.stringify(newMessage));
      setMessages((prev) => [...prev, newMessage]);
      setInput("");
    }
  };

  const renderMessage = ({ item }: { item: any }) => {
    const isMyMessage = item.sender === "me";
    return (
      <View
        style={[
          styles.messageContainer,
          isMyMessage ? styles.myMessage : styles.otherMessage,
        ]}
      >
        <Text style={styles.messageText}>{item.text}</Text>
        <Text style={styles.timestamp}>{item.timestamp}</Text>
      </View>
    );
  };

  return (
    <>
      <Stack.Screen
        options={{
          title: "Chat",
          headerShown: true,
          headerStyle: { backgroundColor: appTheme.dark.background },
          headerTintColor: appTheme.dark.text,
          headerTitleStyle: { fontWeight: "bold" },
          headerBackVisible: false,
        }}
      />
      <LinearGradient
        colors={[appTheme.dark.background, "#1E1B29"]}
        style={styles.gradientContainer}
      >
        <KeyboardAvoidingView
          style={styles.container}
          behavior={Platform.OS === "ios" ? "padding" : "height"}
        >
          {/* Status Bar for mobile */}
          <StatusBar barStyle="light-content" />

          {/* Header */}
          {/* <View style={styles.headerContainer}>
          <Text style={styles.title}>RandomTalk</Text>
          <Text style={styles.subTitle}>
            Meet random users from anywhere, instantly!
          </Text>
        </View> */}

          {/* Chat List */}
          <FlatList
            data={messages}
            renderItem={renderMessage}
            keyExtractor={(item) => item.id}
            style={styles.chatContainer}
            contentContainerStyle={{ paddingBottom: 20 }}
          />

          {/* Message Input */}
          <View style={styles.inputContainer}>
            <TextInput
              style={styles.input}
              onChangeText={setInput}
              value={input}
              placeholder="Type a message..."
              placeholderTextColor="#AAA"
              accessible
              accessibilityLabel="Message input"
              accessibilityHint="Type your message here"
            />
            <TouchableOpacity
              style={styles.sendButton}
              onPress={handleSend}
              accessible
              accessibilityLabel="Send message"
              accessibilityHint="Sends the typed message"
            >
              <Text style={styles.sendButtonText}>Send</Text>
            </TouchableOpacity>
          </View>
        </KeyboardAvoidingView>
      </LinearGradient>
    </>
  );
}
