package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/TylerGrey/lotte_server/lib/mysql"
	"github.com/go-kit/kit/log"
)

type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = &serializedLogger{Logger: logger}
	}

	logger.Log("RDS_APP_DB_NAME", os.Getenv("RDS_USER_DB_NAME"))

	userDB, dbErr := mysql.InitializeDatabase(os.Getenv("RDS_USER_DB_NAME"))
	if dbErr != nil {
		logger.Log("InitializeDatabase Error", dbErr.Error(), userDB)
		panic(dbErr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", getUsers)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	logger.Log(http.ListenAndServe(":8080", accessControl(mux)))
}

func getUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("{\"hello\":1}")
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
