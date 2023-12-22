package group

import (
	"botgo/handler/group/modules/test"
	"botgo/pkg/command"
	"botgo/sdk/dto"
	"botgo/sdk/openapi"
	"context"
	"log"
	"strings"
)

var cmd = command.New("", "", "", nil)

func init() {
	cmd.AddSubCommands(
		command.New("测试组", "test", "用来测试一些基础功能", nil).AddSubCommands(
			command.New("随机数", "random", "回复一串随机数", test.Random),
			command.New("当前时间", "time", "回复当前时间", test.Time),
			command.New("消息信息", "msg", "类似于ping，展示一些debug信息", test.Msg),
		),
		command.New("测试组2", "test2", "用来测试一些基础功能", nil).AddSubCommands(
			command.New("随机数", "random", "回复一串随机数", test.Random),
			command.New("当前时间", "time", "回复当前时间", test.Time),
			command.New("消息信息", "msg", "类似于ping，展示一些debug信息", test.Msg),
		),
	)
}

func HandleGroupAtMessage(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
	groupId := data.GroupId
	userId := data.Author.UserId
	content := strings.TrimSpace(data.Content)
	msgId := data.MsgId

	log.Printf("GroupId: %v, UserId: %v, MsgId: %v, Content: %v", groupId, userId, msgId, content)

	ctx := context.WithValue(context.Background(), "msg", data)
	resp := cmd.Exec(ctx, content)
	if resp == nil {
		return nil
	}
	switch resp.(type) {
	case string:
		_, err := openapi.DefaultImpl.PostGroupMessage(context.Background(), groupId, &dto.GroupMessageToCreate{
			Content: resp.(string),
			MsgID:   msgId,
		})
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
