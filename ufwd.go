package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Version stores current version number.
const Version = "0.1.1"

// Config structure represents app's configuration.
type Config struct {
	BufSize  uint   // Maximum buffer size.
	Proto    string // Protocol to use (udp, udp4, udp6).
	BindAddr string // Address to bind.
	DestAddr string // Address to connect to.
	Debug    bool   // Enable debug mode?
}

// conf stores global configuration.
var conf = new(Config)

// showVersion when set to true makes the app to print current version and exit.
var showVersion = false

// debugf is a default debugging function. If debug mode is enabled it will be
// replaced with log.Printf function.
var debugf = func(f string, a ...interface{}) {}

func init() {
	log.SetFlags(0)

	flag.UintVar(&conf.BufSize, "b", 1024, "buffer maximum size")
	flag.StringVar(&conf.Proto, "p", "udp4", "which protocol to use (udp4 or udp6)")
	flag.BoolVar(&conf.Debug, "d", false, "enable debug mode")
	flag.BoolVar(&showVersion, "v", false, "display version number and exit")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 2 {
		conf.BindAddr = flag.Arg(0)
		conf.DestAddr = flag.Arg(1)
	}
}

func main() {
	// If ran with -v flag, display version number and exit.
	if showVersion {
		fmt.Println(Version)
		return
	}

	// Debug mode enabled? Set log.Printf as debug function.
	if conf.Debug {
		debugf = log.Printf
		debugf("D -- debug mode enabled")
	}

	// Validate buffer size, bind address and destination.
	if conf.BufSize == 0 {
		log.Fatalf("E -- bufsize must be greater than 0")
	}
	if conf.BindAddr == "" {
		log.Fatalf("E -- missing bind address")
	}
	if conf.DestAddr == "" {
		log.Fatalf("E -- missing connect address")
	}

	debugf("D -- provided configuration: %+v", conf)

	// Establish connection with destination server.
	conn, err := dstConn(conf)
	if err != nil {
		log.Fatalf("E -- %v", err)
	}
	defer conn.Close()

	// Start local server and forwarding loop.
	log.Fatalf("E -- %v", srvLoop(conf, conn))
}

// usage displays application help screen.
func usage() {
	fmt.Print(strings.TrimLeft(Help, "\n"))
}

// dstConn connects with given c.DestAddr and returns established connection
// as io.Wroter. If any error occurs then returned connection will be nil
// and error will be returned as second parameter.
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

// srvLoop starts server bound to address giben in c.BindAddr and listens for
// packets on it. Whenever a packet comes in it's immediately written to the
// writer passed as a second argument. If any error occurs the function exits
// and returns error information.
func srvLoop(c *Config, out io.Writer) error {
	addr, err := net.ResolveUDPAddr(c.Proto, c.BindAddr)
	if err != nil {
		return err
	}
	sock, err := net.ListenUDP(c.Proto, addr)
	if err != nil {
		return err
	}

	defer func() {
		log.Printf("I -- shutting down the server")
		sock.Close()
	}()

	log.Printf("I -- about to start a bridge from %s to %s", c.BindAddr, c.DestAddr)

	var buf []byte = make([]byte, c.BufSize)
	for {
		n, err := sock.Read(buf[0:])
		if err != nil {
			return err
		}

		debugf("D -- received packet of size %d", n)
		debugf("D -- > %s", buf[:n])

		w, err := out.Write(buf[:n])
		if err != nil {
			return err
		}

		debugf("D -- forwarded packet of size %d", w)

		if n != w {
			err = fmt.Errorf("corrupted packet, written %d out of %d bytes", w, n)
			return err
		}
	}

	return nil
}

const Help = `
Usage:

  ufwd [-b SIZE] [-d] [-h] [-p PROTOCOL] [-v] BIND CONNECT

Start UDP forwaring server on address BIND that forwards all incoming
connections to address CONNECT.

Options:

  -b SIZE      Maximum buffer size in bytes. Default: 1024.
  -d           Enables debug mode.
  -h           Shows help screen.
  -p PROTOCOL  Which protocol to use. Can be udp, udp4 or udp6. Default: udp4.
  -v           Print version number and exit.

For more information check 'man ufwd'.
`
