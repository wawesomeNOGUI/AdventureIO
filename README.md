AdventureIO

I made this as a project to practice concurrent programming and gamemaking. I chose to remake Adventure by Warren Robinett into an online multiplayer game. I really enjoy collaborative games and thought it would be fun to see if adventure could work with multiple players. Turns out it does! It is a very chaotic, but fun experience to run around exploring the world together with a few other players, teaming up to fight dragons, dodge bats, and make your way to the end of the game together.

## Technology
I use webRTC with a fully reliable datachannel and an ordered, no resend datachannel to communicate between the clients and server.


In the server, there is a goroutine for every client connection, room updates, game state serialization, and broadcasting updates. Because Adventure is split into distinct rooms, it would also be possible to have a goroutine handling updates for each individual room. This concurrency necessitated careful controls over goroutine memory access. I utilized mutexes for controlling access to shared memory such as modifying entity positions, and for moving entities between rooms. I also utilized go channels to synchronize the initiation of the player dataChannel update functions in the client connection goroutine.
