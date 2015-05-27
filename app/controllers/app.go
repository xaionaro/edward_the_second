package controllers;

import (
	"github.com/revel/revel";
	"golang.org/x/net/websocket";

	"os/exec";
);

type App struct {
	*revel.Controller;
}

func (c App) Index() revel.Result {
	return c.Render();
}

func (c App) Websocket(ws *websocket.Conn) revel.Result {
	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			// Otherwise, say something.
			revel.INFO.Printf("msg: %v\n", msg);
			exec.Command("i2csend", msg).Run();
		}
	}
	return nil
}

