package server

import (
	"errors"
	"fmt"
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/db"
	"gitlab.com/gps2.0/errcheck"
	"gitlab.com/gps2.0/models"
	"go.mongodb.org/mongo-driver/bson"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var clients []net.Conn
var count = 0

func HandleConnection(conn net.Conn) {

	errorChan := make(chan error)
	dataChan := make(chan []byte)

	go readWrapper(conn, dataChan, errorChan)
	go connCheckForShutdown(conn)

	for {
		select {
		case data := <-dataChan:

			log.Printf("[SERVER} Client %s sent: %s", conn.RemoteAddr(), string(data))
			//gtplDevice := ParseGTPLData(string(data))
			//
			//fmt.Println(gtplDevice)

		case err := <-errorChan:
			log.Println("[SERVER] An error occured:", err.Error())
			return
		}
	}
}

var (
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection    = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection
	collectionMutex             = &sync.Mutex{}
)

func connCheckForShutdown(conn net.Conn) error {
	var (
		n    int
		err  error
		buff [1]byte
	)

	sconn, ok := conn.(syscall.Conn)

	if !ok {
		return nil
	}

	rc, err := sconn.SyscallConn()

	if err != nil {
		return err
	}

	rerr := rc.Read(func(fd uintptr) bool {
		n, err = syscall.Read(int(fd), buff[:])
		return true
	})

	switch {
	case rerr != nil:
		return rerr
	case n == 0 && err == nil:
		return io.EOF
	case n > 0:
		return errors.New("unexpected read from socket")
	case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
		return nil
	default:
		return err
	}
}

func readWrapper(conn net.Conn, dataChan chan []byte, errorChan chan error) {
	for {
		buf := make([]byte, 5*1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			errorChan <- err
			return
		}
		dataChan <- buf[:reqLen]
	}
}

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

func removeClient(conn net.Conn) {
	log.Printf("[SERVER] Client %s disconnected", conn.RemoteAddr())
	count--
	conn.Close()
	//remove client from clients here
}

func InsertGTPLDataIntoMongo(gtplDevice *models.GTPLDevice) error {
	locationHistoriesCollection, ctx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	collectionMutex.Lock()

	_, err := locationHistoriesCollection.InsertOne(ctx, gtplDevice)

	collectionMutex.Unlock()
	errcheck.CheckError(err)

	return err
}

func UpdateGTPLDataIntoMongo(gtpldevice *models.GTPLDevice) error {

	vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	cursor, err := vehicleDetailsCollection.Find(ctx, bson.M{"deviceid": gtpldevice.DeviceID}, &options)

	if cursor.Next(ctx) {

		collectionMutex.Lock()

		_, err = vehicleDetailsCollection.ReplaceOne(ctx, bson.M{"deviceid": gtpldevice.DeviceID}, &gtpldevice)

		collectionMutex.Unlock()
		errcheck.CheckError(err)

		return err
	}

	collectionMutex.Lock()

	_, err = vehicleDetailsCollection.InsertOne(ctx, gtpldevice)

	collectionMutex.Unlock()
	errcheck.CheckError(err)

	return err
}
