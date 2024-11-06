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
    discovered BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE doors (
    id INTEGER PRIMARY KEY,
    room_id INTEGER NOT NULL, 
    direction TEXT CHECK (direction IN ('north', 'south', 'west', 'east')),
    locked BOOLEAN NOT NULL DEFAULT false,
    key_id INTEGER NULL,  
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (key_id) REFERENCES items(id)
);

-- +goose Down
DROP TABLE doors;
DROP TABLE rooms;
DROP TABLE items;