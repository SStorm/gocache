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
)

type SetOperation struct {
	key string
	flags, timeout int
	numBytes int
	body []byte
}

func ParseSet(buf *[]byte) (*SetOperation, error) {

	split := bytes.Fields(*buf)
	if len(split) < 4 {
		return nil, errors.New("CLIENT_ERROR bad command line format")
	}
	key := string(split[0])
	flags, err := strconv.Atoi(string(split[1]))
	if err != nil {
		return nil, errors.New("CLIENT_ERROR bad command line format")
	}
	timeout, err := strconv.Atoi(string(split[2]))
	if err != nil {
		return nil, errors.New("CLIENT_ERROR bad command line format")
	}
	numBytes, err := strconv.Atoi(string(split[3]))
	if err != nil {
		return nil, errors.New("CLIENT_ERROR bad command line format")
	}

	return &SetOperation{key, flags, timeout, numBytes, nil}, nil
}

