package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/an-halim/url-shortner-grpc/entity"
	"github.com/an-halim/url-shortner-grpc/handler"
	pb "github.com/an-halim/url-shortner-grpc/proto/url_service/v1"
	"github.com/an-halim/url-shortner-grpc/repository"
	"github.com/an-halim/url-shortner-grpc/service"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	grpcServers "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func responseHeaderMatcher(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}

func main() {
	dsn := "postgresql://postgres:root@localhost:5432/golang_advance"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	log.Print("Connected to database : ", gormDB)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// auto migrate
	err = gormDB.AutoMigrate(&entity.Url{})

	if err != nil {
		log.Fatalln(err)
	}

	urlRepo := repository.NewUrlRepository(gormDB)
	urlService := service.NewUrlService(urlRepo, redisClient)
	urlHandler := handler.NewUrlHandler(*urlService)

	grpcServer := grpcServers.NewServer()
	pb.RegisterUrlServiceServer(grpcServer, urlHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Running grpc server in port :50051")
		_ = grpcServer.Serve(lis)
	}()
	time.Sleep(1 * time.Second)

	// run grpc gateway
	conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(responseHeaderMatcher),
	)

	if err := pb.RegisterUrlServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	gwServer := gin.Default()

	gwServer.Group("/*{grpc_gateway}").Any("", gin.WrapH(gwmux))
	log.Println("Running grpc gateway in port :8000")
	_ = gwServer.Run(":8000")
}
