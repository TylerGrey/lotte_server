package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/TylerGrey/lotte_server/lib/mysql"
	"github.com/TylerGrey/lotte_server/model/user"
	userApi "github.com/TylerGrey/lotte_server/service/user"
	"github.com/go-kit/kit/log"
	"github.com/sebest/xff"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	logger.Log("DB_HOST", os.Getenv("RDS_HOSTNAME"), "PORT", os.Getenv("RDS_PORT"))
	logger.Log("REDIS_HOST", os.Getenv("ELASTIC_CACHE_HOST"), "PORT", os.Getenv("ELASTIC_CACHE_PORT"))
}

func main() {
	ctx := context.Background()
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	// DB 초기화
	userDB, err := mysql.InitializeDatabase(os.Getenv("RDS_USER_DB_NAME"))
	if err != nil {
		logger.Log("USER_DB_ERROR", err.Error())
		panic(err)
	}
	userRepo := user.NewUserRepository(userDB)

	// API 생성
	var userService userApi.Service
	{
		userService = userApi.NewService(ctx, logger, userRepo)
		userService = userApi.NewLoggingService(log.With(logger, "api", "user"), userService)
	}

	// Endpoint 생성
	userSignUpEndpoint := userApi.MakeSignUpEndpoint(userService)

	userEndpoints := userApi.Endpoints{
		SignUpEndpoint: userSignUpEndpoint,
	}

	userRoute := userApi.MakeHTTPHandler(userEndpoints, logger)

	xffmw, _ := xff.Default()
	mux := http.NewServeMux()
	mux.Handle("/api/user/", xffmw.Handler(userRoute))
	http.Handle("/", accessControl(mux))

	logger.Log(http.ListenAndServe(*httpAddr, nil))
}

func accessControl(h http.Handler) http.Handler {
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
