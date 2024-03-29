package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	//"os"
	//"math"
	//"reflect"
	//"encoding/binary"

	"github.com/gorilla/websocket"

	"github.com/pion/webrtc/v3"
	"github.com/wawesomeNOGUI/webrtcGameTemplate/signal"
)

// Game Updates Map
// var Updates sync.Map
var NumberOfPlayers int

// Concurrent Safe Maps of Datachannels for broadcasting or sending messages between players
// DataChannelContainer in messaging.go
var reliableChans DataChannelContainer = DataChannelContainer{chans: make(map[string]*webrtc.DataChannel)}
var unreliableChans DataChannelContainer = DataChannelContainer{chans: make(map[string]*webrtc.DataChannel)}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var api *webrtc.API

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade: ", err)
		return
	}
	defer c.Close()
	fmt.Println("User connected from: ", c.RemoteAddr())

	//===========This Player's Variables===================
	var playerTag string
	var playerPtr *Player
	var room *Room
	// var pRoomMutex sync.Mutex

	defer func() {
		playerPtr.room.Entities.DeleteEntity(playerTag) 
	}()

	// lets UDP chan onOpen and reliable chan onOpen know that the player has been fully setup in the OnICEConnectionStateChange
	playerReady := make(chan bool)
	//===========WEBRTC====================================
	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	//Setup dataChannel to act like UDP with ordered messages (no retransmits)
	//with the DataChannelInit struct
	var udpPls webrtc.DataChannelInit
	var retransmits uint16 = 0

	//DataChannel will drop any messages older than
	//the most recent one received if ordered = true && retransmits = 0
	//This is nice so we can always assume client
	//side that the message received from the server
	//is the most recent update, and not have to
	//implement logic for handling old messages
	var ordered = true

	udpPls.Ordered = &ordered
	udpPls.MaxRetransmits = &retransmits

	// Create a datachannel with label 'UDP' and options udpPls
	dataChannel, err := peerConnection.CreateDataChannel("UDP", &udpPls)
	if err != nil {
		panic(err)
	}

	//Create a reliable datachannel with label "TCP" for all other communications
	reliableChannel, err := peerConnection.CreateDataChannel("TCP", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			//Store a new x and y for this player
			NumberOfPlayers++
			playerTag = strconv.Itoa(NumberOfPlayers)
			fmt.Println(playerTag)

			playerPtr = newPlayer(playerTag, 50, 50)

			// start player in default room
			v, ok := Rooms.Load("r1")
			if !ok {
				fmt.Println("Couldn't find room")
			}
			room = v.(*Room)

			playerPtr.room = room
		
			//Store a pointer to a Player Struct in the default room
			room.Entities.StoreEntity(playerTag, playerPtr)

			// go func() {
			// 	for {
			// 		// room = <-tmpPlayer.roomChangeChan //WallCheck will first set room to nil so player can't move while changing rooms
			// 		tmp := <-tmpPlayer.roomChangeChan 
			// 		pRoomMutex.Lock()
			// 		room = tmp
			// 		pRoomMutex.Unlock()
			// 	}
			// }()

			// send two playerreadies, one for each datachannel we're opening
			playerReady <- true
			playerReady <- true

			fmt.Println("stored player")
		} else if connectionState == webrtc.ICEConnectionStateDisconnected || connectionState == webrtc.ICEConnectionStateClosed {
			playerPtr.room.Entities.DeleteEntity(playerTag) 
			fmt.Println("Deleted Player")

			reliableChans.DeletePlayerChan(playerTag)
			unreliableChans.DeletePlayerChan(playerTag)

			err := peerConnection.Close() //deletes all references to this peerconnection in mem and same for ICE agent (ICE agent releases the "closed" status)
			if err != nil {               //https://www.w3.org/TR/webrtc/#dom-rtcpeerconnection-close
				fmt.Println(err)
			}
		}
	})

	//====================No retransmits, ordered dataChannel=======================
	// Register channel opening handling
	dataChannel.OnOpen(func() {
		<-playerReady
		unreliableChans.AddPlayerChan(playerTag, dataChannel)
		/*
		for {
			time.Sleep(time.Millisecond * 50) //50 milliseconds = 20 updates per second
			//20 milliseconds = ~60 updates per second

			//fmt.Println(UpdatesString)
			// Send the message as text so we can JSON.parse in javascript
			sendErr := dataChannel.SendText(UpdatesString)
			if sendErr != nil {
				fmt.Println("data send err", sendErr)
				break
			}
		}
		*/
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		playerPtr.mu.Lock()
		defer playerPtr.mu.Unlock()

		playerPtr.room.Entities.mu.Lock()
		defer playerPtr.room.Entities.mu.Unlock()
		
		//can use non concurrent safe methods on entities below here

		roomNumIndex := strings.Index(string(msg.Data), ",")
		if roomNumIndex == -1 {
			return
		}
		getRoomFromMsg := string(msg.Data[:roomNumIndex])
		if getRoomFromMsg != playerPtr.room.roomKey  {
			return
		}

		msg.Data = msg.Data[roomNumIndex + 1:]

		if msg.Data[0] == 'X' { //88 = "X"
			if playerPtr.owner != nil {
				return
			}

			x, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			// Edge of screen
			if x < 0 {
				x = 0
			} else if x > 160 - playerPtr.width {
				x = 160 - playerPtr.width
			}	

			// Move Owned Item (playerPtr.held is an EntitiyInterface)
			if playerPtr.held != nil {
				playerPtr.held.Update(x - playerPtr.X, 0)
			}
			
			playerPtr.X = x
		} else if msg.Data[0] == 'Y' { //89 = "Y"		
			if playerPtr.owner != nil {
				return
			}

			y, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			if y < 0 {
				y = 0
			} else if y > 105 - playerPtr.height {
				y = 105 - playerPtr.height
			}

			// Move Owned Item (playerPtr.held is an EntitiyInterface)
			if playerPtr.held != nil {
				playerPtr.held.Update(0, y - playerPtr.Y)
			}

			playerPtr.Y = y
		}
	})

	//==============================================================================

	//=========================Reliable DataChannel=================================
	// Register channel opening handling
	reliableChannel.OnOpen(func() {
		<-playerReady
		reliableChans.AddPlayerChan(playerTag, reliableChannel)

		//Send Client their playerTag so they know who they are in the Updates Array
		sendErr := reliableChannel.SendText("T" + playerTag)
		if sendErr != nil {
			panic(err)
		}
	})

	// Register message handling (Data all served as a bytes slice []byte)
	// for user controls
	reliableChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		playerPtr.mu.Lock()
		defer playerPtr.mu.Unlock()

		playerPtr.room.Entities.mu.Lock()
		defer playerPtr.room.Entities.mu.Unlock()
		
		//can use non concurrent safe methods on entities below here

		roomNumIndex := strings.Index(string(msg.Data), ",")
		if roomNumIndex == -1 {
			return
		}
		getRoomFromMsg := string(msg.Data[:roomNumIndex])
		if getRoomFromMsg != playerPtr.room.roomKey {
			return
		}

		msg.Data = msg.Data[roomNumIndex + 1:]

		if msg.Data[0] == 'X' { //88 = "X"
			if playerPtr.owner != nil {
				return
			}

			x, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			// Edge of screen
			if x < 0 {
				x = 0
			} else if x > 160 - playerPtr.width {
				x = 160 - playerPtr.width
			}	

			// Move Owned Item (playerPtr.held is an EntitiyInterface)
			if playerPtr.held != nil {
				playerPtr.held.Update(x - playerPtr.X, 0)
			}
			
			playerPtr.X = x
		} else if msg.Data[0] == 'Y' { //89 = "Y"
			if playerPtr.owner != nil {
				return
			}

			y, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			if y < 0 {
				y = 0
			} else if y > 105 - playerPtr.height {
				y = 105 - playerPtr.height
			}

			// Move Owned Item (playerPtr.held is an EntitiyInterface)
			if playerPtr.held != nil {
				playerPtr.held.Update(0, y - playerPtr.Y)
			}

			playerPtr.Y = y
		} else if msg.Data[0] == 'D' {
			//dropped item
			if playerPtr.held != nil {
				playerPtr.room.Entities.entities[playerPtr.held.Key()] = playerPtr.held
			    playerPtr.held.SetOwner(nil)
				playerPtr.held.SetRoom(room)
				playerPtr.held = nil
			}
		} else if msg.Data[0] == 'P' && playerPtr.held == nil {
			// Picked Up Item
			msg.Data = msg.Data[1:]
			dataSlice := strings.Split(string(msg.Data), ",")
			
			sX := dataSlice[0]
			sY := dataSlice[1]
			sDir := dataSlice[2]

			hitX, err := strconv.ParseFloat(sX, 64)
			if err != nil {
				fmt.Println(err)
			}
			hitY, err := strconv.ParseFloat(sY, 64) 
			if err != nil {
				fmt.Println(err)
			}

			gotItem, _ := playerPtr.room.Entities.nonConcurrentSafeTryPickUpItem(playerPtr, hitX, hitY)

			if gotItem {
				if sDir == "" {
					sDir = "l"
					playerPtr.held.SetX(playerPtr.X - 10)
					playerPtr.held.SetY(playerPtr.Y)
				}
				if sDir[0] == 'l' {
					playerPtr.held.SetX(playerPtr.held.GetX() - 4)
				} else if sDir[0] == 'r' {
					playerPtr.held.SetX(playerPtr.held.GetX() + 4)
				}

				if strings.Contains(sDir, "u") {
					playerPtr.held.SetY(playerPtr.held.GetY() - 4)
				} else if strings.Contains(sDir, "d") {
					playerPtr.held.SetY(playerPtr.held.GetY() + 4)
				}

				// Send this player the item key and offset so they can render it with no delay clientside
				str := "I" + playerPtr.held.Key() + "," + fmt.Sprintf("%.1f", playerPtr.held.GetX()-playerPtr.X)  + "," + fmt.Sprintf("%.1f", playerPtr.held.GetY()-playerPtr.Y)
				reliableChans.SendToPlayer(playerTag, str)
			}
		} else if msg.Data[0] == 'U' { // player sent username
			reliableChans.Broadcast("U" + playerTag + "," + string(msg.Data[1:]))
		}
	})

	//==============================================================================

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	fmt.Println(*peerConnection.LocalDescription())

	//Send the SDP with the final ICE candidate to the browser as our offer
	err = c.WriteMessage(1, []byte(signal.Encode(*peerConnection.LocalDescription()))) //write message back to browser, 1 means message in byte format?
	if err != nil {
		fmt.Println("write:", err)
	}

	//Wait for the browser to send an answer (its SDP)
	msgType, message, err2 := c.ReadMessage() //ReadMessage blocks until message received
	if err2 != nil {
		fmt.Println("read:", err)
	}

	answer := webrtc.SessionDescription{}

	signal.Decode(string(message), &answer) //set answer to the decoded SDP
	fmt.Println(answer, msgType)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	//=====================Trickle ICE==============================================
	//Make a new struct to use for trickle ICE candidates
	var trickleCandidate webrtc.ICECandidateInit
	var leftBracket uint8 = 123 //123 = ascii value of "{"

	for {
		_, message, err2 := c.ReadMessage() //ReadMessage blocks until message received
		if err2 != nil {
			fmt.Println("read:", err)
		}

		//If staement to make sure we aren't adding websocket error messages to ICE
		if message[0] == leftBracket {
			//Take []byte and turn it into a struct of type webrtc.ICECandidateInit
			//(declared above as trickleCandidate)
			err := json.Unmarshal(message, &trickleCandidate)
			if err != nil {
				fmt.Println("errorUnmarshal:", err)
			}

			fmt.Println(trickleCandidate)

			err = peerConnection.AddICECandidate(trickleCandidate)
			if err != nil {
				fmt.Println("errorAddICE:", err)
			}
		}

	}

}

