package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/GeertJohan/go.rice"
	"github.com/fortytw2/embercrest/api"
	"github.com/fortytw2/embercrest/api/util"
	"github.com/fortytw2/embercrest/datastore/pgsql"
	"github.com/fortytw2/embercrest/gc"
	"github.com/fortytw2/embercrest/gc/redisq"
	"github.com/fortytw2/embercrest/web"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	gc := buildGC()

	router.GET("/", web.Homepage)
	router.GET("/docs", web.Docs)

	// public API
	router.GET("/classes", api.Classes(gc))
	router.GET("/tiles", api.Tiles(gc))
	router.GET("/version", api.Version)
	router.GET("/leaderboard", api.Leaderboard(gc))

	// authenticated API
	router.GET("/matches", api.Match(gc))

	router.POST("/queue", api.Queue(gc))
	router.DELETE("/queue", api.CancelQueue(gc))

	router.GET("/game/:id/match", api.Match(gc))
	router.POST("/game/:id/match/concede", api.Concede(gc))
	router.POST("/game/:id/draft", api.Draft(gc))
	router.POST("/game/:id/turn", api.Turn(gc))

	// user management endpoints
	router.POST("/user/new", api.NewUser(gc))
	router.POST("/user/forgot_pass", api.ForgotPass(gc))

	router.NotFound = func(w http.ResponseWriter, r *http.Request) {
		util.JSONError(w, errors.New("404 not found"), http.StatusNotFound)
	}
	// staticfiles
	router.ServeFiles("/static/*filepath", rice.MustFindBox("static").HTTPBox())

	go gc.Matchmake()
	go gc.EnforceTimeLimits()

	// launch the http server
	log.Println("ec: now listening on port", os.Getenv("PORT"))
	err := http.ListenAndServe(os.Getenv("PORT"), httpLogger(router))
	if err != nil {
		panic(err)
	}
}

// cleanly log all HTTP requests
func httpLogger(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		router.ServeHTTP(w, req)
		finishTime := time.Now()
		elapsedTime := finishTime.Sub(startTime)
		log.Println(req.Method, req.URL, elapsedTime)
	})
}

func buildGC() *gc.GC {
	ds, err := pgsql.NewDatastore()
	if err != nil {
		log.Fatalln("embercrest: fatal error building datastore", err)
	}

	return &gc.GC{
		Datastore: ds,
		Q:         redisq.NewQueue(getRedisPool()),
	}
}

func getRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("REDIS_PORT"))
			if err != nil {
				log.Fatalln("ec: failed to connect to redis @", os.Getenv("REDIS_PORT"))
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
