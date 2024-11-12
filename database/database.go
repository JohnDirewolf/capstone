package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/JohnDirewolf/hatrock_dungeon/shared/types"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbURL string

const dbURLSuffix string = "/game_database?sslmode=disable"

var heart *sql.DB

// --------------------- Database functions -------------------- //
// Many database functions return an error but we just log the error and do not catch the error in the calling functions. But the error is there if we do want it.
func Init() error {
	var err error

	//Get the database connection string from .env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	heart, err = sql.Open("postgres", dbURL+dbURLSuffix)
	if err != nil {
		log.Printf("Database, Init: Error in connecting to database: %v", err)
		return err
	}

	// Verify we can connect
	err = heart.Ping()
	if err != nil {
		heart.Close()
		log.Printf("Database, Init: Failed to ping database: %v", err)
		return err
	}
	return nil
}

func Clear() error {
	_, err := heart.Exec("DELETE FROM creatures;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing creatures: %v", err)
	}
	_, err = heart.Exec("DELETE FROM doors;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing doors: %v", err)
	}
	_, err = heart.Exec("DELETE FROM items;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing items: %v", err)
	}
	_, err = heart.Exec("DELETE FROM rooms;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing rooms: %v", err)
	}
	return err
}

func Close() error {
	err := heart.Close()
	if err != nil {
		log.Printf("Database, Close: Error closing database: %v", err)
	}
	return err
}

