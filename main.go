package main

import (
	"flag"
	"fmt"
	"mqtt-mock-aws-iot/lib/logger"
	"mqtt-mock-aws-iot/lib/mock"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

var (
	PemCa    = flag.String("pem", "", "iot server ca pem")
	PemCrt   = flag.String("crt", "", "mqtt client pem.crt")
	PemKey   = flag.String("key", "", "mqtt client pem.key")
	EndPoint = flag.String("endpoint", "", "connection endpoint")
	DataNum  = flag.Int("datanum", 0, "generate or use datanum")
)

func main() {
	flag.Parse()
	var entities []mock.MockEntity
	if *PemCa == "" || *PemCrt == "" || *PemKey == "" || *EndPoint == "" {
		panic("All Arguments are required.")
	}
	entities = make([]mock.MockEntity, 0, 2)
	entities = append(entities, &mock.MockEntitySample{ClientId: "111"})
	entities = append(entities, &mock.MockEntitySample{ClientId: "222"})
	logger := createLogger()
	mc := mock.NewMockCluster(*PemCa, *PemCrt, *PemKey, *EndPoint, logger)
	mc.AddParallel(entities)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	logger.Info(fmt.Sprintf("SIGNAL %d received", <-quit))
	mc.Close()
	time.Sleep(3 * time.Second)
}

func createLogger() logger.LoggerIf {
	return logger.NewLoggerZerolog(func(l *zerolog.Logger) *zerolog.Logger {
		mod := l.Level(zerolog.InfoLevel)
		return &mod
	})
}
