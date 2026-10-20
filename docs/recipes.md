# Recipes

## Local encrypt / decrypt

```bash
katame protect --in secret.txt
katame decrypt --in secret.txt.katame --out secret.out --key secret.txt.key
```

## Password mode

```bash
katame protect --in secret.txt --password-file pw.txt --kdf argon2id
katame decrypt --in secret.txt.katame --out secret.out --password-file pw.txt
```

## Hybrid multi-recipient

```bash
katame keygen --out alice --algo x25519-aes-256-gcm
katame keygen --out bob --algo x25519-aes-256-gcm
katame protect --in secret.txt --recipient alice.pub --recipient bob.pub
katame decrypt --in secret.txt.katame --out secret.out --privkey bob.key
```

## Key-source URIs

```bash
export KATAME_DATA_KEY="b64:$(base64 -w0 secret.key.raw)"
katame decrypt --in secret.katame --out secret.out --key-source env://KATAME_DATA_KEY
```

## Transparency signing

```bash
katame sign --in artifact.katame --signing-priv id.key --transparency
katame verify --in artifact.katame --pubkey id.pub --check-transparency
```

## Service mode sidecar

```bash
katame serve --listen unix:/run/katame.sock --auth-disabled
# peer process talks HTTP over the unix socket
```

## Tamper-evident repo

```bash
katame repo init --dir ./vp-repo --signing-key repo.key --store-bundles
katame repo add --dir ./vp-repo --bundle secret.katame --signing-key repo.key --copy
katame repo verify --dir ./vp-repo --verify-key repo.pub
```

## Audit report

```bash
export KATAME_AUDIT_LOG=./audit.jsonl
katame protect --in a.txt
katame audit verify --log ./audit.jsonl
katame report --audit-log ./audit.jsonl --format json --out report.json
```

## Completions & doctor

```bash
katame completion bash > katame.bash
katame doctor
```
