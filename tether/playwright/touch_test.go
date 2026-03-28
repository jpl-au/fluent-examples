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

// TestTouchSwipeViaJS dispatches synthetic touch events to verify
// the swipe handler fires correctly.
func TestTouchSwipeViaJS(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/touch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Dispatch a synthetic swipe (touchstart + touchend with offset).
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
		t.Errorf("swipe result not visible: %v", err)
	}
}
