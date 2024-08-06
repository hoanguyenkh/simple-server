package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var uniqueNumber = struct {
	sync.RWMutex
	m map[string]bool
}{
	m: make(map[string]bool),
}

func genUniqueNumber() string {
	for {
		// Generate a new random big.Int
		n, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 256))
		if err != nil {
			log.Println("Error generating random number:", err)
			continue
		}

		// Convert to string
		numStr := n.String()

		// Check if it's unique
		uniqueNumber.Lock()
		if _, exists := uniqueNumber.m[numStr]; !exists {
			uniqueNumber.m[numStr] = true
			uniqueNumber.Unlock()
			return numStr
		}
		uniqueNumber.Unlock()
	}
}

// serveWs handles the websocket connection
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error msg:", err)
			break
		}
		log.Printf("recv: %s", message)

		un := genUniqueNumber()
		err = ws.WriteMessage(websocket.TextMessage, []byte(un))
		if err != nil {
			log.Println("Error writing message: ", err)
			break
		}
	}

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
