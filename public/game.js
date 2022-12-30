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

var dragonSprite = new Image();
    dragonSprite.src = "sprites/dragon.gif";

var dragonSpriteMouthOpen = new Image();
    dragonSpriteMouthOpen.src = "sprites/dragonMouthOpen.gif";

var doorGrateSprite = new Image();
    doorGrateSprite.src = "sprites/doorGrate.gif";

var spriteMap = {
  "sword": swordSprite,
  "key": keySprite,
  "bat": batSprite0,
  "drg": dragonSprite,
  "dG": doorGrateSprite
}

var spriteAnimationInterval = setInterval(function(){
  //Bat Animation
  if (spriteMap["bat"] == batSprite0) {
    spriteMap["bat"] = batSprite1;
  } else {
    spriteMap["bat"] = batSprite0;
  }

  //Dragon animation
  if (spriteMap["drg"] == dragonSprite) {
    spriteMap["drg"] = dragonSpriteMouthOpen;
  } else {
    spriteMap["drg"] = dragonSprite;
  }
  
}, 250)
//==========================================================

//======================Rooms===============================
var r1 = new Image();
    r1.src = "sprites/rooms/r1.gif";

var r2 = new Image();
    r2.src = "sprites/rooms/r2.gif";

var r3 = new Image();
    r3.src = "sprites/rooms/r3.gif";

var r4 = new Image();
    r4.src = "sprites/rooms/r4.gif";

var r5 = new Image();
    r5.src = "sprites/rooms/r5.gif";

var roomSprites = {
  "r1": r1,
  "r2": r2,
  "r3": r3,
  "r4": r4,
  "r5": r5
}
//==========================================================

// Player Vars
// var pX = 50;
// var pY = 50;
var speed = 0.75;
var pColor = hexToRgb(borderColor);
var imgData;

// World Vars
var items = [];

var hitX;
var hitY;
var hitDirection = "";
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

var wHImgData;
var wallHitDirection = "";
function checkForPixelPerfectWallHit(color) {
  wHImgData = ctx.getImageData(pX-1, pY-1, 6, 6);
  var colorRGB = color.match(/\d+/g);   // gets rgb values from css "rgb(xxx, xxx, xxx)"
  var i = 0;

  wallHitDirection = "";

  // check for left hit
  for (var i = 1; i < wHImgData.height-1; i++) {
    if (colorRGB[0] == wHImgData.data[i*4*wHImgData.width] && colorRGB[1] == wHImgData.data[i*4*wHImgData.width+1] && colorRGB[2] == wHImgData.data[i*4*wHImgData.width+2]) {
      wallHitDirection += "l";
      break;
    }
  }

  // check for right hit
  for (var i = 1; i < wHImgData.height-1; i++) {
    if (colorRGB[0] == wHImgData.data[(wHImgData.width-1)*4 + i*4*wHImgData.width] && colorRGB[1] == wHImgData.data[(wHImgData.width-1)*4 + i*4*wHImgData.width+1] && colorRGB[2] == wHImgData.data[(wHImgData.width-1)*4 + i*4*wHImgData.width+2]) {
      wallHitDirection += "r";
      break;
    }
  }

  // check for up hit
  for (var i = 1; i < wHImgData.width-1; i++) {
    if (colorRGB[0] == wHImgData.data[i*4] && colorRGB[1] == wHImgData.data[i*4+1] && colorRGB[2] == wHImgData.data[i*4+2]) {
      wallHitDirection += "u";
      break;
    }
  }

  // check for down hit
  for (var i = 1; i < wHImgData.width-1; i++) {
    if (colorRGB[0] == wHImgData.data[i*4 + wHImgData.data.length-wHImgData.width*4] && colorRGB[1] == wHImgData.data[i*4+wHImgData.data.length-wHImgData.width*4+1] && colorRGB[2] == wHImgData.data[i*4+wHImgData.data.length-wHImgData.width*4+2]) {
      wallHitDirection += "d";
      break;
    }
  }
}

// ran after receive playerTag
function askForUserName() {
  ctx.clearRect(0, 0, width, height);
  
  let userName = prompt("Inpit User Name:", "Fizz Buzz");

  if (userName == null || userName == "Fizz Buzz") {
    askForUserName();
  } else {
    TCPChan.send(room +",U" + userName);
    userNameMap[playerTag] = userName;
    animate(step);
  }
}

