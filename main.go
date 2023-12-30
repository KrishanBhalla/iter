package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/KrishanBhalla/iter/internal/websocket"
	"github.com/KrishanBhalla/iter/models"

	// "github.com/gocolly/colly"
	// "github.com/gocolly/colly/queue"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	hmacKey             = "secret-hmac-key"
	userPwPepper        = "secret-pepper"
	similarityThreshold = 0.8
	port                = ":8080"
)

func isProd() bool {
	return false
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	// on first run only
	// scrapeAudleyTravel()

	r := mux.NewRouter()

	services, err := models.NewServices(
		models.WithUser(hmacKey, userPwPepper),
		models.WithContent(similarityThreshold),
	)
	defer services.Close()
	// err = services.DestructiveReset()
	must(err)

	setupRoutes(r, *services)

	// Middleware
	// b, err := rand.Bytes(32)
	// csrfMw := csrf.Protect(b, csrf.Secure(isProd()))
	// userMw := middleware.User{
	// 	UserService: services.User,
	// }
	// requireUserMw := middleware.RequireUser{
	// 	User: userMw,
	// }

	// Listen And Serve
	fmt.Println("Listening on", port)
	// Apply this to every request
	err = http.ListenAndServe(port, r) //csrfMw(r))
	must(err)
}

func setupRoutes(r *mux.Router, services models.Services) {
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
	r.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		countries, err := services.Content.Countries()
		if err != nil {
			log.Println(err)
			return
		}
		jsonData, err := json.Marshal(countries)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}).Methods("GET")
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

// func scrapeAudleyTravel() {

// 	baseDomain := "https://www.audleytravel.com"

// 	c := colly.NewCollector(
// 		colly.AllowedDomains("www.audleytravel.com"),
// 	)
// 	// Find and print all links
// 	topLevelDestinations := make([]string, 0)
// 	finalDestinations := make([]string, 0)
// 	content := make([]models.Content, 0)
// 	c.OnHTML(".dest-list__links > li > a[href]", func(e *colly.HTMLElement) {
// 		link := e.Attr("href")
// 		topLevelDestinations = append(topLevelDestinations, baseDomain+link+"/places-to-go#list")
// 	})

// 	c.OnHTML(".block-link", func(e *colly.HTMLElement) {
// 		link := e.Attr("href")
// 		finalDestinations = append(finalDestinations, baseDomain+link)
// 	})

// 	countryRegex := regexp.MustCompile(`[\w\-]+`)
// 	placeRegex := regexp.MustCompile(`[\w\-]+$`)
// 	c.OnHTML(".readmore", func(e *colly.HTMLElement) {
// 		texts := make([]string, 0)
// 		text := "Title: Intro\n Description: "
// 		e.ForEach("*", func(i int, el *colly.HTMLElement) {
// 			if el.Name == "h3" {
// 				texts = append(texts, text)
// 				text = "Title: " + el.Text + "\n Description: "
// 			} else if el.Name == "p" {
// 				text += el.Text + " \n "
// 			}
// 		})
// 		tmp := make([]models.Content, len(texts))
// 		for i, text := range texts {
// 			url := e.Request.URL.String()
// 			tmp[i] = models.Content{
// 				URL:      url,
// 				Country:  hyphenCaseToUpperCase(countryRegex.FindString(url[len(baseDomain):])),
// 				Location: hyphenCaseToSentenceCase(placeRegex.FindString(url[len(baseDomain):])),
// 				Content:  text,
// 			}
// 		}
// 		content = append(content, tmp...)
// 	})

// 	c.OnRequest(func(request *colly.Request) {
// 		fmt.Println("Visiting", request.URL.String())
// 	})

// 	err := c.Visit("https://www.audleytravel.com/destinations")
// 	if err != nil {
// 		fmt.Println("Error visiting https://www.audleytravel.com/destinations:", err)
// 	}

// 	destQueue, _ := queue.New(16, &queue.InMemoryQueueStorage{MaxSize: 1000}) // tried up to 8 threads

// 	// Uncomment only to test
// 	// topLevelDestinations = []string{baseDomain + "/usa" + "/places-to-go#list"}
// 	for _, dest := range topLevelDestinations {
// 		destQueue.AddURL(dest)
// 	}
// 	destQueue.Run(c)

// 	resultQueue, _ := queue.New(16, &queue.InMemoryQueueStorage{MaxSize: 1000}) // tried up to 8 threads

// 	for _, dest := range finalDestinations {
// 		resultQueue.AddURL(dest)
// 	}
// 	resultQueue.Run(c)

// 	services, err := models.NewServices(
// 		models.WithUser(hmacKey, userPwPepper),
// 		models.WithContent(0.8),
// 	)
// 	defer services.Close()

// 	for _, res := range content {
// 		fmt.Println("")
// 		services.Content.Update(&res)
// 	}
// }

// func getLinks(whiteList func(*colly.Collector), callback func(*colly.HTMLElement), querySelector, baseDomain string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	c := colly.NewCollector(whiteList)
// 	c.OnHTML(querySelector, callback)
// }

// func hyphenCaseToSentenceCase(str string) string {
// 	str = deHyphen(str)
// 	return cases.Title(language.BritishEnglish).String(str)
// }

// func hyphenCaseToUpperCase(str string) string {
// 	str = deHyphen(str)
// 	return cases.Upper(language.BritishEnglish).String(str)
// }

// func deHyphen(str string) string {
// 	return strings.Join(strings.Split(str, "-"), " ")
// }
