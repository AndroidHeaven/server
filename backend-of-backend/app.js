var express = require('express');
var bodyParser = require('body-parser');
var multer = require('multer');

var STATE_READY = 0;
var STATE_CREATING = 1;
var STATE_RUNNING = 2;

var app = express();
var state = STATE_READY;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));

app.use(multer({
  inMemory: true
}));

app.post('/create', function(req, res) {
  console.log(req.files);
  if (state != STATE_READY) {
    console.log("something terrible happened");
  }
  var zip = req.files.z;
  // make zip into ipa
  console.log(zip);

  res.send(zip);
});

app.listen(process.env.PORT || 3000);