//Render using the NTSC Atari 2600 pallete
//https://en.wikipedia.org/wiki/List_of_video_game_console_palettes#Atari_2600
//(inspect element to get the hex color values from the atari color table)
var render = function () {
  if (Updates == undefined || previousUpdate == undefined) {
    return;
  }

  // Draw default background color
  ctx.fillStyle = "#b0b0b0";
  ctx.fillRect(0, 0, width, height);

  //Draw room picture (room defined in index.html)
  ctx.drawImage(roomSprites[room], 0, 0);

  //Special room drawings
  if (room == "r2") {
    drawColorText(ctx, "ITS DANGEROUS TO GO", "#ececec", 14, 12, 7, 8, 8);
    drawColorText(ctx, "ALONE! TAKE THIS", "#ececec", 22, 22, 7, 8, 8);
  }
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
      } else {
        var d = Math.sqrt(Math.pow(Updates[playerTag].X - pX, 2) + Math.pow(Updates[playerTag].Y - pY, 2));
        if (Updates[playerTag].BeingHeld != "") {
          var x;
          var y;
          if (previousUpdate[playerTag] != undefined) {
            x = smoothstep(previousUpdate[playerTag].X, Updates[playerTag].X, t);
            y = smoothstep(previousUpdate[playerTag].Y, Updates[playerTag].Y, t);
          } else {
            x = Updates[playerTag].X
            y = Updates[playerTag].Y
          }
          pX = Math.round(x);
          pY = Math.round(y);          
        }
        //Local Player
        ctx.fillStyle = pColor;
        ctx.fillRect(pX, pY, 4, 4);
      }

      // Draw all usernames
      if (userNameMap[key] != undefined && keysDown[17] != undefined) {  // have to hold down ctrl to see usernames
        ctx.fillStyle = "#000000A0";
        ctx.font = "0.5px";

        if (key == playerTag) {
          ctx.fillText(userNameMap[key], pX - 2, pY);
        } else {
          ctx.fillText(userNameMap[key], Math.round(x) - 2, Math.round(y));
        }
      }
    } else if (previousUpdate.hasOwnProperty(key)) { // its an item or entity
      if (key == ownedItemXYOffset[0]) {
        // Draw local player's item
        //ctx.drawImage(swordSprite, pX + Math.round(ownedItemXYOffset[0]), pY + Math.round(ownedItemXYOffset[1]));
        drawColorSprite(ctx, spriteMap[Updates[key].K], "#FF00FF", pX + Math.round(ownedItemXYOffset[1]), pY + Math.round(ownedItemXYOffset[2]));
      } else {
        var x = smoothstep(previousUpdate[key].X, Updates[key].X, t);
        var y = smoothstep(previousUpdate[key].Y, Updates[key].Y, t);

        ctx.drawImage(spriteMap[Updates[key].K], Math.round(x), Math.round(y));
        // ctx.drawImage(spriteMap[Updates[key].Kind], Updates[key].X, Updates[key].Y);
      }
    }
  }

  // If local player not holding item do Item hit detection
  // if item goes inside player, pick up item
  if (checkForPixelPerfectHit()) {
    TCPChan.send(room +",P" + hitX + "," + hitY + "," + hitDirection);
  }
}

var prevRoom = "r1";
var update = function() {
  prevRoom = room;

  //keyPress();

  //Check for hitting edge of screen
  if(pX < 0){
      //Wall
      pX = 0;
  }else if(pX > canvas.width - 4){
      //Wall
      pX = canvas.width - 4;
  }
  
  if(pY < 0){
      //Wall
      pY = 0;
  }else if(pY > canvas.height - 4){
      //Wall
      pY = canvas.height - 4;
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

var keyCheckInterval = setInterval(function() {
  for (var i = 0; i < 3; i++) {
    keyPress();
  }
}, 10)

var keyPress = function() {
  checkForPixelPerfectWallHit(hexToRgb(wallColor));

  for(var key in keysDown) {
    var value = Number(key);

    if (value == 37 && !wallHitDirection.includes("l")) {   //37 = left
      pX = Math.round(pX - speed);
      UDPChan.send(room +",X" + pX);
    } else if (value == 39 && !wallHitDirection.includes("r")) {  //39 = right
      pX = Math.round(pX + speed);
      UDPChan.send(room +",X" + pX);
    } else if (value == 40 && !wallHitDirection.includes("d")) {  //40 = down
      pY = Math.round(pY + speed);
      UDPChan.send(room +",Y" + pY);
    } else if (value == 38 && !wallHitDirection.includes("u")) {  //38 = up
      pY = Math.round(pY - speed);
      UDPChan.send(room +",Y" + pY);
    }
  }
};

window.addEventListener("keydown", function (event) {
  if (Updates == undefined) {
    return;
  }

  // Single sends
  if (event.keyCode == 32 && !keysDown[32] && Updates[playerTag].BeingHeld == "") {  // space
    ownedItemXYOffset = [0, 0, 0]
    TCPChan.send(room +",D");  // drop item
  } else if (event.keyCode == 37 && !keysDown[37]) { // left
    // pX = Math.round(pX - speed);
    // TCPChan.send("X" + pX);
  } else if (event.keyCode == 39 && !keysDown[39]) { // right
    // pX = Math.round(pX + speed);
    // TCPChan.send("X" + pX);
  } else if (event.keyCode == 40 && !keysDown[40]) { // down
    // pY = Math.round(pY + speed);
    // TCPChan.send("Y" + pY);
  } else if (event.keyCode == 38 && !keysDown[38]) { // up
    // pY = Math.round(pY - speed);
    // TCPChan.send("Y" + pY);
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
    UDPChan.send(room +",X" + pX);
    UDPChan.send(room +",X" + pX);
  } else if (event.keyCode == 39) { // right
    UDPChan.send(room +",X" + pX);
    UDPChan.send(room +",X" + pX);
  } else if (event.keyCode == 40) { // down
    UDPChan.send(room +",Y" + pY);
    UDPChan.send(room +",Y" + pY);
  } else if (event.keyCode == 38) { // up
    UDPChan.send(room +",Y" + pY);
    UDPChan.send(room +",Y" + pY);
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
