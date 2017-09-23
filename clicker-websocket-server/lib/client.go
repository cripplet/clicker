package cc_websocket_server

import (
	"fmt"
	"github.com/cripplet/clicker/lib"
	"github.com/gorilla/websocket"
	"net/http"
)

var client_upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	connectionStatus ClientConnectionStatus
	connection       *websocket.Conn
	gameID           string
	game             *cookie_clicker.GameStateStruct
}

func (self *Client) execute(request *CommandRequest, response *CommandResponse) {
	request.validate(&(response.Error))
	if response.Error.ErrorCode != ERROR_TYPE_SUCCESS {
		return
	}
	COMMAND_DISPATCH_TABLE[SupportedCommand{
		object: request.Object,
		hasID:  request.ID != "",
		method: request.Method,
	}](self, request, response)
}

func (self *Client) run() {
	defer self.connection.Close()

	for {
		request := CommandRequest{}
		response := CommandResponse{}
		err := self.connection.ReadJSON(&request)
		if err != nil {
			// TODO(cripplet): If connection stable, return err, otherwise call self.game.Stop() and exit function.
			return
		} else {
			self.execute(&request, &response)
		}
		self.connection.WriteJSON(&response)
	}
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := client_upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	c := Client{
		connection: conn,
	}
	go c.run()
}
