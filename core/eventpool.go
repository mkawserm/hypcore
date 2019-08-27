package core

import (
	"golang.org/x/sys/unix"
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"
)

type EventPool struct {
	fd            int
	connectionMap map[int]net.Conn

	lock *sync.RWMutex

	eventQueueSize int
	waitingTime    int
}

func MakeEventPool() (*EventPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventPool{
		fd:             fd,
		lock:           &sync.RWMutex{},
		connectionMap:  make(map[int]net.Conn),
		eventQueueSize: 100,
		waitingTime:    100,
	}, nil

}

func MakeCustomEventPool(eventQueueSize int, waitingTime int) (*EventPool, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventPool{
		fd:             fd,
		lock:           &sync.RWMutex{},
		connectionMap:  make(map[int]net.Conn),
		eventQueueSize: eventQueueSize,
		waitingTime:    waitingTime,
	}, nil

}

func (e *EventPool) GetConnection(cid int) (net.Conn, bool) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	v, ok := e.connectionMap[cid]
	return v, ok
}

func (e *EventPool) AddConnection(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.fd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		})

	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.connectionMap[fd] = conn
	if len(e.connectionMap)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connectionMap))
	}

	return nil
}

func (e *EventPool) RemoveConnection(conn net.Conn) error {
	fd := WebsocketFileDescriptor(conn)

	err := unix.EpollCtl(e.fd,
		syscall.EPOLL_CTL_DEL,
		fd, nil)

	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	delete(e.connectionMap, fd)

	if len(e.connectionMap)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connectionMap))
	}

	return nil
}

func (e *EventPool) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, e.eventQueueSize)
	n, err := unix.EpollWait(e.fd, events, e.waitingTime)

	if err != nil {
		return nil, err
	}

	e.lock.RLock()
	defer e.lock.RUnlock()

	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := e.connectionMap[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func (e *EventPool) TotalActiveWebSocketConnections() int {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return len(e.connectionMap)
}

func WebsocketFileDescriptor(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
