package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"github.com/0x17de/uoproxy/packets"
	"github.com/0x17de/uoproxy/uocrypt"
)

var (
	GameHost  string
	GamePort  int
	LobbyHost string
	LobbyPort int
)

type Session struct {
	C           net.Conn
	Uo          net.Conn
	isGame      bool
	ClilocTable map[uint]string
}

func streamReader(c net.Conn) chan []byte {
	o := make(chan []byte)
	go func() {
		defer close(o)
		for {
			data := make([]byte, 3072)
			count, err := c.Read(data)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				log.Printf("Socker error: %s\n", err)
				break
			}
			if count <= 0 {
				log.Printf("Shutdown socket. 0len data\n")
				break
			} else {
				o <- data[0:count]
			}
		}
	}()
	return o
}

func byteReader(in chan []byte) chan byte {
	o := make(chan byte)
	go func() {
		defer close(o)
		for {
			b, ok := <-in
			if !ok {
				return
			}
			for _, c := range b {
				o <- c
			}
		}
	}()
	return o
}

func packetReader(s *Session, clientToServer bool, in chan byte) chan packets.Packeter {
	o := make(chan packets.Packeter)
	go func() {
		defer close(o)
		if !s.isGame {
			if clientToServer {
				// read version
				p := &packets.BytePacket{Len: 4}
				p.Read(in)
				o <- p
			}
		}
		for {
			cmd, ok := <-in
			if !ok {
				return
			}
			info := packets.PacketMap[int(cmd)]
			if info.Handler != nil {
				p := reflect.New(reflect.Indirect(reflect.ValueOf(info.Handler)).Type()).Interface().(packets.Packeter)
				p.Base().Info = info
				p.Base().ClientToServer = clientToServer
				p.Read(in)
				o <- p
			} else {
				p := &packets.UnknownPacket{BasePacket: packets.BasePacket{Info: info, ClientToServer: clientToServer}}
				p.Read(in)
				o <- p
			}
		}
	}()
	return o
}

func (s *Session) runGame() {
	uoStream := streamReader(s.Uo)
	uoBytes := uocrypt.Decompress(byteReader(uoStream))
	uoToCStream := packetReader(s, false, uoBytes)
	cStream := streamReader(s.C)
	cBytes := byteReader(cStream)
	s.Uo.Write([]byte{<-cBytes})
	s.Uo.Write([]byte{<-cBytes})
	s.Uo.Write([]byte{<-cBytes})
	s.Uo.Write([]byte{<-cBytes})
	cToUoStream := packetReader(s, true, cBytes)

sessionLoop:
	for {
		select {
		case p, ok := <-uoToCStream: // from server
			if !ok {
				break sessionLoop
			}
			switch pp := p.(type) {
			case *packets.GenericCommandPacket:
				log.Printf("UO => C: generic subcmd 0x%X", pp.SubCmd)
			case *packets.HealthPacket: // ignore
			case *packets.MegaClilocPacket:
				setLastMegaCliloc(pp)
				/*
					s := pp.Server
					for _, c := range s.Clilocs {
						fmt.Printf("Item %X.%X|%X: %s\n", s.Serial1, s.Serial2, c.Id, c.Text)
					}
				*/
			case *packets.MovePacket: // ignore
			case *packets.PingPacket: // ignore
			case *packets.RevisionPacket: // ignore
			case *packets.SendSpeechPacket:
				log.Printf("UO => C: speaking 0x%X: %s", pp.Target, pp.Text)
			case *packets.UnknownPacket:
				log.Printf("UO => C: cmd 0x%X of len %d", pp.Info.Id, pp.Len)
			default:
				log.Printf("UO => C: %T", p)
			}
			s.C.Write(uocrypt.Compress(p.Bytes()))
		case p, ok := <-cToUoStream: // from client
			if !ok {
				break sessionLoop
			}
			switch pp := p.(type) {
			case *packets.GenericCommandPacket:
				if pp.SubCmd == 0x24 {
					break
				}
				log.Printf("C => UO: generic subcmd 0x%X", pp.SubCmd)
			case *packets.MovePacket: // ignore
			case *packets.MoveRequestPacket: // ignore
			case *packets.PingPacket: // ignore
			case *packets.PostLoginPacket:
				log.Printf("C => UO: login %s:%s %X", pp.User, pp.Pass, pp.Key)
			case *packets.UnknownPacket:
				log.Printf("C => UO: cmd 0x%X of len %d", pp.Info.Id, pp.Len)
			default:
				log.Printf("C => UO: %T", p)
			}
			s.Uo.Write(p.Bytes())
		}
	}
	log.Printf("Game session loop ended\n")

	s.C.Close()
	s.Uo.Close()
}
func (s *Session) runLobby() {
	uoStream := streamReader(s.Uo)
	uoBytes := byteReader(uoStream)
	uoToCStream := packetReader(s, false, uoBytes)
	cStream := streamReader(s.C)
	cBytes := byteReader(cStream)
	cToUoStream := packetReader(s, true, cBytes)

sessionLoop:
	for {
		select {
		case p, ok := <-uoToCStream:
			if !ok {
				break sessionLoop
			}
			switch pp := p.(type) {
			case *packets.UnknownPacket:
				log.Printf("UO => C: cmd 0x%X of len %d", pp.Info.Id, pp.Len)
			case *packets.ServerListPacket:
				pp.Servers[0].Name += " (LOCAL)"
			case *packets.ServerRedirectPacket:
				GameHost = net.IPv4(pp.Address[0], pp.Address[1], pp.Address[2], pp.Address[3]).String()
				GamePort = pp.Port
				fmt.Printf("Original host: %s:%d\n", GameHost, GamePort)
				pp.Address[0] = 127
				pp.Address[1] = 0
				pp.Address[2] = 0
				pp.Address[3] = 1
				pp.Port = 1061
			default:
				log.Printf("UO => C: %T", p)
			}
			s.C.Write(p.Bytes())
		case p, ok := <-cToUoStream:
			if !ok {
				break sessionLoop
			}
			switch pp := p.(type) {
			case *packets.UnknownPacket:
				log.Printf("C => UO: cmd 0x%X of len %d", pp.Info.Id, pp.Len)
			default:
				log.Printf("C => UO: %T", p)
			}
			s.Uo.Write(p.Bytes())
		}
	}
	log.Printf("Lobby session loop ended\n")

	s.C.Close()
	s.Uo.Close()
}

func main() {
	var (
		isServer bool
	)
	flag.BoolVar(&isServer, "server", false, "enable server mode")
	flag.Parse()
	initUi()

	loginServer, _ := net.Listen("tcp", "127.0.0.1:1060")
	gameServer, _ := net.Listen("tcp", "127.0.0.1:1061")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			c, _ := loginServer.Accept()
			fmt.Printf("LOBBY CONNECT!")
			uo, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", LobbyHost, LobbyPort))
			s := &Session{
				C:      c,
				Uo:     uo,
				isGame: false,
			}
			go s.runLobby()
		}
	}()
	go func() {
		defer wg.Done()
		for {
			c, _ := gameServer.Accept()
			fmt.Printf("GAME CONNECT! %T", c)
			uo, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", GameHost, GamePort))
			s := &Session{
				C:      c,
				Uo:     uo,
				isGame: true,
			}
			go s.runGame()
		}
	}()
	runUi()
	wg.Wait()
}
