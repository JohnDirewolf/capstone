package database

import (
	"database/sql"
	"log"

	"github.com/JohnDirewolf/capstone/shared/types"

	_ "github.com/lib/pq"
)

const dbURL string = "postgres://postgres:postgres@localhost:5432/game_database?sslmode=disable"

var heart *sql.DB

func Init() error {
	var err error
	heart, err = sql.Open("postgres", dbURL)
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
	_, err := heart.Exec("DELETE FROM doors;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing doors: %v", err)
	}
	_, err = heart.Exec("DELETE FROM rooms;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing rooms: %v", err)
	}
	_, err = heart.Exec("DELETE FROM items;")
	if err != nil {
		log.Printf("Database, Clear: Error clearing items: %v", err)
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

// ////////// Item functions //////////
func InsertItem(itemInfo types.ItemData) error {
	_, err := heart.Exec(`
		INSERT INTO items (
			id, 
			name,
			article,
			description,
			curLocation
		) VALUES ($1, $2, $3, $4, $5);`,
		itemInfo.Id,
		itemInfo.Name,
		itemInfo.Article,
		itemInfo.Description,
		itemInfo.CurLocation)
	if err != nil {
		log.Printf("Database, InsertItem: Error inserting values for item: %v", err)
	}
	return err
}

func GetItemByID(itemID int) (types.ItemData, error) {
	var item types.ItemData

	itemRecord := heart.QueryRow("SELECT id, name, article, description, curLocation FROM items WHERE id=$1;", itemID)

	err := itemRecord.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.CurLocation)
	if err != nil {
		log.Printf("Database, GetItem: Error getting data for item: %v", err)
		return types.ItemData{}, err
	}

	return item, nil
}

func GetItemByName(itemName string) (types.ItemData, error) {
	var item types.ItemData

	itemRecord := heart.QueryRow("SELECT id, name, article, description, curLocation FROM items WHERE name=$1;", itemName)

	err := itemRecord.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.CurLocation)
	if err != nil {
		log.Printf("Database, GetItem: Error getting data for item: %v", err)
		return types.ItemData{}, err
	}

	return item, nil
}

func GetItemsByLocation(location int) ([]types.ItemData, error) {
	var items []types.ItemData
	var item types.ItemData
	rows, err := heart.Query("SELECT id, name, article, description, curLocation FROM items WHERE curLocation=$1;", location)
	if err != nil {
		log.Printf("Database, GetItemsByLocation: Error getting list of items: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&item.Id, &item.Name, &item.Article, &item.Description, &item.CurLocation)
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
	_, err := heart.Exec("UPDATE items SET Curlocation=$1 WHERE id=$2;", location, itemId)
	if err != nil {
		log.Printf("Database, MoveItemToLocation: Error updating location for item: %v", err)
	}
	return err
}

func DoesUserHaveKey() bool {
	//All this does is check if the user has the Golden Key.
	var curLocation int

	keyRecord := heart.QueryRow("SELECT curLocation FROM items WHERE id=1;")

	err := keyRecord.Scan(&curLocation)
	if err != nil {
		log.Printf("Database, DoesUserHaveKey: Error getting data for key: %v", err)
		return false
	}
	if curLocation == -1 {
		return true
	}
	return false
}

// ////////// Room functions //////////
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

	doorRows, err := heart.Query(`SELECT id, room_id, direction, locked, key_id FROM doors WHERE room_id=$1`, roomID)
	for doorRows.Next() {
		err := doorRows.Scan(&door.Id, &door.RoomId, &door.Direction, &door.Locked, &door.KeyId)
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
