package main

import (
	"database/sql"
	"log"

	"github.com/ifo16u375/tp_db/internal/user/delivery"
	"github.com/ifo16u375/tp_db/internal/user/repository"
	"github.com/ifo16u375/tp_db/internal/user/usecase"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main()  {

	e := echo.New()

	conn, _ := sql.Open("postgres", "host=localhost port=5432 dbname=forum_db sslmode=disable")
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	userRep := repository.NewUserRepository(conn)
	userUCase := usecase.NewUserUsecase(userRep)
	_ = delivery.NewUserHandler(e, userUCase)

	log.Fatal(e.Start(":8000"))
}