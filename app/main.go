package main

import (
	// "net"
	// "strconv"
	"database/sql"
	
	"fmt"

	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	core "message/core"
)

func init() {
	viper.SetConfigFile(`/home/regate/regate-services/message/app/config.json`)
	// viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Println(loc)
	}
	host := viper.GetString(`database.host`)
	port := viper.GetInt(`database.port`)
	user := viper.GetString(`database.user`)
	password := viper.GetString(`database.pass`)
	dbname := viper.GetString(`database.name`)
	time.Local = loc
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()
	core.Init(db)

}
