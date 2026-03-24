package playwright_test

import (
	"testing"
)

// TestDragDropPageRenders verifies the drag-and-drop demo loads and
// shows both zones with items.
func TestDragDropPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/dragdrop/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	alpha := page.GetByText("Alpha")
	if err := expect(alpha).ToBeVisible(); err != nil {
		t.Fatalf("Alpha item not visible: %v", err)
	}

	charlie := page.GetByText("Charlie")
	if err := expect(charlie).ToBeVisible(); err != nil {
		t.Fatalf("Charlie item not visible: %v", err)
	}

	zoneA := page.GetByText("Zone A")
	if err := expect(zoneA).ToBeVisible(); err != nil {
		t.Fatalf("Zone A not visible: %v", err)
	}

	zoneB := page.GetByText("Zone B")
	if err := expect(zoneB).ToBeVisible(); err != nil {
		t.Fatalf("Zone B not visible: %v", err)
	}
}

// TestDragDropMoveItem drags an item from Zone A to Zone B and
// verifies the item appears in the target AND is gone from the
// source. The second check catches the bug where dataTransfer is
// empty and the server silently does nothing.
func TestDragDropMoveItem(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/dragdrop/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	sourceZone := page.Locator("[data-tether-data-zone='left']")
	targetZone := page.Locator("[data-tether-data-zone='right']")

	// Alpha (id=1) starts in Zone A.
	source := page.Locator("[data-tether-data-id='1']")
	if err := expect(source).ToBeVisible(); err != nil {
		t.Fatalf("source item not visible: %v", err)
	}

	// Verify Alpha is in Zone A before the drag.
	alphaInA := sourceZone.Locator("[data-tether-data-id='1']")
	if err := expect(alphaInA).ToBeVisible(); err != nil {
		t.Fatalf("Alpha should start in Zone A: %v", err)
	}

	// Perform drag and drop.
	if err := source.DragTo(targetZone); err != nil {
		t.Fatalf("drag to: %v", err)
	}

	// Alpha should now be in Zone B.
	alphaInB := targetZone.Locator("[data-tether-data-id='1']")
	if err := expect(alphaInB).ToBeVisible(); err != nil {
		t.Errorf("Alpha not found in Zone B after drag: %v", err)
	}

	// Alpha should be GONE from Zone A. This is the critical check -
	// without it, a no-op server response (empty dataTransfer) would
	// pass because the item stays in its original zone.
	alphaStillInA := sourceZone.Locator("[data-tether-data-id='1']")
	if err := expect(alphaStillInA).Not().ToBeVisible(); err != nil {
		t.Errorf("Alpha should no longer be in Zone A: %v", err)
	}
}

// TestDragDropDataTransfer verifies the DnD JS correctly populates
// dataTransfer during dragstart. This tests the actual event handler
// code path rather than relying on Playwright's DragTo simulation.
func TestDragDropDataTransfer(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/dragdrop/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Dispatch a real dragstart event via JS and read back what the
	// DnD extension wrote to dataTransfer.
	result, err := page.Evaluate(`() => {
		var source = document.querySelector('[data-tether-data-id="1"]');
		if (!source) return {error: "source not found"};
		var dt = new DataTransfer();
		var ev = new DragEvent('dragstart', {dataTransfer: dt, bubbles: true});
		source.dispatchEvent(ev);
		var raw = dt.getData('application/tether');
		if (!raw) return {error: "dataTransfer empty"};
		return JSON.parse(raw);
	}`)
	if err != nil {
		t.Fatalf("evaluate dragstart: %v", err)
	}

	data, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}

	if errMsg, has := data["error"]; has {
		t.Fatalf("dragstart failed: %s", errMsg)
	}

	id, _ := data["id"].(string)
	if id != "1" {
		t.Errorf("dataTransfer id = %q, want %q", id, "1")
	}
}

// TestDragDropCrossSession verifies that dragging an item in one tab
// updates another tab via Group.Broadcast.
func TestDragDropCrossSession(t *testing.T) {
	srv := startApp(t, serverMode())
	sender, cleanupSender := newPage(t)
	defer cleanupSender()
	receiver, cleanupReceiver := newPage(t)
	defer cleanupReceiver()

	if _, err := sender.Goto(srv + "/dragdrop/"); err != nil {
		t.Fatalf("sender goto: %v", err)
	}
	waitForConnected(t, sender)

	if _, err := receiver.Goto(srv + "/dragdrop/"); err != nil {
		t.Fatalf("receiver goto: %v", err)
	}
	waitForConnected(t, receiver)

	// Drag Bravo (id=2) from Zone A to Zone B on the sender.
	source := sender.Locator("[data-tether-data-id='2']")
	target := sender.Locator("[data-tether-data-zone='right']")
	if err := source.DragTo(target); err != nil {
		t.Fatalf("drag to: %v", err)
	}

	// Verify the receiver sees Bravo in Zone B.
	receiverTarget := receiver.Locator("[data-tether-data-zone='right']")
	movedItem := receiverTarget.Locator("[data-tether-data-id='2']")
	if err := expect(movedItem).ToBeVisible(); err != nil {
		t.Errorf("receiver did not see Bravo in Zone B: %v", err)
	}

	// Verify the receiver sees Bravo gone from Zone A.
	receiverSource := receiver.Locator("[data-tether-data-zone='left']")
	bravoStillInA := receiverSource.Locator("[data-tether-data-id='2']")
	if err := expect(bravoStillInA).Not().ToBeVisible(); err != nil {
		t.Errorf("receiver should not see Bravo in Zone A: %v", err)
	}
}