// ////////// Creature functions //////////
func InsertCreature(creatureInfo types.CreatureData) error {
	_, err := heart.Exec(`
		INSERT INTO creatures (
			id, 
			name,
			type,
			description,
			is_alive,
			vanquished_by,
			cur_location,
			guards
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		creatureInfo.Id,
		creatureInfo.Name,
		creatureInfo.Type,
		creatureInfo.Description,
		creatureInfo.IsAlive,
		creatureInfo.VanquishedBy,
		creatureInfo.CurLocation,
		creatureInfo.Guards)
	if err != nil {
		log.Printf("Database, InsertCreature: Error inserting values for creature: %v", err)
	}
	return err
}

func GetCreatureInLocation(location int) types.CreatureData {
	// Currently a room can only have 1 creature
	var creature types.CreatureData

	creatureRecord := heart.QueryRow("SELECT id, name, type, description, is_alive, vanquished_by, cur_location, guards FROM creatures WHERE cur_location=$1;", location)

	err := creatureRecord.Scan(&creature.Id, &creature.Name, &creature.Type, &creature.Description, &creature.IsAlive, &creature.VanquishedBy, &creature.CurLocation, &creature.Guards)
	if err != nil {
		//ErrNoRows is not a real error, it means there is no creature in the room. We do not log it, we just return the empty creature. Other errors are logged.
		if err != sql.ErrNoRows {
			log.Printf("Database, GetCreatureInLocation: Error getting creature information.: %v", err)
		}
		return types.CreatureData{}
	}
	return creature
}

func VanquishCreature(creatureId int) error {
	_, err := heart.Exec("UPDATE creatures SET is_alive=false WHERE id=$1;", creatureId)
	if err != nil {
		log.Printf("Database, VanquishCreature: Error vanquishing creature.: %v", err)
	}
	return err
}

// -------------------- Item functions -------------------- //
func InsertItem(itemInfo types.ItemData) error {
	_, err := heart.Exec(`
		INSERT INTO items (
			id, 
			name,
			article,
			description,
			type,
			cur_location
		) VALUES ($1, $2, $3, $4, $5, $6);`,
		itemInfo.Id,
		itemInfo.Name,
		itemInfo.Article,
		itemInfo.Description,
		itemInfo.Type,
		itemInfo.CurLocation)
	if err != nil {
		log.Printf("Database, InsertItem: Error inserting values for item: %v", err)
	}
	return err
}

func GetItemByID(itemID int) (types.ItemData, error) {
	var item types.ItemData

	itemRecord := heart.QueryRow("SELECT id, name, article, description, type, cur_location FROM items WHERE id=$1;", itemID)

	err := itemRecord.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.Type, &item.CurLocation)
	if err != nil {
		log.Printf("Database, GetItemByID: Error getting data for item: %v", err)
		return types.ItemData{}, err
	}

	return item, nil
}

func GetItemByName(itemName string) (types.ItemData, error) {
	var item types.ItemData

	itemRecord := heart.QueryRow("SELECT id, name, article, description, type, cur_location FROM items WHERE name=$1;", itemName)

	err := itemRecord.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.Type, &item.CurLocation)
	if err != nil {
		log.Printf("Database, GetItemByID: Error getting data for item: %v", err)
		return types.ItemData{}, err
	}

	return item, nil
}

func GetItemsByLocation(location int) ([]types.ItemData, error) {
	var items []types.ItemData
	var item types.ItemData
	rows, err := heart.Query("SELECT id, name, article, description, type, cur_location FROM items WHERE cur_location=$1;", location)
	if err != nil {
		log.Printf("Database, GetItemsByLocation: Error getting list of items in locatoin: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.Type, &item.CurLocation)
		if err != nil {
			log.Printf("Database, GetItemsByLocation: Error building list of items: %v", err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func MoveItemToLocation(itemId int, location int) error {
	//location -1 is the Player's Inventory.
	_, err := heart.Exec("UPDATE items SET cur_location=$1 WHERE id=$2;", location, itemId)
	if err != nil {
		log.Printf("Database, MoveItemToLocation: Error moving item to location: %v", err)
	}
	return err
}

func DoesUserHaveItem(itemID int) bool {
	//All this does is check if the user has an item in their inventory.
	var curLocation int

	itemRecord := heart.QueryRow("SELECT cur_location FROM items WHERE id=$1;", itemID)

	err := itemRecord.Scan(&curLocation)
	if err != nil {
		log.Printf("Database, DoesUserHaveItem: Error finding if user has item: %v", err)
		return false
	}
	//Player inventory is room -1
	if curLocation == -1 {
		return true
	}
	return false
}

// -------------------- Door functions --------------------- //
func UnlockDoor(playerLocation int) error {
	//Get the doors from the player location and get the one that is locked. While we have only 1 locked door,
	//this generalizes it if we what to change the door or add more locked doors
	var doorId int

	doorRecord := heart.QueryRow("SELECT id FROM doors WHERE room_id=$1 AND locked=TRUE;", playerLocation)
	err := doorRecord.Scan(&doorId)
	if err != nil {
		log.Printf("Database, UnlockDoor: Error getting locked door in player loction: %v", err)
		return err
	}

	//We could make sure the character has the correct key, but currently there is only 1 key and locked door and that logic is decided when the action is given.
	_, err = heart.Exec("UPDATE doors SET locked=false WHERE id=$1;", doorId)
	if err != nil {
		log.Printf("Database, UnlockDoor: Error unlocking door: %v", err)
	}
	return err
}

func LockDoor(doorId int, keyId int) error {
	//There is currently just the one door, but we accept a doorId from maze if we add more locked doors.
	//We set locked to true and also the item_id of the key.
	_, err := heart.Exec("UPDATE doors SET locked=true, key_id=$1 WHERE id=$2;", keyId, doorId)
	if err != nil {
		log.Printf("Database, LockDoor: Error locking door: %v", err)
	}
	return err
}

func GuardDoor(doorId int) error {
	//We set the door the creature is guarding.
	_, err := heart.Exec("UPDATE doors SET guarded=true WHERE id=$1;", doorId)
	if err != nil {
		log.Printf("Database, GuardDoor: Error guarding door: %v", err)
	}
	return err
}

func UnguardDoor(doorId int) error {
	_, err := heart.Exec("UPDATE doors SET guarded=false WHERE id=$1;", doorId)
	if err != nil {
		log.Printf("Database, UnguardDoor: Error unguarding door: %v", err)
	}
	return err
}

// -------------------- Room functions -------------------- //
func GetRoom(roomID int) (types.RoomData, error) {
	var room types.RoomData
	room.Doors = make(map[string]types.DoorData)
	var door types.DoorData
	// Temporary variables for scanning
	//var northExists, northLocked, southExists, southLocked, westExists, westLocked, eastExists, eastLocked bool
	//var northKeyID, southKeyID, westKeyID, eastKeyID sql.NullInt16

	roomRecord := heart.QueryRow(`SELECT id, title, description, discovered FROM rooms WHERE id=$1;`, roomID)
	err := roomRecord.Scan(&room.Id, &room.Title, &room.Description, &room.Discovered)
	if err != nil {
		log.Printf("Database, GetRoom: Error getting main data for room: %v", err)
		return types.RoomData{}, err
	}

	doorRows, err := heart.Query(`SELECT id, room_id, direction, locked, guarded, key_id FROM doors WHERE room_id=$1`, roomID)
	for doorRows.Next() {
		err := doorRows.Scan(&door.Id, &door.RoomId, &door.Direction, &door.Locked, &door.Guarded, &door.KeyId)
		if err != nil {
			log.Printf("Database, GetRoom: Error getting list of doors: %v", err)
			return types.RoomData{}, err
		}
		room.Doors[door.Direction] = door
	}
	return room, nil
}

func GetDiscoveredRooms() ([]int, error) {
	var rooms []int
	var roomID int
	rows, err := heart.Query("SELECT id FROM rooms WHERE discovered = TRUE;")
	if err != nil {
		log.Printf("Database, GetDiscoveredRooms: Error getting list of rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&roomID)
		if err != nil {
			log.Printf("Database, GetDiscoveredRooms: Error building list of room IDs: %v", err)
			return nil, err
		}
		rooms = append(rooms, roomID)
	}

	return rooms, nil
}

func InsertRoom(roomInfo types.RoomData) error {
	var door types.DoorData
	//Insert the main room data.
	_, err := heart.Exec(`INSERT INTO rooms (id, title, description, discovered) 
						  VALUES ($1, $2, $3, $4);`,
		roomInfo.Id, roomInfo.Title, roomInfo.Description, roomInfo.Discovered)
	if err != nil {
		log.Printf("Database, InsertRoom: Error inserting main values for room: %v", err)
	}
	//Insert the door data for the room.
	for _, door = range roomInfo.Doors {
		_, err := heart.Exec(`INSERT INTO doors (id, room_id, direction, locked, key_id)
							  VALUES ($1, $2, $3, $4, $5);`,
			door.Id, door.RoomId, door.Direction, door.Locked, door.KeyId)
		if err != nil {
			log.Printf("Database, InsertRoom: Error inserting information for a door: %v", err)
		}
	}

	return err
}

func DiscoverRoom(roomID int) error {
	_, err := heart.Exec("UPDATE rooms SET discovered=TRUE WHERE id=$1;", roomID)
	if err != nil {
		log.Printf("Database, DiscoverRoom: Error updating discovered for room: %v", err)
	}
	return err
}
