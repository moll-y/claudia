package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetPrefix("claudia: ")
	// Set up channel on which to send signal notifications. We must use a
	// buffered channel or risk missing the signal if we're not ready to
	// receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	address := "/tmp/claudia.sock"
	l, err := net.Listen("unix", address)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			log.Print(err)
		}
	}()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// This error "ErrClosed" is returned when the
				// listener is closed during shutdownw. In this
				// case we just call "return" to exit this
				// goroutine.
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Print(err)
				continue
			}

			func() {
				defer func() {
					if err := conn.Close(); err != nil {
						log.Print(err)
					}
				}()

				buf := make([]byte, 512)
				n, err := conn.Read(buf)
				if err != nil {
					log.Print(err)
					return
				}

				done, word, err := parse(buf[:n])
				if err != nil {
					log.Print(err)
					return
				}

				// Waybar expects the exec-script to output its
				// data in JSON format. This should look like
				// this:
				//   {
				//     "alt": "$alt",
				//     "text": "$text",
				//     "class": "$class",
				//     "tooltip": "$tooltip",
				//     "percentage": $percentage
				//   }
				if done {
					fmt.Printf("{\"text\":\"π %c\",\"class\":\"done\"}\n", word[0])
				} else {
					fmt.Printf("{\"text\":\"π %c\"}\n", word[0])
				}
			}()
		}
	}()

	log.Printf("listening to socket: %s", address)
	fmt.Println("{\"text\":\"π I\"}")
	// Block until any signal is received.
	log.Print(<-c)
}
