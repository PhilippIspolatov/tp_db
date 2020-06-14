package main

import (
	"log"
	"github.com/PhilippIspolatov/tp_db/internal/db"
	"github.com/sirupsen/logrus"

	forumDelivery "github.com/PhilippIspolatov/tp_db/internal/forum/delivery"
	forumRepo "github.com/PhilippIspolatov/tp_db/internal/forum/repository"
	forumUcase "github.com/PhilippIspolatov/tp_db/internal/forum/usecase"
	postDelivery "github.com/PhilippIspolatov/tp_db/internal/post/delivery"
	postRepo "github.com/PhilippIspolatov/tp_db/internal/post/repository"
	postUcase "github.com/PhilippIspolatov/tp_db/internal/post/usecase"
	serviceDelivery "github.com/PhilippIspolatov/tp_db/internal/service/delivery"
	serviceRepo "github.com/PhilippIspolatov/tp_db/internal/service/repository"
	serviceUcase "github.com/PhilippIspolatov/tp_db/internal/service/usecase"
	threadDelivery "github.com/PhilippIspolatov/tp_db/internal/thread/delivery"
	threadRepo "github.com/PhilippIspolatov/tp_db/internal/thread/repository"
	threadUcase "github.com/PhilippIspolatov/tp_db/internal/thread/usecase"
	userDelivery "github.com/PhilippIspolatov/tp_db/internal/user/delivery"
	userRepo "github.com/PhilippIspolatov/tp_db/internal/user/repository"
	userUcase "github.com/PhilippIspolatov/tp_db/internal/user/usecase"
	voteDelivery "github.com/PhilippIspolatov/tp_db/internal/vote/delivery"
	voteRepo "github.com/PhilippIspolatov/tp_db/internal/vote/repository"
	voteUcase "github.com/PhilippIspolatov/tp_db/internal/vote/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

func main() {

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339_nano} latency=${latency_human}\n",
	}))

	// conn, _ := sql.Open("postgres", "host=localhost port=5432 database=docker user=docker password=docker sslmode=disable")
	// if err := conn.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	db, err := db.NewDataBase("./db.json")

	if err != nil {
		logrus.Info(err)
		return
	}

	defer func() {
		err := db.Conn.Close

		if err != nil {
			logrus.Info(err)
			return
		}
	}()

	userRep := userRepo.NewUserRepository(db.Conn)
	forumRep := forumRepo.NewForumRepository(db.Conn)
	threadRep := threadRepo.NewThreadRepository(db.Conn)
	postRep := postRepo.NewPostRepository(db.Conn)
	voteRep := voteRepo.NewVoteRepository(db.Conn)
	serviceRep := serviceRepo.NewServiceRepository(db.Conn)

	userUCase := userUcase.NewUserUsecase(userRep)
	forumUCase := forumUcase.NewForumUsecase(forumRep)
	threadUCase := threadUcase.NewThreadUsecase(threadRep)
	postUCase := postUcase.NewPostUsecase(postRep)
	voteUCase := voteUcase.NewVoteUsecase(voteRep)
	serviceUCase := serviceUcase.NewServiceUsecase(serviceRep)

	_ = userDelivery.NewUserHandler(e, userUCase, forumUCase)
	_ = forumDelivery.NewForumHandler(e, forumUCase, userUCase)
	_ = threadDelivery.NewThreadHandler(e, threadUCase, userUCase, forumUCase)
	_ = postDelivery.NewPostHandler(e, postUCase, threadUCase, userUCase, forumUCase)
	_ = voteDelivery.NewVoteHandler(e, voteUCase, threadUCase)
	_ = serviceDelivery.NewServiceHandler(e, serviceUCase)

	log.Fatal(e.Start(":5000"))


}