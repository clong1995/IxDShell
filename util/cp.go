package util

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

func CurrDir() (string, error) {
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return currDir, nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
