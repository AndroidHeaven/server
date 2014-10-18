var express = require('express');
var multer = require('multer');
var bodyParser = require('body-parser');
var sa = require('superagent');
var mongoose = require('mongoose');
var fs = require('fs');
var FormData = require('form-data');
var db = mongoose.connect('mongodb://192.168.59.103:27017');

var app = express();

app.use(bodyParser.urlencoded({ extended: false }));

app.get('/', function(req, res) {
  res.send("hello \n");
});

app.use(multer({
  dest: './uploads/',
  rename: function (fieldname, filename) {
    return filename.replace(/\W+/g, '-').toLowerCase() + Date.now()
  }
}))

app.get('/success', function(req, res) {
  res.send("success! \n");
});

var AppSchema = new mongoose.Schema({
  ipa: Buffer,
  md5: { type: String, Unique: true},
  apk: Buffer
});

var App = db.model("App", AppSchema);

app.all('*', function(req, res, next) {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'PUT, GET, POST, DELETE, OPTIONS');
  res.header('Access-Control-Allow-Headers', 'Content-Type');
  next();
});

app.post('/upload', function(req, res) {
  console.log("f",req.files.file);

  var zip = req.files.file.path;

  sa
  .post('localhost:3000/create')
  .attach('file', zip)
  .end(function(error, res){
    console.log("resp: ",res);
    fs.unlinkSync(zip);
  });

  // // pass zip to worker, on callback save ipa
  // var app = new App({ ipa: compiledIPA });
  // app.save(function(err, done) {
  //   //worker takes done.id and gives back apk in callback, saves to app.apk
  // });
  // res.redirect('/success');
});

app.listen(process.env.PORT || 3001);