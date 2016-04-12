local press = require("press")

function routes()
  press.setRoute("hello", "world")
  press.setRoute("goodbye", "moon")
end

