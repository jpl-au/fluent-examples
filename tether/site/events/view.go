package events

import (
	"strconv"
	"time"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/dropdown"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/option"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/composite/viewport"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
	"github.com/jpl-au/fluent-examples/tether/components/simple/spinner"
)

// Render builds the events and forms page, demonstrating every bind
// event type: click, submit, input, change, keydown, focus, blur,
// viewport, throttle, custom events, and typed/bound parsing.
func Render(s State) node.Node {
	return page.New(
		panel.Card("Click Events", "Click the button - each click sends an event to the server and increments the counter. Because this page is stateless HTTP, the count is carried with each event.", "bind.Click + bind.EventData", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Click me", "events.click",
					bind.EventData("count", strconv.Itoa(s.ClickCount)),
				),
				span.Text("Clicked "+strconv.Itoa(s.ClickCount)+" times"),
			),
		),

		panel.Card("Input with Debounce", "Type in the field - the server only receives the value after you stop typing for 300ms, reducing unnecessary requests.", "bind.Input + bind.Debounce", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.TextValue("text", s.InputValue, "Type something..."),
					bind.OnInput("events.input"),
					bind.Debounce(300*time.Millisecond),
				),
				layout.Container(inputResult(s.InputValue)).Dynamic("input-result"),
			),
		),

		panel.Card("Form Submit", "Enter a name and submit. The button disables while the server processes the request, then shows the result below.", "bind.Submit + bind.Disable", panel.AllTransports,
			bind.Apply(field.Inline(
				field.Group(field.Label("name", "Name"), field.Text("name", "Enter a name...")),
				button.Submit("Submit", bind.Disable("Submitting...")),
			),
				bind.OnSubmit("events.submit"),
			),
			layout.Container(submitResult(s.SubmitResult, s.SubmitError)).Dynamic("submit-result"),
		),

		panel.Card("AutoFocus", "Submit with an empty email field - the server validates it, returns an error, and automatically refocuses the input so you can correct it immediately. AutoFocus is applied conditionally in Render: the element only carries the autofocus attribute when state says there is an error, so the focus is server-driven rather than hard-wired in markup.", "bind.AutoFocus", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.Inline(
					field.Group(field.Label("autofocus-email", "Email"), autoFocusInput(s)),
					button.Submit("Submit"),
				),
					bind.OnSubmit("events.autofocus"),
				),
				layout.Container(autoFocusResult(s.AutoFocusResult, s.AutoFocusError)).Dynamic("autofocus-result"),
			),
		),

		panel.Card("Change Events", "Pick a colour from the dropdown - the server receives the selected value when the selection changes.", "bind.Change", panel.AllTransports,
			layout.Stack(
				bind.Apply(dropdown.New(option.Option("", "Select a colour..."), option.Option("red", "Red"), option.Option("green", "Green"), option.Option("blue", "Blue")).Name("colour"),
					bind.OnChange("events.change"),
				),
				layout.Container(colourResult(s.ChangeValue)).Dynamic("colour-result"),
			),
		),

		panel.Card("Keyboard Events", "Click into the field and press Enter - only the Enter key triggers the server event. Other keys are filtered out.", "bind.KeyDown + bind.FilterKey", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.Text("key", "Press Enter..."),
					bind.OnKeyDown("events.keydown"),
					bind.FilterKey("Enter"),
				),
				layout.Container(valueResult("Last key", s.LastKey, "Press Enter above")).Dynamic("key-result"),
			),
		),

		panel.Card("Focus & Blur", "Click into the field to trigger a focus event; click away to trigger blur. Each fires a separate server event so the handler knows exactly when the element gains or loses attention.", "bind.Focus · bind.Blur", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.Text("focus-demo", "Click in and out of this field..."),
					bind.OnFocus("events.focus"),
					bind.OnBlur("events.blur"),
				),
				layout.Container(valueResult("Status", s.FocusBlurResult, "Click into the field above")).Dynamic("focus-result"),
			),
		),

		panel.Card("Confirm Dialog", "Click the button - the browser shows a confirmation dialog. The event only fires if you confirm.", "bind.Confirm", panel.AllTransports,
			button.DangerAction("Delete Everything", "events.confirm",
				bind.Confirm("Are you sure you want to delete everything?"),
			),
		),

		panel.Card("Throttle", "Click rapidly - the server only receives one event per second, no matter how fast you click.", "bind.Throttle", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Rapid Click", "events.throttle",
					bind.EventData("count", strconv.Itoa(s.ThrottleHits)),
					bind.Throttle(time.Second),
				),
				span.Text("Server received "+strconv.Itoa(s.ThrottleHits)+" events (1s throttle)"),
			),
		),

		panel.Card("Event Data", "Click to send extra metadata with the event. The server receives item-id and category as key-value pairs alongside the click.", "bind.EventData", panel.AllTransports,
			layout.Stack(
				button.PrimaryAction("Send with Metadata", "events.data",
					bind.EventData("item-id", "42"),
					bind.EventData("category", "demo"),
				),
				layout.Container(valueResult("Result", s.EventDataResult, "Click the button to see event data")).Dynamic("data-result"),
			),
		),

		panel.Card("Raw Data Attributes", "bind.Data sets a plain data-* attribute on the element - useful for CSS selectors or JS hooks. Unlike bind.EventData which sends data with the event, bind.Data just renders an HTML attribute. Click the button to see the server read the attribute value from the event data.", "bind.Data", panel.AllTransports,
			layout.Stack(
				button.PrimaryAction("Button with data-status", "events.raw-data",
					bind.Data("data-status", "active"),
				),
				layout.Container(valueResult("Result", s.RawDataResult, "Click to see the data-* attribute value")).Dynamic("raw-data-result"),
			),
		),

		panel.Card("Typed Extraction", "Enter a quantity, a price, and pick an urgency, then submit. The handler uses ev.Int for the integer, ev.Float64 for the decimal price, and ev.Bool for the dropdown - no manual strconv needed. Each returns a typed value directly from the event data.", "ev.Int · ev.Float64 · ev.Bool", panel.AllTransports,
			layout.Stack(
				bind.Apply(form.New(
					field.Group(field.Label("qty", "Quantity"), field.TextWithID("qty", "qty", "e.g. 5")),
					field.Group(field.Label("price", "Price"), field.TextWithID("price", "price", "e.g. 9.99")),
					field.Group(
						field.Label("urgent", "Urgent?"),
						dropdown.New(
							option.Option("", "Select..."),
							option.Option("true", "Yes"),
							option.Option("false", "No"),
						).ID("urgent").Name("urgent"),
					),
					button.Submit("Submit"),
				),
					bind.OnSubmit("events.typed"),
				),
				layout.Container(valueResult("Parsed", s.TypedResult, "Submit the form to see parsed values")).Dynamic("typed-result"),
			),
		),

		panel.Card("Struct Binding", "Fill in the fields and submit. The handler calls ev.Bind to deserialise all form values into a struct at once using struct tags - no repeated ev.Get calls, no manual type conversion.", "ev.Bind", panel.AllTransports,
			layout.Stack(
				bind.Apply(form.New(
					field.Group(field.Label("bind-name", "Name"), field.TextWithID("bind-name", "name", "Your name")),
					field.Group(field.Label("bind-email", "Email"), field.TextWithID("bind-email", "email", "your@email.com")),
					button.Submit("Submit"),
				),
					bind.OnSubmit("events.bind"),
				),
				layout.Container(valueResult("Bound", s.BindResult, "Submit the form to see bound struct values")).Dynamic("bind-result"),
			),
		),

		panel.Card("Loading Indicator", "Click the button - a spinner appears while the server processes the request (1 second simulated delay).", "bind.Indicator", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Load Something", "events.indicator",
					bind.Indicator("#indicator-spinner"),
				),
				spinner.New().ID("indicator-spinner"),
			),
		),

		panel.Card("Arbitrary DOM Events", "Double-click the button - bind.On lets you listen for any DOM event, not just the built-in ones.", "bind.On", panel.AllTransports,
			layout.Stack(
				button.Primary("Double-click me", bind.Event("dblclick", "events.custom")),
				layout.Container(valueResult("Result", s.CustomEventResult, "Double-click the button above")).Dynamic("custom-result"),
			),
		),

		panel.Card("Form Reset", "Type a message and submit - the form fields clear automatically after a successful submission.", "bind.Reset", panel.AllTransports,
			bind.Apply(field.Inline(
				field.Group(field.Label("message", "Message"), field.Text("message", "Type a message...")),
				button.Submit("Send & Reset"),
			),
				bind.OnSubmit("events.reset"),
				bind.Reset(),
			),
			layout.Container(resetResult(s.ResetResult)).Dynamic("reset-result"),
		),

		panel.Card("Viewport Trigger", "Scroll to the bottom of the list - each time the sentinel enters the viewport the server loads five more items. This is the infinite scroll pattern: bind.Viewport on the trailing sentinel triggers automatic pagination. The current page is carried in EventData so the stateless handler knows which batch to append.", "bind.Viewport · bind.EventData", panel.AllTransports,
			viewportList(s.ViewportPage),
		),

		panel.Card("Paste Event", "Paste text into the input below. The pasted content is sent to the server via bind.OnPaste and displayed in the result. The pasted text arrives in ev.Value().", "bind.OnPaste", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.Text("paste-input", "Paste something here..."), bind.OnPaste("events.paste")),
				pasteResult(s.PasteResult),
			),
		),

		panel.Card("Context Menu", "Right-click the box below. The browser's default context menu is suppressed via bind.PreventDefault and the event is sent to the server instead.", "bind.Event · bind.PreventDefault", panel.AllTransports,
			layout.Stack(
				bind.Apply(
					div.New(result.Block("Right-click me")),
					bind.Event("contextmenu", "events.contextmenu"),
					bind.PreventDefault(),
				),
				contextMenuResult(s.ContextMenuResult),
			),
		),

		panel.Card("Client-Side Validation", "The name field is required. Try submitting with an empty field - the browser shows a validation tooltip without a server round-trip. bind.Required, bind.MinLength, bind.MaxLength, and bind.Pattern use the native constraint validation API.", "bind.Required", panel.AllTransports,
			layout.Stack(
				bind.Apply(field.Inline(
					field.Group(
						field.Label("validated-name", "Name"),
						bind.Apply(field.Text("validated-name", "Required field..."), bind.Required("This field is required")),
					),
					button.Submit("Submit"),
				), bind.OnSubmit("events.validated")),
				validatedResult(s.ValidatedResult),
			),
		),

		panel.Card("Content Editable", "Click the text below to edit it inline. When you click away (blur), the edited text is sent to the server. bind.Editable sets contenteditable and fires the action on blur.", "bind.Editable", panel.AllTransports,
			layout.Stack(
				bind.Apply(
					span.Text("Click me to edit this text").Class("signal-panel"),
					bind.Editable("events.editable"),
				),
				editableResult(s.EditableResult),
			),
		),
	)
}

