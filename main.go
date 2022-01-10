package main

import (
	"log"
	"os"
	"os/signal"

	"morshed/controllers"
	"morshed/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/recover"

	_ "github.com/GoAdminGroup/go-admin/adapter/iris"             // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	_ "github.com/GoAdminGroup/themes/sword"                      // ui theme

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"

	"morshed/goadmin/pages"
	"morshed/goadmin/tables"
)

func main() {
	app := iris.New()
	// Serve our front-end and its assets.
	app.HandleDir("/", iris.Dir("./app/public"))
	template.AddComp(chartjs.NewChart())

	eng := engine.Default()
	if err := eng.AddConfigFromJSON("./goadmin/config.json").
		AddGenerators(tables.Generators).
		Use(app); err != nil {
		panic(err)
	}
	eng.HTML("GET", "/karakon", pages.GetDashBoard)
	eng.HTMLFile("GET", "/karakon/hello", "./goadmin/html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	app.HandleDir("/goadmin/uploads", "./goadmin/uploads", iris.DirOptions{
		IndexName: "/index.html",
		//Gzip:      false,
		ShowList:  false,
	})

	go func() {
		_ = app.Run(iris.Addr(":80"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()

	// Note, it's buffered, so make sure it's closed so it can flush any buffered contents.
	ac := accesslog.File("./access.log")
	defer ac.Close()

	app.UseRouter(ac.Handler)
	app.UseRouter(recover.New())
	// Group routes and mvc apps based on /api path prefix.
	api := app.Party("/api")
	{
		// Group based on /api/counter path prefix.
		counterAPI := api.Party("/counter")
		// Optionally, a <trick> to keep the `m` local variable
		// unaccessible outside of this block's scope. That
		// way you can register many mvc apps for different Parties
		// with a "m" variable.
		// Alternatively you can use the mvc.Configure function :)

		// Register a new MVC Application to the counterAPI Party.
		m := mvc.New(counterAPI)
		m.Register(
			// Register a static dependency (static because it doesn't accept an iris.Context,
			// only one instance of that it's used). Helps us to keep a global counter across
			// clients requests.
			services.NewGlobalCounter(),
			// Register a dynamic dependency (GetFields accepts an iris.Context,
			// it binds a new instance on every request). Helps us to
			// set custom fields based on the request handler.
			accesslog.GetFields,
		)
		// Register our controller.
		m.Handle(new(controllers.Counter))
	}

	// GET http://localhost:8080/api/counter
	// POST http://localhost:8080/api/counter/increment
	app.Listen(":8080")
}
