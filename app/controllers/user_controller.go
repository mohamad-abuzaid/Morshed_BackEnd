// file: controllers/user_controller.go

package controllers

import (
	"morshed/data/models"
	"morshed/domain/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// UserController is our /auth controller.
// UserController is responsible to handle the following requests:
// GET  			/auth/register
// POST 			/auth/register
// GET 				/auth/login
// POST 			/auth/login
// GET 				/auth/me
// All HTTP Methods /auth/logout
type UserController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

const userIDKey = "UserID"

func (c *UserController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *UserController) logout() {
	c.Session.Destroy()
}

var registerStaticView = mvc.View{
	Name: "auth/register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister handles GET: http://localhost:8080/auth/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return registerStaticView
}

// PostRegister handles POST: http://localhost:8080/auth/register.
func (c *UserController) PostRegister() mvc.Result {
	// get firstname, username and password from the form.
	var (
		firstname 		= c.Ctx.FormValue("firstname")
		username  		= c.Ctx.FormValue("username")
		password  		= c.Ctx.FormValue("password")
		dob  			= "10/12/2009"
		address  		= "California"
		description 	= "My Test User"
	)

	// create the new user, the password will be hashed by the service.
	u, err := c.Service.CreateUser(password, models.User{
		Username:  username,
		Firstname: firstname,
		Dob: dob,
		Address: address,
		Description: description,
	})

	// set the user's id to this session even if err != nil,
	// the zero id doesn't matters because .getCurrentUserID() checks for that.
	// If err != nil then it will be shown, see below on mvc.Response.Err: err.
	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		// if not nil then this error will be shown instead.
		Err: err,
		// redirect to /auth/me.
		Path: "/auth/me",
		// When redirecting from POST to GET request you -should- use this HTTP status code,
		// however there're some (complicated) alternatives if you
		// search online or even the HTTP RFC.
		// Status "See Other" RFC 7231, however iris can automatically fix that
		// but it's good to know you can set a custom code;
		// Code: 303,
	}
}

var loginStaticView = mvc.View{
	Name: "auth/login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/auth/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.logout()
	}

	return loginStaticView
}

// PostLogin handles POST: http://localhost:8080/auth/login.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	attrs := map[string]interface{}{"username": username, "password": password}
	u, err := c.Service.GetByAttrs(attrs)

	if err != nil {
		return mvc.Response{
			Path: "/auth/register",
		}
	}

	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		Path: "/auth/me",
	}
}

// GetMe handles GET: http://localhost:8080/auth/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/auth/login"}
	}

	u, err := c.Service.GetByID(c.getCurrentUserID())
	if err != nil {
		// if the  session exists but for some reason the user doesn't exist in the "database"
		// then logout and re-execute the function, it will redirect the client to the
		// /auth/login page.
		c.logout()
		return c.GetMe()
	}

	return mvc.View{
		Name: "auth/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/auth/logout.
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}

	c.Ctx.Redirect("/auth/login")
}
