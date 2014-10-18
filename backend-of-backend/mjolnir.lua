local hotkey = require "mjolnir.hotkey"
local application = require "mjolnir.application"
local window = require "mjolnir.window"

hotkey.bind({"cmd", "alt", "ctrl"}, "C", function ()
	local app = application.applicationsforbundleid("com.apple.iphonesimulator")[1]
	if not app then return end -- If the simulator isn't running, this function is noop.
	local win = app:mainwindow()
	local f = win:frame()
	f.y = 0
	f.x = 0
	win:setframe(f)
end)
