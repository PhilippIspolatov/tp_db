package delivery

import (
	"net/http"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/user"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUcase user.Usecase
}

func NewUserHandler(router *echo.Echo, uUC user.Usecase) *UserHandler {
	uh := &UserHandler{
		userUcase: uUC,
	}

	router.GET("/user/:nickname/profile", uh.GetUser())
	router.POST("/user/:nickname/create", uh.CreateUser())
	router.POST("/user/:nickname/profile", uh.UpdateUser())

	return uh
}

func (uh *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		u, err := uh.userUcase.GetUser(nickname)

		if err != nil {
			logrus.Info(err)
			return c.JSON(http.StatusNotFound, tools.Error{
				ErrMsg: tools.ErrNotFound.Error(),
			})
		}
		return c.JSON(http.StatusOK, u)
	}
}

func (uh *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		u := &models.User{Nickname:nickname}
		if err := c.Bind(u); err != nil {
			logrus.Info(err)
			return err
		}

		if ConflictData, err := uh.userUcase.CreateUser(u); err != nil {
			logrus.Info(err)
			if err == tools.ErrConflict {
				return c.JSON(http.StatusConflict, ConflictData)
			}
			return c.JSON(http.StatusBadRequest, tools.Error{
				ErrMsg: tools.BadRequest.Error(),
			})
		}
		return c.JSON(http.StatusCreated, u)
	}
}

func (uh *UserHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		u := &models.User{Nickname:nickname}

		if err := c.Bind(u); err != nil {
			logrus.Info(err)
			return err
		}

		if err := uh.userUcase.UpdateUser(u); err != nil {
			if err == tools.ErrConflict {
				logrus.Info(err)
				return c.JSON(http.StatusConflict, tools.Message{
					Message: "Conflict"})
			}
			if err == tools.ErrNotFound {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message: "Not found"})
			}
			return c.JSON(http.StatusBadRequest, tools.Error{
				ErrMsg: tools.BadRequest.Error(),
			})
		}
		return c.JSON(http.StatusOK, u)
	}
}
