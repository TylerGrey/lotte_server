package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"

	"github.com/TylerGrey/lotte_server/db"
	"github.com/TylerGrey/lotte_server/lib/mysql"
	boardApi "github.com/TylerGrey/lotte_server/service/board"
	userApi "github.com/TylerGrey/lotte_server/service/user"
	"github.com/go-kit/kit/log"
	"github.com/sebest/xff"
)

var logger log.Logger

func init() {
	// Logger 초기화
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	logger.Log("DB_HOST", os.Getenv("RDS_HOSTNAME"), "PORT", os.Getenv("RDS_PORT"))
	logger.Log("REDIS_HOST", os.Getenv("ELASTIC_CACHE_HOST"), "PORT", os.Getenv("ELASTIC_CACHE_PORT"))
}

func main() {
	ctx := context.Background()
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")

	// DB 초기화
	dbCli, err := mysql.InitializeDatabase(os.Getenv("RDS_DB_NAME"))
	if err != nil {
		logger.Log("USER_DB_ERROR", err.Error())
		panic(err)
	}
	userRepo := db.NewUserRepository(dbCli)
	boardRepo := db.NewBoardRepository(dbCli)

	// API 설정
	var userService userApi.Service
	{
		userService = userApi.NewService(ctx, logger, userRepo)
		userService = userApi.NewLoggingService(log.With(logger, "api", "user"), userService)
	}

	var boardService boardApi.Service
	{
		boardService = boardApi.NewService(ctx, logger, boardRepo)
		boardService = boardApi.NewLoggingService(log.With(logger, "api", "board"), boardService)
	}

	// Endpoint 설정
	userSignUpEndpoint := userApi.MakeSignUpEndpoint(userService)
	userSignInEndpoint := userApi.MakeSignInEndpoint(userService)

	userEndpoints := userApi.Endpoints{
		SignUpEndpoint: userSignUpEndpoint,
		SignInEndpoint: userSignInEndpoint,
	}

	boardListEndpoint := boardApi.MakeListEndpoint(boardService)
	var boardAddEndpoint endpoint.Endpoint
	{
		boardAddEndpoint = boardApi.MakeAddEndpoint(boardService)
		boardAddEndpoint = boardApi.MakeAuthVerifyMiddleware()(boardAddEndpoint)
	}

	boardEndpoints := boardApi.Endpoints{
		ListEndpoint: boardListEndpoint,
		AddEndpoint:  boardAddEndpoint,
	}

	// 핸들러 설정
	userRoute := userApi.MakeHTTPHandler(userEndpoints, logger)
	boardRoute := boardApi.MakeHTTPHandler(boardEndpoints, logger)

	xffmw, _ := xff.Default()
	mux := http.NewServeMux()
	mux.Handle("/api/user/", xffmw.Handler(userRoute))
	mux.Handle("/api/board/", xffmw.Handler(boardRoute))
	http.Handle("/", accessControl(mux))

	logger.Log(http.ListenAndServe(*httpAddr, nil))
}

func accessControl(h http.Handler) http.Handler {
	// 크로스 도메인 설정
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Access-Control-Allow-Headers, DeviceInfo, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
