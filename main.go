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

			//Store a slice for player x, y, and other data (player{x float64, y float64, holding item})
			Updates.Store(playerTag, Player{50, 50, Item{}})

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
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	//==============================================================================

	//=========================Reliable DataChannel=================================
	// Register channel opening handling
	reliableChannel.OnOpen(func() {

		//Send Client their playerTag so they know who they are in the Updates Array
		sendErr := reliableChannel.SendText("T" + playerTag)
		if sendErr != nil {
			panic(err)
		}

		//add this channel's pointer to list so can broadcast messages to all players
		reliableChans.AddPlayerChan(playerTag, reliableChannel)
		unreliableChans.AddPlayerChan(playerTag, dataChannel)

	})

	// Register message handling (Data all served as a bytes slice []byte)
	// for user controls
	reliableChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		//fmt.Printf("Message from DataChannel '%s': '%s'\n", reliableChannel.Label(), string(msg.Data))

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

			tmpPlayer.Y = y
			Updates.Store(playerTag, tmpPlayer)
		} else if msg.Data[0] == 'D' {
			//dropped item
			if tmpPlayer.Held.Kind != "" {
				tmpItem := tmpPlayer.Held
				tmpItem.X = tmpPlayer.X + tmpItem.X * 2
				tmpItem.Y = tmpPlayer.Y + tmpItem.Y * 2
				tmpItem.Owner = ""

				strayItems.StoreItem(tmpItem.Kind, tmpItem)

				tmpPlayer.Held = Item{}
				Updates.Store(playerTag, tmpPlayer)
			}
		} else if msg.Data[0] == 'P' {
			// Picked Up Item
			sX := string(msg.Data[1:strings.Index(string(msg.Data), ",")])
			sY := string(msg.Data[strings.Index(string(msg.Data), ",") + 1:])

			hitX, err := strconv.ParseFloat(sX, 64)
			if err != nil {
				fmt.Println(err)
			}
			hitY, err := strconv.ParseFloat(sY, 64) 
			if err != nil {
				fmt.Println(err)
			}

			itemHere, itemKey := strayItems.TryPickUpItem(&ownedItems, hitX, hitY)
			
			if itemHere {
				fmt.Println(itemKey)
				fmt.Println(strayItems)
				fmt.Println(ownedItems)
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
func sendGameStateUnreliableLoop(m *sync.Map) {
	for {
		time.Sleep(time.Millisecond * 50) //50 milliseconds = 20 updates per second

		tmpMap := make(map[string]interface{})
		m.Range(func(k, v interface{}) bool {
			tmpMap[k.(string)] = v.(Player)
			return true
		})

		for k, v := range strayItems.GetItems() {
			tmpMap[k] = v
		}

		jsonTemp, err := json.Marshal(tmpMap)
		if err != nil {
			panic(err)
		}

		unreliableChans.Broadcast(string(jsonTemp))
	}
}

// Game Vars

// will contain items that can be picked up by players (mutex)
// ItemContainer & Item defined in types.go
var strayItems ItemContainer = ItemContainer{items: make(map[string]Item)} 

// will contain items with the key being the item, and each item has an Owner tag set to the playerTag who owns it
var ownedItems ItemContainer = ItemContainer{items: make(map[string]Item)} 

func initGameVars() {
	strayItems.StoreItem("sword", Item{20, 20, "", "sword"})
}

// All server orchestrated game logic
func gameLoop() {
	for {
		time.Sleep(time.Millisecond * 15)

		Updates.Range(func(k, v interface{}) bool {
			/*
			for ki, vi := range strayItems.GetItems() {
				// d := math.Sqrt(math.Pow(v.(Player).X - vi.(Item).X - 5, 2) + math.Pow(v.(Player).Y - vi.(Item).Y - 2, 2))
				dX := vi.X + 5 - v.(Player).X + 2 
				dY := vi.Y - 2 - v.(Player).Y + 2

				if math.Abs(dX) < 8 && math.Abs(dY) < 4 {
					// pick up item
					tmpItem := vi
					strayItems.DeleteItem(ki)

					tmpItem.X = dX * 2  // to offset Item from player
					tmpItem.Y = dY * 2
					tmpItem.Owner = k.(string); 

					tmpPlayer := v.(Player)
					tmpPlayer.Held = tmpItem;

					//Updates.Store(k, tmpPlayer)
					//ownedItems.Store(ki, tmpItem)

					break
				}
			}
			*/

			return true   
			// return false	// If f returns false, range stops the iteration. 
		})
	}
}

func main() {

	go sendGameStateUnreliableLoop(&Updates)
	go gameLoop()

	initGameVars()

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
	settingEngine.SetNAT1To1IPs([]string{}, webrtc.ICECandidateTypeHost)

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
