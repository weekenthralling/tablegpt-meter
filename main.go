package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/tablegpt_meter/config"
	"github.com/tablegpt_meter/models"
	"github.com/tablegpt_meter/server"
	"github.com/tablegpt_meter/store"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	ratelimit "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	token "github.com/tablegpt_meter/proto/token"
)

var (
	grpcPort  = flag.String("grpc-port", "50051", "The port for the gRPC server")
	probeAddr = flag.String("probe-address", ":8080", "The address for the health probe")
)

func main() {

	flag.Parse()

	cfg := config.InitConfig()

	var storeImpl store.TokenStore
	ctx := context.Background()

	switch cfg.Type {
	case "redis":
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisConfig.Addr,
			Password: cfg.RedisConfig.Password,
			DB:       cfg.RedisConfig.DB,
		})
		storeImpl = store.NewRedisStore(rdb)
		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Redis connection failed: %v", err)
		}
		log.Printf("Redis connected: %s", pong)

	case "postgres":
		// Create a GORM DB instance using pgx connection pool
		db, err := gorm.Open(postgres.Open(cfg.DBConfig.DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("PostgreSQL connection failed: %v", err)
		}
		storeImpl = store.NewPostgresStore(db)
		log.Println("PostgreSQL connected.")

		// Auto migrate tables
		err = db.AutoMigrate(&models.UserTokens{}, &models.UsedToken{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database tables migrated.")

	default:
		log.Fatalf("Unsupported DB type: %s", cfg.Type)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", *grpcPort)

	grpcServer := grpc.NewServer()

	// Create and register the TokenServiceServer
	tokenServiceServer := server.NewTokenServiceServer(storeImpl)
	token.RegisterTokenServiceServer(grpcServer, tokenServiceServer)

	// Create and register the TokenLimitServiceServer
	tokenLimitServiceServer := server.NewTokenLimitServiceServer(storeImpl)
	ratelimit.RegisterRateLimitServiceServer(grpcServer, tokenLimitServiceServer)

	// Create a simple HTTP server for health checks
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		handler := NewRecordHandler(tokenServiceServer)
		handler(w, r) // Call the handler returned from NewRecordHandler
	})

	// Start the health check server in a new goroutine
	go func() {
		log.Printf("Starting health check server on %s", *probeAddr)
		if err := http.ListenAndServe(*probeAddr, nil); err != nil {
			log.Fatalf("failed to start health check server: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func writeHttpStatus(writer http.ResponseWriter, code int) {
	http.Error(writer, http.StatusText(code), code)
}

func NewRecordHandler(svc token.TokenServiceServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req token.RecordTokenUsageRequest

		ctx := context.Background()

		body, err := io.ReadAll(request.Body)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			writeHttpStatus(writer, http.StatusBadRequest)
			return
		}

		if err := protojson.Unmarshal(body, &req); err != nil {
			log.Fatalf("error: %s", err.Error())
			writeHttpStatus(writer, http.StatusBadRequest)
			return
		}

		resp, err := svc.RecordTokenUsage(ctx, &req)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			writeHttpStatus(writer, http.StatusBadRequest)
			return
		}

		jsonResp, err := protojson.Marshal(resp)
		if err != nil {
			log.Fatalf("error marshaling proto3 to json: %s", err.Error())
			writeHttpStatus(writer, http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(jsonResp)
	}
}
