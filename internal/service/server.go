package service

import (
	"context"
	"strconv"

	"server-template/internal/biz"
	"server-template/internal/conf"

	pb "server-template/api/server"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerService struct {
	pb.UnimplementedServerServer
	user *biz.UserBiz
	log  *log.Helper
	cfg  *conf.Config
}

func NewServerService(
	uc *biz.UserBiz, logger log.Logger, conf *conf.Config,
) *ServerService {
	return &ServerService{
		user: uc,
		log:  log.NewHelper(log.With(logger, "module", "service")),
		cfg:  conf,
	}
}

func (b *ServerService) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserReply, error) {
	id, err := b.user.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		Id: strconv.FormatInt(id, 10),
	}, nil
}

func (s *ServerService) Ping(
	ctx context.Context, in *emptypb.Empty,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
