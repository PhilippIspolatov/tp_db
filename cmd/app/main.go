package main

import (
	"database/sql"
	"log"

	forumDelivery "github.com/ifo16u375/tp_db/internal/forum/delivery"
	forumRepo "github.com/ifo16u375/tp_db/internal/forum/repository"
	forumUcase "github.com/ifo16u375/tp_db/internal/forum/usecase"
	postDelivery "github.com/ifo16u375/tp_db/internal/post/delivery"
	postRepo "github.com/ifo16u375/tp_db/internal/post/repository"
	postUcase "github.com/ifo16u375/tp_db/internal/post/usecase"
	serviceDelivery "github.com/ifo16u375/tp_db/internal/service/delivery"
	serviceRepo "github.com/ifo16u375/tp_db/internal/service/repository"
	serviceUcase "github.com/ifo16u375/tp_db/internal/service/usecase"
	threadDelivery "github.com/ifo16u375/tp_db/internal/thread/delivery"
	threadRepo "github.com/ifo16u375/tp_db/internal/thread/repository"
	threadUcase "github.com/ifo16u375/tp_db/internal/thread/usecase"
	userDelivery "github.com/ifo16u375/tp_db/internal/user/delivery"
	userRepo "github.com/ifo16u375/tp_db/internal/user/repository"
	userUcase "github.com/ifo16u375/tp_db/internal/user/usecase"
	voteDelivery "github.com/ifo16u375/tp_db/internal/vote/delivery"
	voteRepo "github.com/ifo16u375/tp_db/internal/vote/repository"
	voteUcase "github.com/ifo16u375/tp_db/internal/vote/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

func main() {

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339_nano} latency=${latency_human}\n",
	}))

	conn, _ := sql.Open("postgres", "host=localhost port=5432 dbname=docker sslmode=disable")
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	userRep := userRepo.NewUserRepository(conn)
	forumRep := forumRepo.NewForumRepository(conn)
	threadRep := threadRepo.NewThreadRepository(conn)
	postRep := postRepo.NewPostRepository(conn)
	voteRep := voteRepo.NewVoteRepository(conn)
	serviceRep := serviceRepo.NewServiceRepository(conn)

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
