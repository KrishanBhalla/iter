package main

import (
	"fmt"
	"net/http"

	"github.com/KrishanBhalla/iter/models"
	"github.com/KrishanBhalla/iter/rand"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	hmacKey      = "secret-hmac-key"
	userPwPepper = "secret-pepper"
	port         = ":3000"
)

func isProd() bool {
	return false
}

func main() {

	r := mux.NewRouter()

	services, err := models.NewServices(
		models.WithUser(hmacKey, userPwPepper),
		models.WithContent(0.8),
	)
	defer services.Close()
	// err = services.DestructiveReset()
	must(err)

	setupRoutes(r)

	// Middleware
	b, err := rand.Bytes(32)
	csrfMw := csrf.Protect(b, csrf.Secure(isProd()))
	// userMw := middleware.User{
	// 	UserService: services.User,
	// }
	// requireUserMw := middleware.RequireUser{
	// 	User: userMw,
	// }

	// Listen And Serve
	fmt.Println("Listening on", port)
	// Apply this to every request
	err = http.ListenAndServe(port, csrfMw(r))
	must(err)
}

func setupRoutes(r *mux.Router) {
// 	// Websocket
// 	pool := websocket.NewChannelPool()
// 	go pool.Start()
// 	serveWsFunc := func(w http.ResponseWriter, r *http.Request) {
// 		serveWs(pool, w, r)
// 	}

	// // Handlers ---------
	// // Static
	// r.Handle("/", staticC.Home)
	// r.Handle("/home", staticC.Home)
	// r.Handle("/contact", staticC.Contact)

	// // Users
	// r.Handle("/signup", usersC.SignupView).Methods("GET")
	// r.HandleFunc("/signup", usersC.CreateOrLogin).Methods("POST")
	// r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	// r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// r.PathPrefix("/ws").HandlerFunc(http.HandlerFunc(serveWsFunc))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
