package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// STRUCTURES
type pageData struct {
	Title   string
	Rooms   template.HTML
	Compass template.HTML
}

type roomData struct {
	discovered bool
	exitNorth  bool
	exitSouth  bool
	exitWest   bool
	exitEast   bool
}

const north string = "north"
const south string = "south"
const west string = "west"
const east string = "east"

var theMaze map[int]roomData
var playerLocation int

// HANDLERS
func Root(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		Title string
	}{
		Title: "Root of Capstone Project!",
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/index.html")
	if err != nil {
		log.Printf("Handler, Root, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = pageTemplate.Execute(w, pageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Root, Error executing page: %v", err)
	}
}

func Game(w http.ResponseWriter, r *http.Request) {
	//This takes the action encoded in r and updates and displays the new updated page.
	//To prevent caching, add the following header to the response "Cache-Control: no-store or Cache-Control: no-cache, no-store, must-revalidate"

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	//Get the action path parameter from the request.
	action := r.URL.Query().Get("action")
	switch action {
	case "start":
		start(w)
	case "end":
		end(w)
	case north, south, west, east:
		move(action, w)
	default:
		//Unknown action
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		var pageBody string = "400: Unknown action."
		w.Write([]byte(pageBody))
	}
}

// Maze
func mazeInit() {
	//create the map
	theMaze = make(map[int]roomData)
	//set up the rooms, set everything to false.
	for i := 0; i < 16; i++ {
		theMaze[i] = roomData{discovered: false, exitNorth: false, exitSouth: false, exitWest: false, exitEast: false}
	}
	//Set the doors for each room as well as discover the first room (02)
	room := theMaze[0]
	room.exitEast = true
	theMaze[0] = room
	room = theMaze[1]
	room.exitNorth = true
	room.exitWest = true
	room.exitEast = true
	theMaze[1] = room
	room = theMaze[2]
	room.discovered = true
	room.exitNorth = true
	room.exitWest = true
	room.exitEast = true
	theMaze[2] = room
	room = theMaze[3]
	room.exitNorth = true
	room.exitWest = true
	theMaze[3] = room
	room = theMaze[4]
	room.exitNorth = true
	room.exitEast = true
	theMaze[4] = room
	room = theMaze[5]
	room.exitNorth = true
	room.exitSouth = true
	room.exitWest = true
	theMaze[5] = room
	room = theMaze[6]
	room.exitNorth = true
	room.exitSouth = true
	theMaze[6] = room
	room = theMaze[7]
	room.exitSouth = true
	theMaze[7] = room
	room = theMaze[8]
	room.exitSouth = true
	theMaze[8] = room
	room = theMaze[9]
	room.exitNorth = true
	room.exitSouth = true
	theMaze[9] = room
	room = theMaze[10]
	room.exitSouth = true
	room.exitEast = true
	theMaze[10] = room
	room = theMaze[11]
	room.exitNorth = true
	room.exitWest = true
	theMaze[11] = room
	room = theMaze[12]
	room.exitEast = true
	theMaze[12] = room
	room = theMaze[13]
	room.exitSouth = true
	room.exitWest = true
	room.exitEast = true
	theMaze[13] = room
	room = theMaze[14]
	room.exitWest = true
	theMaze[14] = room
	room = theMaze[15]
	room.exitSouth = true
	theMaze[15] = room
}

func generateKnownMap() template.HTML {
	//This runs through the map and all rooms showing discovered have their container added to the string.
	var knownMap strings.Builder
	for i := 0; i < 16; i++ {
		if theMaze[i].discovered {
			fmt.Fprintf(&knownMap, "<div class='room%d'><img src='/assets/images/r%d.png' alt='Maze Room' width='200' height='200' /></div>\n", i, i)
		}
	}
	//Add the Player
	playerLocationTop := ((15 - playerLocation) / 4) * 200
	playerLocationLeft := (playerLocation % 4) * 200
	//log.Printf("playerLocationLeft: %d", playerLocationLeft)
	fmt.Fprintf(&knownMap, "<div style='position:absolute;top:%dpx;left:%dpx'><img src='/assets/images/player.png' alt='Player!' width='200' height='200' /></div>\n", playerLocationTop, playerLocationLeft)
	return template.HTML(knownMap.String())
}

