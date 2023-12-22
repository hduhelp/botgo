package test

import (
	"botgo/sdk/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func Random(ctx context.Context, args []string) any {
	return uuid.New().String()
}

func Time(ctx context.Context, args []string) any {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Msg(ctx context.Context, args []string) any {
	data := ctx.Value("msg").(*dto.WSGroupATMessageData)
	groupId := data.GroupId
	userId := data.Author.UserId
	content := strings.TrimSpace(data.Content)
	msgId := data.MsgId

	return fmt.Sprintf("UserId: %v\nGroupId: %v\nMsgId: %v\nContent: %v", userId, groupId, msgId, content)
}
