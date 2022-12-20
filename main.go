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
var Updates sync.Map
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

		//3 = ICEConnectionStateConnected
		if connectionState == 3 {
			//Store a new x and y for this player
			NumberOfPlayers++
			playerTag = strconv.Itoa(NumberOfPlayers)
			fmt.Println(playerTag)

			//Store a pointer to a Player Struct
			v, ok := Rooms.Load("r1")
			if !ok {
				fmt.Println("Couldn't find room")
			}
			v.(*Room).Entities.StoreEntity(playerTag, newPlayer(playerTag, 50, 50))

		} else if connectionState == 5 || connectionState == 6 || connectionState == 7 {
			Updates.Delete(playerTag)
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
		playerStruct, ok := Updates.Load(playerTag)
		if ok == false {
			fmt.Println("Uh oh")
		}
		tmpPlayer := playerStruct.(Player)  // need to create a temporary copy to edit: https://stackoverflow.com/questions/17438253/accessing-struct-fields-inside-a-map-value-without-copying

		if msg.Data[0] == 'X' { //88 = "X"
			x, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			// Walls
			if x < 2 {
				x = 2
			} else if x > 154 {
				x = 154
			}	

			// Move Owned Item
			if tmpPlayer.Held != "" {
				tmpItem := ownedItems.LoadItem(tmpPlayer.Held)
				tmpItem.X += x - tmpPlayer.X
				ownedItems.StoreItem(tmpPlayer.Held, tmpItem)
			}
			
			tmpPlayer.X = x
			Updates.Store(playerTag, tmpPlayer)
		} else if msg.Data[0] == 'Y' { //89 = "Y"
			y, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			if y < 2 {
				y = 2
			} else if y > 99 {
				y = 99
			}

			// Move Owned Item
			if tmpPlayer.Held != "" {
				tmpItem := ownedItems.LoadItem(tmpPlayer.Held)
				tmpItem.Y += y - tmpPlayer.Y 
				ownedItems.StoreItem(tmpPlayer.Held, tmpItem)
			}

			tmpPlayer.Y = y
			Updates.Store(playerTag, tmpPlayer)
		}
	})

	//==============================================================================

	//=========================Reliable DataChannel=================================
	// Register channel opening handling
	reliableChannel.OnOpen(func() {
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
		// fmt.Printf("Message from DataChannel '%s': '%s'\n", reliableChannel.Label(), string(msg.Data))

		playerStruct, ok := Updates.Load(playerTag)
		if ok == false {
			fmt.Println("Uh oh")
		}
		tmpPlayer := playerStruct.(Player)  // need to create a temporary copy to edit: https://stackoverflow.com/questions/17438253/accessing-struct-fields-inside-a-map-value-without-copying

		if msg.Data[0] == 'X' { //88 = "X"
			x, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			// Walls
			if x < 2 {
				x = 2
			} else if x > 154 {
				x = 154
			}	

			// Move Owned Item
			if tmpPlayer.Held != "" {
				tmpItem := ownedItems.LoadItem(tmpPlayer.Held)
				tmpItem.X += x - tmpPlayer.X
				ownedItems.StoreItem(tmpPlayer.Held, tmpItem)
			}
			
			tmpPlayer.X = x
			Updates.Store(playerTag, tmpPlayer)
		} else if msg.Data[0] == 'Y' { //89 = "Y"
			y, err := strconv.ParseFloat(string(msg.Data[1:]), 64)
			if err != nil {
				fmt.Println(err)
			}

			if y < 2 {
				y = 2
			} else if y > 99 {
				y = 99
			}

			// Move Owned Item
			if tmpPlayer.Held != "" {
				tmpItem := ownedItems.LoadItem(tmpPlayer.Held)
				tmpItem.Y += y - tmpPlayer.Y 
				ownedItems.StoreItem(tmpPlayer.Held, tmpItem)
			}

			tmpPlayer.Y = y
			Updates.Store(playerTag, tmpPlayer)
		} else if msg.Data[0] == 'D' {
			//dropped item
			if tmpPlayer.Held != "" {
				tmpItem := ownedItems.DeleteItem(tmpPlayer.Held)
				k := tmpPlayer.Held
				tmpItem.Owner = ""
				tmpPlayer.Held = ""

				strayItems.StoreItem(k, tmpItem)
				Updates.Store(playerTag, tmpPlayer)
			}
		} else if msg.Data[0] == 'P' && tmpPlayer.Held == "" {
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

			gotItem, itemKey := strayItems.TryPickUpItem(&ownedItems, playerTag, hitX, hitY)

			if gotItem {
				tmpPlayer.Held = itemKey

				tmpItem := ownedItems.LoadItem(itemKey)

				if sDir == "" {
					sDir = "l"
					tmpItem.X = tmpPlayer.X - 10
					tmpItem.Y = tmpPlayer.Y
				}
				if sDir[0] == 'l' {
					tmpItem.X -= 4
				} else if sDir[0] == 'r' {
					tmpItem.X += 4
				}

				if strings.Contains(sDir, "u") {
					tmpItem.Y -= 4
				} else if strings.Contains(sDir, "d") {
					tmpItem.Y += 4
				}

				ownedItems.StoreItem(itemKey, tmpItem)
				Updates.Store(playerTag, tmpPlayer)

				// Send this player the item offset so they can render it with no delay clientside
				str := "I" + fmt.Sprintf("%.1f", tmpItem.X-tmpPlayer.X)  + "," + fmt.Sprintf("%.1f", tmpItem.Y-tmpPlayer.Y)
				reliableChans.SendToPlayer(playerTag, str)
			}
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
func sendGameStateUnreliableLoop() {
	for {
		time.Sleep(time.Millisecond * 50) //50 milliseconds = 20 updates per second

		Rooms.Range(func(rk, rv interface{}) bool {
			switch z := rv.(type) {
			case *Room:
				// here z is a pointer to a Room
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
	//InitializeEntities(&entities)
}

// All server orchestrated game logic
func gameLoop() {
	for {
		time.Sleep(time.Millisecond * 16)  // 16 ms is a little faster than 60 updates per second

		// Update Rooms
		Rooms.Range(func(k, v interface{}) bool {
			switch z := v.(type) {
			case *Room:
				// here z is a pointer to a Room
				z.updateFunc(z)
			default:
				// no match; here z has the same type as v (interface{})
			}	
			return true
		})
	}
}

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
