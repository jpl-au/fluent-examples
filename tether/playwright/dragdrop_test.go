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

// TestDragDropMoveItem drags an item from Zone A to Zone B using
// Playwright's drag-and-drop API and verifies the item appears in
// the target zone after the server processes the drop event.
func TestDragDropMoveItem(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/dragdrop/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Alpha starts in Zone A (left zone).
	source := page.Locator("[data-tether-data-id='1']")
	if err := expect(source).ToBeVisible(); err != nil {
		t.Fatalf("source item not visible: %v", err)
	}

	// Drop target is Zone B (right zone).
	target := page.Locator("[data-tether-data-zone='right']")
	if err := expect(target).ToBeVisible(); err != nil {
		t.Fatalf("target zone not visible: %v", err)
	}

	// Perform drag and drop.
	if err := source.DragTo(target); err != nil {
		t.Fatalf("drag to: %v", err)
	}

	// After the drop, the server moves the item and broadcasts.
	// Alpha should now be inside Zone B. Verify by checking the
	// item is a child of the right zone.
	movedItem := target.Locator("[data-tether-data-id='1']")
	if err := expect(movedItem).ToBeVisible(); err != nil {
		t.Errorf("Alpha not found in Zone B after drag: %v", err)
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

	// Open both pages and wait for connections.
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

	// Verify the receiver also sees Bravo in Zone B.
	receiverTarget := receiver.Locator("[data-tether-data-zone='right']")
	movedItem := receiverTarget.Locator("[data-tether-data-id='2']")
	if err := expect(movedItem).ToBeVisible(); err != nil {
		t.Errorf("receiver did not see Bravo in Zone B: %v", err)
	}
}
