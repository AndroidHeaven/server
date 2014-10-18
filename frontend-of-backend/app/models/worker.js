var _ = require('lodash');
var mongoose = require('mongoose');
var Schema = mongoose.Schema;

var ObjectId = mongoose.Types.ObjectId;

var WorkerSchema = new Schema({
	state: { type: String }, //creator, runner, ready
	uri: { type: String }
});

WorkerSchema.method('isReady', function(){
	if(this.state == "ready") {
		return true;
	}
	return false;
});

WorkerSchema.method('createFiles', function(app, callback) {

})

module.exports = mongoose.model('Worker', WorkerSchema);