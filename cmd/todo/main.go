package main

import (
	"flag"
	"log"

	"github.com/AlekseyAytov/skillrock-todo/internal/app"
	"github.com/AlekseyAytov/skillrock-todo/internal/models/master"
	"github.com/AlekseyAytov/skillrock-todo/store/pg"
)

var (
	flagServerSocket string // адрес и номер порта хоста на котором запущен проект
	flagDatabaseDSN  string // строка для подключения к базе данных
)

func init() {
	flag.StringVar(&flagServerSocket, "a", ":8091", "host:port of target HTTP address")
	flag.StringVar(&flagDatabaseDSN, "d", "postgres://{user}:{password}@{hostname}:{port}/{database-name}", "connection string to database")
}

func main() {
	flag.Parse()

	store, err := pg.NewDBStorage(flagDatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	app := app.NewToDoAPI(master.NewTaskMaster(store))
	log.Fatal(app.StartServer(flagServerSocket))
}
