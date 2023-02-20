// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"errors"

	ws "github.com/illacloud/builder-backend/internal/websocket"
)

func SignalCooperateAttach(hub *ws.Hub, message *ws.Message) error {
	currentClient := hub.Clients[message.ClientID]

	// attach components
	inRoomUsers := hub.GetInRoomUsersByRoomID(currentClient.APPID)
	displayNames := make([]string, 0)
	for _, displayNameInterface := range message.Payload {
		displayName, assertCorrectly := displayNameInterface.(string)
		if !assertCorrectly {
			return errors.New("user input assert failed with signal cooperate attach.")
		}
		displayNames = append(displayNames, displayName)
	}
	inRoomUsers.AttachComponent(currentClient.ExportMappedUserIDToString(), displayNames)

	// broadcast attached components users
	message.SetBroadcastType(ws.BROADCAST_TYPE_ATTACH_COMPONENT)
	message.RewriteBroadcast()
	message.SetBroadcastPayload(inRoomUsers.FetchAllAttachedUsers())
	hub.BroadcastToRoomAllClients(message, currentClient)

	return nil
}
