package internal

import (
	"context"
	"math/rand"
	"shortlink/database"
	pb "shortlink/proto"
	"time"

	"gorm.io/gorm"
)

var (
	long  = map[string]string{}
	short = map[string]string{}
)

const template = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Randomable interface {
	String10() string
}

type BaseHandle interface {
	Find(dest interface{}, left, right interface{})
	Create(value interface{})
}

type Server struct {
	pb.UnimplementedLinkBuilderServer
	Storage    string
	HandleDB   BaseHandle
	HandleRand Randomable
}

type HandleDB struct {
	DB *gorm.DB
}

type HandleRand struct{}

func (h *HandleDB) Find(dest interface{}, key, value interface{}) {
	h.DB.Find(dest, key, value)
}

func (h *HandleDB) Create(value interface{}) {
	h.DB.Create(value)
}

func (h *HandleRand) String10() string {
	out_string := make([]byte, 10)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		out_string[i] = template[r.Intn(len(template))]
	}
	return string(out_string)
}

func (s *Server) Post(ctx context.Context, in *pb.LongLink) (*pb.ShortLink, error) {
	if in.LongLink == "" {
		return &pb.ShortLink{ShortLink: "Empty request field"}, nil
	}
	if s.Storage == "inmemory" {
		if _, ok := long[in.LongLink]; !ok {
			for {
				tmp := s.HandleRand.String10()
				if _, ok := short[tmp]; !ok {
					short[tmp] = in.LongLink
					long[in.LongLink] = tmp
					break
				}
			}
		}
		return &pb.ShortLink{ShortLink: long[in.LongLink]}, nil
	}
	var result database.Mapping
	s.HandleDB.Find(&result, "long=?", in.LongLink)
	if result.Short == "" {
		for {
			tmp := s.HandleRand.String10()
			s.HandleDB.Find(&result, "short=?", tmp)
			if result.Short == tmp {
				continue
			}
			result = database.Mapping{Short: tmp, Long: in.LongLink}
			s.HandleDB.Create(result)
			break
		}
	}
	return &pb.ShortLink{ShortLink: result.Short}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.ShortLink) (*pb.LongLink, error) {
	if in.ShortLink == "" {
		return &pb.LongLink{LongLink: "Empty request field"}, nil
	}
	if s.Storage == "inmemory" {
		if _, ok := short[in.ShortLink]; !ok {
			return &pb.LongLink{LongLink: "Link is not exist"}, nil
		}
		return &pb.LongLink{LongLink: short[in.ShortLink]}, nil
	}
	var result database.Mapping
	s.HandleDB.Find(&result, "short=?", in.ShortLink)
	if result.Long == "" {
		return &pb.LongLink{LongLink: "Link is not exist"}, nil
	}
	return &pb.LongLink{LongLink: result.Long}, nil
}
