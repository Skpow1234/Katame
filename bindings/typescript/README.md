# @katame/wasm

WebAssembly bindings for [Katame](https://github.com/Skpow1234/Katame)
— encrypt, decrypt, sign, and verify `.katame` bundles directly in the
browser or in Node.js, with no native code required at install time.

## Installation

```bash
npm install @katame/wasm
```

Two files ship with the package:

- `katame.wasm`         — the Go-compiled WASM module (~8 MB).
- `katame-loader.js`    — a tiny ESM loader that pulls in Go's
  `wasm_exec.js` and instantiates the module.

Build from the Katame source tree with:

```bash
GOOS=js GOARCH=wasm go build -o katame.wasm ./cmd/katame-wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
```

## Quick start (Node 20+)

```javascript
import { loadKatame } from "@katame/wasm";

const Katame = await loadKatame(
  new URL("./katame.wasm", import.meta.url)
);

console.log("Version:", Katame.version());

// Encrypt with a generated key.
const plaintext = new TextEncoder().encode("hello, sdk");
const enc = Katame.protect({ plaintext });
if (!enc.ok) throw new Error(enc.error);
console.log("manifest:", enc.manifest);

// Decrypt.
const dec = Katame.decrypt({
  bundle: enc.bundle,
  key: enc.generatedKey,
});
if (!dec.ok) throw new Error(dec.error);
console.log(new TextDecoder().decode(dec.plaintext));
```

## Quick start (browser)

```html
<script type="module">
  import { loadKatame } from "./katame-loader.js";
  const Katame = await loadKatame("./katame.wasm");
  document.body.textContent = "Katame " + Katame.version();
</script>
```

## API

See `katame.d.ts` for the full TypeScript surface. The shape is the
same as the C-shared and Python wrappers:

- `version() -> string`
- `protect({plaintext, key?, password?, cipher?, inputName?, aad?})`
- `decrypt({bundle, key?, password?, aad?})`
- `inspect(bundleBytes)`
- `sign({bundle, privateKeyPEM, algo?})`
- `verify({bundle, publicKeyPEM})`

Every call returns either `{ok: true, ...}` on success or
`{ok: false, error: "..."}` on failure.

## Supported modes (v1.0)

- Symmetric key (auto-generated or supplied).
- Password-based (argon2id / scrypt / pbkdf2-sha256).
- Hybrid recipients, KMS wrap, compression, Shamir split (see `ProtectOptions` / `DecryptOptions` in `katame.d.ts`).
- Detached signatures (ed25519 / ecdsa / rsa-pss / ml-dsa / slh-dsa).
- Bundle inspection.

Sigstore transparency uploads remain CLI-oriented for now.

## License

MIT.
