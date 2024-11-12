# Hatrock Dungeon Explorer!

## Description
This is an example of Go programming a server and database driven program as a very simple fantasy web based game. 
In the game you explore the "Hatrock Dungeon" finding a magic key and a magic sword to defeat a goblin and escape the dungeon's maze. 
This program creates a local server using ServeMux as a file server and handler for game actions.
A PostgreSQL database is then used to store game data.
All art assets were created by myself.

I wanted to create a fun little game that would pull together web server handling, database interfacing, and my own artistic ability into a fun web app in the Go language.

## Installation & Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/JohnDirewolf/capstone
   cd capstone
   ```

2. Install PostgreSQL if needed (see below.)

3. Install dependencies:
        ```bash
        go mod download
        ```

4. Create a .env file based on .env.example and configure your environment variables. (See below.)

5. Running the Application
   To run the application:
        ```bash
        go run .
        ```

## Required:
* Postgres (PostgreSQL) 14.13

### Install PostgreSQL 14.13 on Ubuntu:
To install PostgreSQL on Ubuntu, run the following commands in your terminal:
    ```bash
    sudo apt update
    sudo apt install postgresql
    ```

### Set up your PostgreSQL database connection string:
1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Open the .env file in your text editor:
    ```bash
    nano .env
    ```

3. Update the DB_URL variable with your PostgreSQL connection string:

postgres://<username>:<password>@<host>:<port>

Example format explanation:
- username: your PostgreSQL username
- password: your PostgreSQL password
- host: localhost (if running locally) or your database host
- port: 5432 (default PostgreSQL port)
- database_name: name of your database (game_database)

4. Save the file (in nano: Ctrl + O, then Ctrl + X to exit)

## Set up the database 'game_database' in PostgreSQL
1. Open your terminal and connect to PostgreSQL:

   ```bash
   psql -U your-username
   ```

2. Create the database by running:
   ```sql
   CREATE DATABASE game_database;
   ```

3. Exit the PostgreSQL prompt:
   ```sql
   \q
   ```

### Run Goose Migration to set up the database:
1. Install Goose if needed. In your terminal run:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Use the connection string you set up in a Goose migration command in the root directory for this program:
   ```bash
   goose -dir sql/schema ppostgres://<username>:<password>@<host>:<port>/game_database up
   ```
## Usage

After starting the app you open a browser to your localhost. There you will see the opening page. It gives the title and some lore. When reading just click start.
You will find your avatar in the first room of the dungeon maze, the little adventurer icon.
To the right of the maze board you will see a navagation compass, your inventory, and a basic room description.
Click on any of the green compass arrows to move in the maze. New rooms will be revealed.
You will find a locked door and a mean goblin blocking you path, you must find the key and sword.
When you enter a room if there is an action to take you will see it above the description.
Be careful with the goblin! If you don't have the sword it may end your adventure!.
Find the final room, get your treasure and leave the dungeon!
