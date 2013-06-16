package store

import (
	"bytes"
	"fmt"
)

type ChunkStore struct {
	minKey int
	maxKey int
	storage Storage
}

type DistributedStorage struct {
	chunks *[]ChunkStore
}

func NewDistributedStorage() *DistributedStorage {
	chunks := make([]ChunkStore, 2)

	chunk1 := ChunkStore{0, 0xFFFFFFFF/2, NewInMemoryStorage()}
	chunk2 := ChunkStore{0xFFFFFFFF/2+1, 0xFFFFFFFF, NewInMemoryStorage()}

	chunks[0] = chunk1;
	chunks[1] = chunk2;

	s := &DistributedStorage{&chunks}

	return s
}

func (s *DistributedStorage) Set(key string, value []byte, flags, timeout int) {
	hash := hashcode(key)

	fmt.Println("Hashcode: ", hash)
	for _, elem := range(*s.chunks) {
		if elem.minKey < hash && elem.maxKey >= hash {
			elem.storage.Set(key, value, flags, timeout)
			break;
		}
		// TODO Handle if key doesn't fit anywhere...
	}
}

func (s *DistributedStorage) Get(key string) []byte {
	hash := hashcode(key)

	for _, elem := range(*s.chunks) {
		if elem.minKey > hash && elem.maxKey <= hash {
			return elem.storage.Get(key)
		}
	}
	return nil
}

func (s *DistributedStorage) Dump() string {
	return ""
}

func (s *DistributedStorage) Stats() string {

	buffer := bytes.NewBufferString("")

	for _, elem := range(*s.chunks) {
		buffer.WriteString(fmt.Sprintf("ChunkStore[%d, %d]\n", elem.minKey, elem.maxKey))
		buffer.WriteString(elem.storage.Stats())
	}
	return buffer.String()
}

// :) best hashcode ever
func hashcode(key string) int {

	hash := 7
	for pos, char := range key {
		hash = hash + pos + int(char)
	}

	return hash
}
