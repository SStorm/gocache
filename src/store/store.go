/**
 * Created with IntelliJ IDEA.
 * User: romanas
 * Date: 01/06/13
 * Time: 18:02
 * To change this template use File | Settings | File Templates.
 */
package store

import (
	"bytes"
)

type Storage interface {
	Set(key, value string, flags, timeout int)
	Get(key string) string
	Dump() string
}

type InMemoryStorage struct {
	storageMap map[string]string
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Set(key, value string, flags, timeout int) {
	s.init()
	s.storageMap[key] = value
}

func (s *InMemoryStorage) Get(key string) string {
	s.init()
	return s.storageMap[key]
}

func (s *InMemoryStorage) Dump() string {
	buffer := bytes.NewBufferString("");
	for k, v := range s.storageMap {
		buffer.WriteString(k + " -> " + v + "\n")
	}
	return buffer.String()
}

func (s *InMemoryStorage) init() {
	if s.storageMap == nil {
		s.storageMap = make(map[string]string)
	}
}

