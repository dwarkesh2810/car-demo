package controllers

import (
	"car_demo/helper"
	"car_demo/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// Car_masterController operations for Car_master
type Car_masterController struct {
	beego.Controller
}

// URLMapping ...
func (c *Car_masterController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Car_master
// @Param	body		body 	models.Car_master	true		"body for Car_master content"
// @Success 201 {int} models.Car_master
// @Failure 403 body is empty
// @router / [post]
func (c *Car_masterController) Post() {
	ctx := c.Ctx.Input.GetData("user")

	uid := ctx.(*models.Users).Id
	log.Print(uid)
	var v models.Car_master
	// json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if err := c.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		c.Ctx.WriteString("Error while parsing form data: " + err.Error())
		return
	}
	imgPath, err := helper.GetFileAndStore(c.Controller, "imageFile", "cars", string(v.CarType))

	if err != nil {

		log.Print("1111111111111111111111111111111111111111111111")
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}

	v.CarImage = imgPath
	v.CreatedAt = int(time.Now().Unix())
	v.UpdatedAt = int(time.Now().Unix())
	v.UserId = uid

	switch v.CarType {
	case models.Hatchback, models.SUV, models.Sedan:
		if _, err := models.AddCar_master(&v); err == nil {
			helper.JsonResponse(c.Controller, http.StatusCreated, 1, v, "")
			return
		} else {
			helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
			return
		}
	default:
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, "invalid car type")
		return
	}

}

// GetOne ...
// @Title Get One
// @Description get Car_master by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Car_master
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Car_masterController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetCar_masterById(id)
	if err != nil {
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	} else {
		helper.JsonResponse(c.Controller, http.StatusOK, 1, v, "")
		return
	}
}

// GetAll ...
// @Title Get All
// @Description get Car_master
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Car_master
// @Failure 403
// @router / [get]
func (c *Car_masterController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, "invalid query key/value pair")
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllCar_master(query, fields, sortby, order, offset, limit)
	if err != nil {
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	} else {
		helper.JsonResponse(c.Controller, http.StatusOK, 1, l, "")
		return
	}
}

// Put ...
// @Title Put
// @Description update the Car_master
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Car_master	true		"body for Car_master content"
// @Success 200 {object} models.Car_master
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Car_masterController) Put() {

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Car_master{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	switch v.CarType {
	case models.Hatchback, models.SUV, models.Sedan:
		if err := models.UpdateCar_masterById(&v); err == nil {
			helper.JsonResponse(c.Controller, http.StatusOK, 1, "OK", "")
			return
		} else {
			helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
			return
		}
	default:
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, "invalid car type")
		return
	}
}

// Delete ...
// @Title Delete
// @Description delete the Car_master
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Car_masterController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteCar_master(id); err == nil {
		helper.JsonResponse(c.Controller, http.StatusOK, 1, "OK", "")
		return
	} else {
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
}
