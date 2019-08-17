package server

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/gpsv2-tn/amqputils"
	"gitlab.com/gpsv2-tn/config"
	"gitlab.com/gpsv2-tn/errorcheck"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn

var (
	// AMQP
	amqpConnection = amqputils.GetAMQPInstance()
	amqpQueue      = config.GetAppConfig().AMQPConfig.AMQPQueue
)

// readTCPClient reads data sent by the device(a TCP client)
// and pushes it to the DB in an overview. Read more documentation below
func HandleConnection(conn net.Conn)  {

	fmt.Println("running goroutines: ", runtime.NumGoroutine())

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed EOF")
				_ = conn.Close()
				count--
			}
		} else {
			publishChannel, err := amqpConnection.Channel()

			q, err := publishChannel.QueueDeclare(amqpQueue, false, false, false, false, nil)

			errorcheck.CheckError(err)

			fmt.Println("ONE CHUNK......")

			buffer := string(buf)

			if strings.Contains(buffer, "GTPL") || strings.Contains(buffer, "BSTPL") {
				dataslice := strings.Split(string(buf), "#")

				for _, record := range dataslice {

					fmt.Println(record)

					err = publishChannel.Publish("", q.Name, false, false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(record),
						})
				}

				_ = publishChannel.Close()

			} else if strings.Contains(buffer, "AVA") || strings.Contains(buffer, "*") {
				dataslice := strings.Split(string(buf), "*")

				for _, record := range dataslice {

					err = publishChannel.Publish("", q.Name, false, false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(record),
						})
				}
				_ = publishChannel.Close()
			}
		}
	}
}

// signalHandler notices termination signals or
// interrupts from the command line. Eg: ctrl-c and exits cleanly
func signalHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		for sig := range sigchan {
			log.Printf("[SERVER] Closing due to Signal: %s", sig)
			log.Printf("[SERVER] Graceful shutdown")

			fmt.Println("Done.")

			// Exit cleanly
			os.Exit(0)
		}
	}()
}
