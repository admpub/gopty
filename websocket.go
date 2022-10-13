package gopty

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/admpub/gopty/interfaces"
	"github.com/admpub/websocket"
)

func PTY2Websocket(ws *websocket.Conn, pty interfaces.Console) {
	buffer := make([]byte, 1024)
	var payload, overflow []byte
	for {
		n, err := pty.Read(buffer)
		if err != nil {
			fmt.Println("[PTY2Websocket] read from pty error: ", err)
			return
		}

		// Empty the overflow from the last read into the payload first.
		payload = append(payload[0:], overflow...)
		overflow = nil
		// Then empty the new buf read into the payload.
		payload = append(payload, buffer[:n]...)

		// Strip out any incomplete utf-8 from current payload into overflow.
		for !utf8.Valid(payload) {
			overflow = append(overflow[:0], append(payload[len(payload)-1:], overflow[0:]...)...)
			payload = payload[:len(payload)-1]
		}

		if len(payload) >= 1 {
			err = ws.WriteMessage(websocket.BinaryMessage, payload[:len(payload)])
			if err != nil {
				fmt.Println("[PTY2Websocket] write to ws error: ", err)
				return
			}
		}

		// Empty the payload.
		payload = nil
	}
}

func Websocket2PTY(ws *websocket.Conn, pty interfaces.Console) {
	for {
		mt, message, err := ws.ReadMessage()
		if mt == -1 || err != nil {
			log.Println("[Websocket2PTY] websocket read error: ", err)
			return
		}
		msg := string(message)
		if strings.HasPrefix(msg, "<RESIZE>") {
			size := msg[len("<RESIZE>"):len(msg)]
			sizeArr := strings.Split(size, ",")
			rows, _ := strconv.Atoi(sizeArr[0])
			cols, _ := strconv.Atoi(sizeArr[1])
			err = pty.SetSize(cols, rows)
			log.Printf("[Websocket2PTY] pty resize window to %d, %d", cols, rows)
			if err != nil {
				log.Println("[Websocket2PTY] pty resize error: ", err)
				return
			}
		} else {
			_, err = pty.Write(message)
			if err != nil {
				log.Println("[Websocket2PTY] pty write error: ", err)
			}
		}
	}
}

func Bash(wsc *websocket.Conn) error {
	pty, err := New(120, 60)
	if err != nil {
		return err
	}
	defer pty.Close()
	args := []string{}
	err = pty.Start(args)
	if err != nil {
		err = fmt.Errorf("open terminal err: %w", err)
		return err
	}

	go PTY2Websocket(wsc, pty)
	// block from close
	Websocket2PTY(wsc, pty)
	return nil
}
