var locomotive = require('locomotive')
  , Controller = locomotive.Controller;

var creatorController = new Controller();
var App = require('../models/app');

creatorController.upload = function() {
	var app = new App({
		ipa = this.param('zip');
	});
}

module.exports = creatorController;