// submitResult renders the submit demo outcome - either a validation
// error, a success message, or an empty placeholder before submission.
func submitResult(res, err string) node.Node {
	if err != "" {
		return field.Error(err)
	}
	if res != "" {
		return result.Success(res)
	}
	return layout.Container()
}

// valueResult renders a labelled result block when the server has
// responded, or a hint paragraph when there is nothing to show yet.
func valueResult(lbl, val, placeholder string) node.Node {
	if val == "" {
		return hint.Text(placeholder)
	}
	return layout.Container(
		result.Label(lbl),
		result.Block(val),
	)
}

// inputResult renders the debounced input result.
func inputResult(val string) node.Node {
	if val == "" {
		return hint.Text("Type above to see the server response")
	}
	return layout.Container(
		result.Label("Server received"),
		result.Block(val),
	)
}

// autoFocusInput conditionally applies bind.AutoFocus so the cursor
// returns to the input after a validation error.
func autoFocusInput(s State) node.Node {
	n := field.TextWithID("autofocus-email", "email", "Enter your email...")
	if s.AutoFocusError != "" {
		return bind.Apply(n, bind.AutoFocus())
	}
	return n
}

// autoFocusResult renders the autofocus demo outcome, following the
// same error/success/empty pattern as submitResult.
func autoFocusResult(res, err string) node.Node {
	if err != "" {
		return field.Error(err)
	}
	if res != "" {
		return result.Success(res)
	}
	return layout.Container()
}

