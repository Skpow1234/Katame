# Recipes

## Local encrypt / decrypt

```bash
vaultpack protect --in secret.txt
vaultpack decrypt --in secret.txt.vpack --out secret.out --key secret.txt.key
```

## Password mode

```bash
vaultpack protect --in secret.txt --password-file pw.txt --kdf argon2id
vaultpack decrypt --in secret.txt.vpack --out secret.out --password-file pw.txt
```

## Hybrid multi-recipient

```bash
vaultpack keygen --out alice --algo x25519-aes-256-gcm
vaultpack keygen --out bob --algo x25519-aes-256-gcm
vaultpack protect --in secret.txt --recipient alice.pub --recipient bob.pub
vaultpack decrypt --in secret.txt.vpack --out secret.out --privkey bob.key
```

## Key-source URIs

```bash
export VPACK_DATA_KEY="b64:$(base64 -w0 secret.key.raw)"
vaultpack decrypt --in secret.vpack --out secret.out --key-source env://VPACK_DATA_KEY
```

## Transparency signing

```bash
vaultpack sign --in artifact.vpack --signing-priv id.key --transparency
vaultpack verify --in artifact.vpack --pubkey id.pub --check-transparency
```

## Service mode sidecar

```bash
vaultpack serve --listen unix:/run/vaultpack.sock --auth-disabled
# peer process talks HTTP over the unix socket
```

## Tamper-evident repo

```bash
vaultpack repo init --dir ./vp-repo --signing-key repo.key --store-bundles
vaultpack repo add --dir ./vp-repo --bundle secret.vpack --signing-key repo.key --copy
vaultpack repo verify --dir ./vp-repo --verify-key repo.pub
```

## Audit report

```bash
export VAULTPACK_AUDIT_LOG=./audit.jsonl
vaultpack protect --in a.txt
vaultpack audit verify --log ./audit.jsonl
vaultpack report --audit-log ./audit.jsonl --format json --out report.json
```

## Completions & doctor

```bash
vaultpack completion bash > vaultpack.bash
vaultpack doctor
```
