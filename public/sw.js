// This is the script for the servide worker ("SW"), 
// which will act as an HTTP server deciding whether to
// (a) handle a request all by itself, or (b) forward it.

// An SW cannot both be a proxy and access the DOM.
// https://elijahm.com/posts/local_first_htmx_part2/
// A service worker is a separate thread in the browser that 
// has special privileges. MDN documentation describes it:
// “Service workers essentially act as proxy servers that sit
// between web apps, the browser, and the network (when avail-
// able).” Now the service worker does have many of the same
// restrictions as a web worker, such as no access to the DOM. 

// Fetch the cript we need, from the EXACT Go version.
// (We have copied it into our own directory.)
importScripts("wasm_exec.js");
// importScripts(
//   "https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js",
// );

addEventListener("install", (event) => {
  event.waitUntil(skipWaiting());
});

addEventListener("activate", (event) => {
  event.waitUntil(clients.claim());
});

registerWasmHTTPListener("main.wasm", { base: "" });

function registerWasmHTTPListener(wasm, { base, args = [] } = {}) {
  let path = new URL(registration.scope).pathname;
  if (base && base !== "") {
    path = `${trimEnd(path, "/")}/${trimStart(base, "/")}`;
  }

  const handlerPromise = new Promise((setHandler) => {
    self.wasmhttp = {
      path,
      setHandler,
    };
  });

  const go = new Go();
  go.argv = [wasm, ...args];
  // Fetch the wasm binary and start it up.
  WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then((
    { instance },
  ) => go.run(instance));

  // Add an event listener for “fetch”: this registers
  // that we are proxying network requests.
  addEventListener("fetch", (e) => {
    const url = new URL(e.request.url);
    const { pathname } = new URL(e.request.url);
    console.log(e);
    console.log("pathname: ", pathname);
    console.log("path: ", path);
    if (
      pathname === "/wasm_exec.js" || pathname == "/sw.js" ||
      pathname === "/start_worker.js"
    ) {
      e.respondWith(fetch(e.request));
      return;
    } else if (url.hostname === "localhost") {
      e.respondWith(handlerPromise.then((handler) => handler(e.request)));
    } else {
      // For requests to other domains, just pass them along to the network
      e.respondWith(fetch(e.request));
    }
    // if (!pathname.startsWith(path)) {
    //   console.log("fallback to network");
    //   return fetch(request);
    // }
  });
}

function trimStart(s, c) {
  let r = s;
  while (r.startsWith(c)) r = r.slice(c.length);
  return r;
}

function trimEnd(s, c) {
  let r = s;
  while (r.endsWith(c)) r = r.slice(0, -c.length);
  return r;
}
