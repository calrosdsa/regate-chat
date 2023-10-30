package usecase

import (
	"context"
	"fmt"
	"log"
	r "message/domain/repository"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


type utilUseCase struct {}

func NewUseCase()r.UtilUseCase{
	return &utilUseCase{}
}

func (u *utilUseCase)SendMessageToKafka(w *kafka.Writer,data interface{},key string){
	json, err := json.Marshal(data)
		if err != nil {
			log.Println("Fail to parse", err)
		}
		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(key),
				Value: json,
			},
		)
		if err != nil {
			u.LogError("SendMessageToKafka","util_usecase",err.Error())
		}
}

func (u *utilUseCase)LogError(method string,file string,err string){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		logrus.Error("error opening file: %v", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Error(err)
}

func (u *utilUseCase)LogInfo(method string,file string,err string){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		logrus.Fatalf("error opening file: %v", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Info(err)
}


func (u *utilUseCase)CustomLog(method string,file string,err string,payload map[string]interface{}){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		log.Println("error opening file", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
		"extra":payload,
    })
    ctx.Error(err)
// 	if u.logger != nil {
// 	ctx := u.logger.WithFields(logrus.Fields{
// 		"method": method,
// 		"file":file,
// 		"extra":payload,
//     })
//     ctx.Error(err)
// }
}

func (u *utilUseCase)LogFatal(method string,file string,err string,payload map[string]interface{}){
	now := time.Now()
	t := fmt.Sprintf("%slog/%s-%s-%s", viper.GetString("path"),strconv.Itoa(now.Year()),now.Month().String(),strconv.Itoa(now.Day()))
	f, errL := os.OpenFile(t, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if errL != nil {
		log.Println("error opening file", errL)
	}
	logrus.SetOutput(f)
	defer func ()  {
		log.Println("closing file")
		f.Close()	
	} ()
	ctx := logrus.WithFields(logrus.Fields{
		"method": method,
		"file":file,
    })
    ctx.Error(err)
}


func (u *utilUseCase)PaginationValues(p int16)(page int16){
	if p == 1 || p == 0 {
		page = 0
	} else {
		page = p - 1
	}
	return
}


func (h *utilUseCase)GetNextPage(results int8,pageSize int8,page int16) (nextPage int16){
	if results == 0{
	   nextPage = 0
   }else if results != pageSize{
	   nextPage = 0
   } else{
	   nextPage = page + 1
   }
   return
}