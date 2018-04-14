/*
 *  Copyright (c) 2017, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rpc

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/grpc_util"
	"github.com/nebulaim/telegramd/mtproto"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/biz/base"
	chat2 "github.com/nebulaim/telegramd/biz/core/chat"
	"github.com/nebulaim/telegramd/biz/core/account"
	"github.com/nebulaim/telegramd/biz/core/user"
)

// messages.getFullChat#3b831c66 chat_id:int = messages.ChatFull;
func (s *MessagesServiceImpl) MessagesGetFullChat(ctx context.Context, request *mtproto.TLMessagesGetFullChat) (*mtproto.Messages_ChatFull, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("MessagesGetFullChat - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	messagesChatFull := mtproto.NewTLMessagesChatFull()
	// &mtproto.TLMessagesChatFull{}
	chatFull := chat2.GetChatFull(request.ChatId)
	peer := &base.PeerUtil{}
	peer.PeerType = base.PEER_CHAT
	peer.PeerId = request.ChatId
	chatFull.SetNotifySettings(account.GetNotifySettings(md.UserId, peer))
	chat := chat2.GetChat(request.ChatId)
	// chat.ParticipantsCount = len(chatFull.GetParticipants().GetChatParticipants().GetParticipants())
	chatUserIdList := make([]int32, 0)
	participants := chatFull.GetParticipants().GetData2().GetParticipants()
	for _, participant := range participants {
		switch participant.GetConstructor() {
		case mtproto.TLConstructor_CRC32_chatParticipantCreator:
			chatUserIdList = append(chatUserIdList, md.UserId)
		case mtproto.TLConstructor_CRC32_chatParticipant:
			chatUserIdList = append(chatUserIdList, participant.GetData2().GetUserId())
		case mtproto.TLConstructor_CRC32_chatParticipantAdmin:
			chatUserIdList = append(chatUserIdList, participant.GetData2().GetUserId())
		}
	}
	chat.SetParticipantsCount(int32(len(participants)))
	users := user.GetUsersBySelfAndIDList(md.UserId, chatUserIdList)

	messagesChatFull.SetUsers(users)
	//for _, u := range users {
	//	if u.GetId() == md.UserId {
	//		u.SetSelf(true)
	//	}
	//	messagesChatFull.Data2.Users = append(messagesChatFull.Data2.Users, u.To_User())
	//}
	messagesChatFull.Data2.Chats = append(messagesChatFull.Data2.Chats, chat.To_Chat())
	messagesChatFull.SetFullChat(chatFull.To_ChatFull())

	glog.Infof("MessagesGetFullChat - reply: %s", logger.JsonDebugData(chatFull))
	return messagesChatFull.To_Messages_ChatFull(), nil
}
