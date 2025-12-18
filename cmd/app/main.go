package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/torrentxok/order_service/internal/cache"
	"github.com/torrentxok/order_service/internal/config"
	"github.com/torrentxok/order_service/internal/http"
	"github.com/torrentxok/order_service/internal/http/handler"
	kafkaConsumer "github.com/torrentxok/order_service/internal/kafka"
	"github.com/torrentxok/order_service/internal/logger"
	"github.com/torrentxok/order_service/internal/repository"
	"github.com/torrentxok/order_service/internal/service"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.New("debug")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("service starting")

	db, err := repository.NewRepository(cfg.DB, log)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer db.Close()

	orderCache := cache.NewLRUCache(cfg.Cache.Size)

	orderService := service.NewOrderService(db, orderCache, log)

	kafkaReader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
	})
	defer kafkaReader.Close()

	consumer := kafkaConsumer.NewConsumer(kafkaReader, orderService, log)

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Error("kafka consumer stopped with error", zap.Error(err))
		}
	}()

	orderHandler := handler.NewOrderHandler(orderService, log)

	httpServer := http.NewServer(
		":"+cfg.Server.Port,
		orderHandler,
		log,
	)

	httpServer.Start()

	<-ctx.Done()

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown error", zap.Error(err))
	}

	log.Info("service stopped gracefully")
}