func generateCompass() template.HTML {
	var compass strings.Builder
	if theMaze[playerLocation].exitNorth {
		fmt.Fprint(&compass, "<div class='n_arrow'><a href='/app?action=north'><img src='/assets/images/green_arrow_n.png' alt='Green Arrow North' width='150' height='200' /></a></div>\n")
	} else {
		fmt.Fprint(&compass, "<div class='n_arrow'><img src='/assets/images/red_arrow_n.png' alt='Red Arrow North' width='150' height='200' /></div>\n")
	}
	if theMaze[playerLocation].exitSouth {
		fmt.Fprint(&compass, "<div class='s_arrow'><a href='/app?action=south'><img src='/assets/images/green_arrow_s.png' alt='Green Arrow South' width='150' height='200' /></a></div>\n")
	} else {
		fmt.Fprint(&compass, "<div class='s_arrow'><img src='/assets/images/red_arrow_s.png' alt='Red Arrow South' width='150' height='200' /></div>\n")
	}
	if theMaze[playerLocation].exitWest {
		fmt.Fprint(&compass, "<div class='w_arrow'><a href='/app?action=west'><img src='/assets/images/green_arrow_w.png' alt='Green Arrow West' width='200' height='150' /></a></div>\n")
	} else {
		fmt.Fprint(&compass, "<div class='w_arrow'><img src='/assets/images/red_arrow_w.png' alt='Red Arrow West' width='200' height='150' /></div>\n")
	}
	if theMaze[playerLocation].exitEast {
		fmt.Fprint(&compass, "<div class='e_arrow'><a href='/app?action=east'><img src='/assets/images/green_arrow_e.png' alt='Green Arrow East' width='200' height='150' /></a></div>\n")
	} else {
		fmt.Fprint(&compass, "<div class='e_arrow'><img src='/assets/images/red_arrow_e.png' alt='Red Arrow East' width='200' height='150' /></div>\n")
	}
	return template.HTML(compass.String())
}

// ACTIONS
func start(w http.ResponseWriter) {
	//Initialize the Maze
	mazeInit()
	//Set player location to starting room, 2
	playerLocation = 2

	startPageData := pageData{
		Title:   "Maze Runner - Start!",
		Rooms:   generateKnownMap(),
		Compass: generateCompass(),
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/maze.html")
	if err != nil {
		log.Printf("Handler, Game, start Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	//log.Printf("Known Map: %s\n", startPageData.Rooms)
	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, start, Error executing page: %v", err)
	}
}

func end(w http.ResponseWriter) {
	startPageData := pageData{
		Title: "Maze Runner End!",
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/end.html")
	if err != nil {
		log.Printf("Handler, Game, end, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, end, Error executing page: %v", err)
	}
}

func move(direction string, w http.ResponseWriter) {
	//Each direction changes the playerLocation by a different value. If there is no direction, then playerLocation does not change.
	switch direction {
	case north:
		if theMaze[playerLocation].exitNorth {
			playerLocation += 4
		}
	case south:
		if theMaze[playerLocation].exitSouth {
			playerLocation -= 4
		}
	case west:
		if theMaze[playerLocation].exitWest {
			playerLocation -= 1
		}
	case east:
		if theMaze[playerLocation].exitEast {
			playerLocation += 1
		}
	}

	//This sets the room to discovered, regardless if the player moved or not.
	room := theMaze[playerLocation]
	room.discovered = true
	theMaze[playerLocation] = room

	startPageData := pageData{
		Title:   "Maze Runner",
		Rooms:   generateKnownMap(),
		Compass: generateCompass(),
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/maze.html")
	if err != nil {
		log.Printf("Handler, Game, move, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, move, Error executing page: %v", err)
	}
}
