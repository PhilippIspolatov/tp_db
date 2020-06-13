package delivery

import (
	"net/http"

	"github.com/ifo16u375/tp_db/internal/service"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/labstack/echo"
)

type ServiceHandler struct {
	serviceUcase service.Usecase
}

func NewServiceHandler(router *echo.Echo, sUC service.Usecase) *ServiceHandler {
	sh := &ServiceHandler{
		serviceUcase: sUC,
	}

	router.POST("/api/service/clear", sh.ClearDB())
	router.GET("/api/service/status", sh.GetStatusDB())

	return sh
}

func (sh *ServiceHandler) ClearDB() echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := sh.serviceUcase.ClearAllDB(); err != nil {
			return c.JSON(http.StatusInternalServerError, tools.Message{
				Message:"server error",
			})
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func (sh *ServiceHandler) GetStatusDB() echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := sh.serviceUcase.GetInfoDB()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, tools.Message{
				Message:"server error",
			})
		}
		return c.JSON(http.StatusOK, s)
	}
}