var express = require('express');
var bodyParser = require('body-parser');
var multer = require('multer');
var fs = require('fs');
var exec = require('child_process').exec;

var STATE_READY = 0;
var STATE_CREATING = 1;
var STATE_RUNNING = 2;

var app = express();
var state = STATE_READY;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));

app.use(multer({
  dest: './uploads/'
}));

app.post('/create', function(req, res) {
  if (state != STATE_READY) {
    console.log("something terrible happened");
    res.send("failure");
    return;
  }
  var zip = req.files.file;

  exec('./compile_ipa.sh '+process.cwd()+'/'+zip.path, {cwd: process.cwd()}, function(error, stdout, stderr) {
    if (error !== null) {
      console.log('exec error: ', error);
    }
    var path = process.cwd()+'/'+'./artifacts/artifact.tar.gz';
    res.download(path, null, function(err){
      if (err) {
        console.log(err);
        res.status(err.status).end();
      }
      else {
        console.log('Sent:', path);
      }
      exec('rm -rf '+process.cwd()+'/artifacts');
    });
  });
});

app.listen(process.env.PORT || 3000);
