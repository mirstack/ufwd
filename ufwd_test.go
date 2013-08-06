package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"
)

var msgs = make(chan string)

func init() {
	log.SetOutput(ioutil.Discard)

	conf = &Config{
		BufSize:  1024,
		Proto:    "udp4",
		BindAddr: ":5566",
		DestAddr: ":5567",
		Debug:    true,
	}

	ready := make(chan bool, 1)

	go func() {
		addr, err := net.ResolveUDPAddr(conf.Proto, conf.DestAddr)
		if err != nil {
			panic(err)
		}
		sock, err := net.ListenUDP(conf.Proto, addr)
		if err != nil {
			panic(err)
		}
		var buf []byte = make([]byte, conf.BufSize)
		ready <- true
		for {
			n, err := sock.Read(buf[0:])
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}
			msgs <- string(buf[:n])
		}
	}()

	<-ready

	go main()
	<-time.After(100 * time.Millisecond)
}

func TestServer(t *testing.T) {
	addr, err := net.ResolveUDPAddr(conf.Proto, conf.BindAddr)
	if err != nil {
		t.Errorf("Expected to resolve the addr, got error: %v", err)
		return
	}
	conn, err := net.DialUDP(conf.Proto, nil, addr)
	if err != nil {
		t.Errorf("Expected to connect with the server, got error: %v", err)
		return
	}

	messages := []string{
		"hello world!",
		"another message to make sure it works",
		"and one more time",
	}

	for _, msg := range messages {
		_, err = conn.Write([]byte(msg))
		if err != nil {
			t.Errorf("Expected to write data, got error: %v", err)
			continue
		}

		select {
		case rmsg := <-msgs:
			if rmsg != msg {
				t.Errorf("Expected to get '%s', got '%s'", msg, rmsg)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("Expected to get a message, got nothing")
		}
	}
}
