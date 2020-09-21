package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

//Server struct definition
type Server struct {
	listener         net.Listener
	quit             chan struct{}
	exited           chan struct{}
	db               runDB
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
}

//NewServer function to create a new instance
func NewServer() *Server {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("--- Failed to create listener", err.Error())
	}
	srv := &Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               createDB(),
		connections:      map[int]net.Conn{},
		connCloseTimeout: 10 * time.Second,
	}

	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	var id int
	fmt.Println("--- Listening for clients")
	for {
		select {
		case <-srv.quit:
			fmt.Println("--- Shutting down the server")
			err := srv.listener.Close()
			if err != nil {
				fmt.Println("--- Error in closing the listener", err.Error())
			}
			if len(srv.connections) > 0 {
				srv.warnConnections(srv.connCloseTimeout)
				<-time.After(srv.connCloseTimeout)
				srv.closeConnections()
			}
			close(srv.exited)
			return
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Println("--- Failed to set listener deadline", err.Error())
			}

			conn, err := tcpListener.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("--- Failed to accept connection", err.Error())
			}

			write(conn, "--- Welcome to RuntimeDB server")

			srv.connections[id] = conn

			go func(connID int) {
				fmt.Println("--- Client with ID: ", connID, " joined")
				srv.handleConn(conn)

				delete(srv.connections, connID)
				fmt.Println("--- Client with id: ", connID, " left")
			}(id)
			id++
		}
	}
}

//--- Write message to the terminal
func write(con net.Conn, s string) {
	_, err := fmt.Fprint(con, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		switch {
		case len(values) == 3 && values[0] == "set":
			srv.db.set(values[1], values[2])
			write(conn, "--- OK, inserted the data")

		case len(values) == 2 && values[0] == "get":
			val, found := srv.db.get(values[1])
			if !found {
				write(conn, fmt.Sprintf("--- Key %s not found"))
			} else {
				write(conn, val)
			}
		case len(values) == 2 && values[0] == "delete":
			srv.db.delete(values[1])
			write(conn, "--- OK, deleted the data")
		case len(values) == 1 && values[0] == "exit":
			if err := conn.Close(); err != nil {
				fmt.Println("--- Could'nt close connection", err.Error())
			}
		default:
			write(conn, fmt.Sprintf("--- UNKNOWN: %s", l))
		}
	}
}

//--- Warn about closing the server
func (srv *Server) warnConnections(timeout time.Duration) {
	for _, conn := range srv.connections {
		write(conn, fmt.Sprintf("--- Host wants to shut down the server in %s", srv.connCloseTimeout.String()))
	}
}

//--- Close the connections
func (srv *Server) closeConnections() {
	fmt.Println("--- Closing all connections")
	for id, conn := range srv.connections {
		err := conn.Close()
		if err != nil {
			fmt.Println("--- Couldn't close connection with ID: ", id, err.Error())
		}
	}
}

func (srv *Server) Stop() {
	fmt.Println("--- Stopping the database server")
	close(srv.quit)
	<-srv.exited
	fmt.Println("--- Saving records to file")
	srv.db.save()
	fmt.Println("--- Database server successfully stopped")
}
