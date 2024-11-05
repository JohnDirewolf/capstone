-- +goose Up
CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL DEFAULT 'Item Name',
    article TEXT NOT NULL DEFAULT 'an ',
    description TEXT NOT NULL DEFAULT 'A Mysterious Object',
    curLocation INTEGER
);

CREATE TABLE rooms (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL DEFAULT 'Maze Runner!',
    description TEXT NOT NULL DEFAULT 'No Description Provided',
    discovered Boolean NOT NULL DEFAULT false,
    doorNorth BOOLEAN NOT NULL DEFAULT false,
    doorNorthLocked BOOLEAN NOT NULL DEFAULT false,
    doorNorthKey_id INTEGER,  
    doorSouth BOOLEAN NOT NULL DEFAULT false,
    doorSouthLocked BOOLEAN NOT NULL DEFAULT false,
    doorSouthKey_id INTEGER, 
    doorWest BOOLEAN NOT NULL DEFAULT false,
    doorWestLocked BOOLEAN NOT NULL DEFAULT false,
    doorWestKey_id INTEGER, 
    doorEast BOOLEAN NOT NULL DEFAULT false,
    doorEastLocked BOOLEAN NOT NULL DEFAULT false,
    doorEastKey_id INTEGER,
    FOREIGN KEY (doorNorthKey_id) REFERENCES items(id),
    FOREIGN KEY (doorSouthKey_id) REFERENCES items(id),
    FOREIGN KEY (doorWestKey_id) REFERENCES items(id),
    FOREIGN KEY (doorEastKey_id) REFERENCES items(id)
);

-- +goose Down
DROP TABLE rooms;
DROP TABLE items;