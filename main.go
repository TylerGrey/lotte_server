package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"

	"github.com/TylerGrey/lotte_server/db"
	"github.com/TylerGrey/lotte_server/lib/mysql"
	reservationApi "github.com/TylerGrey/lotte_server/service/reservation"
	roomApi "github.com/TylerGrey/lotte_server/service/room"
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
	roomRepo := db.NewRoomRepository(dbCli)
	reservationRepo := db.NewReservationRepository(dbCli)

	// API 설정
	var userService userApi.Service
	{
		userService = userApi.NewService(ctx, logger, userRepo)
		userService = userApi.NewLoggingService(log.With(logger, "api", "user"), userService)
	}

	var roomService roomApi.Service
	{
		roomService = roomApi.NewService(ctx, logger, roomRepo)
		roomService = roomApi.NewLoggingService(log.With(logger, "api", "room"), roomService)
	}

	var reservationService reservationApi.Service
	{
		reservationService = reservationApi.NewService(ctx, logger, reservationRepo)
		reservationService = reservationApi.NewLoggingService(log.With(logger, "api", "reservation"), reservationService)
	}

	// Endpoint 설정
	userSignUpEndpoint := userApi.MakeSignUpEndpoint(userService)
	userSignInEndpoint := userApi.MakeSignInEndpoint(userService)
	userListEndpoint := userApi.MakeListEndpoint(userService)

	userEndpoints := userApi.Endpoints{
		SignUpEndpoint: userSignUpEndpoint,
		SignInEndpoint: userSignInEndpoint,
		ListEndpoint:   userListEndpoint,
	}

	roomListEndpoint := roomApi.MakeListEndpoint(roomService)
	var roomAddEndpoint endpoint.Endpoint
	{
		roomAddEndpoint = roomApi.MakeAddEndpoint(roomService)
		roomAddEndpoint = roomApi.MakeAuthVerifyMiddleware()(roomAddEndpoint)
	}

	roomEndpoints := roomApi.Endpoints{
		ListEndpoint: roomListEndpoint,
		AddEndpoint:  roomAddEndpoint,
	}

	reservationListEndpoint := reservationApi.MakeListEndpoint(reservationService)
	var reservationAddEndpoint endpoint.Endpoint
	{
		reservationAddEndpoint = reservationApi.MakeAddEndpoint(reservationService)
		reservationAddEndpoint = reservationApi.MakeAuthVerifyMiddleware()(reservationAddEndpoint)
	}

	reservationEndpoints := reservationApi.Endpoints{
		ListEndpoint: reservationListEndpoint,
		AddEndpoint:  reservationAddEndpoint,
	}

	// 핸들러 설정
	userRoute := userApi.MakeHTTPHandler(userEndpoints, logger)
	roomRoute := roomApi.MakeHTTPHandler(roomEndpoints, logger)
	reservationRoute := reservationApi.MakeHTTPHandler(reservationEndpoints, logger)

	xffmw, _ := xff.Default()
	mux := http.NewServeMux()
	mux.Handle("/api/user/", xffmw.Handler(userRoute))
	mux.Handle("/api/room/", xffmw.Handler(roomRoute))
	mux.Handle("/api/reservation/", xffmw.Handler(reservationRoute))
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
