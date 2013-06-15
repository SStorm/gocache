all: 
	go clean
	rm -fr bin out pkg
	go build github.com/sstorm/gocache/store
	go install github.com/sstorm/gocache

clean:
	go clean
	rm -fr bin out pkg
