/**
 * Created with IntelliJ IDEA.
 * User: romanas
 * Date: 02/06/13
 * Time: 17:40
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"strconv"
	"bytes"
	"errors"
	"unicode"
)

func ParseSet(buf *[]byte) (key string, flags, timeout int, numBytes int, err error) {

	split := bytes.Fields(buf)
	if len(split) < 4 {
		return "", 0, 0, "", errors.New("CLIENT_ERROR bad command line format")
	}
	key = string(split[0])
	flags, err := strconv.Atoi(string(split[1]))
	if err != nil {
		return "", 0, 0, "", errors.New("CLIENT_ERROR bad command line format")
	}
	timeout, err := strconv.Atoi(string(split[2]))
	if err != nil {
		return "", 0, 0, "", errors.New("CLIENT_ERROR bad command line format")
	}
	numBytes, err := strconv.Atoi(string(split[3]))
	if err != nil {
		return "", 0, 0, "", errors.New("CLIENT_ERROR bad command line format")
	}

	return key, flags, timeout, numbytes, nil
}

