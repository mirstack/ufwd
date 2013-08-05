package main

import (
	"net"
	"log"
	"io"
	"fmt"
	"flag"
)

const Version = "0.0.1"

type Config struct {
	BufSize  uint
	Proto    string
	Addr     string
	DestAddr string
	Debug    bool
}

var (
	conf        = new(Config)
	showVersion = false
)

func init() {
	log.SetFlags(0)

	flag.UintVar(&conf.BufSize, "bufsize", 1024, "read/write buffer maximum size")
	flag.StringVar(&conf.Proto, "proto", "udp4", "which protocol to use (udp4 or udp6)")
	flag.StringVar(&conf.Addr, "addr", "", "bind address")
	flag.StringVar(&conf.DestAddr, "dest", "", "destination address")
	flag.BoolVar(&conf.Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&showVersion, "version", false, "display version number and exit")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Printf("ufwd v%s\n", Version)
		return
	}

	if conf.BufSize == 0 {
		log.Fatalf("ERROR: bufsize must be greater than 0")
	}

	conn, err := dstConn(conf)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	defer conn.Close()
	log.Fatalf("ERROR: %v", srvLoop(conf, conn))
}

func debugf(format string, args ...interface{}) {
	if conf.Debug {
		log.Printf("DEBUG: " + format, args...)
	}
}

func dstConn(c *Config) (io.WriteCloser, error) {
	addr, err := net.ResolveUDPAddr(c.Proto, c.DestAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(c.Proto, nil, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func srvLoop(c *Config, out io.Writer) error {
	addr, err := net.ResolveUDPAddr(c.Proto, c.Addr)
	if err != nil {
		return err
	}
	sock, err := net.ListenUDP(c.Proto, addr)
	if err != nil {
		return err
	}

	defer func() {
		log.Printf("INFO: shutting down the server")
		sock.Close()
	}()

	log.Printf("INFO: about to start bridge from %s to %s", c.Addr, c.DestAddr)
	var buf []byte = make([]byte, c.BufSize)

	for {
		n, err := sock.Read(buf[0:])
		if err != nil {
			return err
		}
		debugf("received packet of size %d", n)
		debugf("> %s", buf[:n])

		w, err := out.Write(buf[:n])
		if err != nil {
			return err
		}
		debugf("passed packet of size %d", w)

		if n != w {
			err = fmt.Errorf("corrupted packet, written %d out of %d bytes", w, n)
			return err
		}
	}

	return nil
}