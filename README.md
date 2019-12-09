# Interactive Fiction Gaming Engine

A framework for Web based games in the spirit of [MUD's](https://en.wikipedia.org/wiki/MUD), [Interactive Fiction](https://en.wikipedia.org/wiki/Interactive_fiction), and Choose Your Own Adventure games. 

## PROGRESS

This is a glorified chatroom right now. See [NOTES](#NOTES) for a list of hopes and dreams.

## DEPENDENCIES

See the `go.mod` for the GOLANG version used during development and any dependencies. The goal is to keep
dependencies to a minimum.

## DEVELOPMENT

See the build and run sections to the get the game servers running. The Game Frontend is available in your browser
at [localhost:8080](localhost:8080) by default. 

### Visual Studio Code

A `launch.json` configuration is provided that will allow you to run the project with a debugger using the `Start Debug` (F5 shortcut by default) and `Start without Debugging` menu items.

### CLI

From the project root:

```bash
go run .
```

## BUILD

```bash
./bin/build.sh
```

The `bin/build.sh` script can take two arguments, operating system and architecture, for cross compiling. The binary will be named with the os and architecture. 

Example:
```bash
./bin/build.sh windows 386

jake@devbox:~/iff$ ls bin/
build.sh  iff-game-server iff-game-server-windows-386
```

OR

Manually build the binary:

```bash
go build -o bin/iff-game-server .
```

## RUN

```bash
./bin/iff-game-server
```

## NOTES

### Mission

- Web Server
    - Serves the client (Frontend) for the game
- Generic Game Client
    - Web-based
    - Text heavy (duh)
    - Controlled with text commands & links
    - Customizable GUI elements that map to commands 
    - Web socket communication with Game Server
- Game Server 
    - Authorative
    - Supports Single & Multiplayer Game Worlds
    - Persistence
    - Web socket communication with Game Client
- Web Based Authoring Tools for creating and extending the game
- Embedded scripting language for most (all?) game logic that can be hot loaded/reloaded during runtime