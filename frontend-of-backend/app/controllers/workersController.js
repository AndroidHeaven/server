var locomotive = require('locomotive')
  , Controller = locomotive.Controller;

var workerController = new Controller();
var Worker = require('../models/worker');

workerController.create = function() {
	// register new worker
	var worker = new Worker({
		state: "ready",
		uri: this.param('uri')
	})
}


module.exports = workerController;
