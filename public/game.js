//Animation and game updates will be triggered once connection made in index.html

var animate = window.requestAnimationFrame || window.webkitRequestAnimationFrame || window.mozRequestAnimationFrame || function (callback) {
        window.setTimeout(callback, 1000/60)
    };

const canvas = document.getElementById("draw");
const ctx = canvas.getContext("2d");
//canvas.width = 160;  //Atari 2600 Resolution
//canvas.height = 192;
canvas.width = 160;  //Adventure Resolution
canvas.height = 105;
var width = canvas.width;
var height = canvas.height;

//Don't Anti-Alias Scaled Images
ctx.imageSmoothingEnabled = false;

//Set border color / width here
// object.style.border = "width style color|initial|inherit" 
//canvas.style.borderStyle = "solid"
//canvas.style.borderWidth = "20px";
var borderColor = "#A4FCD4";
//canvas.style.borderColor = "#A4FCD4";

function hexToRgb(hex) {
  var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return "rgb(" + parseInt(result[1], 16) + "," + parseInt(result[2], 16) + "," + parseInt(result[3], 16) + ")";
}

//Scale canvas to fit user's screen
//canvas.style.width = ""+ window.innerWidth - 40 +"px";
//canvas.style.height = ""+ window.innerHeight - 40 +"px";
canvas.style.width = ""+ window.innerWidth +"px";
canvas.style.height = ""+ window.innerHeight +"px";

var keysDown = {};     //This creates var for keysDown event

//===================Sprites================================
var swordSprite = new Image();
    swordSprite.src = "sprites/sword.gif";  // 10 x 5 pixels

var keySprite = new Image();
    keySprite.src = "sprites/key.gif";

var batSprite0 = new Image();
    batSprite0.src = "sprites/bat_0.png";

var batSprite1 = new Image();
    batSprite1.src = "sprites/bat_1.png";

var spriteMap = {
  "sword": swordSprite,
  "key": keySprite,
  "bat": batSprite0
}

var spriteAnimationInterval = setInterval(function(){
  //Bat Animation
  if (spriteMap["bat"] == batSprite0) {
    spriteMap["bat"] = batSprite1;
  } else {
    spriteMap["bat"] = batSprite0;
  }
  
}, 250)
//==========================================================

// Player Vars
var pX = 50;
var pY = 50;
var speed = 0.75;
var pColor = hexToRgb(borderColor);
var imgData;

// World Vars
var items = [];

var hitX;
var hitY;
var hitDirection;
function checkForPixelPerfectHit() {
  imgData = ctx.getImageData(pX, pY, 4, 4);
  var pColorRGB = pColor.match(/\d+/g);   // gets rgb values from css "rgb(xxx, xxx, xxx)"
  var i = 0;

  for (var y = 0; y < imgData.height; y++) {
    for (var x = 0; x < imgData.width; x++) {
      if (pColorRGB[0] != imgData.data[i] && pColorRGB[1] != imgData.data[i+1] && pColorRGB[2] != imgData.data[i+2] ) {
        hitX = pX + x;
        hitY = pY + y;

        hitDirection = "";
        if (keysDown[37]) {  // left
          hitDirection = "l";
        } else if (keysDown[39]) { // right
          hitDirection = "r";
        }

        if (keysDown[40]) { // down
          hitDirection += "d";
        } else if (keysDown[38]) {  // up
          hitDirection += "u";
        }

        return true;
      } 

      i += 4;
    }
  }

  return false;
}

//Render using the NTSC Atari 2600 pallete
//https://en.wikipedia.org/wiki/List_of_video_game_console_palettes#Atari_2600
//(inspect element to get the hex color values from the atari color table)
var render = function () {
  if (Updates == undefined || previousUpdate == undefined) {
    return;
  }

  ctx.fillStyle = "#b0b0b0";
  ctx.fillRect(0, 0, width, height);

   //Draw Players And Items
  t += interpolateInc;

  for (var key in Updates) {     //Updates defined in index.html  
    if (Number(key)) {  // then its a player
      if (key != playerTag && previousUpdate.hasOwnProperty(key))  {
        var x = smoothstep(previousUpdate[key].X, Updates[key].X, t);
        var y = smoothstep(previousUpdate[key].Y, Updates[key].Y, t);

        // Draw Player
        ctx.fillStyle = "#ecb0e0";
        ctx.fillRect(Math.round(x), Math.round(y), 4, 4);
      } else if (key == playerTag) {
        //Local Player
        ctx.fillStyle = pColor;
        ctx.fillRect(pX, pY, 4, 4);
      }
    } else if (previousUpdate.hasOwnProperty(key)) { // its an item or entity
      if (Updates[key].hasOwnProperty("Owner") && Updates[key].Owner == playerTag) {
        // Draw local player's item
        //ctx.drawImage(swordSprite, pX + Math.round(ownedItemXYOffset[0]), pY + Math.round(ownedItemXYOffset[1]));
        drawColorSprite(ctx, spriteMap[Updates[key].Kind], "#FF00FF", pX + Math.round(ownedItemXYOffset[0]), pY + Math.round(ownedItemXYOffset[1]));
      } else {
        var x = smoothstep(previousUpdate[key].X, Updates[key].X, t);
        var y = smoothstep(previousUpdate[key].Y, Updates[key].Y, t);

        ctx.drawImage(spriteMap[Updates[key].Kind], Math.round(x), Math.round(y));
        // ctx.drawImage(spriteMap[Updates[key].Kind], Updates[key].X, Updates[key].Y);

      }
    }
  }

  // If local player not holding item do Item hit detection
  // if item goes inside player, pick up item
  if (Updates[playerTag].Held == "" && checkForPixelPerfectHit()) {
    TCPChan.send("P" + hitX + "," + hitY + "," + hitDirection);
  }
}

