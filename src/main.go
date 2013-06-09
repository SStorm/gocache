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
	var setInProgress *SetOperation
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
		fmt.Println("Reading: ", string(buffer[0:n-1]))
		prefix := string(buffer[0:3])

		if setInProgress != nil {
			err := handleSetBody(setInProgress, buffer)
			if err != nil {
				c.Write([]byte("CLIENT_ERROR bad data chunk"))
				continue
			}
			c.Write([]byte("STORED\n"))
			setInProgress = nil
			continue;
		}

		switch prefix {
		case "set", "add":
			setInProgress, err = handleSet(c, buffer[4:n])
			if setInProgress == nil && err == nil {
				c.Write([]byte("STORED\r\n"))
			}
		case "get":
			res, err := handleGet(buffer[4:n])
			if err != nil {
				c.Write([]byte(err.Error() + "\r\n"))
			}
			if len(res) > 0 {
				info := fmt.Sprintf("%s %s %d %d\r\n", "VALUE", buffer[4:n-2], 0, len(res))
				c.Write([]byte(info))
				c.Write(res)
			}
			c.Write([]byte("\r\nEND\r\n"))
		case "dum":
			c.Write([]byte("==== DUMP ====\n"))
			c.Write([]byte(handleDump()))
			c.Write([]byte("==== END DUMP ====\n"))
		default:
			c.Write([]byte("Unknown command\n"))
		}
	}
}

func handleSet(c net.Conn, buf []byte) (*SetOperation, error) {
	fmt.Println("Handle set")

	oper, err := ParseSet(&buf)

	if err != nil {
		c.Write([]byte("CLIENT_ERROR bad command line format\n"))
		return nil, errors.New("CLIENT_ERROR bad command line format")
	}

	fmt.Println("Numbytes: ", oper.numBytes, "ReadSoFar:", oper.readSoFar)
	if oper.numBytes <= oper.readSoFar {
		dataStorage.Set(oper.key, oper.body[0:oper.numBytes], oper.flags, oper.timeout)
		return nil, nil
	}

	return oper, nil
}

func handleSetBody(oper *SetOperation, buf []byte) error {
	continuedOperation, err := ParseSetContinue(oper, &buf)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	dataStorage.Set(continuedOperation.key, continuedOperation.body, continuedOperation.flags,
		continuedOperation.timeout)
	return nil
}

func handleGet(buf []byte) ([]byte, error) {
	fmt.Println("Handle get")
	var key string
	split := bytes.Fields(buf)
	if len(split) != 1 {
		return nil, errors.New("Only key must be specified in get operations")
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
