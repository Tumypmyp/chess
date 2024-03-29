// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: proto/game/game.proto

package game

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GameID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID int64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *GameID) Reset() {
	*x = GameID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameID) ProtoMessage() {}

func (x *GameID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameID.ProtoReflect.Descriptor instead.
func (*GameID) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{0}
}

func (x *GameID) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type MoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameID   int64  `protobuf:"varint,1,opt,name=gameID,proto3" json:"gameID,omitempty"`
	PlayerID int64  `protobuf:"varint,2,opt,name=playerID,proto3" json:"playerID,omitempty"`
	Text     string `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *MoveRequest) Reset() {
	*x = MoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveRequest) ProtoMessage() {}

func (x *MoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveRequest.ProtoReflect.Descriptor instead.
func (*MoveRequest) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{1}
}

func (x *MoveRequest) GetGameID() int64 {
	if x != nil {
		return x.GameID
	}
	return 0
}

func (x *MoveRequest) GetPlayerID() int64 {
	if x != nil {
		return x.PlayerID
	}
	return 0
}

func (x *MoveRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type NewGameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayersID []int64 `protobuf:"varint,1,rep,packed,name=playersID,proto3" json:"playersID,omitempty"`
}

func (x *NewGameRequest) Reset() {
	*x = NewGameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewGameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewGameRequest) ProtoMessage() {}

func (x *NewGameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewGameRequest.ProtoReflect.Descriptor instead.
func (*NewGameRequest) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{2}
}

func (x *NewGameRequest) GetPlayersID() []int64 {
	if x != nil {
		return x.PlayersID
	}
	return nil
}

type GameStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string         `protobuf:"bytes,1,opt,name=Description,proto3" json:"Description,omitempty"`
	Keyboard    []*ArrayButton `protobuf:"bytes,2,rep,name=Keyboard,proto3" json:"Keyboard,omitempty"`
}

func (x *GameStatus) Reset() {
	*x = GameStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameStatus) ProtoMessage() {}

func (x *GameStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameStatus.ProtoReflect.Descriptor instead.
func (*GameStatus) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{3}
}

func (x *GameStatus) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GameStatus) GetKeyboard() []*ArrayButton {
	if x != nil {
		return x.Keyboard
	}
	return nil
}

type ArrayButton struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Buttons []*Button `protobuf:"bytes,1,rep,name=Buttons,proto3" json:"Buttons,omitempty"`
}

func (x *ArrayButton) Reset() {
	*x = ArrayButton{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArrayButton) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArrayButton) ProtoMessage() {}

func (x *ArrayButton) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArrayButton.ProtoReflect.Descriptor instead.
func (*ArrayButton) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{4}
}

func (x *ArrayButton) GetButtons() []*Button {
	if x != nil {
		return x.Buttons
	}
	return nil
}

type Button struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text         string `protobuf:"bytes,1,opt,name=Text,proto3" json:"Text,omitempty"`
	CallbackData string `protobuf:"bytes,2,opt,name=CallbackData,proto3" json:"CallbackData,omitempty"`
}

func (x *Button) Reset() {
	*x = Button{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_game_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Button) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Button) ProtoMessage() {}

func (x *Button) ProtoReflect() protoreflect.Message {
	mi := &file_proto_game_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Button.ProtoReflect.Descriptor instead.
func (*Button) Descriptor() ([]byte, []int) {
	return file_proto_game_game_proto_rawDescGZIP(), []int{5}
}

func (x *Button) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Button) GetCallbackData() string {
	if x != nil {
		return x.CallbackData
	}
	return ""
}

var File_proto_game_game_proto protoreflect.FileDescriptor

var file_proto_game_game_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x67, 0x61, 0x6d,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x18, 0x0a, 0x06, 0x47, 0x61,
	0x6d, 0x65, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x49, 0x44, 0x22, 0x55, 0x0a, 0x0b, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x2e, 0x0a, 0x0e, 0x4e,
	0x65, 0x77, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03,
	0x52, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x49, 0x44, 0x22, 0x5d, 0x0a, 0x0a, 0x47,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x08, 0x4b,
	0x65, 0x79, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x67, 0x61, 0x6d, 0x65, 0x2e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e,
	0x52, 0x08, 0x4b, 0x65, 0x79, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x22, 0x35, 0x0a, 0x0b, 0x41, 0x72,
	0x72, 0x61, 0x79, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x07, 0x42, 0x75, 0x74,
	0x74, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x2e, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x52, 0x07, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e,
	0x73, 0x22, 0x40, 0x0a, 0x06, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12,
	0x22, 0x0a, 0x0c, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x44,
	0x61, 0x74, 0x61, 0x32, 0x92, 0x01, 0x0a, 0x04, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x04,
	0x4d, 0x6f, 0x76, 0x65, 0x12, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x28, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0c, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x2e, 0x47, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x1a, 0x10, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x47,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2d, 0x0a, 0x07, 0x4e, 0x65, 0x77,
	0x47, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4e, 0x65, 0x77, 0x47,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x67, 0x61,
	0x6d, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_game_game_proto_rawDescOnce sync.Once
	file_proto_game_game_proto_rawDescData = file_proto_game_game_proto_rawDesc
)

func file_proto_game_game_proto_rawDescGZIP() []byte {
	file_proto_game_game_proto_rawDescOnce.Do(func() {
		file_proto_game_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_game_game_proto_rawDescData)
	})
	return file_proto_game_game_proto_rawDescData
}

var file_proto_game_game_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_game_game_proto_goTypes = []interface{}{
	(*GameID)(nil),         // 0: game.GameID
	(*MoveRequest)(nil),    // 1: game.MoveRequest
	(*NewGameRequest)(nil), // 2: game.NewGameRequest
	(*GameStatus)(nil),     // 3: game.GameStatus
	(*ArrayButton)(nil),    // 4: game.ArrayButton
	(*Button)(nil),         // 5: game.Button
	(*empty.Empty)(nil),    // 6: google.protobuf.Empty
}
var file_proto_game_game_proto_depIdxs = []int32{
	4, // 0: game.GameStatus.Keyboard:type_name -> game.ArrayButton
	5, // 1: game.ArrayButton.Buttons:type_name -> game.Button
	1, // 2: game.Game.Move:input_type -> game.MoveRequest
	0, // 3: game.Game.Status:input_type -> game.GameID
	2, // 4: game.Game.NewGame:input_type -> game.NewGameRequest
	6, // 5: game.Game.Move:output_type -> google.protobuf.Empty
	3, // 6: game.Game.Status:output_type -> game.GameStatus
	0, // 7: game.Game.NewGame:output_type -> game.GameID
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_game_game_proto_init() }
func file_proto_game_game_proto_init() {
	if File_proto_game_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_game_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_game_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MoveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_game_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewGameRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_game_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_game_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArrayButton); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_game_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Button); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_game_game_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_game_game_proto_goTypes,
		DependencyIndexes: file_proto_game_game_proto_depIdxs,
		MessageInfos:      file_proto_game_game_proto_msgTypes,
	}.Build()
	File_proto_game_game_proto = out.File
	file_proto_game_game_proto_rawDesc = nil
	file_proto_game_game_proto_goTypes = nil
	file_proto_game_game_proto_depIdxs = nil
}
