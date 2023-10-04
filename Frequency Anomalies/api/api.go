package api

import (
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"military/transmitter"
	"os"
	"time"
)

type GRPCServer struct {
	transmitter.UnimplementedTransmitterServer
}

func getParams() (int32, float64) {
	sec := rand.NewSource(time.Now().UnixNano())
	r := rand.New(sec)
	return int32(r.Intn(20) - 10), float64(r.Intn(12)+3) / 10
}

func (s *GRPCServer) getFreq(expectedValue int32, standartDiviation float64) float64 {
	return rand.NormFloat64()*standartDiviation + float64(expectedValue)
}

func SaveData(id string, expVal int32, stdDiv float64) {
	file, err := os.OpenFile("requests.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Unable to open file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s %d %.1f\n", id, expVal, stdDiv))
}

func (s *GRPCServer) Transmitter(rq *transmitter.TransmitterRequest, st transmitter.Transmitter_TransmitterServer) error {
	expectedValue, standartDiviation := getParams()
	sessionId := uuid.NewString()
	fmt.Println(sessionId, expectedValue, standartDiviation)
	SaveData(sessionId, expectedValue, standartDiviation)
	for i := 0; ; i++ {
		if err := st.Send(&transmitter.TransmitterResponce{
			SessionId: sessionId,
			Frequency: s.getFreq(expectedValue, standartDiviation),
			Timestamp: timestamppb.New(time.Now())}); err != nil {
			return err
		}
	}
}
