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

	testName := "First"
	t.Run(testName, func(t *testing.T) {
		rand10 := "abcdefgjik"
		pShort := pb.ShortLink{ShortLink: rand10}
		pLong := pb.LongLink{LongLink: "longlink"}
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
}

func TestMockPostExistLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseHandle := NewMockBaseHandle(ctrl)
	tests := []struct {
		testName string
		rand10   string
		data     database.Mapping
		s        Server
		pShort   pb.ShortLink
		pLong    pb.LongLink
	}{
		{
			testName: "First",
			s:        Server{Storage: "database", HandleDB: baseHandle},
			pShort:   pb.ShortLink{ShortLink: "shortlink1"},
			pLong:    pb.LongLink{LongLink: "longlink1"},
		}, {
			testName: "Second",
			s:        Server{Storage: "database", HandleDB: baseHandle},
			pShort:   pb.ShortLink{ShortLink: "shortlink2"},
			pLong:    pb.LongLink{LongLink: "longlink2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			baseHandle.EXPECT().Find(&tt.data, "long=?", tt.pLong.LongLink).SetArg(0, database.Mapping{Short: tt.pShort.ShortLink, Long: tt.pLong.LongLink}).Return()
			result, _ := tt.s.Post(context.Background(), &tt.pLong)
			if result.ShortLink != tt.pShort.ShortLink {
				t.Errorf("want: %v, got: %v\n", tt.pShort.ShortLink, result.ShortLink)
			}
		})
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
	// data := database.Mapping{}
	pShort := pb.ShortLink{ShortLink: "shortlink"}
	pLong := pb.LongLink{LongLink: "longlink"}

	gomock.InOrder(
		baseHandle.EXPECT().Find(&database.Mapping{}, "long=?", pLong.LongLink).Return(),

		randomable.EXPECT().String10().Return(pShort.ShortLink),
		baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some long link"}, "short=?", pShort.ShortLink).Return(),

		randomable.EXPECT().String10().Return(pShort.ShortLink),
		baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some long link"}, "short=?", pShort.ShortLink).Return(),

		randomable.EXPECT().String10().Return(rand10),
		baseHandle.EXPECT().Find(&database.Mapping{}, "short=?", rand10).Return(),

		baseHandle.EXPECT().Create(database.Mapping{Short: rand10, Long: pLong.LongLink}).Return(),
	)
	result, _ := s.Post(context.Background(), &pLong)
	if result.ShortLink != rand10 {
		t.Errorf("want: %v, got: %v\n", rand10, result.ShortLink)
	}
}
