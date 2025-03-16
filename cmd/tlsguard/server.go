package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"tlsguard/pkg/cert"

	"github.com/charmbracelet/log"
)

func runServer() {
	proxycfg := flag.Arg(0)

	ports := strings.Split(proxycfg, ":")
	if len(ports) != 2 {
		log.Errorf("invalid usage: expected <secure_port>:<insecure_port>, but got %v", proxycfg)
		return
	}

	sourcePort, err := strconv.Atoi(ports[0])
	if err != nil {
		log.Errorf("invalid usage: secure_port must be a number, but got %v", ports[0])
		return
	}

	targetPort, err := strconv.Atoi(ports[1])
	if err != nil {
		log.Errorf("invalid usage: insecure_port must be a number, but got %v", ports[1])
		return
	}

	listener, err := cert.GenTls(fmt.Sprintf(":%v", sourcePort))
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Server started on port %v", sourcePort)
	for {
		sourceConn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			continue
		}

		go func() {
			log.Infof("Connection established from %v", sourceConn.RemoteAddr())

			log.Info("Connecting to target endpoint...")

			targetConn, err := net.DialTimeout("tcp", fmt.Sprintf(":%v", targetPort), time.Second*15)
			if err != nil {
				log.Error(err)
				return
			}

			log.Info("Connection established")

			err = sourceConn.SetDeadline(time.Now().Add(time.Second * 15))
			if err != nil {
				log.Error(err)
				return
			}

			err = targetConn.SetDeadline(time.Now().Add(time.Second * 15))
			if err != nil {
				log.Error(err)
				return
			}

			go func() {
				for {
					buf := make([]byte, 1024)
					n, err := sourceConn.Read(buf)
					if err != nil {
						log.Error(err)
						return
					}

					err = sourceConn.SetReadDeadline(time.Now().Add(time.Second * 15))
					if err != nil {
						log.Error(err)
						return
					}

					_, err = targetConn.Write(buf[:n])
					if err != nil {
						log.Error(err)
						return
					}

					err = targetConn.SetWriteDeadline(time.Now().Add(time.Second * 15))
					if err != nil {
						log.Error(err)
						return
					}
				}
			}()

			go func() {
				for {
					buf := make([]byte, 1024)
					n, err := targetConn.Read(buf)
					if err != nil {
						log.Error(err)
						return
					}

					err = targetConn.SetReadDeadline(time.Now().Add(time.Second * 15))
					if err != nil {
						log.Error(err)
						return
					}

					_, err = sourceConn.Write(buf[:n])
					if err != nil {
						log.Error(err)
						return
					}

					err = sourceConn.SetWriteDeadline(time.Now().Add(time.Second * 15))
					if err != nil {
						log.Error(err)
						return
					}
				}
			}()

			log.Info("Tunnel initiated")
		}()
	}
}
