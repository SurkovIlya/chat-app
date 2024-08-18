package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	chatserver "github.com/SurkovIlya/chat-app/internal/chat_server"
	"github.com/SurkovIlya/chat-app/internal/server"
	"github.com/SurkovIlya/chat-app/internal/storage/pg"
	st "github.com/SurkovIlya/chat-app/pkg/postgres"
)

const port = "8080"

func main() {
	pgParams := st.DBParams{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	conn, err := st.Connect(pgParams)
	if err != nil {
		panic(err)
	}

	storage := pg.New(st.New(conn))

	chs := chatserver.New(storage)

	srv := server.New(port, chs, storage)

	go func() {
		if err := srv.Run(); err != nil {
			log.Panicf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Println("app Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("app Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Panicf("error occured on server shutting down: %s", err.Error())
	}
}
