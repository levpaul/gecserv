## IdolScape

This is a toy project to experience developing a UDP game server for an action MMO style game. The idea is to have MMO type mechanics mixed with a fast paced "real-time" action PVP mechanics. Basically mix FPS netcode with runescape.

### Overall Architecture

- Each play connects via WebRTC/UDP
- Players sync game state updates with server (60tick)
- Players send input commands to server
- Server processes input commands at 60fps
- Server sends partial updates to each client at 60fps
- If server has no input from client then server "predicts" client action instead, by repeating previous input
- On subsequent old packets coming in, server can replay them

### Details

 - Client sends inputs with gametick #
 - If client action includes hit-tests for attacks, server rolls back target to client's gametick for hit testing
 - If client hits based on what "they saw" then rewrite game state to accept hit
 
### Add
 - Add a maximum rollback limit for attackers - they will need to compensate themselves after such a limit
 