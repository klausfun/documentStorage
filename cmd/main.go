package main

import (
	"documentStorage"
	"documentStorage/pkg/handler"
	"documentStorage/pkg/repository"
	"documentStorage/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	redisClient, err := repository.NewRedisClient(
		repository.RedisConfig{
			Addr:     viper.GetString("redis.hostAndPort"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       viper.GetInt("redis.dbname"),
		},
	)
	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}

	repos := repository.NewRepository(db, redisClient)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8081", nil); err != nil {
			logrus.Fatalf("failed to start metrics server: %s", err.Error())
		}
	}()

	srv := new(documentStorage.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
