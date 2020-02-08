## WebRTC MMO Server

This is a toy project to experience developing a UDP game server for an action MMO style game. The idea is to have MMO type mechanics mixed with a fast paced "real-time" action PVP mechanics. Basically mix FPS netcode with runescape.

# TODO: 
 7 Feb
 - Impl propagation system
 - Impl player movement system
 - Write benchmarking software


## BIG ARCH CONCRETE

 - Input/output with players is handled via event bus
 - Game state is managed via ECS design

## BIG ARCH WOOD

 - Single map open world - divided into sectors for horizontal scaling
 - "Interest" zones for player updates made by dividing map sectors to locales and taking surroundings
 - Client/Server prediction + Small Lag compensation UDP net design
 - Separate servers for non-gameworld systems, like chat
 
## BIG ARCH VAPEN8TION

 - Persistence layer (probably dynamo + redis caches?)
 - Are _all_ C in ECS flatbuffer managed?

# Network

### Brief Architecture

- Each play connects via WebRTC/UDP
- Players sync game state updates with server (60tick? or 20tick?)
- Players send input commands to server
- Server processes input commands at 60fps
- Server sends partial updates to each client at 60fps
- If server has no input from client then server "predicts" client action instead, by repeating previous input
- On subsequent old packets coming in, server can replay them

### Brief Details

 - Client sends inputs with gametick #
 - If client action includes hit-tests for attacks, server rolls back target to client's gametick for hit testing
 - If client hits based on what "they saw" then rewrite game state to accept hit
