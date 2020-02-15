package debug

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/levpaul/gecserv/internal/fb"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"net"
	"os"
)

const debugSocket = "/opt/idolscape/debug.sock"

type client struct {
	conn net.Conn
}

var count float64
var tiny = math.Float64frombits(1)

func StartDebugServer() {
	os.Remove(debugSocket)
	socket, err := net.Listen("unix", debugSocket)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal().Msg("Failed to accept socket connection")
		}
		log.Info().Strs("Addresses", []string{conn.LocalAddr().String(), conn.RemoteAddr().String()}).Msg("Accepted new client")

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
	//conn, err := ingest.DebugGetLiveConnection()
	//if err != nil {
	//	c.conn.Write([]byte("No active connections found\n"))
	//	return
	//}
	//c.conn.Write([]byte("Adding new player\n"))
	//conn.SendNewPlayerState(genRandomPlayer())
	log.Print("Echo: ", input)
}

func genRandomPlayer() []byte {
	mss := new(fb.MessageT)
	mss.Data = &fb.GameMessageT{
		Type: fb.GameMessageMapUpdate,
		Value: &fb.MapUpdateT{
			Seq: 12342345,
			Logins: []*fb.PlayerT{{
				Posx: rand.Float32() * 10,
				Posy: rand.Float32() * 10,
				Sid:  count,
				Col:  fb.ColorBlue,
			}},
			Logouts: nil,
			Psyncs:  nil,
		},
	}

	count += tiny

	b := flatbuffers.NewBuilder(1024)
	b.Finish(mss.Pack(b))

	return b.FinishedBytes()
}
