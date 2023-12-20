package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KrishanBhalla/iter/internal/websocket"
	"github.com/KrishanBhalla/iter/models"
	"github.com/KrishanBhalla/iter/rand"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	hmacKey      = "secret-hmac-key"
	userPwPepper = "secret-pepper"
	port         = ":8080"
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
	// Websocket
	r.PathPrefix("/ws").HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	}))

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

}

func serveWs(w http.ResponseWriter, r *http.Request) {

	log.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%+V\n", err)
	}
	client := &websocket.Client{
		ID:   uuid.NewString(),
		Conn: conn,
	}
	client.Read()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
