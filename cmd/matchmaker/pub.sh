#!/bin/bash

# This script is used to test publish a chat notification to NATS Stream.
# It will publish a chat notification to the NATS Stream.

# Usage:
#   ./pub.sh <user-id> <age> <[]interests]>

# Example:
#   ./pub.sh 1 20

# Check if the number of arguments is correct
if [ "$#" -ne 2 ]; then
    echo "Usage: ./pub.sh <user-id> <age>"
    exit 1
fi

# Get the user id
USER_ID=$1

# Publish a chat notification with matching preferences
nats pub randomtalk.notifications.chat.users."${USER_ID}".connected '{"id": "'"${USER_ID}"'","type": "chat_session_started","specversion": "1.0","source": "randomtalk.chat","subject": "1","time": "2020-08-01T00:00:00Z","data": {"user_id": "'"${USER_ID}"'","user_age": 20, "user_gender": "GENDER_MALE", "user_preferences": {"min_age": 20, "max_age": 30, "interests": ["music", "sports"]}}}'

# Publish a chat notification to match user 1
# nats pub randomtalk.chat.notifications.sessions.1.started '{"id": "2","type": "chat_session_started","specversion": "1.0","source": "randomtalk.chat","subject": "2","time": "2020-08-01T00:00:00Z","data": {"user_id": "user-2","user_age": 20, "user_gender": "FEMALE_GENDER", "user_preferences": {"min_age": 18, "max_age": 20, "interests": ["music", "sports"]}}}'

# # Publish a chat notification to NOT match user 1 nor user 2
# nats pub randomtalk.chat.notifications.sessions.1.started '{"id": "3","type": "chat_session_started","specversion": "1.0","source": "randomtalk.chat","subject": "3","time": "2020-08-01T00:00:00Z","data": {"user_id": "user-3","user_age": 35, "user_preferences": {"min_age": 30, "max_age": 40, "interests": ["gaming", "movies"]}}}'

# # Publish a chat notification to match user 3
# nats pub randomtalk.chat.notifications.sessions.1.started '{"id": "4","type": "chat_session_started","specversion": "1.0","source": "randomtalk.chat","subject": "4","time": "2020-08-01T00:00:00Z","data": {"user_id": "user-4","user_age": 35, "user_preferences": {"min_age": 30, "max_age": 40, "interests": ["gaming", "movies"]}}}'
