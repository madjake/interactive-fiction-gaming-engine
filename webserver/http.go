package webserver

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/game")
}

// NewWebServer creates a simple web server to serve the FE client to a browser
func NewWebServer(addr *string) {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  

var IFGame = function () {
    this.ws = false;
};

IFGame.prototype.connect = function () {
    if (this.ws) {
        return false;
    }
    
    this.ws = new WebSocket("{{.}}");
    
    this.ws.onopen = function(evt) {
        print("Connected to the game server.");
    }

    this.ws.onclose = function(evt) {
        print("Your connection to the game server has ended.");
        this.ws = null;
    }
    
    this.ws.onmessage = function(evt) {
        print(evt.data);
    }
    
    this.ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }

    return false;
};

IFGame.prototype.disconnect = function () {
    if (!this.ws) {
        return false;
    }
  
    this.ws.close();
    return true;
};

IFGame.prototype.send = function (msgEvent) {
    if (!this.ws) {
        return false;
    }

    //print("SEND: " + msgEvent.input);
    this.ws.send(msgEvent.input);

    return true;
}

var print = function(message) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var d = document.createElement("div");
    d.innerHTML = message;
    output.appendChild(d);
};

window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var game = new IFGame();
    game.connect();

    document.getElementById("send").onclick = function(evt) {
        game.send({input: input.value});
        evt.preventDefault();
    };
});
</script>
</head>
<body>

<div id="output"></div>
<form>
<input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</body>
</html>
`))
