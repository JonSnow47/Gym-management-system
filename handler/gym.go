/*
 * Revision History:
 *     Initial: 2018/05/24        Chen Yanchen
 */

package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"github.com/JonSnow47/Gym-management-system/common"
	"github.com/JonSnow47/Gym-management-system/model"
)

type groundHandler struct{}

var Ground *groundHandler

// 新建场地
func (*groundHandler) New(c echo.Context) error {
	var (
		err error
		req struct {
			Name string `json:"name" validate:"required,min=1,max=20"`
		}
	)

	if err = c.Bind(&req); err != nil {
		c.Logger().Error("[Bind]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrParam))
	}

	if err = c.Validate(&req); err != nil {
		c.Logger().Error("[Validate]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrValidate))
	}

	err = model.GymService.New(req.Name)
	if err != nil {
		c.Logger().Error("[New account]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrMongoDB))
	}

	return c.JSON(http.StatusOK, Resp(common.RespSuccess, nil))
}

// 场地信息
func (*groundHandler) Info(c echo.Context) error {
	var (
		err error
		req struct {
			Id int
		}
	)

	if err = c.Bind(&req); err != nil {
		c.Logger().Error("[Bind]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrParam))
	}

	g, err := model.GymService.Info(req.Id)
	if err != nil {
		if err == mgo.ErrNotFound {
			c.Logger().Error("[Info]", err)
			return c.JSON(http.StatusOK, Resp(common.ErrNotFound))
		}
		c.Logger().Error("[Info]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrMongoDB))
	}

	return c.JSON(http.StatusOK, Resp(common.RespSuccess, g))
}

// 场地列表
func (*groundHandler) List(c echo.Context) error {

	g, err := model.GymService.List()
	if err != nil {
		c.Logger().Error("[List]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrMongoDB))
	}

	return c.JSON(http.StatusOK, Resp(common.RespSuccess, g))
}

// 修改场地状态
func (*groundHandler) ModifyState(c echo.Context) error {
	var (
		err error
		req struct {
			Id int
		}
	)

	if err = c.Bind(&req); err != nil {
		c.Logger().Error("[Bind]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrParam))
	}

	err = model.GymService.State(req.Id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return c.JSON(http.StatusOK, Resp(common.ErrNotFound))
		}
		c.Logger().Error("[ModifyState]", err)
		return c.JSON(http.StatusOK, Resp(common.ErrMongoDB))
	}

	return c.JSON(http.StatusOK, Resp(common.RespSuccess, nil))
}
