package controllers

import (
	"morshed/data/engine/service"
	"morshed/data/engine/sql"
	"morshed/data/models"
	"morshed/helpers"

	"github.com/kataras/iris/v12"
)

// CategoryHandler is the http mux for categories.
type CategoryHandler struct {
	// [...options]

	service *service.CategoryService
}

// NewCategoryHandler returns the main controller for the categories API.
func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

// GetByID fetches a single record from the database and sends it to the client.
// Method: GET.
func (h *CategoryHandler) GetByID(ctx iris.Context) {
	id := ctx.Params().GetInt64Default("id", 0)

	var cat models.Category
	err := h.service.GetByID(ctx.Request().Context(), &cat, id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.MwriteEntityNotFound(ctx)
			return
		}

		helpers.Mdebugf("CategoryHandler.GetByID(id=%d): %v", id, err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	ctx.JSON(cat)
}

/*

type (
	List struct {
		Data  interface{} `json:"data"`
		Order string      `json:"order"`
		Next  Range       `json:"next,omitempty"`
		Prev  Range       `json:"prev,omitempty"`
	}

	Range struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit`
	}
)
*/

// List lists a set of records from the database.
// Method: GET.
func (h *CategoryHandler) List(ctx iris.Context) {
	q := ctx.Request().URL.Query()
	opts := sql.ParseListOptions(q)

	// initialize here in order to return an empty json array `[]` instead of `null`.
	categories := models.Categories{}
	err := h.service.List(ctx.Request().Context(), &categories, opts)
	if err != nil && err != sql.ErrNoRows {
		helpers.Mdebugf("CategoryHandler.List(DB) (limit=%d offset=%d where=%s=%v): %v",
			opts.Limit, opts.Offset, opts.WhereColumn, opts.WhereValue, err)

			helpers.MwriteInternalServerError(ctx)
		return
	}

	ctx.JSON(categories)
}

// Create adds a record to the database.
// Method: POST.
func (h *CategoryHandler) Create(ctx iris.Context) {
	var cat models.Category
	if err := ctx.ReadJSON(&cat); err != nil {
		return
	}

	id, err := h.service.Insert(ctx.Request().Context(), cat)
	if err != nil {
		if err == sql.ErrUnprocessable {
			ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("CategoryHandler.Create(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	// Send 201 with body of {"id":$last_inserted_id"}.
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{cat.PrimaryKey(): id})
}

// Update performs a full-update of a record in the database.
// Method: PUT.
func (h *CategoryHandler) Update(ctx iris.Context) {
	var cat models.Category
	if err := ctx.ReadJSON(&cat); err != nil {
		return
	}

	affected, err := h.service.Update(ctx.Request().Context(), cat)
	if err != nil {
		if err == sql.ErrUnprocessable {
			ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "required fields are missing"))
			return
		}

		helpers.Mdebugf("CategoryHandler.Update(DB): %v", err)
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
func (h *CategoryHandler) PartialUpdate(ctx iris.Context) {
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

		helpers.Mdebugf("CategoryHandler.PartialUpdate(DB): %v", err)
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
func (h *CategoryHandler) Delete(ctx iris.Context) {
	id := ctx.Params().GetInt64Default("id", 0)

	affected, err := h.service.DeleteByID(ctx.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.MwriteEntityNotFound(ctx)
			return
		}

		helpers.Mdebugf("CategoryHandler.Delete(DB): %v", err)
		helpers.MwriteInternalServerError(ctx)
		return
	}

	status := iris.StatusOK // StatusNoContent
	if affected == 0 {
		status = iris.StatusNotModified
	}

	ctx.StatusCode(status)
}

// Products.

// ListProducts lists products of a Category.
// Example: from cheap to expensive:
// http://localhost:8080/category/3/products?offset=0&limit=30&by=price&order=asc
// Method: GET.
func (h *CategoryHandler) ListProducts(ctx iris.Context) {
	id := ctx.Params().GetInt64Default("id", 0)

	// NOTE: could add cache here too.

	q := ctx.Request().URL.Query()
	opts := sql.ParseListOptions(q).Where("category_id", id)
	opts.Table = "products"
	if opts.OrderByColumn == "" {
		opts.OrderByColumn = "updated_at"
	}

	var products models.Products
	err := h.service.List(ctx.Request().Context(), &products, opts)
	if err != nil && err != sql.ErrNoRows {
		helpers.Mdebugf("CategoryHandler.ListProducts(DB) (table=%s where=%s=%v limit=%d offset=%d): %v",
			opts.Table, opts.WhereColumn, opts.WhereValue, opts.Limit, opts.Offset, err)

			helpers.MwriteInternalServerError(ctx)
		return
	}

	ctx.JSON(products)
}

// InsertProducts assigns new products to a Category (accepts a list of products).
// Method: POST.
func (h *CategoryHandler) InsertProducts(productService *service.ProductService) iris.Handler {
	return func(ctx iris.Context) {
		categoryID := ctx.Params().GetInt64Default("id", 0)

		var products []models.Product
		if err := ctx.ReadJSON(&products); err != nil {
			return
		}

		for i := range products {
			products[i].CategoryID = categoryID
		}

		inserted, err := productService.BatchInsert(ctx.Request().Context(), products)
		if err != nil {
			if err == sql.ErrUnprocessable {
				ctx.StopWithJSON(iris.StatusUnprocessableEntity, newError(iris.StatusUnprocessableEntity, ctx.Request().Method, ctx.Path(), "required fields are missing"))
				return
			}

			helpers.Mdebugf("CategoryHandler.InsertProducts(DB): %v", err)
			helpers.MwriteInternalServerError(ctx)
			return
		}

		if inserted == 0 {
			ctx.StatusCode(iris.StatusNotModified)
			return
		}

		// Send 201 with body of {"inserted":$inserted"}.
		ctx.StatusCode(iris.StatusCreated)
		ctx.JSON(iris.Map{"inserted": inserted})
	}
}
