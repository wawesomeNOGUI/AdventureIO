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
canvas.style.borderStyle = "solid"
canvas.style.borderWidth = "20px";
canvas.style.borderColor = "#A4FCD4";

//Scale canvas to fit user's screen
canvas.style.width = ""+ window.innerWidth - 40 +"px";
canvas.style.height = ""+ window.innerHeight - 40 +"px";

var keysDown = {};     //This creates var for keysDown event

//===================Sprites================================
var swordSprite = new Image();
    swordSprite.src = "sword.gif";
//==========================================================

// Player Vars
var pX = 50;
var pY = 50;
var speed = 0.75;

// World Vars
var items = [];

//Render using the NTSC Atari 2600 pallete
//https://en.wikipedia.org/wiki/List_of_video_game_console_palettes#Atari_2600
//(inspect element to get the hex color values from the atari color table)
var render = function () {

   ctx.fillStyle = "#000000";
   ctx.fillRect(0, 0, width, height);

   //Draw Players
   t += interpolateInc;

   for (var key in Updates) {     //Updates defined in index.html
      if (Updates.hasOwnProperty(key) && key != playerTag && previousUpdate != undefined && previousUpdate.hasOwnProperty(key))  {
        var x = smoothstep(previousUpdate[key].X, Updates[key].X, t);
        var y = smoothstep(previousUpdate[key].Y, Updates[key].Y, t);

        // Draw Player
        ctx.fillStyle = "#ecb0e0";
        ctx.fillRect(Math.floor(x), Math.floor(y), 4, 4);

        // Draw Held Item
        if (Updates[key].Held.Kind == "sword") {
          ctx.drawImage(swordSprite, Math.floor(x + Updates[key].Held.X), y + Math.floor(Updates[key].Held.Y));
        }
      } else if (Updates.hasOwnProperty(key) && key == playerTag) {
        //Local Player
        ctx.fillStyle = canvas.style.borderColor;
        ctx.fillRect(Math.floor(pX), Math.floor(pY), 4, 4);

          if (Updates[playerTag].Held.Kind == "sword") {
            ctx.drawImage(swordSprite, Math.floor(pX + 2.5), Math.floor(pY + 2.5));
          }
      }
   }

   

   //Draw items laying on ground
   ctx.drawImage(swordSprite, 20, 20);
};


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

var keyPress = function() {
  for(var key in keysDown) {
    var value = Number(key);

    if (value == 37) {   //37 = left
      pX = pX - speed;
      TCPChan.send("X" + pX);
    } else if (value == 39) {  //39 = right
      pX = pX + speed;
      TCPChan.send("X" + pX);
    } else if (value == 40) {  //40 = down
      pY = pY + speed;
      TCPChan.send("Y" + pY);
    } else if (value == 38) {  //38 = up
      pY = pY - speed;
      TCPChan.send("Y" + pY);
    }
  }
};

var spaceKeyDown = false;
window.addEventListener("keydown", function (event) {
  keysDown[event.keyCode] = true;

  if (event.keyCode == 32 && !spaceKeyDown) {
    spaceKeyDown = true;
    TCPChan.send("D");  // drop item
  }
});

window.addEventListener("keyup", function (event) {
  delete keysDown[event.keyCode];

  if (event.keyCode == 32 && spaceKeyDown) {
    spaceKeyDown = false;
  }
});

window.addEventListener("resize", function(){
  canvas.style.width = ""+ window.innerWidth - 40 +"px";
  canvas.style.height = ""+ window.innerHeight - 40 +"px";
});

window.addEventListener("contextmenu", function(e){
  e.preventDefault();   //stops right click bringing up a menu
});
