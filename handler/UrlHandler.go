package handler

import (
	"context"

	"github.com/an-halim/url-shortner-grpc/entity"
	pb "github.com/an-halim/url-shortner-grpc/proto/url_service/v1"
	"github.com/an-halim/url-shortner-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type urlHandler struct {
	pb.UnimplementedUrlServiceServer
	urlService service.UrlService
}

func NewUrlHandler(urlService service.UrlService) *urlHandler {
	return &urlHandler{urlService: urlService}
}

func (h *urlHandler) Create(ctx context.Context, req *pb.CreateUrlRequest) (*pb.CreateUrlResponse, error) {
	shortUrl := entity.Url{
		Original: req.GetOriginal(),
	}

	created, err := h.urlService.Short(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUrlResponse{
		Id:        int32(created.ID),
		Original:  created.Original,
		Short:     created.ShortUrl,
		CreatedAt: timestamppb.New(created.CreatedAt),
	}, nil
}

func (h *urlHandler) GetByShort(ctx context.Context, req *pb.GetUrlRequest) (*pb.Empty, error) {
	url, err := h.urlService.GetByShort(ctx, req.GetShort())
	if err != nil {
		return nil, err
	}

	header := metadata.Pairs("Location", url.Original)
	grpc.SendHeader(ctx, header)

	return &pb.Empty{}, nil
}
