var tmpCanvas = document.createElement("canvas");
var tmpCtx = tmpCanvas.getContext("2d");

function drawColorSprite(myContext, sprite, color, x, y) {
    tmpCanvas.width = sprite.width;
    tmpCanvas.height = sprite.height;

    tmpCtx.drawImage(sprite, 0, 0);

    //Stuff only kept where drawn pixels overlap with already non background colored pixels
    tmpCtx.globalCompositeOperation = "source-atop";
  
    //Color for image
    tmpCtx.fillStyle = color;
    tmpCtx.fillRect(0, 0, tmpCanvas.width, tmpCanvas.height);
  
    //Then draw to user's canvas
    myContext.drawImage(tmpCanvas, x, y);
  }