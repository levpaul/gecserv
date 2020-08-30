# Performance Matters


## Critical Optimisations Required

 - Currently when the propagator sends a full sync to a player it pulls all players from the relevant interestmaps and sends the PlayerT types as an event bus message to be handled by netpub. Netpub subsequently packs all the players into a flatbuffer and sends them. This is done for every player needing a full sync. Some flatbuffer caching of players per interestmap can be done here to save on recomputes for full syncs. Partial updates will suffer the same issue but at with smaller player lists but at a larger player scale.


### Things to perf test:

 - swap flatbuffers for messagepack
 - pass players as structs not pointers for player updates (from propagator.go to netpub.go)
 