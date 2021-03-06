package store

import (
	"bytes"
	"fmt"
)

type Storage interface {
	Set(key string, value []byte, flags, timeout int)
	Get(key string) []byte
	Dump() string
	Stats() string
}

type InMemoryStorage struct {
	storageMap map[string][]byte
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Set(key string, value []byte, flags, timeout int) {
	s.init()
	s.storageMap[key] = value
}

func (s *InMemoryStorage) Get(key string) []byte {
	s.init()
	return s.storageMap[key]
}

func (s *InMemoryStorage) Dump() string {
	buffer := bytes.NewBufferString("")
	for k, v := range s.storageMap {
		buffer.WriteString(k + " -> " + string(v) + "\n")
	}
	return buffer.String()
}

func (s *InMemoryStorage) Stats() string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("Total objects: %d \n", len(s.storageMap)))
	return buffer.String()
}

func (s *InMemoryStorage) init() {
	if s.storageMap == nil {
		s.storageMap = make(map[string][]byte)
	}
}
