package playwright_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	tether "github.com/jpl-au/tether"
	pw "github.com/playwright-community/playwright-go"
)

// Exercise the real runtime and DOM with controlled transport delivery and
// animation scheduling. The browser still performs all event delegation,
// morphing, history and storage operations.
func auditClient(t *testing.T, content string, attrs ...string) pw.Page {
	t.Helper()
	// Serve the selected module's embedded runtime, including when GOWORK=off.
	mux := http.NewServeMux()
	mux.Handle("/_tether/", http.StripPrefix("/_tether/", tether.ServeClient()))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!doctype html><div data-tether-root data-tether-endpoint="/" data-tether-session="ORIGINAL" data-tether-transport="sse" data-tether-flash-duration="137" %s><span id="count">0</span>%s</div><script src="/_tether/idiomorph.min.js"></script><script src="/_tether/tether.js"></script>`, strings.Join(attrs, " "), content)
	})
	base := startServer(t, RealHTTP1, mux)
	page, cleanup := newPage(t)
	t.Cleanup(cleanup)
	if err := page.AddInitScript(pw.Script{Content: pw.String(`
window.__frames = new Map(); window.__frameID = 0; window.__requests = []; window.__timeouts = []; window.__beacons = []; window.__listeners = [];
const realAddListener = window.addEventListener.bind(window);
window.addEventListener = (type, ...args) => { __listeners.push(type); return realAddListener(type, ...args); };
navigator.sendBeacon = (url, body) => { __beacons.push({url, body}); return true; };
window.requestAnimationFrame = fn => { const id = ++__frameID; __frames.set(id, fn); return id; };
window.cancelAnimationFrame = id => __frames.delete(id);
window.__flush = () => { const fs = Array.from(__frames.values()); __frames.clear(); fs.forEach(fn => fn(performance.now())); };
const realTimeout = window.setTimeout;
window.setTimeout = (fn, ms, ...args) => { __timeouts.push(ms); return realTimeout(fn, ms, ...args); };
window.fetch = (url, opts) => { __requests.push({url: String(url), ...opts}); return Promise.resolve(new Response("ticket", {status: 200})); };
window.EventSource = class { constructor() { window.__source = this; this.readyState = 1; queueMicrotask(() => this.onopen && this.onopen()); } close() { this.readyState = 2; } };
window.__deliver = msg => __source.onmessage({data: JSON.stringify({type: "update", ...msg})});
Object.defineProperty(navigator, "clipboard", {value: {writeText: () => Promise.resolve()}});
`)}); err != nil {
		t.Fatal(err)
	}
	if _, err := page.Goto(base); err != nil {
		t.Fatal(err)
	}
	if _, err := page.WaitForFunction(`window.__source && document.querySelector('[data-tether-root]').getAttribute('data-tether-state') === 'connected'`, nil); err != nil {
		t.Fatal(err)
	}
	return page
}

func auditEval(t *testing.T, page pw.Page, script string) any {
	t.Helper()
	result, err := page.Evaluate(script)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestAuditClientSessionAssignmentBeforeFrame(t *testing.T) {
	page := auditClient(t, "")
	got := auditEval(t, page, `() => { __deliver({session: 'REPLACEMENT'}); Tether.sendEvent('click', 'test', {}); return __requests.at(-1).headers['Tether-Session']; }`)
	if got != "REPLACEMENT" {
		t.Errorf("event used stale session: %v", got)
	}
}

func TestAuditClientQueuesRetryablePost(t *testing.T) {
	for _, status := range []int{429, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			page := auditClient(t, "", `data-tether-background-sync="true"`)
			auditEval(t, page, fmt.Sprintf(`() => {
window.fetch = () => Promise.resolve(new Response('', {status:%d}));
Tether.sendEvent('click', 'retry', {});
}`, status))
			got := auditEval(t, page, `async () => {
const db = await new Promise((resolve,reject) => { const r=indexedDB.open('tether-events',1); r.onupgradeneeded=()=>r.result.createObjectStore('queue',{autoIncrement:true}); r.onsuccess=()=>resolve(r.result); r.onerror=reject; });
let rows=[]; const end=Date.now()+1000;
do { rows=await new Promise(resolve => { const r=db.transaction('queue').objectStore('queue').getAll(); r.onsuccess=()=>resolve(r.result); }); if(rows.length) break; await new Promise(resolve=>setTimeout(resolve,20)); } while(Date.now()<end);
db.close(); return rows.map(x=>JSON.parse(x.payload).action).join(',');
}`)
			if got != "retry" {
				t.Errorf("retryable POST was not retained: %v", got)
			}
		})
	}
}

func TestAuditEventReplayPreservesOrder(t *testing.T) {
	for _, worker := range []bool{false, true} {
		t.Run(map[bool]string{false: "client", true: "worker"}[worker], func(t *testing.T) {
			page := auditClient(t, "", `data-tether-background-sync="true"`)
			if worker {
				if _, err := page.AddScriptTag(pw.PageAddScriptTagOptions{URL: pw.String("/_tether/tether-worker.js")}); err != nil {
					t.Fatal(err)
				}
			}
			auditEval(t, page, `async () => {
window.__db = await new Promise((resolve,reject) => { const r=indexedDB.open('tether-events',1); r.onupgradeneeded=()=>r.result.createObjectStore('queue',{autoIncrement:true}); r.onsuccess=()=>resolve(r.result); r.onerror=reject; });
await new Promise(resolve => { const tx=__db.transaction('queue','readwrite'); for(let i=1;i<=3;i++)tx.objectStore('queue').add({sessionID:'ORIGINAL', endpoint:location.origin+'/', payload:JSON.stringify({action:String(i)}), ts:Date.now()}); tx.oncomplete=resolve; });
window.__attempts=[]; window.__reply=null;
window.fetch=(url,opts) => { if(String(url).includes('tether=ticket'))return Promise.resolve(new Response('ticket')); __attempts.push(JSON.parse(opts.body).action); return new Promise(resolve => { __reply=resolve; }); };
}`)
			trigger := `dispatchEvent(new PageTransitionEvent('pageshow', {persisted:true}))`
			if worker {
				trigger = `void (window.__replay = drainAndReplay().catch(() => {}))`
			}
			auditEval(t, page, trigger)
			if _, err := page.WaitForFunction(`__attempts.length > 0`, nil); err != nil {
				t.Fatal(err)
			}
			if got := auditEval(t, page, `__attempts.join(',')`); got != "1" {
				t.Errorf("replay sent later events before acknowledgement: %v", got)
			}
			auditEval(t, page, `__reply(new Response('',{status:429}))`)
			// A new transaction follows response processing and observes any deletion.
			got := auditEval(t, page, `async () => { await new Promise(resolve=>setTimeout(resolve,50)); return new Promise(resolve=>{ const r=__db.transaction('queue').objectStore('queue').getAll(); r.onsuccess=()=>resolve(r.result.map(x=>JSON.parse(x.payload).action).join(',')); }); }`)
			if got != "1,2,3" {
				t.Errorf("overload discarded queued events: %v", got)
			}
			auditEval(t, page, `() => { window.__attempts=[]; window.fetch=(url,opts)=>{ if(String(url).includes('tether=ticket')) return Promise.resolve(new Response('ticket')); __attempts.push(JSON.parse(opts.body).action); return Promise.resolve(new Response(null,{status:204})); }; }`)
			auditEval(t, page, trigger)
			got = auditEval(t, page, `async () => {
const end=Date.now()+2000; let count;
do { count=await new Promise(resolve=>{ const r=__db.transaction('queue').objectStore('queue').count(); r.onsuccess=()=>resolve(r.result); }); if(!count)break; await new Promise(resolve=>setTimeout(resolve,20)); } while(Date.now()<end);
return {remaining:count, sent:__attempts.join(',')};
}`)
			result := got.(map[string]any)
			if result["remaining"] != 0 || result["sent"] != "1,2,3" {
				t.Errorf("successful retry did not drain in order: %v", result)
			}
		})
	}
}

func TestAuditWorkerRetriesSessionsIndependently(t *testing.T) {
	page := auditClient(t, "")
	if _, err := page.AddScriptTag(pw.PageAddScriptTagOptions{URL: pw.String("/_tether/tether-worker.js")}); err != nil {
		t.Fatal(err)
	}
	got := auditEval(t, page, `async () => {
const db=await openEventDB();
await new Promise(resolve=>{const tx=db.transaction('queue','readwrite'); [['A','a1'],['A','a2'],['B','b1']].forEach(([sessionID,action])=>tx.objectStore('queue').add({sessionID,endpoint:location.origin+'/',payload:JSON.stringify({action}),ts:Date.now()}));tx.oncomplete=resolve;});
const sent=[]; window.fetch=(url,opts)=>{const action=JSON.parse(opts.body).action;sent.push(action);return Promise.resolve(new Response(null,{status:action==='a1'?503:204}));};
let retry=false;try{await drainAndReplay();}catch{retry=true;}
const remaining=await new Promise(resolve=>{const r=db.transaction('queue').objectStore('queue').getAll();r.onsuccess=()=>resolve(r.result.map(x=>JSON.parse(x.payload).action).join(','));});db.close();
return {sent:sent.sort().join(','),remaining,retry};
}`)
	result := got.(map[string]any)
	if result["sent"] != "a1,b1" || result["remaining"] != "a1,a2" || result["retry"] != true {
		t.Errorf("one session blocked another session's replay: %v", result)
	}
}

func TestAuditClientHiddenUpdatesStayOrdered(t *testing.T) {
	page := auditClient(t, "")
	got := auditEval(t, page, `() => {
__deliver({patches:[{key:'count', html:'<span id="count">visible</span>'}]});
Object.defineProperty(document, 'hidden', {configurable:true, value:true}); document.dispatchEvent(new Event('visibilitychange'));
for (let i = 0; i < 100; i++) __deliver({patches:[{key:'count', html:'<span id="count">'+i+'</span>'}]});
return {count:document.getElementById('count').textContent, queued:__frames.size};
}`)
	values := got.(map[string]any)
	if values["count"] != "99" || values["queued"] != 0 {
		t.Errorf("hidden updates retained or reordered: %v", values)
	}
	if got := auditEval(t, page, `() => { Object.defineProperty(document, 'hidden', {value:false}); document.dispatchEvent(new Event('visibilitychange')); __flush(); return document.getElementById('count').textContent; }`); got != "99" {
		t.Errorf("old frame overwrote newer state: %v", got)
	}
}

func TestAuditClientEffectDoesNotDiscoverDOM(t *testing.T) {
	page := auditClient(t, "")
	got := auditEval(t, page, `() => { const root=document.querySelector('[data-tether-root]'); let scans=0; for (const name of ['querySelector','querySelectorAll']) { const original=root[name]; root[name]=function(...args){scans++;return original.apply(this,args)}; } __deliver({toast:'hello'}); __flush(); return scans; }`)
	if got != 0 {
		t.Errorf("effect-only update rescanned the DOM %v times", got)
	}
}

func TestAuditClientConfiguredFlashDuration(t *testing.T) {
	page := auditClient(t, `<button id="copy" data-tether-copy="#count" data-tether-flash-text="Copied">Copy</button>`)
	auditEval(t, page, `document.getElementById('copy').click()`)
	if _, err := page.WaitForFunction(`__timeouts.length > 0`, nil); err != nil {
		t.Fatal(err)
	}
	if got := auditEval(t, page, `__timeouts.includes(137)`); got != true {
		t.Errorf("configured Flash duration ignored: timers=%v", auditEval(t, page, `__timeouts`))
	}
}

func TestAuditClientLinksRespectPreventDefault(t *testing.T) {
	for _, own := range []bool{false, true} {
		t.Run(fmt.Sprint(own), func(t *testing.T) {
			attrs := ""
			if own {
				attrs = ` data-tether-event-click="handled" data-tether-prevent-default`
			}
			page := auditClient(t, `<a id="link" href="/other" data-tether-link`+attrs+`>link</a>`)
			if !own {
				auditEval(t, page, `document.getElementById('link').addEventListener('click', e => e.preventDefault())`)
			}
			auditEval(t, page, `document.getElementById('link').click()`)
			if got := auditEval(t, page, `__requests.filter(r => r.body && JSON.parse(r.body).type === 'navigate').length`); got != 0 {
				t.Errorf("prevented click still navigated: %v", got)
			}
		})
	}
}

func TestAuditClientPageCacheLifecycle(t *testing.T) {
	page := auditClient(t, "")
	if got := auditEval(t, page, `__listeners.includes('beforeunload')`); got != false {
		t.Error("runtime installs a permanent beforeunload listener")
	}
	if got := auditEval(t, page, `() => { window.__oldSource=__source; dispatchEvent(new PageTransitionEvent('pagehide',{persisted:true})); return __beacons.length === 0 && __oldSource.readyState === 2; }`); got != true {
		t.Error("cached page did not close transport while retaining its session")
	}
	auditEval(t, page, `dispatchEvent(new PageTransitionEvent('pageshow',{persisted:true}))`)
	if _, err := page.WaitForFunction(`__source !== __oldSource`, nil, pw.PageWaitForFunctionOptions{Timeout: pw.Float(2000)}); err != nil {
		t.Error("cached page did not reconnect")
	}
	auditEval(t, page, `dispatchEvent(new PageTransitionEvent('pagehide',{persisted:false}))`)
	if got := auditEval(t, page, `__beacons.length === 1 && __beacons[0].body === 'ORIGINAL'`); got != true {
		t.Error("departing page did not destroy its session")
	}
}

func TestAuditClientDelayedTransitionCannotReorderUpdates(t *testing.T) {
	page := auditClient(t, "", `data-tether-view-transitions`)
	auditEval(t, page, `() => {
window.__updates=0; document.addEventListener('tether:update', () => __updates++);
document.startViewTransition = fn => { window.__lateDOM=fn; return {skipTransition(){}, updateCallbackDone:new Promise(resolve => window.__lateDone=resolve)}; };
__deliver({patches:[{key:'count',html:'<span id="count">first</span>'}]}); __flush();
for(let i=0;i<300;i++) __deliver({patches:[{key:'count',html:'<span id="count">'+i+'</span>'}]});
Object.defineProperty(document,'hidden',{configurable:true,value:true}); document.dispatchEvent(new Event('visibilitychange'));
}`)
	if got := auditEval(t, page, `document.getElementById('count').textContent === '299' && __frames.size === 0 && __updates === 301`); got != true {
		t.Fatal("queued transition did not preserve ordered, bounded delivery")
	}
	auditEval(t, page, `__lateDOM(); __lateDone()`)
	if got := auditEval(t, page, `document.getElementById('count').textContent === '299' && __updates === 301`); got != true {
		t.Error("late transition replayed stale DOM or duplicate effects")
	}
}

func TestAuditClientBatchesConditionalEvaluation(t *testing.T) {
	page := auditClient(t, `<div id="condition" data-tether-bind-show-when='[[0,"a"],[0,"b"],[3,"and"]]'>visible</div>`)
	got := auditEval(t, page, `() => {
let scans=0; const query=document.querySelectorAll.bind(document);
document.querySelectorAll = selector => { if(selector.includes('data-tether-bind-show-when')) scans++; return query(selector); };
__deliver({signals:{a:true,b:false}}); __flush();
return {scans, hidden:document.getElementById('condition').style.display === 'none'};
}`)
	result := got.(map[string]any)
	if result["scans"] != 1 || result["hidden"] != true {
		t.Errorf("conditional evaluation repeated or incorrect: %v", result)
	}
}
