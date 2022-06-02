package routers

import (
	"git.gonkar.com/gonkar/infra-cmd/app"
	"git.gonkar.com/gonkar/infra-cmd/controllers"
	"git.gonkar.com/gonkar/infra-cmd/middlewares"

	"github.com/gorilla/mux"
)

// Router - Initialize the router
var MainRouter = mux.NewRouter()

func routerPaths() {

	MainRouter.Use(app.JwtAuthentication)

	// Router - Path registration
	MainRouter.HandleFunc("/webhook/mattermost", middlewares.Chain(controllers.Webhook, middlewares.Logging())).Methods("POST", "OPTIONS")

}
