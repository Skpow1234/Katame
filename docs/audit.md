# Audit trail, provenance, and reports

## Audit log

Use `--audit-log <path>` or `VAULTPACK_AUDIT_LOG` to append JSONL for every operation. Entries include timestamp, operation, paths, key fingerprint, user, hostname, success/error, plus a **hash chain** (`prev_hash` / `entry_hash`) and OCSF-style class tags.

```bash
export VAULTPACK_AUDIT_LOG=/var/log/vaultpack-audit.jsonl
vaultpack protect --in data.csv --out data.vpack

vaultpack audit export --format csv --since 2026-01-01 --operation protect
vaultpack audit export --log /var/log/vaultpack-audit.jsonl --format json --key-fingerprint abc123
vaultpack audit verify --log audit.jsonl
```

## Provenance (SLSA-style)

```bash
vaultpack attest --in bundle.vpack --out provenance.json
vaultpack attest --in bundle.vpack --embed   # store provenance.json inside the bundle
```

## Merkle seal

```bash
vaultpack seal --dir ./bundles/ --out root.txt
vaultpack verify-seal --dir ./bundles/ --root $(cat root.txt)
```

## Compliance reports

```bash
vaultpack report --audit-log audit.jsonl --format json --out report.json
vaultpack report --audit-log audit.jsonl --format csv --out report.csv
```

Reports include totals, successes/failures, counts by operation and OCSF class, first/last timestamp, and hash-chain status.
