//Thanks to http://members.bitstream.net/marksim/atarimac/fonts.html
//For the atari chunky font!
var fontImg = new Image();
  fontImg.src = "./chunkyFont.png";

const tempCanvas = document.createElement("canvas");
var tempCtx = tempCanvas.getContext("2d");
tempCanvas.width = 16;
tempCanvas.height = 16;

function drawText(myContext, string, x, y, spacing = 16, scaleX = 16, scaleY = 16){
  for(var i = 0; i<string.length; i++){
    if(string[i] == " "){
      myContext.fillStyle = "#00000000"; //Transparent
      myContext.fillRect(x+(i*spacing), y, scaleX, scaleY);
      continue;
    }
    myContext.drawImage(fontImg, fontMap[string[i]][0], fontMap[string[i]][1], fontMap[string[i]][2], fontMap[string[i]][3], x+(i*spacing), y, scaleX, scaleY);
  }
}

function drawColorText(myContext, string, color, x, y, spacing = 16, scaleX = 16, scaleY = 16){
  tempCanvas.width = string.length*(spacing);
  tempCanvas.height = scaleY;

  for(var i = 0; i<string.length; i++){
    if(string[i] == " "){
      myContext.fillStyle = "#00000000"; //Transparent
      myContext.fillRect(x+(i*spacing), y, scaleX, scaleY);
      continue;
    }


    tempCtx.drawImage(fontImg, fontMap[string[i]][0], fontMap[string[i]][1], fontMap[string[i]][2], fontMap[string[i]][3], i*spacing, 0, scaleX, scaleY);
  }
  //Stuff only kept where drawn pixels overlap with already non background colored pixels
  tempCtx.globalCompositeOperation = "source-in";

  //Color for image
  tempCtx.fillStyle = color;
  tempCtx.fillRect(0, 0, tempCanvas.width, tempCanvas.height);

  //Then draw to user's canvas
  myContext.drawImage(tempCanvas, x, y);
}

var fontMap = {
  A:[0,0,16,16],
  B:[16,0,16,16],
  C:[32,0,16,16],
  D:[48,0,16,16],
  E:[64,0,16,16],
  F:[80,0,16,16],
  G:[96,0,16,16],
  H:[112,0,16,16],
  I:[128,0,16,16],
  J:[144,0,16,16],
  K:[160,0,16,16],
  L:[176,0,16,16],
  M:[192,0,16,16],
  N:[108,0,16,16],
  O:[224,0,16,16],
  P:[240,0,16,16],
  Q:[256,0,16,16],
  R:[272,0,16,16],
  S:[288,0,16,16],
  T:[304,0,16,16],
  U:[320,0,16,16],
  V:[336,0,16,16],
  W:[352,0,16,16],
  X:[368,0,16,16],
  Y:[384,0,16,16],
  Z:[400,0,16,16],
  a:[0,16,16,16],
  b:[16,16,16,16],
  c:[32,16,16,16],
  d:[48,16,16,16],
  e:[64,16,16,16],
  f:[80,16,16,16],
  g:[96,16,16,16],
  h:[112,16,16,16],
  i:[128,16,16,16],
  j:[144,16,16,16],
  k:[160,16,16,16],
  l:[176,16,16,16],
  m:[192,16,16,16],
  n:[108,16,16,16],
  o:[224,16,16,16],
  p:[240,16,16,16],
  q:[256,16,16,16],
  r:[272,16,16,16],
  s:[288,16,16,16],
  t:[304,16,16,16],
  u:[320,16,16,16],
  v:[336,16,16,16],
  w:[352,16,16,16],
  x:[368,16,16,16],
  y:[384,16,16,16],
  z:[400,16,16,16],
  "1":[0,32,16,16],
  "2":[16,32,16,16],
  "3":[32,32,16,16],
  "4":[48,32,16,16],
  "5":[64,32,16,16],
  "6":[80,32,16,16],
  "7":[96,32,16,16],
  "8":[112,32,16,16],
  "9":[128,32,16,16],
  "0":[144,32,16,16],
  "(":[160,32,16,16],
  "!":[176,32,16,16],
  "@":[192,32,16,16],
  "$":[208,32,16,16],
  "%":[224,32,16,16],
  "&":[240,32,16,16],
  "*":[256,32,16,16],
  "+":[272,32,16,16],
  "?":[288,32,16,16],
  ";":[304,32,16,16],
  ":":[320,32,16,16],
  ")":[336,32,16,16],
  "-":[352,32,16,16],
}
