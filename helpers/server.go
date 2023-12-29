package helpers

import (
	"fmt"
	"net"
)

const DEFAULT_PORT = 28899

type Credential struct {
}

func StoreToken(cred *Credential) error {
	return nil
}
func handleCallbackRequest(c net.Conn) (*Credential, error) {
	return &Credential{}, nil
}

func StartCallbackServer(port int) error {
	var listenerPort int
	if port != 0 {
		listenerPort = port
	} else {
		listenerPort = DEFAULT_PORT
	}

	var listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", listenerPort))
	if err != nil {
		return err
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			var credential, err = handleCallbackRequest(conn)
			if err != nil {
				return err
			}
			StoreToken(credential)
			break
		} else {
			return err
		}
	}
	return nil
}
