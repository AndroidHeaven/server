var mongoose = require('mongoose');

module.exports = function() {
  switch (this.env) {
    case 'development':
      mongoose.connect('mongodb://192.168.59.103:27017');
      break;
    case 'production':
      mongoose.connect('mongodb://mongodb.example.com/prod');
      break;
  }
}