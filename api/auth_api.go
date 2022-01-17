// Package api contains the handlers for our HTTP Endpoints.
package api

import (
	"time"

	"morshed/app/controllers"
	"morshed/data/engine/sql"
	"morshed/data/repositories"
	middleware "morshed/domain/middlewares"
	"morshed/domain/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// Router accepts any required dependencies and returns the main server's handler.
func Router(db sql.Database, secret string) func(iris.Party) {
	return func(r iris.Party) {
		r.Use(requestid.New())

		signer := jwt.NewSigner(jwt.HS256, secret, 15*time.Minute)
		r.Get("/token", writeToken(signer))

		verify := jwt.NewVerifier(jwt.HS256, secret).Verify(nil)
		r.Use(verify)
		// Generate a token for testing by navigating to
		// http://localhost:8080/token endpoint.
		// Copy-paste it to a ?token=$token url parameter or
		// open postman and put an Authentication: Bearer $token to get
		// access on create, update and delete endpoinds.

		var (
			userRepository = repositories.NewUserRepository(db)
			userService    = services.NewUserService(userRepository)

			productRepository = repositories.NewProductRepository(db)
			productService    = services.NewProductService(productRepository)
		)

		/////////////////// User /////////////////////

		// "/user" based mvc application.
		sessManager := sessions.New(sessions.Config{
			Cookie:  "morshed_cookie",
			Expires: 120 * time.Hour,
		})
		user := mvc.New(r.Party("/auth"))
		user.Register(
			userService,
			sessManager.Start,
		)
		user.Handle(new(controllers.UserController))

		/////////////////// Users /////////////////////

		// "/users" based mvc application.
		users := mvc.New(r.Party("/users"))
		// Add the basic authentication(admin:password) middleware
		// for the /users based requests.
		users.Router.Use(middleware.BasicAuth)
		// Bind the "userService" to the UserController's Service (interface) field.
		users.Register(userService)
		users.Handle(new(controllers.UsersController))

		/////////////////// Product /////////////////////

		prod := mvc.New(r.Party("/product"))
		prod.Register(
			productService,
		)
		prod.Handle(new(controllers.ProductController))
	}
}

func writeToken(signer *jwt.Signer) iris.Handler {
	return func(ctx iris.Context) {
		claims := jwt.Claims{
			Issuer:   "https://iris-go.com",
			Audience: []string{requestid.Get(ctx)},
		}

		token, err := signer.Sign(claims)
		if err != nil {
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		ctx.Write(token)
	}
}
