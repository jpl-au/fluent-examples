package playwright_test

import "testing"

// TestTouchPageRenders verifies the touch gestures demo loads.
func TestTouchPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	swipe := page.GetByText("Swipe here on a touch device")
	if err := expect(swipe).ToBeVisible(); err != nil {
		t.Fatalf("swipe area not visible: %v", err)
	}

	longPress := page.GetByText("Long-press here on a touch device")
	if err := expect(longPress).ToBeVisible(); err != nil {
		t.Fatalf("long-press area not visible: %v", err)
	}
}

// TestTouchSwipeLeft dispatches synthetic touch events to verify
// the swipe left handler fires correctly.
func TestTouchSwipeLeft(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Dispatch a synthetic swipe left (large negative X offset).
	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-swipe]');
		if (!el) return;
		el.dispatchEvent(new TouchEvent('touchstart', {
			touches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
			bubbles: true
		}));
		el.dispatchEvent(new TouchEvent('touchend', {
			changedTouches: [new Touch({identifier: 0, target: el, clientX: 50, clientY: 200})],
			bubbles: true
		}));
	}`)
	if err != nil {
		t.Fatalf("dispatch touch: %v", err)
	}

	result := page.GetByText("Swiped left")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("swipe left result not visible: %v", err)
	}
}

// TestTouchSwipeRight dispatches a rightward swipe and verifies
// the server receives direction "right".
func TestTouchSwipeRight(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-swipe]');
		if (!el) return;
		el.dispatchEvent(new TouchEvent('touchstart', {
			touches: [new Touch({identifier: 0, target: el, clientX: 50, clientY: 200})],
			bubbles: true
		}));
		el.dispatchEvent(new TouchEvent('touchend', {
			changedTouches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
			bubbles: true
		}));
	}`)
	if err != nil {
		t.Fatalf("dispatch touch: %v", err)
	}

	result := page.GetByText("Swiped right")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("swipe right result not visible: %v", err)
	}
}

// TestTouchSwipeUp dispatches an upward swipe and verifies the
// server receives direction "up".
func TestTouchSwipeUp(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-swipe]');
		if (!el) return;
		el.dispatchEvent(new TouchEvent('touchstart', {
			touches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
			bubbles: true
		}));
		el.dispatchEvent(new TouchEvent('touchend', {
			changedTouches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 50})],
			bubbles: true
		}));
	}`)
	if err != nil {
		t.Fatalf("dispatch touch: %v", err)
	}

	result := page.GetByText("Swiped up")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("swipe up result not visible: %v", err)
	}
}

// TestTouchSwipeDown dispatches a downward swipe and verifies the
// server receives direction "down".
func TestTouchSwipeDown(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-swipe]');
		if (!el) return;
		el.dispatchEvent(new TouchEvent('touchstart', {
			touches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 50})],
			bubbles: true
		}));
		el.dispatchEvent(new TouchEvent('touchend', {
			changedTouches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
			bubbles: true
		}));
	}`)
	if err != nil {
		t.Fatalf("dispatch touch: %v", err)
	}

	result := page.GetByText("Swiped down")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("swipe down result not visible: %v", err)
	}
}

// TestTouchLongPress dispatches a synthetic long-press and verifies
// the server fires the longpress event.
func TestTouchLongPress(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Simulate a long press: touchstart, wait 600ms (>500ms threshold),
	// then touchend at the same position (within 10px tolerance).
	_, err = page.Evaluate(`() => {
		return new Promise(resolve => {
			var el = document.querySelector('[data-tether-longpress]');
			if (!el) { resolve(); return; }
			el.dispatchEvent(new TouchEvent('touchstart', {
				touches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
				bubbles: true
			}));
			setTimeout(() => {
				el.dispatchEvent(new TouchEvent('touchend', {
					changedTouches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
					bubbles: true
				}));
				resolve();
			}, 600);
		});
	}`)
	if err != nil {
		t.Fatalf("dispatch long press: %v", err)
	}

	result := page.GetByText("Long press detected!")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("long press result not visible: %v", err)
	}
}

// TestTouchSwipeCancelledByShortDistance verifies that a touch
// movement under the 30px threshold does not register as a swipe.
func TestTouchSwipeCancelledByShortDistance(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Move only 15px horizontally - well under the 30px threshold.
	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-swipe]');
		if (!el) return;
		el.dispatchEvent(new TouchEvent('touchstart', {
			touches: [new Touch({identifier: 0, target: el, clientX: 200, clientY: 200})],
			bubbles: true
		}));
		el.dispatchEvent(new TouchEvent('touchend', {
			changedTouches: [new Touch({identifier: 0, target: el, clientX: 185, clientY: 200})],
			bubbles: true
		}));
	}`)
	if err != nil {
		t.Fatalf("dispatch touch: %v", err)
	}

	// The hint text should still be showing - no swipe result.
	hint := page.GetByText("Swipe to see the direction.")
	if err := expect(hint).ToBeVisible(); err != nil {
		t.Errorf("expected hint to remain visible after short swipe: %v", err)
	}
}
