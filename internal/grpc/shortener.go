package shortener

import (
	"context"
	errorsEntity "mfuss/internal/entity/errors"
	auth "mfuss/internal/grpc/auth"
	"mfuss/internal/service"
	"mfuss/internal/utilits"
	pb "mfuss/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServer
	service *service.Service
}

func NewGRPCServer(s *service.Service) *GRPCServer {
	return &GRPCServer{service: s}
}

// Ping - обработчик для проверки связи с хранилищем.
func (s GRPCServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Ping(ctx); err != nil {
		return nil, status.Error(codes.Internal, "database ping failed")
	}

	return &emptypb.Empty{}, nil
}

// Short - обработчик для создания короткой ссылки.
func (s GRPCServer) Short(ctx context.Context, req *pb.ShortRequest) (*pb.ShortResponse, error) {
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}
	shortURL, err := s.service.SaveURL(ctx, req.Url, user.String())

	if err != nil {
		_, ok := err.(errorsEntity.URLConflict)

		if !ok {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := utilits.CheckURL(shortURL); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "short error: %v", err)
		}

		return nil, status.Errorf(codes.AlreadyExists, "short error: %v", err)
	} else {
		if err := utilits.CheckURL(shortURL); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "short error: %v", err)
		}
		return &pb.ShortResponse{
			Id:       user.String(),
			Url:      req.Url,
			ShortUrl: shortURL,
		}, nil
	}

}

// Get - обработчик, который получает полную ссылку из id короткой.
func (s GRPCServer) Get(ctx context.Context, l *pb.GetRequest) (*pb.GetResponse, error) {
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}
	if len(l.Id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "id length should be greater than 0")
	}
	url, err := s.service.GetShortURL(ctx, l.Id, user.String())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	if url.IsDeleted {
		return nil, status.Error(codes.Unavailable, "link is deleted")
	}

	return &pb.GetResponse{
		Id:       l.Id,
		Url:      url.Origin,
		ShortUrl: url.ResultURL,
	}, nil
}

// GetLinks - обработчик возвращающий все ссылки принадлежащие текущему пользователю.
func (s GRPCServer) GetLinks(ctx context.Context, _ *emptypb.Empty) (*pb.GetLinksResponse, error) {
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	links := s.service.GetAllURLs(ctx, user.String())

	b := &pb.GetLinksResponse{}
	for _, link := range links {
		b.Links = append(b.Links, &pb.GetLinksResponse_Link{
			Id:       link.ID,
			Url:      link.Origin,
			ShortUrl: link.ResultURL,
		})
	}

	return b, nil
}

// BatchShort - обработчик для создания пачки коротких ссылок.
func (s GRPCServer) BatchShort(ctx context.Context, in *pb.BatchShortRequest) (*pb.BatchShortResponse, error) {
	user, err := auth.GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %v", err)
	}

	res := &pb.BatchShortResponse{}

	for _, link := range in.Links {
		shortURL, err := s.service.SaveURL(ctx, link.Url, user.String())
		if err != nil {
			continue
		}
		res.Links = append(res.Links, &pb.BatchShortResponse_Link{
			Id:            user.String(),
			Url:           link.Url,
			ShortUrl:      shortURL,
			CorrelationId: link.CorrelationId,
		})
	}

	return res, nil
}

// Delete - обработчик для удаления ссылок пользователя.
func (s GRPCServer) Delete(ctx context.Context, b *pb.DeleteRequest) (*emptypb.Empty, error) {

	ids := make([]string, 0)
	for _, link := range b.Ids {
		if len(link) == 0 {
			continue
		}
		ids = append(ids, link)
	}

	s.service.Queue.Write(ids)

	return &emptypb.Empty{}, nil
}

// GetStats - обработчик, который возвращает статистику сервера.
func (s GRPCServer) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.GetStatsResponse, error) {
	stats, err := s.service.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &pb.GetStatsResponse{
		Links: uint64(stats.URLs),
		Users: uint64(stats.Users),
	}, nil
}
