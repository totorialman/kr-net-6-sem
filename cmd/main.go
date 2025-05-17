package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/totorialman/kr-net-6-sem/internal/consts"
	"github.com/totorialman/kr-net-6-sem/internal/handlers"
	"github.com/totorialman/kr-net-6-sem/internal/kafka"
	"github.com/totorialman/kr-net-6-sem/internal/storage"
	"github.com/totorialman/kr-net-6-sem/internal/utils"
)

func main() {
	// запуск consumer-а
	go func() {
		maxRetries := 5
		retryInterval := 5 * time.Second
	
		for i := 0; i < maxRetries; i++ {
			if err := kafka.ReadFromKafka(); err != nil {
				fmt.Println(err)
				// Если ошибка, ждем 5 секунд и пробуем снова
				time.Sleep(retryInterval)
			} else {
				// Если подключение успешно
				break
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(consts.KafkaReadPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				storage.ScanStorage(utils.SendReceiveRequest)
			}
		}
	}()

	// создание роутера
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	r.HandleFunc("/receive", handlers.HandleTransfer).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/send", handlers.HandleSend).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// запуск http сервера
	srv := http.Server{
		Handler:           r,
		Addr:              ":8010",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Server stopped")
		}
	}()
	fmt.Println("Server started")

	// graceful shutdown
	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}
}
