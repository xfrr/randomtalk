// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: randomtalk/matchmaking/v1/matchmaking_service.proto

package matchpb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Gender int32

const (
	Gender_GENDER_UNSPECIFIED Gender = 0
	Gender_GENDER_MALE        Gender = 1
	Gender_GENDER_FEMALE      Gender = 2
)

// Enum value maps for Gender.
var (
	Gender_name = map[int32]string{
		0: "GENDER_UNSPECIFIED",
		1: "GENDER_MALE",
		2: "GENDER_FEMALE",
	}
	Gender_value = map[string]int32{
		"GENDER_UNSPECIFIED": 0,
		"GENDER_MALE":        1,
		"GENDER_FEMALE":      2,
	}
)

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_enumTypes[0].Descriptor()
}

func (Gender) Type() protoreflect.EnumType {
	return &file_randomtalk_matchmaking_v1_matchmaking_service_proto_enumTypes[0]
}

func (x Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Gender.Descriptor instead.
func (Gender) EnumDescriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{0}
}

type FindMatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId           string            `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserName         string            `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	UserAge          int32             `protobuf:"varint,3,opt,name=user_age,json=userAge,proto3" json:"user_age,omitempty"`
	UserGender       Gender            `protobuf:"varint,4,opt,name=user_gender,json=userGender,proto3,enum=randomtalk.matchmaking.v1.Gender" json:"user_gender,omitempty"`
	UserLocation     *LatLng           `protobuf:"bytes,5,opt,name=user_location,json=userLocation,proto3" json:"user_location,omitempty"`
	MatchPreferences *MatchPreferences `protobuf:"bytes,6,opt,name=match_preferences,json=matchPreferences,proto3" json:"match_preferences,omitempty"`
}

func (x *FindMatchRequest) Reset() {
	*x = FindMatchRequest{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindMatchRequest) ProtoMessage() {}

func (x *FindMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindMatchRequest.ProtoReflect.Descriptor instead.
func (*FindMatchRequest) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{0}
}

func (x *FindMatchRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *FindMatchRequest) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *FindMatchRequest) GetUserAge() int32 {
	if x != nil {
		return x.UserAge
	}
	return 0
}

func (x *FindMatchRequest) GetUserGender() Gender {
	if x != nil {
		return x.UserGender
	}
	return Gender_GENDER_UNSPECIFIED
}

func (x *FindMatchRequest) GetUserLocation() *LatLng {
	if x != nil {
		return x.UserLocation
	}
	return nil
}

func (x *FindMatchRequest) GetMatchPreferences() *MatchPreferences {
	if x != nil {
		return x.MatchPreferences
	}
	return nil
}

type FindMatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MatchId string `protobuf:"bytes,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
}

func (x *FindMatchResponse) Reset() {
	*x = FindMatchResponse{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindMatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindMatchResponse) ProtoMessage() {}

func (x *FindMatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindMatchResponse.ProtoReflect.Descriptor instead.
func (*FindMatchResponse) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{1}
}

func (x *FindMatchResponse) GetMatchId() string {
	if x != nil {
		return x.MatchId
	}
	return ""
}

type GetMatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MatchId string `protobuf:"bytes,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
}

