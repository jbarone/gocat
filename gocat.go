package main

import (
    "flag"
    "fmt"
    "io"
    "net"
    "os"
    "path/filepath"
    "strings"
)

func Read(conn net.Conn, quit chan<- bool) {
    _, err := io.Copy(conn, os.Stdin)
    if err != nil {
        quit <- true
    }
    quit <- true
}

func Write(conn net.Conn, quit chan<- bool) {
    _, err := io.Copy(os.Stdout, conn)
    if err != nil {
        quit <- true
    }
    quit <- true
}

func Dial(host string, port int) error {

    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        return err
    }
    defer conn.Close()

    quit := make(chan bool)

    go Read(conn, quit)
    go Write(conn, quit)

    <- quit

    return err
}

func main() {
    flag.Usage = func() {
        fmt.Println(
            strings.Replace(
                "usage: $name [-options] hostname port",
                "$name",
                filepath.Base(os.Args[0]),
                -1))
    }
    flag.Parse()

    if flag.NArg() != 2 {
        flag.Usage()
        return
    }
    dialPort := 0
    fmt.Sscanf(flag.Arg(1), "%d", &dialPort)
    Dial(flag.Arg(0), dialPort)
}
