# Key rotation & rewrap

Each rotation clears any prior signature on the manifest (the manifest changes) — re-sign with `katame sign` afterward.

| Command | Payload re-encrypted? | DEK changes? | Use case |
| ------- | --------------------- | ------------ | -------- |
| `rewrap` | no | no | Rotate the **KMS key** wrapping the DEK |
| `rotate-key` | **yes** | yes | DEK compromise or full recipient revocation |
| `add-recipient` | no | no | Wrap existing DEK for more public keys |
| `remove-recipient` | no | no | Drop a recipient entry (combine with `rotate-key` for full revoke) |

```bash
katame rewrap --in data.katame \
  --kms-provider aws \
  --from-kms-key-id arn:aws:kms:us-east-1:111111111111:key/old \
  --to-kms-key-id   arn:aws:kms:us-east-1:111111111111:key/new

katame rotate-key --in data.katame --old-key data.key --new-key-out data.new.key
katame rotate-key --in data.katame --old-password "old" --new-password "new"
katame rotate-key --in data.katame --kms-provider aws
katame rotate-key --in data.katame --old-privkey alice.key --recipient alice.pub --recipient bob.pub

katame add-recipient --in data.katame --privkey alice.key --recipient bob.pub --recipient carol.pub
katame remove-recipient --in data.katame --recipient bob.pub
```

Every rotation appends a `rotated_from` entry:

- `operation` — `rewrap` | `rotate-key` | `add-recipient` | `remove-recipient`
- `rotated_at` — RFC 3339
- `bundle_hash` — SHA-256 of previous bundle bytes
- `notes` — context (old/new KMS IDs, recipient counts, …)

`katame inspect` shows the full chain.
