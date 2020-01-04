package epoll

import (
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type EventsCollector struct {
	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

func New() (*EventsCollector, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	return &EventsCollector{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]net.Conn),
	}, nil
}

func (ec *EventsCollector) Add(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := websocketFD(conn)
	err := unix.EpollCtl(ec.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})

	if err != nil {
		return err
	}

	ec.lock.Lock()
	defer ec.lock.Unlock()

	ec.connections[fd] = conn

	if len(ec.connections)%100 == 0 {
		log.Printf("number of connections: %v", len(ec.connections))
	}

	return nil
}

func (ec *EventsCollector) Remove(conn net.Conn) error {
	fd := websocketFD(conn)
	err := unix.EpollCtl(ec.fd, syscall.EPOLL_CTL_DEL, fd, nil)

	if err != nil {
		return err
	}

	ec.lock.Lock()
	defer ec.lock.Unlock()

	delete(ec.connections, fd)

	if len(ec.connections)%100 == 0 {
		log.Printf("number of connections: %v", len(ec.connections))
	}

	return nil
}

func (ec *EventsCollector) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(ec.fd, events, 100)

	if err != nil {
		return nil, err
	}

	ec.lock.RLock()
	defer ec.lock.RUnlock()

	var connections []net.Conn

	for i := 0; i < n; i++ {
		conn := ec.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}

	return connections, nil
}

func websocketFD(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
