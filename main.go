package main

import (
	"botgo/sdk"
	"botgo/sdk/dto"
	"botgo/sdk/event"
	"botgo/sdk/token"
	"botgo/sdk/websocket"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	botToken := token.BotToken(102084832, "", string(token.TypeBot))
	api := sdk.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Printf("%+v, err:%v", ws, err)
	}

	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var atMessage event.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		fmt.Println(event, data)
		return nil
	}

	//var groupMessage event.GroupMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupMessageData) error {
	//	fmt.Println(event, data)
	//	newMsg := &dto.GroupMessageToCreate{
	//		Content: "干嘛",
	//		MsgID:   data.MsgId,
	//		MsgType: 0,
	//	}
	//	api.PostGroupMessage(ctx, data.GroupId, newMsg)
	//	return nil
	//}

	var groupAtMessage event.GroupAtMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		groupId := data.GroupId
		//userId := data.Author.UserId
		content := strings.TrimSpace(data.Content)
		msgId := data.MsgId

		resp, err := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: "https://avatars.githubusercontent.com/u/114216659?s=200&v=4", SrvSendMsg: false})
		if err != nil {
			newMsg := &dto.GroupMessageToCreate{
				Content: "图片上传失败",
				MsgID:   msgId,
				MsgType: 0,
			}
			api.PostGroupMessage(ctx, groupId, newMsg)
			return nil
		}

		newMsg := &dto.GroupMessageToCreate{
			Content: content,
			Media:   &dto.FileInfo{FileInfo: resp.FileInfo},
			MsgID:   msgId,
			MsgType: 7,
		}
		api.PostGroupMessage(ctx, groupId, newMsg)
		return nil
	}

	intent := websocket.RegisterHandlers(atMessage, groupAtMessage)

	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	sdk.NewSessionManager().Start(ws, botToken, &intent)
}
