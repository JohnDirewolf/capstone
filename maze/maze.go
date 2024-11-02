package maze

//This handles the maze and compass specific functions.

import (
	"fmt"
	"html/template"
	"strings"
)

type roomData struct {
	discovered bool
	exitNorth  bool
	exitSouth  bool
	exitWest   bool
	exitEast   bool
}

// Global variables
var theMaze map[int]roomData
var playerLocation int

// Direction constants
const North string = "north"
const South string = "south"
const West string = "west"
const East string = "east"

func Init() {
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

	//Set player location to starting room, 2
	playerLocation = 2
}

func Move(direction string) {
	//Each direction changes the playerLocation by a different value. If there is no direction, then playerLocation does not change.
	switch direction {
	case North:
		if theMaze[playerLocation].exitNorth {
			playerLocation += 4
		}
	case South:
		if theMaze[playerLocation].exitSouth {
			playerLocation -= 4
		}
	case West:
		if theMaze[playerLocation].exitWest {
			playerLocation -= 1
		}
	case East:
		if theMaze[playerLocation].exitEast {
			playerLocation += 1
		}
	}

	//This sets the room to discovered, regardless if the player moved or not.
	room := theMaze[playerLocation]
	room.discovered = true
	theMaze[playerLocation] = room
}

func GenerateKnownMap() template.HTML {
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

func GenerateCompass() template.HTML {
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
