# Dummy KEM Plugin (Example)

A minimal Katame plugin that registers the KEM scheme **`dummy-kem`**. For demonstration only — the "encapsulation" is a trivial XOR and is **not secure**.

## Build

```bash
go build -o dummy-kem .
# On Windows: dummy-kem.exe
```

## Use

1. Put the built executable in a directory and set the plugin dir:

   ```bash
   export KATAME_PLUGIN_DIR=/path/to/dir   # or set plugin_dir in ~/.katame.yaml
   ```

2. Generate a key pair:

   ```bash
   katame keygen --out demo --algo dummy-kem
   ```

3. Encrypt for the recipient:

   ```bash
   katame protect --in file.txt --recipient demo.pub
   ```

4. Decrypt with the private key:

   ```bash
   katame decrypt --in file.txt.katame --out file.txt --privkey demo.key
   ```

See the [Plugin Author Guide](../../docs/plugins.md) for the full contract and how to implement your own schemes.
