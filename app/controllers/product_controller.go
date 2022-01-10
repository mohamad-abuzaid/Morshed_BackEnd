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
	Ctx iris.Context
	Service services.ProductService
}

// GetByID fetches a single record from the database and sends it to the client.
// Method: GET.
func (c *ProductController) GetBy(id int64) (product models.Product, found bool) {
	u, found := c.Service.GetByID(id)
	if !found {
		// this message will be binded to the
		// main.go -> app.OnAnyErrorCode -> NotFound -> shared/error.html -> .Message text.
		c.Ctx.Values().Set("message", "Product couldn't be found!")
	}
	return u, found // it will throw/emit 404 if found == false.
}

// List lists a set of records from the database.
// Method: GET.
func (h *ProductController) List(ctx iris.Context) {
	key := ctx.Request().URL.RawQuery

	products := []byte("[]")
	err := h.cache.List(ctx.Request().Context(), key, &products)
	if err != nil && err != sql.ErrNoRows {
		helpers.Mdebugf("ProductHandler.List(DB) (%s): %v",
			key, err)

			helpers.MwriteInternalServerError(ctx)
		return
	}

	ctx.ContentType("application/json")
	ctx.Write(products)
}

// Create adds a record to the database.
// Method: POST.
func (h *ProductController) Create(ctx iris.Context) {
	var product models.Product
	if err := ctx.ReadJSON(&product); err != nil {
		return
	}

	id, err := h.service.Insert(ctx.Request().Context(), product)
	if err != nil {
		if err == sql.ErrUnprocessable {
			ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("ProductHandler.Create(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	// Send 201 with body of {"id":$last_inserted_id"}.
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{product.PrimaryKey(): id})
}

// Update performs a full-update of a record in the database.
// Method: PUT.
func (h *ProductController) Update(ctx iris.Context) {
	var product models.Product
	if err := ctx.ReadJSON(&product); err != nil {
		return
	}

	affected, err := h.service.Update(ctx.Request().Context(), product)
	if err != nil {
		if err == sql.ErrUnprocessable {
			ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("ProductHandler.Update(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	status := iris.StatusOK
	if affected == 0 {
		status = iris.StatusNotModified
	}

	ctx.StatusCode(status)
}

// PartialUpdate is the handler for partially update one or more fields of the record.
// Method: PATCH.
func (h *ProductController) PartialUpdate(ctx iris.Context) {
	id := ctx.Params().GetInt64Default("id", 0)

	var attrs map[string]interface{}
	if err := ctx.ReadJSON(&attrs); err != nil {
		return
	}

	affected, err := h.service.PartialUpdate(ctx.Request().Context(), id, attrs)
	if err != nil {
		if err == sql.ErrUnprocessable {
			ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "unsupported value(s)"))
			return
		}

		helpers.Mdebugf("ProductHandler.PartialUpdate(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	status := iris.StatusOK
	if affected == 0 {
		status = iris.StatusNotModified
	}

	ctx.StatusCode(status)
}

// Delete removes a record from the database.
// Method: DELETE.
func (h *ProductController) Delete(ctx iris.Context) {
	id := ctx.Params().GetInt64Default("id", 0)

	affected, err := h.service.DeleteByID(ctx.Request().Context(), id)
	if err != nil {
		helpers.Mdebugf("ProductHandler.Delete(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	status := iris.StatusOK // StatusNoContent
	if affected == 0 {
		status = iris.StatusNotModified
	}

	ctx.StatusCode(status)
}
