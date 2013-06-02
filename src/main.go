/**
 * Created with IntelliJ IDEA.
 * User: romanas
 * Date: 01/06/13
 * Time: 17:08
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"fmt"
	"net"
	"log"
	"log/syslog"
	"store"
	"bytes"
	"errors"
	"strings"
)

var Logger *log.Logger

// 10 kb
const RECEIVE_BUFFER_LENGTH = 1024 * 10

var connectionCount int;

var dataStorage store.Storage

func main() {
	initLogging();

	fmt.Println("Starting goCache");
	listener, err := net.Listen("tcp", "0.0.0.0:11212");
	if err != nil {
		fmt.Println("Failed to start listener:", err.Error());
		return;
	}
	fmt.Println("Started");

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err);
		}
		go handleConn(c);
	}
}

func handleConn(c net.Conn) {
	incrementConnectionCount()
	fmt.Println("New connection, total count: ", connectionCount);
	setInProgress := false
	for {
		buffer := make([]byte, RECEIVE_BUFFER_LENGTH);
		n, err := c.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read: ", err);
			c.Close();
			decrementConnectionCount()
			fmt.Println("Client disconnected, total count: ", connectionCount);
			return;
		}
		fmt.Println("Bytes read:", n)
		prefix := string(buffer[0:3])
		switch prefix {
		case setInProgress:

		case "set", "add":
			err := handleSet(buffer[4:n])
			if err != nil {
				c.Write([]byte(err.Error() + "\n"))
				return;
			}
			c.Write([]byte("END\n"))
		case "get":
			res, err := handleGet(buffer[4:n])
			if err != nil {
				c.Write([]byte(err.Error() + "\n"))
			}
			if res != "" {
				c.Write([]byte(res + "\n"))
			}
			c.Write([]byte("END\n"))
		case "dum":
			c.Write([]byte("==== DUMP ====\n"))
			c.Write([]byte(handleDump()))
			c.Write([]byte("==== END DUMP ====\n"))
		default:
			c.Write([]byte("Unknown command\n"))
		}
	}
}

func handleSet(buf []byte) error {
	fmt.Println("Handle set")

	key, flags, timeout, numBytes, err := ParseSet(&buf)

	if err != nil {
		return errors.New("CLIENT_ERROR bad command line format")
	}

	return nil
}

func handleSetValue(buf []byte, key string, flags, timeout int) {

}

func handleGet(buf []byte) (string, error) {
	fmt.Println("Handle get")
	var key string
	split := bytes.Fields(buf)
	if len(split) != 1 {
		return "", errors.New("Only key must be specified in get operations")
	}
	key = string(split[0])
	return dataStorage.Get(key), nil
}

func handleDump() string {
	return dataStorage.Dump()
}

func initLogging() {
	var err error;
	Logger, err = syslog.NewLogger(syslog.LOG_INFO, 0);
	if err != nil {
		fmt.Println("Failed to initialize logger")
	}
}

func init() {
	fmt.Println("Init storage")
	dataStorage = store.NewStorage()
}

func incrementConnectionCount() {
	connectionCount = connectionCount + 1
}

func decrementConnectionCount() {
	connectionCount = connectionCount - 1
}