// viewportList builds the infinite-scroll demo: a growing list of
// items with a sentinel element that triggers the next page load
// via bind.OnViewport when it enters the viewport.
func viewportList(pg int) node.Node {
	const pageSize = 5
	total := (pg + 1) * pageSize
	items := make([]node.Node, total)
	for i := range total {
		items[i] = viewport.Itemf("Item %d", i+1)
	}
	sentinel := bind.Apply(viewport.Sentinel().Dynamic("viewport-sentinel"),
		bind.OnViewport("events.viewport"),
		bind.EventData("page", strconv.Itoa(pg)),
	)
	return layout.Container(
		viewport.List(items...),
		sentinel,
	).Dynamic("viewport-list")
}

// colourResult renders the selected colour in its own colour.
func colourResult(val string) node.Node {
	if val == "" {
		return hint.Text("Pick a colour above")
	}
	switch val {
	case "red":
		return result.Danger(val)
	case "green":
		return result.Success(val)
	case "blue":
		return result.Blue(val)
	default:
		return result.Block(val)
	}
}

// resetResult renders the bind.Reset demo outcome.
func resetResult(val string) node.Node {
	if val == "" {
		return layout.Container()
	}
	return result.Success(val)
}

// pasteResult renders the paste demo outcome.
func pasteResult(val string) node.Node {
	if val == "" {
		return hint.Text("Paste into the input to see the server response.")
	}
	return result.Success(val)
}

// contextMenuResult renders the context menu demo outcome.
func contextMenuResult(val string) node.Node {
	if val == "" {
		return hint.Text("Right-click the box above.")
	}
	return result.Success(val)
}

// validatedResult renders the validation demo outcome.
func validatedResult(val string) node.Node {
	if val == "" {
		return hint.Text("Fill in the field and submit.")
	}
	return result.Success(val)
}

// editableResult renders the contenteditable demo outcome.
func editableResult(val string) node.Node {
	if val == "" {
		return hint.Text("Click the text above to edit, then click away.")
	}
	return result.Success(val)
}
