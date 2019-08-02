const crypto = require("crypto");
const express = require('express');

const app = express();
const router = express.Router();

// Socket.io
const ioserver = require('http').Server(app);
const io = require('socket.io').listen(ioserver);
ioserver.listen(4321);

// Spawn
const spawn = require('child_process').spawn;

// Webserver
const views = __dirname + '/views/';
const files = __dirname + '/files/';
const port = 1234;

router.use(function (req,res,next) {
  console.log('/' + req.method + " " + req.path);
  next();
});

router.get('/', function(req,res){
  res.sendFile(views + 'index.html');
});

// Socket.io
io.on('connection', function(socket){
  var socketId = socket.id;
  core.logSocket(socketId, 'connection');

  // Events
  socket.on('requestConvert', function (data) {

    // Init
    // TODO: Input validation
    core.logSocket(socketId, 'requestConvert' + JSON.stringify(data));

    // Start conversion
    var fileName = core.getRandomFilename();
    var proc = spawn('go',  ['run', '../Backend/StlGenerator.go', '-file=' + fileName]);
    //var proc = spawn('go',  ['run', '../Backend/TiffDecoder.go', '-neLat=50.1', '-neLng=10.1', '-swLat=49.9', '-swLng=9.9', '-model=surface', '-cropping=sqr', '-length=50', '-heightFactor=20.0', '-name=' + fileName]);

    // Server-side only
    proc.stderr.on('data', (data) => {
      console.log(data.toString());
    });

    proc.stdout.setEncoding('utf8');
    proc.stdout.on('data', function (data) {
      var str = data.toString()
      var lines = str.split(/\n/);
      var message = lines.join("")
      core.logSocket(socketId, 'proc/data ' + message);

      for (line of lines) {
        if (line) {
          io.sockets.to(socketId).emit('convertUpdate', line); // + ';0;0'
        }
      }
    });

    proc.on('close', function (code) {
      core.logSocket(socketId, 'proc/close ' + code);

      if (code == 0) {
        io.sockets.to(socketId).emit('convertSuccess', files + fileName);
      } else {
        io.sockets.to(socketId).emit('convertFailed', 'Error ' + code);
      }
    });
  });
  
  socket.on('disconnect', function () {
    console.log('/io/disconnect ' + socketId);
  });
});

// Application
var core = (function Core() {

  // Private member

  return {
    // Public member
    logSocket: function( socketId, args ) {
      console.log("/io/" + socketId + "/" + args);
    },
    getRandomFilename: function() {
      return "model-" + crypto.randomBytes(10).toString('hex') + ".stl";
    }
  };

})();

// Init
app.use(express.static(views));
app.use('/', router);

app.listen(port, function () {
  console.log('App listening on port ' + port + "!")
})