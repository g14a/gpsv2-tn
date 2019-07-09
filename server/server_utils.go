package server

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/gpsv2-kudankulam/amqputils"
	"gitlab.com/gpsv2-kudankulam/config"
	"gitlab.com/gpsv2-kudankulam/errorcheck"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn
func HandleConnection(conn net.Conn) {

	var wg sync.WaitGroup

	wg.Add(1)
	go readTCPClient(conn, &wg)
	wg.Wait()

}

var (
	// AMQP
	amqpConnection = amqputils.GetAMQPInstance()
	amqpQueue      = config.GetAppConfig().AMQPConfig.AMQPQueue
)

// readTCPClient reads data sent by the device(a TCP client)
// and pushes it to the DB in an overview. Read more documentation below
func readTCPClient(conn net.Conn, wg *sync.WaitGroup) {

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

	defer wg.Done()

	ch, err := amqpConnection.Channel()
	errorcheck.CheckError(err)

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 5*1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed EOF")
				_ = conn.Close()
				count--
			}
		} else {

			q, err := ch.QueueDeclare(amqpQueue, false, false, false, false, nil)

			errorcheck.CheckError(err)

			err = ch.Publish("", q.Name, false, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(buf),
				})
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
