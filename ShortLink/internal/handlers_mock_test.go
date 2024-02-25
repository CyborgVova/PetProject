package internal

import (
	"context"
	"shortlink/database"
	pb "shortlink/proto"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMockPostNewEntry(t *testing.T) {
	ctrl1 := gomock.NewController(t)
	ctrl2 := gomock.NewController(t)
	defer ctrl1.Finish()
	defer ctrl2.Finish()

	baseHandle := NewMockBaseHandle(ctrl1)
	randomable := NewMockRandomable(ctrl2)

	testName := "First"
	t.Run(testName, func(t *testing.T) {
		rand10 := "abcdefgjik"
		pLong := pb.LongLink{LongLink: "longlink"}
		pShort := pb.ShortLink{ShortLink: rand10}
		data := database.Mapping{Short: pShort.ShortLink, Long: pLong.LongLink}
		s := Server{Storage: "database", HandleDB: baseHandle, HandleRand: randomable}
		empty := database.Mapping{Short: "", Long: ""}

		baseHandle.EXPECT().Find(&empty, "long=?", data.Long).Return()
		randomable.EXPECT().String10().Return(rand10)
		baseHandle.EXPECT().Find(&empty, "short=?", rand10).Return()
		baseHandle.EXPECT().Create(data).Return()

		result, _ := s.Post(context.Background(), &pLong)
		if result.ShortLink != data.Short {
			t.Errorf("want: %v, got: %v", data.Long, result)
		}
	})

	// test := struct {
	// 	testName string
	// 	rand10   string
	// 	data     database.Mapping
	// 	s        Server
	// 	pLong    proto.LongLink
	// 	pShort   proto.ShortLink
	// }{
	// 	"First",
	// 	"abcdefgjik",
	// 	database.Mapping{Short: "shortlink", Long: "longlink"},
	// 	Server{Storage: "database", DB: nil},
	// 	// proto.LongLink{LongLink: ""},
	// 	// proto.ShortLink{ShortLink: "Empty request field"},
	// 	proto.LongLink{LongLink: "longlink"},
	// 	proto.ShortLink{ShortLink: "shortlink"},
	// }

	// t.Run(testName, func(t *testing.T) {
	// d := database.Mapping{Short: "", Long: ""}
	// baseHandle.EXPECT().Find(&d, "long=?", data.Long).Return()
	// randomable.EXPECT().String10().Return(rand10)
	// baseHandle.EXPECT().Find(&data, "short=?", data.Short).Return()
	// baseHandle.EXPECT().Create(&data).Return()
	// s.Post(context.Background(), &pLong)
	// })
	// if result.ShortLink != test.pShort.ShortLink {
	// 	t.Errorf("want: %v, got: %v", test.data.Long, result)
	// }
}
