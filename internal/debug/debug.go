package debug

import (
	"github.com/rs/zerolog/log"
	"net"
	"os"
)

const debugSocket = "/opt/idolscape/debug.sock"

type client struct {
	conn net.Conn
}

func StartDebugServer() {
	os.Remove(debugSocket)
	socket, err := net.Listen("unix", debugSocket)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer socket.Close()

	for {
		log.Info().Msg("Listning")
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal().Msg("Failed to accept socket connection")
		}

		newClient := &client{conn: conn}
		go newClient.handle()
	}
}

func (c *client) handle() {
	for {
		buf := make([]byte, 512)
		nr, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		c.handleInputLine(string(data))
	}
}

func (c *client) handleInputLine(input string) {
	_, err := c.conn.Write([]byte("Right back at you: " + input))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
