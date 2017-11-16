// This program extends part 9.
//
// It creates a web server listening on localhost:8080.
// It registers a websocket handler on "/msg".
// It registers a root handler on "/" that will serve a simple GUI.
// It sends messages on "/msg" to be displayed and receives messages from
// "/msg" to be sent to connected peers.
//
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/sbinet/whispering-gophers/util"
	"golang.org/x/net/websocket"
)

var (
	peerAddr = flag.String("peer", "", "peer host:port")
	httpAddr = flag.String("http", "localhost:8080", "http [ip]:port")
	self     string
)

type Message struct {
	ID   string
	Addr string
	Body string
}

func main() {
	flag.Parse()

	l, err := util.Listen()
	if err != nil {
		log.Fatal(err)
	}
	self = l.Addr().String()
	log.Println("Listening on", self)

	go gui.run()
	go dial(*peerAddr)
	go readInput()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}
}

var gui = &GUI{conns: make(map[*websocket.Conn]bool)}

type GUI struct {
	conns map[*websocket.Conn]bool
	mu    sync.RWMutex
}

func (gui *GUI) run() {
	log.Println("Serving", *httpAddr)
	http.HandleFunc("/", rootHandler)
	http.Handle("/msg", websocket.Handler(msgHandler))
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Add adds a websocket to the list of connections this GUI handles.
func (gui *GUI) Add(c *websocket.Conn) {
	gui.mu.Lock()
	defer gui.mu.Unlock()
	if gui.conns[c] {
		return
	}
	gui.conns[c] = true
}

// Remove removes a websocket from the list of connections this GUI handles
func (gui *GUI) Remove(c *websocket.Conn) {
	gui.mu.Lock()
	defer gui.mu.Unlock()
	delete(gui.conns, c)
}

// Display sends m to all the websocket connections this GUI handles
// so it can be displayed.
func (gui *GUI) Display(m Message) {
	gui.mu.RLock()
	defer gui.mu.RUnlock()
	m.Body = html.EscapeString(m.Body)
	for c := range gui.conns {
		err := websocket.JSON.Send(c, m)
		if err != nil {
			log.Print(err)
		}
	}
}

var peers = &Peers{m: make(map[string]chan<- Message)}

type Peers struct {
	m  map[string]chan<- Message
	mu sync.RWMutex
}

// Add creates and returns a new channel for the given peer address.
// If an address already exists in the registry, it returns nil.
func (p *Peers) Add(addr string) <-chan Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.m[addr]; ok {
		return nil
	}
	ch := make(chan Message)
	p.m[addr] = ch
	return ch
}

// Remove deletes the specified peer from the registry.
func (p *Peers) Remove(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.m, addr)
}

// List returns a slice of all active peer channels.
func (p *Peers) List() []chan<- Message {
	p.mu.RLock()
	defer p.mu.RUnlock()
	l := make([]chan<- Message, 0, len(p.m))
	for _, ch := range p.m {
		l = append(l, ch)
	}
	return l
}

func broadcast(m Message) {
	for _, ch := range peers.List() {
		select {
		case ch <- m:
		default:
			// Okay to drop messages sometimes.
		}
	}
	gui.Display(m)
}

func serve(c net.Conn) {
	defer c.Close()
	d := json.NewDecoder(c)
	for {
		var m Message
		err := d.Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}

		if Seen(m.ID) {
			continue
		}

		fmt.Printf("%#v\n", m)
		broadcast(m)
		go dial(m.Addr)
	}
}

func readInput() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := Message{
			ID:   util.RandomID(),
			Addr: self,
			Body: s.Text(),
		}
		Seen(m.ID)
		broadcast(m)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

func dial(addr string) {
	if addr == self {
		return // Don't try to dial self.
	}

	ch := peers.Add(addr)
	if ch == nil {
		return // Peer already connected.
	}
	defer peers.Remove(addr)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(addr, err)
		return
	}
	defer c.Close()

	e := json.NewEncoder(c)
	for m := range ch {
		err := e.Encode(m)
		if err != nil {
			log.Println(addr, err)
			return
		}
	}
}

var seen = struct {
	sync.Mutex
	m map[string]bool
}{
	m: make(map[string]bool),
}

// Seen returns true if the specified id has been seen before.
// If not, it returns false and marks the given id as "seen".
func Seen(id string) bool {
	seen.Lock()
	ok := seen.m[id]
	seen.m[id] = true
	seen.Unlock()
	return ok
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var data = struct {
		Self string
	}{
		Self: self,
	}
	err := rootTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func msgHandler(ws *websocket.Conn) {
	defer ws.Close()

	gui.Add(ws)
	defer gui.Remove(ws)

	for {
		var txt string
		err := websocket.Message.Receive(ws, &txt)
		if err != nil {
			log.Println(err)
			break
		}
		m := Message{
			ID:   util.RandomID(),
			Addr: self,
			Body: txt,
		}
		Seen(m.ID)
		broadcast(m)
	}
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
    <title>WhisperNet</title>
	<style>
		:host {
			display: block;
			box-sizing: border-box;
			text-align: center;
			margin: 5px;
			max-width: 250px;
			min-width: 200px;
		}
		body {
			font-family: 'Roboto', 'Helvetica Neue', Helvetica, Arial, sans-serif;
		}
		#box {
			overflow: auto;
			border-style: solid;
			width: 60%;
			max-height: 70%;
			padding: 10px;
		}
	</style>
	<script type="text/javascript">
		var sock = null;

		function update(data) {
			var box = document.getElementById("box");
			box.innerHTML += "<b>" + data.Addr + ":</b> " + data.Body + "<br>";
			box.scrollTop = box.scrollHeight;
		};

		window.onload = function() {
			sock = new WebSocket("ws://"+location.host+"/msg");
			sock.onmessage = function(event) {
				var msg = JSON.parse(event.data);
				update(msg);
			};

			var button = document.getElementById("button");
			button.addEventListener("click", function(event){
				var text = document.getElementById("textbox").value;
				sock.send(text);
			});
		};

		function sendMsg(txt, event) {
			var KEY_ENTER = 13;
			if (event.keyCode == KEY_ENTER) {
				var button = document.getElementById("button");
				button.click();
				txt.value = "";
			}
		}
	</script>
</head>
<body>

	<h2 id="self">{{.Self}}</h2>
	<div id="box" style="height:300px;"></div>
	<br>
    <input type="text" placeholder="message" id="textbox" onkeypress="sendMsg(this, event)" autofocus>
    <button id="button">Send</button>
</body>
</html>
`))
