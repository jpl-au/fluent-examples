// Package main runs the tether-wasm todo list example. It is a
// normal tether stateful app - the only difference from a JS-client
// example is that the HTML page loads a WASM bootstrap instead of
// tether.js.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/h1"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	"github.com/jpl-au/tether/mode"
	"github.com/jpl-au/tether/wire"
)

//go:embed static
var staticEmbed embed.FS

// Todo represents a single item in the todo list.
type Todo struct {
	ID   int
	Text string
	Done bool
}

// State holds the session state for a single connected client.
type State struct {
	Todos  []Todo
	NextID int
}

func main() {
	staticFS, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		log.Fatal(err)
	}
	assets := &tether.Asset{
		FS:       staticFS,
		Prefix:   "/static/",
		Precache: []string{"app.css"},
	}

	app := tether.App{
		DevMode: true,
		Assets:  []*tether.Asset{assets},
		Client: tether.Client{
			Runtime: tether.Runtime.WASM("/static/client.wasm"),
		},
	}

	handler := tether.Stateful(app, tether.StatefulConfig[State]{
		Name:       "todo",
		Mode:       mode.ServerSentEvents,
		WireFormat: wire.CBOR,

		InitialState: func(_ *http.Request) State {
			return State{NextID: 1}
		},
		Render: render,
		Handle: handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether WASM - Todo"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	slog.Info("starting tether-wasm todo", "addr", ":8080")
	if err := tether.ListenAndServe(":8080", mux, handler); err != nil {
		log.Fatal(err)
	}
}

// handle processes client events and returns the updated state.
func handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "add":
		text, _ := ev.Get("text")
		if text == "" {
			return s
		}
		s.Todos = append(s.Todos, Todo{
			ID:   s.NextID,
			Text: text,
			Done: false,
		})
		s.NextID++

	case "toggle":
		id, err := ev.Int("id")
		if err != nil {
			return s
		}
		for i := range s.Todos {
			if s.Todos[i].ID == id {
				s.Todos[i].Done = !s.Todos[i].Done
				break
			}
		}

	case "delete":
		id, err := ev.Int("id")
		if err != nil {
			return s
		}
		for i, t := range s.Todos {
			if t.ID == id {
				s.Todos = append(s.Todos[:i], s.Todos[i+1:]...)
				break
			}
		}
	}

	return s
}

// render builds the todo list UI from the current state. Each todo
// item carries a Dynamic key so the WASM client can target it with
// patches.
func render(s State) node.Node {
	return div.New(
		h1.Text("Todo List"),
		todoForm(),
		todoList(s.Todos),
	).Class("todo-app").Dynamic("app")
}

// todoForm renders the input form for adding new todos. The form
// collects the text input value and sends it with the "add" action.
func todoForm() node.Node {
	return bind.Apply(form.New(
		input.New().Name("text").Placeholder("What needs doing?").ID("todo-input"),
		button.Text("Add"),
	).Class("todo-form"),
		bind.OnSubmit("add"),
		bind.Reset(),
	).Dynamic("todo-form")
}

// todoList renders either the list of todos or an empty state message.
func todoList(todos []Todo) node.Node {
	if len(todos) == 0 {
		return p.Text("No todos yet. Add one above.").Class("empty-state").Dynamic("todo-list")
	}

	items := make([]node.Node, len(todos))
	for i, t := range todos {
		items[i] = todoItem(t)
	}
	return ul.New(items...).Class("todo-list").Dynamic("todo-list")
}

// todoItem renders a single todo with toggle and delete buttons. The
// todo ID is passed as event data so the server knows which item to
// act on.
func todoItem(t Todo) node.Node {
	idStr := strconv.Itoa(t.ID)
	key := fmt.Sprintf("todo-%d", t.ID)

	textClass := "todo-text"
	if t.Done {
		textClass = "todo-text todo-done"
	}

	toggleLabel := "Done"
	if t.Done {
		toggleLabel = "Undo"
	}

	return li.New(
		span.Text(t.Text).Class(textClass),
		bind.Apply(button.Text(toggleLabel),
			bind.OnClick("toggle"),
			bind.EventData("id", idStr),
		),
		bind.Apply(button.Text("Delete").Class("delete-btn"),
			bind.OnClick("delete"),
			bind.EventData("id", idStr),
		),
	).Class("todo-item").Dynamic(key)
}
