package maze

//This handles the maze and compass specific functions.

import (
	"fmt"
	"html/template"
	"strings"

	"database/sql"

	"github.com/JohnDirewolf/capstone/database"
	"github.com/JohnDirewolf/capstone/shared/types"
)

// Global variables
var playerLocation int

// These are room constants for the player inventory, start room and goal room.
const (
	playerInventory = -1
	StartRoom       = 2
	GoalRoom        = 12
)

const (
	NorthDoor = iota // 0
	SouthDoor        //1
	WestDoor         // 2
	EastDoor         // 3
)

func Init() {
	//rest the database to the initial state.
	//For database functions, while they will often return an err, this is not used here. Errs are simply logged by the function that raises the error, but we return it if we want to use it.
	//Clear the database
	database.Clear()
	var room types.RoomData

	//Room playerInventory (-1) is a container for the user's inventory, this allows items to reference an existing room even if with the player.
	room = types.RoomData{
		Id:          playerInventory,
		Title:       "Player Inventory",
		Description: "Items the Player has in their inventory",
		Discovered:  false,
	}
	database.InsertRoom(room)

	//Set up rooms to initial state and the room's doors, but all doors are unlocked. Locks reference key items and so after the items are created we lock the door(s) we want locked.
	room = types.RoomData{
		Id:          0,
		Title:       "Bark Room",
		Description: "The room is covered in various types of tree bark. The smell is musky. The only exist is the way you entered to the East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          1,
		Title:       "Leaf Room",
		Description: "You find a room overgrown with leaves. Pushing through the folliage you find exits to the North, West, and East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          2,
		Title:       "Stone Room",
		Description: "You are in a room made of cleaved stones. You see doors to the North, West, and East.",
		Discovered:  true,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          3,
		Title:       "Grass Room",
		Description: "This room is filled with grass and smells natural and clean. You can make out doors to the North and West.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          4,
		Title:       "Dirt Room",
		Description: "You enter a room that is empty, the foor being ony dirt. There are exits to the North and East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          5,
		Title:       "Water Room",
		Description: "As you enter you fall into a pool of warm water. Swimming about you find exits to the North, South and West.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          6,
		Title:       "Oil Room",
		Description: "The air is hard to breath here as the room is knee deep in black oil. Wading through the room you find doors to the North and South. The Northern door is covered in gold inlay.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          7,
		Title:       "Wood Room",
		Description: "All around is worked wood creating a cosy cabin feel. There is only a door to the South.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          8,
		Title:       "Lava Room",
		Description: "The heat in this room is nearly unbarable as the floor is mostly lava. There seems to be no way forward, only the exit to the South.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          9,
		Title:       "Boil Room",
		Description: "The air is full of steam and the sound of roiling water. A simple bridge over the boiling water allows exit to the North and South.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          10,
		Title:       "Mud Room",
		Description: "You find a huge expanse of cracked mud, desperate for water. You can see exits South and East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          11,
		Title:       "Rust Room",
		Description: "There is the sound of clanking metals and steam flowing through old pipes in this room filed with rusted machinery. There are doors to the North and West.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["north"] = types.DoorData{Id: (room.Id * 4) + NorthDoor, RoomId: room.Id, Direction: "north"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          12,
		Title:       "Copper Room - Goal!",
		Description: "Huzzah! Entering this room made of copper and metal you see a large portal open. The way out! You can go back to the maze to the East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          13,
		Title:       "Gold Room",
		Description: "The room glitters with gold in all shapes and sizes, then you realize it is just fool's gold. You see doors through the faux horde to the South, West and East.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	room.Doors["east"] = types.DoorData{Id: (room.Id * 4) + EastDoor, RoomId: room.Id, Direction: "east"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          14,
		Title:       "Magma Room",
		Description: "A powerful heat radiates from this room, magma slowly shifing. There is no other exit then the door you came in to the West.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["west"] = types.DoorData{Id: (room.Id * 4) + WestDoor, RoomId: room.Id, Direction: "west"}
	database.InsertRoom(room)

	room = types.RoomData{
		Id:          15,
		Title:       "Granite Room",
		Description: "You enter a quarry of granite and stone. You only see a door back South.",
		Discovered:  false,
		Doors:       make(map[string]types.DoorData),
	}
	room.Doors["south"] = types.DoorData{Id: (room.Id * 4) + SouthDoor, RoomId: room.Id, Direction: "south"}
	database.InsertRoom(room)

	//Add items.
	database.InsertItem(types.ItemData{
		Id:          1,
		Name:        "Golden Key",
		Article:     "a ",
		Description: "A delicate golden ward key with a ruby in the bow.",
		Type:        "key",
		CurLocation: 8,
	})
	database.InsertItem(types.ItemData{
		Id:          2,
		Name:        "Magic Sword",
		Article:     "a ",
		Description: "A sword that glows blue. Along the blad are runes that says \"Goblin Scourge\"",
		Type:        "sword",
		CurLocation: 15,
	})
	//This sets up a lucky coin in the player's default inventory as item 0. It's not used in the maze.
	database.InsertItem(types.ItemData{
		Id:          0,
		Name:        "Lucky Coin",
		Article:     "a ",
		Description: "A ancient coin you found long ago and you feel has brought you luck.",
		Type:        "coin",
		CurLocation: playerInventory,
	})

	//All items created, so now lock our door(s).
	database.LockDoor((6*4)+NorthDoor, 1)

	//Creatures have a description that is a sentence used to describe the character and what they are doing.
	database.InsertCreature(types.CreatureData{
		Id:           1,
		Name:         "Wartz",
		Type:         "goblin",
		Description:  "A powerful looking goblin with a huge club is guarding the bridge.",
		IsAlive:      true,
		VanquishedBy: sql.NullInt16{Int16: 2, Valid: true},
		CurLocation:  9,
		Guards:       sql.NullInt16{Int16: (9 * 4) + NorthDoor, Valid: true},
	})
	//Set the door they are guarding
	database.GuardDoor((9 * 4) + NorthDoor)

	//Set player location to starting room
	playerLocation = StartRoom
}

func Move(direction types.UrlAction) types.SpecialStatus {
	//Each direction changes the playerLocation by a different value. If there is no direction, then playerLocation does not change.
	var special types.SpecialStatus
	//Get a RoomData from the Database
	room, _ := database.GetRoom(playerLocation)
	//Check for locked or guarded doors.
	if doorData, ok := room.Doors[string(direction)]; ok {
		if doorData.Locked {
			special.IsLocked = true
		} else if doorData.Guarded {
			special.IsGuarded = true
		} else {
			switch direction {
			case types.North:
				playerLocation += 4
			case types.South:
				playerLocation -= 4
			case types.West:
				playerLocation -= 1
			case types.East:
				playerLocation += 1
			}
		}
	}

	//This sets the room the player is in to discovered, regardless if the player moved or not.
	database.DiscoverRoom(playerLocation)

	return special
}

func GetItems() {
	//move all the items from the room into the player's inventory.
	itemList, _ := database.GetItemsByLocation(playerLocation)
	for _, item := range itemList {
		database.MoveItemToLocation(item.Id, playerInventory)
	}
}

func UseKey() {
	database.UnlockDoor(playerLocation)
}

func Attack() bool {
	//Check if the player has the proper vanquishing item for the creatures, if not they are defeated
	//Get the current creature from the player's Location
	creatureInfo := database.GetCreatureInLocation(playerLocation)
	//Check if the player has the vanquishing item and return result
	if database.DoesUserHaveItem(int(creatureInfo.VanquishedBy.Int16)) {
		//Player has the item and defeats the creature, set to vanquished and unguard the door.
		database.VanquishCreature(creatureInfo.Id)
		//If the creature is guarding a door, unguard it. (For if we have a creature that is just a foe but not a guard, currently there is only the goblin guard.)
		if creatureInfo.Guards.Valid {
			database.UnguardDoor(int(creatureInfo.Guards.Int16))
		}
		return true
	}
	return false
}

func GenerateKnownMap() template.HTML {
	//This runs through the map and all rooms showing discovered have their container added to the string.
	var knownMap strings.Builder
	roomsDiscovered, _ := database.GetDiscoveredRooms()
	for _, roomID := range roomsDiscovered {
		fmt.Fprintf(&knownMap, "<div class='room%d'><img src='/assets/images/r%d.jpg' alt='Maze Room' width='200' height='200' /></div>\n", roomID, roomID)
	}
	//Add the Player
	playerLocationTop := (((15 - playerLocation) / 4) * 200) + 25
	playerLocationLeft := ((playerLocation % 4) * 200) + 66
	//log.Printf("playerLocationLeft: %d", playerLocationLeft)
	fmt.Fprintf(&knownMap, "<div style='position:absolute;top:%dpx;left:%dpx'><img src='/assets/images/player.png' alt='Player!' width='67' height='150' /></div>\n", playerLocationTop, playerLocationLeft)
	return template.HTML(knownMap.String())
}

func GetPageInfo(special types.SpecialStatus) types.PageData {
	var compass strings.Builder
	var inventory strings.Builder
	var description strings.Builder
	var action strings.Builder
	var instructions template.HTML
	room, _ := database.GetRoom(playerLocation)
	//log.Printf("Room returned to GenerateCompass is: %v", room)
	if _, ok := room.Doors["north"]; ok {
		compass.WriteString(`<div class="n-arrow"><a href="/app?action=north"><img class="n-arrow" src="/assets/images/green_arrow_n.png" alt="Green Arrow North" width="100" height="100" /></a></div>`)
	} else {
		compass.WriteString(`<div class="n-arrow"><img class="n-arrow" src="/assets/images/red_arrow_n.png" alt="Red Arrow North" width="100" height="100" /></div>`)
	}
	if _, ok := room.Doors["south"]; ok {
		compass.WriteString(`<div class="s-arrow"><a href="/app?action=south"><img class="s-arrow" src="/assets/images/green_arrow_s.png" alt="Green Arrow South" width="100" height="100" /></a></div>`)
	} else {
		compass.WriteString(`<div class="s-arrow"><img class="s-arrow" src="/assets/images/red_arrow_s.png" alt="Red Arrow South" width="100" height="100" /></div>`)
	}
	if _, ok := room.Doors["west"]; ok {
		compass.WriteString(`<div class="w-arrow"><a href="/app?action=west"><img class="w-arrow" src="/assets/images/green_arrow_w.png" alt="Green Arrow West" width="100" height="100" /></a></div>`)
	} else {
		compass.WriteString(`<div class="w-arrow"><img class="w-arrow" src="/assets/images/red_arrow_w.png" alt="Red Arrow West" width="100" height="100" /></div>`)
	}
	if _, ok := room.Doors["east"]; ok {
		compass.WriteString(`<div class="e-arrow"><a href="/app?action=east"><img class="e-arrow" src="/assets/images/green_arrow_e.png" alt="Green Arrow East" width="100" height="100" /></a></div>`)
	} else {
		compass.WriteString(`<div class="e-arrow"><img class="e-arrow" src="/assets/images/red_arrow_e.png" alt="Red Arrow East" width="100" height="100" /></div>`)
	}

	//Show player's Inventory.
	itemsInInventory, _ := database.GetItemsByLocation(playerInventory)
	for index, item := range itemsInInventory {
		if index >= 1 {
			inventory.WriteString("<br />")
		}
		inventory.WriteString(item.Name)
	}

	//Description - First we check for the special cases for page information.
	//If the is the start we set extra text and the instructions set.
	if special.IsStart {
		description.WriteString(`The entrace slams shut behind you. You will have to look for a different exit!<br />`)
		instructions = getInstructions()
	}

	//Next if we moved and door was locked then we print warning text or add an action if the user has the key.
	//If we are in a room with a locked door (room 9) then check if the user has the key (1) and either add acton to use key or tell them the door is locked
	if special.IsLocked {
		//Check if the player has the golden key.
		if database.DoesUserHaveItem(1) {
			description.WriteString(`<span class="locked">Locked! Try the key!</span><br />`)
			//Add use action
			action.WriteString(`<div class="action"><a class="action" href="/app?action=use"><span class="action">Use Golden Key</span></a></div>`)
		} else {
			description.WriteString(`<span class="locked">Locked! The door you tried is locked.<br />Perhaps you can find a key?</span><br />`)
		}
	}

	//Once the key is used, we tell the user it worked.
	if special.Unlocked {
		description.WriteString(`<span class="unlocked">The key turns and unlocks the door!</span><br />`)
	}

	if special.IsGuarded {
		fmt.Fprintf(&description, `<span class="warning">Your foe blocks your path!</span><br />`)
	}

	if special.Vanquished {
		description.WriteString(`<span class="vanquished">You vanquished your foe!</span><br />`)
	}

	//Get The basic room Description
	description.WriteString(room.Description)

	//Checking for items in the room. If there are item(s) in the room we add an action, get item, and the item(s) are added to the description.
	itemsInLocation, _ := database.GetItemsByLocation(playerLocation)
	if len(itemsInLocation) >= 1 {
		description.WriteString("<br />In the room you see: ")
		for index, item := range itemsInLocation {
			if index >= 1 {
				description.WriteString(", ")
			}
			description.WriteString(item.Article)
			description.WriteString(item.Name)
		}
		//Action to get items
		action.WriteString(`<div class="action"><a class="action" href="/app?action=get"><span class="action">`)
		if len(itemsInLocation) == 1 {
			//Single item in room.
			action.WriteString(`Get `)
			action.WriteString(itemsInLocation[0].Name)
		} else {
			//More then one item in the room.
			action.WriteString(`Get all items.`)
		}
		action.WriteString(`</span></a></div>`)
	}

	//Checking for creature in the room. If there is a creature in the room we add an action, and add the creature descrition and room sepcific information to the description.
	creature := database.GetCreatureInLocation(playerLocation)
	if (creature != types.CreatureData{}) {
		//First check if the creature is already dead, skip and say so.
		if creature.IsAlive {
			fmt.Fprintf(&description, "<br />%s", creature.Description)
			//If the player does not have the ability to vanquish the creature, warn them. But they can still attack.
			if !database.DoesUserHaveItem(int(creature.VanquishedBy.Int16)) {
				fmt.Fprintf(&description, `<br /><span class="warning">Careful! You may not be able to vanquish the %s!</span>`, creature.Type)
			}
			//Action to attack items
			fmt.Fprintf(&action, `<div class="action"><a class="action" href="/app?action=attack"><span class="action">Attack %s!</span></a></div>`, creature.Type)
		} else {
			fmt.Fprintf(&description, "<br />You see a dead %s here.", creature.Type)
		}
	}

	//If we are in the exit room we have an action to escape the maze.
	if playerLocation == GoalRoom {
		action.WriteString(`<div class="action"><a class="action" href="/app?action=end"><span class="action">Escape the Maze!</span></a></div>`)
	}

	pageData := types.PageData{
		Title:        room.Title,
		Rooms:        GenerateKnownMap(),
		Compass:      template.HTML(compass.String()),
		Inventory:    template.HTML(inventory.String()),
		Description:  template.HTML(description.String()),
		Action:       template.HTML(action.String()),
		Instructions: instructions,
	}

	return pageData
}

func getInstructions() template.HTML {
	//This sends the Instructions HTML when we want to display them, specifically on the start of the game.
	return template.HTML(`
		<div class="instructions">
			<p class="instructions">Click on the Compass Arrow in the direction you wish to go!</p>
		</div>`)
}
