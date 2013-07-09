package app

import (
	http "net/http"

	controllers "github.com/drbawb/babou/app/controllers"
	filters "github.com/drbawb/babou/app/filters"

	web "github.com/drbawb/babou/lib/web"

	mux "github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	web.Router = r

	// Shorthand for controllers
	home := controllers.NewHomeController()
	login := controllers.NewLoginController()

	// Shows public homepage, redirects to private site if valid session can be found.
	r.HandleFunc("/",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(home, "index")).Name("homeIndex")

	// Displays a login form.
	r.HandleFunc("/login",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "index")).
		Methods("GET").
		Name("loginIndex")

	r.HandleFunc("/login",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "session")).
		Methods("POST").
		Name("loginSession")

	r.HandleFunc("/logout",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "logout")).
		Methods("GET").
		Name("loginDelete")

	r.HandleFunc("/download",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "download")).
		Methods("GET").
		Name("loginDownload")

	// Displays a registration form
	r.HandleFunc("/register",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "new")).Methods("GET").Name("loginNew")
	// Handles a new user's registration request.
	r.HandleFunc("/register",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "create")).Methods("POST").Name("loginCreate")

	// Catch-All: Displays all public assets.
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		web.DisableDirectoryListing(http.FileServer(http.Dir("assets/")))))

	return r
}

//TODO: should move some of this to a library package.
func handleRedirect(redirect *web.RedirectPath, response http.ResponseWriter, request *http.Request) {
	if redirect.NamedRoute != "" {
		url, err := web.Router.Get(redirect.NamedRoute).URL()
		if err != nil {
			http.Error(response, string("While trying to redirect you to another page the server encountered an error. Please reload the homepage"),
				500)
		}

		http.Redirect(response, request, url.Path, 302)
	}
}
