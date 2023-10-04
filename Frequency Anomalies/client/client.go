package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"military/foo"
	"military/transmitter"
	"os"
	"sync"
	"time"
)


type Anomalie struct {
	Id         uint32
	Session_id string
	Frequency  float64
	Time       time.Time
}

func Base() *gorm.DB {
	dsn := "host=localhost user=USERNAME password=PASSWORD dbname=DBNAME port=PORT sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&Anomalie{})
	return db
}

func main() {
	os.Remove("report.txt")
	var k float64
	flag.Float64Var(&k, "k", 3.5, "-k 'float number' for apply coefficient")
	flag.Parse()
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	client := transmitter.NewTransmitterClient(conn)
	tr := transmitter.TransmitterRequest{}
	stream, err := client.Transmitter(context.Background(), &tr)
	if err != nil {
		log.Fatalln(err)
	}

	var pool = sync.Pool{
		New: func() interface{} { return []float64{} },
	}
	var arrf = []float64{}
	var mean, sd float64
	db := Base()
	for i := 1; ; i++ {
		arrf = pool.Get().([]float64)
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		arrf = append(arrf, resp.GetFrequency())
		if i < 1001 {
			mean, sd = foo.Job(int32(i), k, arrf)
			if i == 1000 {
				foo.SaveData(int32(i+1), mean, sd)
			}
		} else {
			if foo.Anomalies(resp.GetFrequency(), mean, sd*k) {
				db.Create(&Anomalie{Session_id: resp.GetSessionId(), Frequency: resp.GetFrequency(), Time: resp.GetTimestamp().AsTime()})
			}
		}
		pool.Put(arrf)
	}
}
