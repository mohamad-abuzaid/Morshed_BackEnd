package controllers

import (
	"morshed/data/engine/sql"
	"morshed/data/models"
	"morshed/domain/services"
	"morshed/helpers"

	"github.com/kataras/iris/v12"
)

// ProductHandler is the http mux for products.
type ProductController struct {
	Ctx     iris.Context
	Service services.ProductService
}

func (c *ProductController) Get() ([]models.Product, error) {
	prods, err := c.Service.GetAll()
	return prods, err
}

// GetByID fetches a single record from the database and sends it to the client.
// Method: GET.
func (c *ProductController) GetBy(id int64) (models.Product, error) {
	prod, err := c.Service.GetByID(id)
	if err != nil {
		// this message will be binded to the
		// main.go -> app.OnAnyErrorCode -> NotFound -> shared/error.html -> .Message text.
		c.Ctx.Values().Set("message", "Product couldn't be found!")
	}
	return prod, err // it will throw/emit 404 if found == false.
}

// Create adds a record to the database.
// Method: POST.
func (h *ProductController) Post() {
	var product models.Product
	if err := h.Ctx.ReadJSON(&product); err != nil {
		return
	}

	id, err := h.Service.Create(product)
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
	h.Ctx.JSON(iris.Map{product.PrimaryKey(): id})
}

// Update performs a full-update of a record in the database.
// Method: PUT.
func (h *ProductController) Put() {
	var product models.Product
	if err := h.Ctx.ReadJSON(&product); err != nil {
		return
	}

	prod, err := h.Service.Update(product)
	if err != nil {
		if err == sql.ErrUnprocessable {
			h.Ctx.StopWithJSON(iris.StatusUnprocessableEntity, helpers.MnewError(iris.StatusUnprocessableEntity, h.Ctx.Request().Method, h.Ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("ProductHandler.Update(DB): %v", err)
		helpers.MwriteInternalServerError(h.Ctx)
		return
	}

	status := iris.StatusOK
	if prod.ID <= 0 {
		status = iris.StatusNotModified
	}

	h.Ctx.StatusCode(status)
}

// PartialUpdate is the handler for partially update one or more fields of the record.
// Method: PATCH.
func (h *ProductController) Patch() {
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

// Delete removes a record from the database.
// Method: DELETE.
func (h *ProductController) Delete() {
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