var update = function() {
  keyPress();

  //Check for hitting wall
  if(pX < 2){
      //Wall
      pX = 2;
  }else if(pX > canvas.width - 6){
      //Wall
      pX = canvas.width - 6;
  }
  
  if(pY < 2){
      //Wall
      pY = 2;
  }else if(pY > canvas.height - 6){
      //Wall
      pY = canvas.height - 6;
  }
};

var step = function() {
  update();
  render();
  animate(step);
};

/*
// For being able to check validity / game state it is nice to send player movement reliably
// Cause then the server can check each update, but it might be fine to use the ordered,
// no resend dataChannel cause it will presumably be faster (no head of line blocking),
// and we can at least check validity of latest move
var sendToServerInterval = setInterval(function(){
  TCPChan.send("X" + pX);
  TCPChan.send("Y" + pY);
}, 40);  // sends updates to server every 40 ms instead of every animation loop
*/

var keyPress = function() {
  for(var key in keysDown) {
    var value = Number(key);

    if (value == 37) {   //37 = left
      pX = Math.round(pX - speed);
      UDPChan.send("X" + pX);
    } else if (value == 39) {  //39 = right
      pX = Math.round(pX + speed);
      UDPChan.send("X" + pX);
    } else if (value == 40) {  //40 = down
      pY = Math.round(pY + speed);
      UDPChan.send("Y" + pY);
    } else if (value == 38) {  //38 = up
      pY = Math.round(pY - speed);
      UDPChan.send("Y" + pY);
    }
  }
};

window.addEventListener("keydown", function (event) {
  // Single sends
  if (event.keyCode == 32 && !keysDown[32]) {  // space
    TCPChan.send("D");  // drop item
  } else if (event.keyCode == 37 && !keysDown[37]) { // left
    pX = Math.round(pX - speed);
    TCPChan.send("X" + pX);
  } else if (event.keyCode == 39 && !keysDown[39]) { // right
    pX = Math.round(pX + speed);
    TCPChan.send("X" + pX);
  } else if (event.keyCode == 40 && !keysDown[40]) { // down
    pY = Math.round(pY + speed);
    TCPChan.send("Y" + pY);
  } else if (event.keyCode == 38 && !keysDown[38]) { // up
    pY = Math.round(pY - speed);
    TCPChan.send("Y" + pY);
  }


  keysDown[event.keyCode] = true;
});

window.addEventListener("keyup", function (event) {
  // send a couple more ordered udp messages for player movement 
  // so server will most likely get the last location the player moved to
  // can't do TCP send here cause if user starts moving again, the fast udp messages might
  // come in first and then the late TCP message of where user stopped will make it look
  // to other players like a player rubberbanded 
  if (event.keyCode == 37) { // left
    UDPChan.send("X" + pX);
    UDPChan.send("X" + pX);
  } else if (event.keyCode == 39) { // right
    UDPChan.send("X" + pX);
    UDPChan.send("X" + pX);
  } else if (event.keyCode == 40) { // down
    UDPChan.send("Y" + pY);
    UDPChan.send("Y" + pY);
  } else if (event.keyCode == 38) { // up
    UDPChan.send("Y" + pY);
    UDPChan.send("Y" + pY);
  }

  delete keysDown[event.keyCode];
});

window.addEventListener("resize", function(){
 // canvas.style.width = ""+ window.innerWidth - 40 +"px";
 // canvas.style.height = ""+ window.innerHeight - 40 +"px";
 canvas.style.width = ""+ window.innerWidth +"px";
 canvas.style.height = ""+ window.innerHeight +"px";
});

window.addEventListener("contextmenu", function(e){
  e.preventDefault();   //stops right click bringing up a menu
});
