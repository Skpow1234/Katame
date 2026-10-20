/**
 * Browser/Node loader for the Katame WASM module.
 *
 * The Go runtime needs its companion `wasm_exec.js` (shipped with the Go
 * SDK at $(go env GOROOT)/lib/wasm/wasm_exec.js as of Go 1.22). Drop a
 * copy alongside this file, then:
 *
 *   import { loadKatame } from "./katame-loader.js";
 *   await loadKatame("./katame.wasm");   // returns globalThis.Katame
 *
 * After the promise resolves the global `Katame` object (typed by
 * katame.d.ts) is ready to use.
 */

/**
 * @param {string|URL} wasmURL Path or URL to katame.wasm.
 * @param {object}     [opts]
 * @param {string|URL} [opts.wasmExecURL=./wasm_exec.js]
 *   Location of Go's wasm_exec.js. Defaults to a sibling of wasmURL.
 * @returns {Promise<import("./katame").Katame>}
 */
export async function loadKatame(wasmURL, opts = {}) {
  if (typeof globalThis.Go !== "function") {
    const execURL = opts.wasmExecURL ?? new URL("./wasm_exec.js", wasmURL);
    // Browser path.
    if (typeof document !== "undefined") {
      await new Promise((resolve, reject) => {
        const s = document.createElement("script");
        s.src = execURL.toString();
        s.onload = () => resolve();
        s.onerror = reject;
        document.head.appendChild(s);
      });
    } else {
      // Node path: dynamic import keeps this file dependency-free.
      await import(execURL.toString());
    }
  }

  // eslint-disable-next-line new-cap
  const go = new globalThis.Go();

  // Node's `fetch` does not support file:// URLs, so for local file inputs
  // (string with a relative path, an absolute path, or a file: URL) we use
  // fs/promises directly. HTTPS URLs go through fetch as usual.
  const isFileLike =
    (typeof wasmURL === "string" && !/^https?:\/\//i.test(wasmURL)) ||
    (wasmURL && wasmURL.protocol === "file:");

  let bytes;
  if (isFileLike && typeof process !== "undefined") {
    const fs = await import("node:fs/promises");
    bytes = await fs.readFile(wasmURL);
  } else {
    const resp = await fetch(wasmURL);
    bytes = new Uint8Array(await resp.arrayBuffer());
  }
  const result = await WebAssembly.instantiate(bytes, go.importObject);

  // run() returns when main() exits — Go's main() uses `select {}` to keep
  // the runtime alive, so we deliberately do not await this promise.
  go.run(result.instance);
  if (typeof globalThis.Katame === "undefined") {
    throw new Error("katame.wasm did not install globalThis.Katame");
  }
  return globalThis.Katame;
}
