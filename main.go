package main

import (
	"log"
	"os"
	"os/signal"

	"morshed/api"
	"morshed/data/datasource"
	"morshed/helpers"

	"github.com/kataras/iris/v12"

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
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")

	// Load the template files.
	tmpl := iris.HTML("./app/views/", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	// Serve our front-end and its assets.
	app.HandleDir("/", iris.Dir("./app/public"))
	template.AddComp(chartjs.NewChart())

	// Note, it's buffered, so make sure it's closed so it can flush any buffered contents.
	ac := accesslog.File("./access.log")
	defer ac.Close()

	app.UseRouter(ac.Handler)
	app.UseRouter(recover.New())

	/////////////////////////////////////////////////
	////////////// ADMIN PANEL /////////////////////

	// Initialize Admin Panel
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
		ShowList: false,
	})

	/////////////////////////////////////////////////
	////////////////// DataSource //////////////////

	// Prepare our repositories and services.
	db, err := datasource.StartMySql(datasource.MySQL)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}

	/////////////////////////////////////////////////
	/////////////////// Routing ////////////////////
	///////////////////////////////////////////////

	secret := helpers.Mgetenv("JWT_SECRET", "EbnJO3bwmX")

	subRouter := api.Router(db, secret)
	app.PartyFunc("/", subRouter)

	/////////////////////////////////////////////////
	//////////////////// RUN ///////////////////////
	///////////////////////////////////////////////

	app.Listen(":80", iris.WithOptimizations)

	go func() {
		_ = app.Run(iris.Addr(":80"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
}
