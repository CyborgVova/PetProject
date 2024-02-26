package internal

import (
	"context"
	"shortlink/database"
	pb "shortlink/proto"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMockPostNewEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseHandle := NewMockBaseHandle(ctrl)
	randomable := NewMockRandomable(ctrl)

	rand10 := "abcdefgjik"
	pShort := pb.ShortLink{ShortLink: rand10}
	pLong := pb.LongLink{LongLink: "longlink"}
	data := &database.Mapping{Short: pShort.ShortLink, Long: pLong.LongLink}
	s := Server{Storage: "database", HandleDB: baseHandle, HandleRand: randomable}
	empty := &database.Mapping{Short: "", Long: ""}

	baseHandle.EXPECT().Find(empty, "long=?", data.Long).Return(empty)
	randomable.EXPECT().String10().Return(rand10)
	baseHandle.EXPECT().Find(empty, "short=?", rand10).Return(empty)
	baseHandle.EXPECT().Create(data).Return(data)

	result, _ := s.Post(context.Background(), &pLong)
	if result.ShortLink != data.Short {
		t.Errorf("want: %v, got: %v", data.Long, result)
	}
}

func TestMockPostExistLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseHandle := NewMockBaseHandle(ctrl)

	data := &database.Mapping{}
	s := Server{Storage: "database", HandleDB: baseHandle}
	pShort := pb.ShortLink{ShortLink: "shortlink1"}
	pLong := pb.LongLink{LongLink: "longlink1"}

	baseHandle.EXPECT().Find(data, "long=?", pLong.LongLink).Return(&database.Mapping{Short: pShort.ShortLink, Long: pLong.LongLink})
	result, _ := s.Post(context.Background(), &pLong)
	if result.ShortLink != pShort.ShortLink {
		t.Errorf("want: %v, got: %v\n", pShort.ShortLink, result.ShortLink)
	}
}

func TestMockPostBadGenerate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseHandle := NewMockBaseHandle(ctrl)
	randomable := NewMockRandomable(ctrl)

	rand10 := "rutyeiwoqp"
	s := Server{
		Storage:    "database",
		HandleDB:   baseHandle,
		HandleRand: randomable,
	}
	pShort := pb.ShortLink{ShortLink: "shortlink"}
	pLong := pb.LongLink{LongLink: "longlink"}

	gomock.InOrder(
		baseHandle.EXPECT().Find(&database.Mapping{}, "long=?", pLong.LongLink).Return(&database.Mapping{}),

		randomable.EXPECT().String10().Return(pShort.ShortLink),
		baseHandle.EXPECT().Find(&database.Mapping{}, "short=?", pShort.ShortLink).Return(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}),

		randomable.EXPECT().String10().Return(pShort.ShortLink),
		baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}, "short=?", pShort.ShortLink).Return(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}),

		randomable.EXPECT().String10().Return(rand10),
		baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}, "short=?", rand10).Return(&database.Mapping{}),
	)
	baseHandle.EXPECT().Create(&database.Mapping{rand10, pLong.LongLink}).Return(&database.Mapping{rand10, pLong.LongLink})
	result, _ := s.Post(context.Background(), &pLong)
	if result.ShortLink != rand10 {
		t.Errorf("want: %v, got: %v\n", rand10, result.ShortLink)
	}
}
