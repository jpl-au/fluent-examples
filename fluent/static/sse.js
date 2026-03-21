// sse.js - Server-Sent Events client for the live log demo.
//
// Opens an EventSource connection to /sse/feed and renders incoming
// log entries into the #sse-log container. The browser handles
// reconnection automatically when the connection drops.

(function () {
  "use strict";

  var log = document.getElementById("sse-log");
  var status = document.getElementById("sse-status");
  var maxEntries = 200;

  var source = new EventSource("/sse/feed");

  source.onopen = function () {
    status.textContent = "connected";
    status.className = "status-connected";
  };

  source.onmessage = function (event) {
    var entry = JSON.parse(event.data);
    appendEntry(entry);
  };

  source.onerror = function () {
    // EventSource reconnects automatically.
    status.textContent = "reconnecting...";
    status.className = "status-reconnecting";
  };

  // appendEntry creates a log line element and appends it to the
  // feed container, auto-scrolling to keep the latest entry visible.
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
})();
