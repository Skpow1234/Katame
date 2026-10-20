# Cloud storage

Katame accepts cloud URIs anywhere a local path is allowed for `--in` / `--out` (and `--dir` / `--out-dir` for batch). Objects are downloaded to a private temp file, processed locally, then re-uploaded — crypto always runs in-process on local data.

| Scheme | Backend | Read | Write | Credentials |
| ------ | ------- | ---- | ----- | ----------- |
| `az://` | Azure Blob | yes | yes | Connection string → `DefaultAzureCredential` |
| `s3://` | AWS S3 | yes | yes | AWS SDK v2 default chain |
| `gs://` | GCS | yes | yes | Application Default Credentials |
| `https://`, `http://` | HTTP | yes | **no** | Optional `KATAME_HTTP_BEARER` or basic auth env |

```bash
katame protect --in az://container/in.csv --out az://container/out.katame --azure-account myaccount
katame protect --in s3://mybucket/in.csv --out s3://mybucket/out.katame --aws-region us-east-1
katame protect --in gs://mybucket/in.csv --out gs://mybucket/out.katame
katame inspect --in https://example.com/release.katame
katame batch-protect --dir s3://mybucket/exports/ --out-dir s3://mybucket/encrypted/ --aws-region us-east-1
```

## URI formats

- `az://container/path/to/blob`
- `s3://bucket/path/to/object`
- `gs://bucket/path/to/object`
- `https://host/path` / `http://host/path`

## Flags

| Flag | Description |
| ---- | ----------- |
| `--azure-account` | Storage account (or `AZURE_STORAGE_ACCOUNT`) |
| `--azure-connection-string` | Connection string (or `AZURE_STORAGE_CONNECTION_STRING`) |
| `--aws-region` | Region for `s3://` (or `AWS_REGION`) |
| `--aws-profile` | Shared config profile (or `AWS_PROFILE`) |
| `--s3-endpoint` | Override endpoint (MinIO, LocalStack) |
| `--s3-path-style` | Path-style addressing |

GCS uses ADC only (no extra flags beyond standard Google env / `gcloud auth application-default login`).