func (x *GetMatchRequest) Reset() {
	*x = GetMatchRequest{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMatchRequest) ProtoMessage() {}

func (x *GetMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMatchRequest.ProtoReflect.Descriptor instead.
func (*GetMatchRequest) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetMatchRequest) GetMatchId() string {
	if x != nil {
		return x.MatchId
	}
	return ""
}

type GetMatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Match *Match `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
}

func (x *GetMatchResponse) Reset() {
	*x = GetMatchResponse{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMatchResponse) ProtoMessage() {}

func (x *GetMatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMatchResponse.ProtoReflect.Descriptor instead.
func (*GetMatchResponse) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetMatchResponse) GetMatch() *Match {
	if x != nil {
		return x.Match
	}
	return nil
}

type MatchPreferences struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gender             Gender   `protobuf:"varint,1,opt,name=gender,proto3,enum=randomtalk.matchmaking.v1.Gender" json:"gender,omitempty"`
	MinAge             int32    `protobuf:"varint,2,opt,name=min_age,json=minAge,proto3" json:"min_age,omitempty"`
	MaxAge             int32    `protobuf:"varint,3,opt,name=max_age,json=maxAge,proto3" json:"max_age,omitempty"`
	MaxDistanceKm      float64  `protobuf:"fixed64,4,opt,name=max_distance_km,json=maxDistanceKm,proto3" json:"max_distance_km,omitempty"`
	Interests          []string `protobuf:"bytes,5,rep,name=interests,proto3" json:"interests,omitempty"`
	MaxWaitTimeSeconds int32    `protobuf:"varint,6,opt,name=max_wait_time_seconds,json=maxWaitTimeSeconds,proto3" json:"max_wait_time_seconds,omitempty"`
}

func (x *MatchPreferences) Reset() {
	*x = MatchPreferences{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MatchPreferences) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchPreferences) ProtoMessage() {}

func (x *MatchPreferences) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchPreferences.ProtoReflect.Descriptor instead.
func (*MatchPreferences) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{4}
}

func (x *MatchPreferences) GetGender() Gender {
	if x != nil {
		return x.Gender
	}
	return Gender_GENDER_UNSPECIFIED
}

func (x *MatchPreferences) GetMinAge() int32 {
	if x != nil {
		return x.MinAge
	}
	return 0
}

func (x *MatchPreferences) GetMaxAge() int32 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

func (x *MatchPreferences) GetMaxDistanceKm() float64 {
	if x != nil {
		return x.MaxDistanceKm
	}
	return 0
}

func (x *MatchPreferences) GetInterests() []string {
	if x != nil {
		return x.Interests
	}
	return nil
}

func (x *MatchPreferences) GetMaxWaitTimeSeconds() int32 {
	if x != nil {
		return x.MaxWaitTimeSeconds
	}
	return 0
}

type LatLng struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The latitude in degrees. It must be in the range [-90.0, +90.0].
	Latitude float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	// The longitude in degrees. It must be in the range [-180.0, +180.0].
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *LatLng) Reset() {
	*x = LatLng{}
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LatLng) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatLng) ProtoMessage() {}

func (x *LatLng) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatLng.ProtoReflect.Descriptor instead.
func (*LatLng) Descriptor() ([]byte, []int) {
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP(), []int{5}
}

func (x *LatLng) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *LatLng) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

var File_randomtalk_matchmaking_v1_matchmaking_service_proto protoreflect.FileDescriptor

var file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDesc = []byte{
	0x0a, 0x33, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c,
	0x6b, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x25, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc9, 0x02, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x12, 0x42, 0x0a, 0x0b,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x21, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x12, 0x46, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x74, 0x4c, 0x6e, 0x67, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x58, 0x0a, 0x11, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b,
	0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73,
	0x52, 0x10, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x73, 0x22, 0x2e, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x49, 0x64, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64,
	0x22, 0x4a, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b,
	0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x22, 0xf8, 0x01, 0x0a,
	0x10, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x12, 0x39, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x21, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07,
	0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6d,
	0x69, 0x6e, 0x41, 0x67, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x41, 0x67, 0x65, 0x12, 0x26,
	0x0a, 0x0f, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6b,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x6d, 0x61, 0x78, 0x44, 0x69, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x4b, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x65,
	0x73, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x65, 0x73, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x15, 0x6d, 0x61, 0x78, 0x5f, 0x77, 0x61, 0x69, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x57, 0x61, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x42, 0x0a, 0x06, 0x4c, 0x61, 0x74, 0x4c, 0x6e,
	0x67, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x2a, 0x44, 0x0a, 0x06, 0x47,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x12, 0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f, 0x4d, 0x41, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f, 0x46, 0x45, 0x4d, 0x41, 0x4c, 0x45, 0x10,
	0x02, 0x32, 0xa5, 0x02, 0x0a, 0x12, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x61, 0x6b, 0x69, 0x6e,
	0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x09, 0x46, 0x69, 0x6e,
	0x64, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x2b, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74,
	0x61, 0x6c, 0x6b, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b,
	0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1b, 0x92, 0x41, 0x02, 0x62, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x3a, 0x01,
	0x2a, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x12, 0x88,
	0x01, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x2a, 0x2e, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61,
	0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b, 0x69, 0x6e, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x92, 0x41, 0x02, 0x62, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x18, 0x12, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x2f, 0x7b,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x7d, 0x42, 0xc4, 0x04, 0x92, 0x41, 0x93, 0x04,
	0x12, 0x85, 0x01, 0x0a, 0x0e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x20, 0x20, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x22, 0x2b, 0x0a, 0x04, 0x78, 0x66, 0x72, 0x72, 0x12, 0x12, 0x68, 0x74, 0x74,
	0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x66, 0x72, 0x6f, 0x6d, 0x65, 0x72, 0x6f, 0x2e, 0x6d, 0x65, 0x1a,
	0x0f, 0x77, 0x6f, 0x72, 0x6b, 0x40, 0x66, 0x72, 0x6f, 0x6d, 0x65, 0x72, 0x6f, 0x2e, 0x6d, 0x65,
	0x2a, 0x42, 0x0a, 0x0a, 0x41, 0x70, 0x61, 0x63, 0x68, 0x65, 0x20, 0x32, 0x2e, 0x30, 0x12, 0x34,
	0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x78, 0x66, 0x72, 0x72, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61,
	0x6c, 0x6b, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x4c, 0x49, 0x43,
	0x45, 0x4e, 0x53, 0x45, 0x32, 0x02, 0x76, 0x31, 0x1a, 0x0f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68,
	0x6f, 0x73, 0x74, 0x3a, 0x35, 0x30, 0x30, 0x30, 0x30, 0x2a, 0x03, 0x01, 0x02, 0x04, 0x32, 0x10,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e,
	0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73,
	0x6f, 0x6e, 0x52, 0x55, 0x0a, 0x03, 0x34, 0x30, 0x33, 0x12, 0x4e, 0x0a, 0x4c, 0x52, 0x65, 0x74,
	0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77, 0x68, 0x65, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72, 0x20, 0x64, 0x6f, 0x65, 0x73, 0x20, 0x6e, 0x6f,
	0x74, 0x20, 0x68, 0x61, 0x76, 0x65, 0x20, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x20, 0x74, 0x6f, 0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x3b, 0x0a, 0x03, 0x34, 0x30, 0x34,
	0x12, 0x34, 0x0a, 0x2a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77, 0x68, 0x65,
	0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x20, 0x64,
	0x6f, 0x65, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x2e, 0x12, 0x06,
	0x0a, 0x04, 0x9a, 0x02, 0x01, 0x07, 0x52, 0x37, 0x0a, 0x03, 0x35, 0x30, 0x30, 0x12, 0x30, 0x0a,
	0x2e, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77, 0x68, 0x65, 0x6e, 0x20, 0x61,
	0x6e, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x6f, 0x63, 0x63, 0x75, 0x72, 0x73, 0x2e, 0x5a,
	0x74, 0x0a, 0x72, 0x0a, 0x06, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x12, 0x68, 0x08, 0x03, 0x28,
	0x04, 0x32, 0x23, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x3a, 0x1f, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x61, 0x75, 0x74,
	0x68, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x1c, 0x0a, 0x1a, 0x0a, 0x04, 0x72, 0x65, 0x61,
	0x64, 0x12, 0x12, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x73, 0x20, 0x72, 0x65, 0x61, 0x64, 0x20, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x62, 0x0c, 0x0a, 0x0a, 0x0a, 0x06, 0x4f, 0x41, 0x75, 0x74, 0x68,
	0x32, 0x12, 0x00, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x78, 0x66, 0x72, 0x72, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescOnce sync.Once
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescData = file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDesc
)

func file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescGZIP() []byte {
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescOnce.Do(func() {
		file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescData)
	})
	return file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDescData
}

var file_randomtalk_matchmaking_v1_matchmaking_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_randomtalk_matchmaking_v1_matchmaking_service_proto_goTypes = []any{
	(Gender)(0),               // 0: randomtalk.matchmaking.v1.Gender
	(*FindMatchRequest)(nil),  // 1: randomtalk.matchmaking.v1.FindMatchRequest
	(*FindMatchResponse)(nil), // 2: randomtalk.matchmaking.v1.FindMatchResponse
	(*GetMatchRequest)(nil),   // 3: randomtalk.matchmaking.v1.GetMatchRequest
	(*GetMatchResponse)(nil),  // 4: randomtalk.matchmaking.v1.GetMatchResponse
	(*MatchPreferences)(nil),  // 5: randomtalk.matchmaking.v1.MatchPreferences
	(*LatLng)(nil),            // 6: randomtalk.matchmaking.v1.LatLng
	(*Match)(nil),             // 7: randomtalk.matchmaking.v1.Match
}
var file_randomtalk_matchmaking_v1_matchmaking_service_proto_depIdxs = []int32{
	0, // 0: randomtalk.matchmaking.v1.FindMatchRequest.user_gender:type_name -> randomtalk.matchmaking.v1.Gender
	6, // 1: randomtalk.matchmaking.v1.FindMatchRequest.user_location:type_name -> randomtalk.matchmaking.v1.LatLng
	5, // 2: randomtalk.matchmaking.v1.FindMatchRequest.match_preferences:type_name -> randomtalk.matchmaking.v1.MatchPreferences
	7, // 3: randomtalk.matchmaking.v1.GetMatchResponse.match:type_name -> randomtalk.matchmaking.v1.Match
	0, // 4: randomtalk.matchmaking.v1.MatchPreferences.gender:type_name -> randomtalk.matchmaking.v1.Gender
	1, // 5: randomtalk.matchmaking.v1.MatchMakingService.FindMatch:input_type -> randomtalk.matchmaking.v1.FindMatchRequest
	3, // 6: randomtalk.matchmaking.v1.MatchMakingService.GetMatch:input_type -> randomtalk.matchmaking.v1.GetMatchRequest
	2, // 7: randomtalk.matchmaking.v1.MatchMakingService.FindMatch:output_type -> randomtalk.matchmaking.v1.FindMatchResponse
	4, // 8: randomtalk.matchmaking.v1.MatchMakingService.GetMatch:output_type -> randomtalk.matchmaking.v1.GetMatchResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_randomtalk_matchmaking_v1_matchmaking_service_proto_init() }
func file_randomtalk_matchmaking_v1_matchmaking_service_proto_init() {
	if File_randomtalk_matchmaking_v1_matchmaking_service_proto != nil {
		return
	}
	file_randomtalk_matchmaking_v1_match_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_randomtalk_matchmaking_v1_matchmaking_service_proto_goTypes,
		DependencyIndexes: file_randomtalk_matchmaking_v1_matchmaking_service_proto_depIdxs,
		EnumInfos:         file_randomtalk_matchmaking_v1_matchmaking_service_proto_enumTypes,
		MessageInfos:      file_randomtalk_matchmaking_v1_matchmaking_service_proto_msgTypes,
	}.Build()
	File_randomtalk_matchmaking_v1_matchmaking_service_proto = out.File
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_rawDesc = nil
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_goTypes = nil
	file_randomtalk_matchmaking_v1_matchmaking_service_proto_depIdxs = nil
}
