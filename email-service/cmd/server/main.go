package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aitu/food-delivery/email-service/internal/config"
	"github.com/aitu/food-delivery/email-service/internal/infrastructure/nats"
	"github.com/aitu/food-delivery/email-service/internal/infrastructure/smtp"
	"github.com/aitu/food-delivery/email-service/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	smtpClient := smtp.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPEmail, cfg.SMTPPass)
	uc := usecase.NewEmailUsecase(smtpClient)

	subscriber, err := nats.NewSubscriber(cfg.NATSURL, uc)
	if err != nil {
		log.Fatal("nats subscriber:", err)
	}
	defer subscriber.Close()

	log.Println("email-service listening for payment.completed on", cfg.NATSURL)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
