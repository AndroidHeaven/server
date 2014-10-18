var _ = require('lodash');
var mongoose = require('mongoose');
var Schema = mongoose.Schema;

var ObjectId = mongoose.Types.ObjectId;
var wo


var AppSchema = new Schema({
	ipa: Buffer,
	apk: Buffer
});

ClientSchema.pre('save', function (next, done){
	var zip = this.ipa;

	// take zip and send to worker
	// get ipa and apk from worker

	next();
});

module.exports = mongoose.model('Apps', AppSchema);