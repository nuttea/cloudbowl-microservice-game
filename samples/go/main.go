package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {

	/// Check HTTP Method
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "only POST method supported")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := playrandom(v)

	// Logic to override random command
	var myx = v.Arena.State["https://cloudbowl-samples-java-quarkus-yngbkt2j3a-uc.a.run.app"].X
	var myy = v.Arena.State["https://cloudbowl-samples-java-quarkus-yngbkt2j3a-uc.a.run.app"].Y
	var myd = v.Arena.State["https://cloudbowl-samples-java-quarkus-yngbkt2j3a-uc.a.run.app"].Direction
	var myh = v.Arena.State["https://cloudbowl-samples-java-quarkus-yngbkt2j3a-uc.a.run.app"].WasHit
	var diffx int
	var diffy int
	var shoot = false
	var shootto string
	log.Printf("MyState: X %v, Y %v, Direction %v, WasHit %v \n", myx, myy, myd, myh)

	for player, state := range v.Arena.State {
		diffx = state.X - myx
		diffy = state.Y - myy
		//log.Printf("Diff: Player %v X %v, Y %v", player, diffx, diffy)

		switch myd {
		case "S":
			if (diffy > 0) && (diffy < 4) && (diffx == 0) {
				shoot = true
				shootto = player
				resp = "T"
			}
		case "N":
			if (diffy < 0) && (diffy > -4) && (diffx == 0) {
				shoot = true
				shootto = player
				resp = "T"
			}
		case "E":
			if (diffx > 0) && (diffx < 4) && (diffy == 0) {
				shoot = true
				shootto = player
				resp = "T"
			}
		case "W":
			if (diffx < 0) && (diffx > -4) && (diffy == 0) {
				shoot = true
				shootto = player
				resp = "T"
			}
		}
	}

	log.Printf("Shoot: %v %v", shootto, shoot)
	fmt.Fprint(w, resp)
}

func playrandom(input ArenaUpdate) (response string) {

	commands := []string{"F", "R", "L", "F", "F", "T", "T", "T", "T", "T", "T", "T", "T", "T", "T"}
	rand := rand2.Intn(15)
	//log.Printf("Rand: %v %v", rand, commands[rand])
	return commands[rand]
}
