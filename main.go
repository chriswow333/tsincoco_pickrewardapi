package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	psql "pickrewardapi/internal/pkg/postgres"

	bankApplication "pickrewardapi/internal/application/bank/v1"

	cardApplication "pickrewardapi/internal/application/card/v1"
	cardService "pickrewardapi/internal/domain/card/service"
	cardStore "pickrewardapi/internal/domain/card/store"

	cardRewardApplication "pickrewardapi/internal/application/card_reward/v1"

	channelApplication "pickrewardapi/internal/application/channel/v1"
	channelService "pickrewardapi/internal/domain/channel/service"
	channelStore "pickrewardapi/internal/domain/channel/store"

	imageApplication "pickrewardapi/internal/application/image/v1"

	payApplication "pickrewardapi/internal/application/pay/v1"
	payService "pickrewardapi/internal/domain/pay/service"
	payStore "pickrewardapi/internal/domain/pay/store"
)

func loadEnvFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func initLogger() {

	environment := os.Getenv("ENV")

	if environment == "prod" {
		log.SetFormatter(&log.JSONFormatter{})

		f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to create logfile" + "log.txt")
			panic(err)
		}

		// Output to stderr instead of stdout, could also be a file.
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)

		// defer f.Close()

	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{
			// ForceColors:     true,
			FullTimestamp:   false,
			TimestampFormat: "2006-01-02 15:04:05",
		})

	}

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

}

func findServerPort() int {
	logPos := "[main][findServerPort]"
	serverPort := os.Getenv("APP_SERVER_PORT")
	if serverPort == "" {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Info("Cannot find APP_SERVER_PORT")
		panic(-1)
	}

	port, err := strconv.Atoi(serverPort)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":         logPos,
			"server.port": serverPort,
		}).Error("Cannot parse APP_SERVER_PORT")
		panic(-1)
	}
	return port
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(psql.NewPsql)

	container.Provide(payService.NewPay)
	container.Provide(payStore.NewPay)

	container.Provide(cardService.NewBank)
	container.Provide(cardStore.NewBank)

	container.Provide(cardService.NewCard)
	container.Provide(cardStore.NewCard)

	container.Provide(cardStore.NewCardReward)
	container.Provide(cardService.NewCardReward)

	container.Provide(cardService.NewFeedbackType)
	container.Provide(cardStore.NewFeedbackType)

	container.Provide(cardService.NewTaskLabel)
	container.Provide(cardStore.NewTaskLabel)

	container.Provide(channelService.NewChannel)
	container.Provide(channelService.NewChannelLabel)
	container.Provide(channelStore.NewChannel)
	container.Provide(channelStore.NewChannelLabel)

	container.Provide(initGrpcServer)
	return container
}

func parseUseTls() bool {
	logPos := "[main][parseUseTls]"

	lowerStr := strings.ToLower(os.Getenv("APP_USE_TLS"))

	if lowerStr == "true" {
		return true
	} else if lowerStr == "false" {
		return false
	} else {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"use.tls": lowerStr,
		}).Error("Cannot parse APP_USE_TLS")
		panic(-1)
	}

}

func initGrpcServer(
	bankService cardService.BankService,
	cardService cardService.CardService,
	channelService channelService.ChannelService,
	channelLabelService channelService.ChannelLabelService,
	cardRewardService cardService.CardRewardService,
	payService payService.PayService,

) *grpc.Server {
	logPos := "[main][initGrpcServer]"

	useTls := parseUseTls()

	var s *grpc.Server

	if useTls {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Info("With TLS server")

		// 加载证书和密钥

		certPath := os.Getenv("APP_TLS_CERT_PATH")
		keyPath := os.Getenv("APP_TLS_KEY_PATH")

		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			log.Fatalf("Failed to load certificate: %v", err)
			panic(-1)
		}

		// 创建 gRPC 服务器配置
		serverCreds := credentials.NewServerTLSFromCert(&cert)
		s = grpc.NewServer(grpc.Creds(serverCreds))

	} else {

		log.WithFields(log.Fields{
			"pos": logPos,
		}).Info("Without TLS server")
		// Create gRPC Server
		s = grpc.NewServer()

	}

	bankApplication.NewBankServer(s, bankService)
	cardApplication.NewCardServer(s, cardService)
	cardRewardApplication.NewCardRewardServer(s, cardRewardService)
	channelApplication.NewChannelServer(s, channelService, channelLabelService)
	payApplication.NewPayServer(s, payService)
	imageApplication.NewImageServer(s)

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("Loaded all domain servers.")

	return s

}

func main() {

	logPos := "[main][main]"

	initLogger()

	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Info("Cannot find ENV_FILE env, alternative get .env.dev")

		envFile = ".env.dev"
	}

	if err := loadEnvFromFile(envFile); err != nil {
		panic(err)
	}

	container := buildContainer()

	port := findServerPort()

	if err := container.Invoke(func(s *grpc.Server) {
		// Create gRPC Server
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Info("starting grpc listener")

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Fatalf("net.Listen failed: %s", err)
		}

		log.WithFields(log.Fields{
			"pos": logPos,
		}).Infof("Starting grpc server port: %d", port)

		if err := s.Serve(lis); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Fatalf("s.Serve failed: %s ", err)
		}

	}); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Errorf("container.Invoke failed: %s", err)
		panic(err)
	}

}