// Sends current game state unreliably to all players
// (seemed better to just put in gameLoop because of weird race condition between SerializeEntities and UpdateEntities kept showing up)
func sendGameStateUnreliableLoop() {
	for {
		time.Sleep(time.Millisecond * 50) //50 milliseconds = 20 updates per second

		Rooms.Range(func(rk, rv interface{}) bool {
			switch z := rv.(type) {
			case *Room:
				//z.sendGameStateUnreliableChan <- true // let room update entities
				// // here z is a pointer to a Room
				s := z.Entities.SerializeEntities()

				for k, _ := range z.Entities.Players() {
					unreliableChans.SendToPlayer(k, s)
				}

			default:
				// no match; here z has the same type as v (interface{})
			}	
			return true
		})
	}
}

// Game Vars

// will contain items that can be picked up by players (mutex)
// ItemContainer & Item defined in types.go
//var strayItems ItemContainer = ItemContainer{items: make(map[string]Item)} 
// will contain items with the key being the item, and each item has an Owner tag set to the playerTag who owns it
//var ownedItems ItemContainer = ItemContainer{items: make(map[string]Item)}

// var entities EntityContainer = EntityContainer{entities: make(map[string]Entity)}
var Rooms sync.Map
func initGameVars() {
	InitializeRooms(&Rooms)
	InitializeEntities(&Rooms)

	fmt.Println("Game Ready. \n\n\n")
}

