package main

import (
	"botgo/handler/group"
	"botgo/sdk"
	"botgo/sdk/event"
	"botgo/sdk/openapi"
	"botgo/sdk/token"
	"botgo/sdk/websocket"
	"context"
	"github.com/jinzhu/configor"
	"log"
	"os"
	"time"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	var cfg Config
	err := configor.New(&configor.Config{Verbose: true}).Load(&cfg, "config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	botToken := token.BotToken(cfg.AppID, cfg.AccessToken, string(token.TypeBot))
	api := sdk.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	openapi.DefaultImpl = api
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Printf("%+v, err:%v", ws, err)
		os.Exit(1)
	}

	intent := websocket.RegisterHandlers(
		event.GroupAtMessageEventHandler(group.HandleGroupAtMessage),
	)

	sdk.NewSessionManager().Start(ws, botToken, &intent)
}
