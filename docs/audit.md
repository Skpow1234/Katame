# Audit trail, provenance, and reports

## Audit log

Use `--audit-log <path>` or `KATAME_AUDIT_LOG` to append JSONL for every operation. Entries include timestamp, operation, paths, key fingerprint, user, hostname, success/error, plus a **hash chain** (`prev_hash` / `entry_hash`) and OCSF-style class tags.

```bash
export KATAME_AUDIT_LOG=/var/log/katame-audit.jsonl
katame protect --in data.csv --out data.katame

katame audit export --format csv --since 2026-01-01 --operation protect
katame audit export --log /var/log/katame-audit.jsonl --format json --key-fingerprint abc123
katame audit verify --log audit.jsonl
```

## Provenance (SLSA-style)

```bash
katame attest --in bundle.katame --out provenance.json
katame attest --in bundle.katame --embed   # store provenance.json inside the bundle
```

## Merkle seal

```bash
katame seal --dir ./bundles/ --out root.txt
katame verify-seal --dir ./bundles/ --root $(cat root.txt)
```

## Compliance reports

```bash
katame report --audit-log audit.jsonl --format json --out report.json
katame report --audit-log audit.jsonl --format csv --out report.csv
```

Reports include totals, successes/failures, counts by operation and OCSF class, first/last timestamp, and hash-chain status.
