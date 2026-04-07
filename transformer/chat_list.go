package transformer

import (
	"pkm/model"
)

type ChatList struct {
	ConversationId string `json:"conversationId,omitempty"` // single
	ChatGroupId    string `json:"chatGroupId,omitempty"`    // group
	Title          string `json:"title,omitempty"`          // username for single chat
	Message        string `json:"message,omitempty"`
	MediaUrl       string `json:"mediaUrl,omitempty"`
	MediaType      string `json:"mediaType,omitempty"`
	GroupChatImg   string `json:"groupChatImg,omitempty"` // group
	Count          int64  `json:"count"`
	LastReadId     string `json:"lastReadId"`
	model.BaseModel
}

// func ToUserChatHistory(m *model.Message, conversationId string, u *model.User, count int64, lastReadId string) *ChatList {
// 	chat := ChatList{
// 		ConversationId: conversationId,
// 		Message:        m.Message,
// 		MediaUrl:       m.MediaUrl,
// 		MediaType:      m.MediaType,
// 		Title:          u.Username,
// 		Count:          count,
// 		LastReadId:     lastReadId,
// 		BaseModel:      m.BaseModel,
// 	}

// 	return &chat
// }

// func ToGroupChatHistory(cg *model.ChatGroup, cgMsg *model.ChatGroupMessages, count int64, lastReadId string) *ChatList {
// 	chat := ChatList{}

// 	if cg != nil {
// 		chat.ChatGroupId = cg.Id
// 		chat.Title = cg.Title
// 		chat.GroupChatImg = cg.GroupChatImg
// 		chat.BaseModel = cg.BaseModel
// 	}

// 	if cgMsg != nil {
// 		chat.Message = cgMsg.Message
// 		chat.MediaUrl = cgMsg.MediaUrl
// 		chat.MediaType = cgMsg.MediaType
// 		chat.BaseModel = cgMsg.BaseModel
// 		chat.Count = count
// 		chat.LastReadId = lastReadId
// 	}

// 	return &chat
// }

// func ToAllChatHistoryLists(cg []*model.ChatGroup, conversationIdMap map[string]string, cgMsgMap map[string]*model.ChatGroupMessages, m []*model.Message, userMap map[string]*model.User, countMap map[string]int64, lastReadIdMap map[string]string) []*ChatList {
// 	totalSize := len(m) + len(cg)
// 	if totalSize == 0 {
// 		return nil
// 	}

// 	o := make([]*ChatList, totalSize)
// 	pool := grpool.NewPool(20, 20)
// 	pool.WaitCount(totalSize)
// 	defer pool.Release()

// 	// Handle personal messages
// 	for i, item := range m {
// 		pool.JobQueue <- func(index int, val *model.Message) func() {
// 			return func() {
// 				defer pool.JobDone()
// 				o[index] = ToUserChatHistory(val, conversationIdMap[val.Id], userMap[val.Id], countMap[val.Id], lastReadIdMap[val.Id])
// 			}
// 		}(i, item)
// 	}

// 	// Handle group messages
// 	for j, group := range cg {
// 		pool.JobQueue <- func(index int, val *model.ChatGroup) func() {
// 			return func() {
// 				defer pool.JobDone()
// 				o[len(m)+index] = ToGroupChatHistory(val, cgMsgMap[val.Id], countMap[val.Id], lastReadIdMap[val.Id])
// 			}
// 		}(j, group)
// 	}

// 	pool.WaitAll()
// 	return o
// }