// All server orchestrated game logic
func gameLoop() {
	for {
		time.Sleep(time.Millisecond * 16)  // 16 ms is a little faster than 60 updates per second

		// Update Rooms
		Rooms.Range(func(k, v interface{}) bool {
			switch z := v.(type) {
			case *Room:
				//z.updateEntitiesChan <- true // let room update entities
				// // here z is a pointer to a Room
				z.updateFunc(z)

				// // then send updates to players in that room
				// s := z.Entities.SerializeEntities()

				// for k, _ := range z.Entities.Players() {
				// 	unreliableChans.SendToPlayer(k, s)
				// }
			default:
				// no match; here z has the same type as v (interface{})
			}	
			return true
		})
	}
}

//var playerUpdatePendingChan chan *Player = make(chan *Player)

// func handlePlayerUpdatesLoop() {
// 	// playerUpdatePendingChan := make(chan *Player)

// 	for {
// 		select {
// 		case p := <-playerUpdatePendingChan:
// 			p.room.Entities.mu.Lock()
// 			p.canUpdate <- true
// 			<-p.updateDone	//wait for player goroutine to finish updates
// 			p.room.Entities.mu.Unlock()
// 		}
// 	}
// }

func main() {
	/*
	dat, err := os.ReadFile("./gameData/itemData.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &itemData); err != nil {
		panic(err)
	}
	*/

	initGameVars()

	go sendGameStateUnreliableLoop()
	go gameLoop()

	// Listen on UDP Port 80, will be used for all WebRTC traffic
	udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: 80,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening for WebRTC traffic at %s\n", udpListener.LocalAddr())

	// Create a SettingEngine, this allows non-standard WebRTC behavior
	settingEngine := webrtc.SettingEngine{}

	//Our Public Candidate is declared here cause we're not using a STUN server for discovery
	//and just hardcoding the open port, and port forwarding webrtc traffic on the router
	settingEngine.SetNAT1To1IPs([]string{"162.200.58.171"}, webrtc.ICECandidateTypeHost)
	// settingEngine.SetNAT1To1IPs([]string{}, webrtc.ICECandidateTypeHost)

	// Configure our SettingEngine to use our UDPMux. By default a PeerConnection has
	// no global state. The API+SettingEngine allows the user to share state between them.
	// In this case we are sharing our listening port across many.
	settingEngine.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpListener))

	// Create a new API using our SettingEngine
	api = webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))

	fileServer := http.FileServer(http.Dir("./public"))
	http.HandleFunc("/echo", echo) //this request comes from webrtc.html
	http.Handle("/", fileServer)

	err = http.ListenAndServe(":80", nil) //Http server blocks
	if err != nil {
		panic(err)
	}
}
