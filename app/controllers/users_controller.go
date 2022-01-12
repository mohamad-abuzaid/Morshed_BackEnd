package controllers

import (
	"morshed/data/engine/sql"
	"morshed/data/models"
	"morshed/domain/services"
	"morshed/helpers"

	"github.com/kataras/iris/v12"
)

// UsersController is our /users API controller.
// GET				/users  | get all
// GET				/users/{id:int64} | get by id
// PUT				/users/{id:int64} | update by id
// DELETE			/users/{id:int64} | delete by id
// Requires basic authentication.
type UsersController struct {
	// Optionally: context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserService
}

// Get returns list of the users.
// Demo:
// curl -i -u admin:password http://localhost:8080/users
//
// The correct way if you have sensitive data:
// func (c *UsersController) Get() (results []viewmodels.User) {
// 	data := c.Service.GetAll()
//
// 	for _, user := range data {
// 		results = append(results, viewmodels.User{user})
// 	}
// 	return
// }
// otherwise just return the datamodels.
func (c *UsersController) Get() ([]models.User, error) {
	users, err := c.Service.GetAll()
	return users, err
}

// GetBy returns a user.
// Demo:
// curl -i -u admin:password http://localhost:8080/users/1
func (c *UsersController) GetBy(id int64) (models.User, error) {
	user, err := c.Service.GetByID(id)
	if err != nil {
		// this message will be binded to the
		// main.go -> app.OnAnyErrorCode -> NotFound -> shared/error.html -> .Message text.
		c.Ctx.Values().Set("message", "User couldn't be found!")
	}
	return user, err // it will throw/emit 404 if found == false.
}

func (h *UsersController) Post() {
	var user models.User
	if err := h.Ctx.ReadJSON(&user); err != nil {
		return
	}

	id, err := h.Service.Create(user)
	if err != nil {
		if err == sql.ErrUnprocessable {
			h.Ctx.StopWithJSON(iris.StatusUnprocessableEntity, helpers.MnewError(iris.StatusUnprocessableEntity, h.Ctx.Request().Method, h.Ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("ProductHandler.Create(DB): %v", err)
		helpers.MwriteInternalServerError(h.Ctx)
		return
	}

	// Send 201 with body of {"id":$last_inserted_id"}.
	h.Ctx.StatusCode(iris.StatusCreated)
	h.Ctx.JSON(iris.Map{user.PrimaryKey(): id})
}

func (h *UsersController) Put() {
	var user models.User
	if err := h.Ctx.ReadJSON(&user); err != nil {
		return
	}

	user, err := h.Service.Update(user)
	if err != nil {
		if err == sql.ErrUnprocessable {
			h.Ctx.StopWithJSON(iris.StatusUnprocessableEntity,
				helpers.MnewError(iris.StatusUnprocessableEntity,
					h.Ctx.Request().Method, h.Ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("ProductHandler.Update(DB): %v", err)
		helpers.MwriteInternalServerError(h.Ctx)
		return
	}

	status := iris.StatusOK
	if user.ID <= 0 {
		status = iris.StatusNotModified
	}

	h.Ctx.StatusCode(status)
}

func (h *UsersController) Patch() {
	id := h.Ctx.Params().GetInt64Default("id", 0)

	var attrs map[string]interface{}
	if err := h.Ctx.ReadJSON(&attrs); err != nil {
		return
	}

	affected, err := h.Service.PatchUpdate(id, attrs)
	if err != nil {
		if err == sql.ErrUnprocessable {
			h.Ctx.StopWithJSON(iris.StatusUnprocessableEntity,
				helpers.MnewError(iris.StatusUnprocessableEntity,
					h.Ctx.Request().Method, h.Ctx.Path(), "unsupported value(s)"))
			return
		}

		helpers.Mdebugf("ProductHandler.PartialUpdate(DB): %v", err)
		helpers.MwriteInternalServerError(h.Ctx)
		return
	}

	status := iris.StatusOK
	if affected == 0 {
		status = iris.StatusNotModified
	}

	h.Ctx.StatusCode(status)
}

func (h *UsersController) Delete() {
	id := h.Ctx.Params().GetInt64Default("id", 0)

	affected, err := h.Service.DeleteByID(id)
	if err != nil {
		helpers.Mdebugf("ProductHandler.Delete(DB): %v", err)
		helpers.MwriteInternalServerError(h.Ctx)
		return
	}

	status := iris.StatusOK // StatusNoContent
	if affected == 0 {
		status = iris.StatusNotModified
	}

	h.Ctx.StatusCode(status)
}
