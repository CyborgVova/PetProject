package internal

import (
	"context"

	"shortlink/database"
	pb "shortlink/proto"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPostEmptyRequest(t *testing.T) {
	tests := []struct {
		name string
		s    Server
	}{
		{
			name: "Inmemory",
			s: Server{
				Storage: "inmemory",
			},
		},
		{
			name: "Database",
			s: Server{
				Storage: "database",
			},
		},
	}
	long := pb.LongLink{LongLink: ""}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := test.s.Post(context.Background(), &long)
			want := "Empty request field"
			if result.ShortLink != want {
				t.Errorf("want: %s, got: %s\n", want, result.ShortLink)
			}
		})
	}
}
func TestMockPostInmemoryNewEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomable := NewMockRandomable(ctrl)

	shortShare := InMemoryMap{
		"short1": "long1",
		"short2": "long2",
		"short3": "long3",
	}
	longShare := InMemoryMap{
		"long1": "short1",
		"long2": "short2",
		"long3": "short3",
	}
	srv := Server{
		Storage:    "inmemory",
		HandleRand: randomable,
		Short:      shortShare,
		Long:       longShare,
	}
	pShort := []pb.ShortLink{
		{ShortLink: "shortlink1"},
		{ShortLink: "shortlink2"},
		{ShortLink: "shortlink3"}}
	pLong := []pb.LongLink{
		{LongLink: "longlink1"},
		{LongLink: "longlink2"},
		{LongLink: "longlink3"}}

	tests := []struct {
		name string
		s    Server
	}{
		{
			name: "First",
			s:    srv,
		}, {
			name: "Second",
			s:    srv,
		}, {
			name: "Third",
			s:    srv,
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			randomable.EXPECT().String10().Return(pShort[i].ShortLink)
			result, _ := test.s.Post(context.Background(), &pLong[i])
			got := result.ShortLink
			want := pShort[i].ShortLink
			if got != want {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}

func TestPostInmemoryExistLink(t *testing.T) {
	shortShare := InMemoryMap{
		"short1": "long1",
		"short2": "long2",
		"short3": "long3",
	}
	longShare := InMemoryMap{
		"long1": "short1",
		"long2": "short2",
		"long3": "short3",
	}
	srv := Server{
		Storage: "inmemory",
		Short:   shortShare,
		Long:    longShare,
	}
	pLong := []pb.LongLink{
		{LongLink: "long1"},
		{LongLink: "long2"},
		{LongLink: "long3"}}

	tests := []struct {
		name string
		s    Server
	}{
		{
			name: "First",
			s:    srv,
		}, {
			name: "Second",
			s:    srv,
		}, {
			name: "Third",
			s:    srv,
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := test.s.Post(context.Background(), &pLong[i])
			got := result.ShortLink
			want := test.s.Long[pLong[i].LongLink]
			if got != want {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}

func TestPostInmemoryBadGenerate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomable := NewMockRandomable(ctrl)

	list := []string{
		"short3",
		"short1",
		"short2",
	}
	shortShare := InMemoryMap{
		"short1": "long1",
		"short2": "long2",
		"short3": "long3",
	}
	longShare := InMemoryMap{
		"long1": "short1",
		"long2": "short2",
		"long3": "short3",
	}
	s := Server{
		Storage:    "inmemory",
		HandleRand: randomable,
		Short:      shortShare,
		Long:       longShare,
	}
	pLong := pb.LongLink{LongLink: "long4"}
	pShort := pb.ShortLink{ShortLink: "short4"}

	for i := 0; i < len(list); i++ {
		randomable.EXPECT().String10().Return(list[i])
	}
	randomable.EXPECT().String10().Return(pShort.ShortLink)

	result, _ := s.Post(context.Background(), &pLong)
	want := longShare[pLong.LongLink]
	got := result.ShortLink
	if got != want {
		t.Errorf("want: %s, got: %s\n", want, got)
	}
}

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

	baseHandle.EXPECT().Find(&database.Mapping{}, "long=?", pLong.LongLink).Return(&database.Mapping{})

	randomable.EXPECT().String10().Return(pShort.ShortLink)
	baseHandle.EXPECT().Find(&database.Mapping{}, "short=?", pShort.ShortLink).Return(&database.Mapping{Short: pShort.ShortLink, Long: "some link"})

	randomable.EXPECT().String10().Return(pShort.ShortLink)
	baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}, "short=?", pShort.ShortLink).Return(&database.Mapping{Short: pShort.ShortLink, Long: "some link"})

	randomable.EXPECT().String10().Return(rand10)
	baseHandle.EXPECT().Find(&database.Mapping{Short: pShort.ShortLink, Long: "some link"}, "short=?", rand10).Return(&database.Mapping{})

	baseHandle.EXPECT().Create(&database.Mapping{Short: rand10, Long: pLong.LongLink}).Return(&database.Mapping{Short: rand10, Long: pLong.LongLink})

	result, _ := s.Post(context.Background(), &pLong)
	if result.ShortLink != rand10 {
		t.Errorf("want: %v, got: %v\n", rand10, result.ShortLink)
	}
}
