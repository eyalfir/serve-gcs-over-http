# serve-gcs-over-http

Inspired by https://fale.io/blog/2018/04/12/an-http-server-to-serve-gcs-files/.

## Configuration

| environmental variable | optional? | default | description |
|-|-|-|-|
| BUCKET | no || GCS bucket to serve. i.e my-bucket |
| SERVICE_ACCOUNT_KEY | no || path to a service account JSON key |
| LOGGING | yes | no | print a log line for every request |
| PORT | yes | 8000 | Port to serve on |

## Example

```bash

gsutil cp ./my-file.txt gs://my-bucket/my-file.txt
docker run -e BUCKET=my-bucket -v /path/to/service/account/file -e SERVICE_ACCOUNT_KEY=/path/to/service/account/file -e PORT=9080 -p 9080:9080 eyalfirst/serve-gcs-over-http

```

Then

```bash
curl http://localhost:9080/my-file.txt
```
