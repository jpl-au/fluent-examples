// ws.js - WebSocket client for the live log demo.
//
// Opens a WebSocket connection to the server's /ws/feed endpoint and
// renders incoming log entries into the #ws-log container. The
// connection status is shown in #ws-status.
//
// Reconnection: if the connection drops, the client retries after a
// short delay with exponential backoff up to 10 seconds.

(function () {
  "use strict";

  var log = document.getElementById("ws-log");
  var status = document.getElementById("ws-status");
  var maxEntries = 200;
  var retryDelay = 1000;

  function connect() {
    // Build the WebSocket URL from the current page origin, swapping
    // http(s) for ws(s) so it works in both dev and production.
    var protocol = location.protocol === "https:" ? "wss:" : "ws:";
    var url = protocol + "//" + location.host + "/ws/feed";
    var ws = new WebSocket(url);

    ws.onopen = function () {
      status.textContent = "connected";
      status.className = "status-connected";
      retryDelay = 1000;
    };

    ws.onmessage = function (event) {
      var entry = JSON.parse(event.data);
      appendEntry(entry);
    };

    ws.onclose = function () {
      status.textContent = "reconnecting...";
      status.className = "status-reconnecting";
      scheduleReconnect();
    };

    ws.onerror = function () {
      status.textContent = "disconnected";
      status.className = "status-disconnected";
    };
  }

  // scheduleReconnect waits with exponential backoff before retrying.
  function scheduleReconnect() {
    setTimeout(connect, retryDelay);
    retryDelay = Math.min(retryDelay * 1.5, 10000);
  }

  // appendEntry creates a log line element and appends it to the feed
  // container, auto-scrolling to keep the latest entry visible. Old
  // entries beyond maxEntries are removed to prevent unbounded growth.
  function appendEntry(entry) {
    var row = document.createElement("div");
    row.className = "log-entry";

    row.innerHTML =
      '<span class="log-time">' + esc(entry.time) + "</span>" +
      '<span class="log-level log-level-' + esc(entry.level) + '">' + esc(entry.level) + "</span>" +
      '<span class="log-message">' + esc(entry.message) + "</span>" +
      '<span class="log-attrs">' + esc(entry.attrs) + "</span>";

    // Newest entries at the top so the latest is always visible.
    log.insertBefore(row, log.firstChild);

    // Trim old entries from the bottom so the DOM stays bounded.
    while (log.children.length > maxEntries) {
      log.removeChild(log.lastChild);
    }
  }

  // esc escapes HTML entities to prevent injection from log content.
  function esc(str) {
    var el = document.createElement("span");
    el.textContent = str;
    return el.innerHTML;
  }

  connect();
})();
