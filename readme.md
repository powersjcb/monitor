# Notes

# Todo
- [ ] execute ping requests using udp
- [ ] intermittently update default gateway ip
- [ ] super lightweight UI/DB for recording and plotting metrics
- [ ] create a deploy script that handle migrations
- [ ] add desktop push notification when pings get terrible
- [ ] bulk upload results
- [ ] update sqlc dependency to support null IPs (track ip from host) 

# Deploy
How to deploy the app to production:
- GO111MODULE=on gcloud app deploy app.yaml

Future work:
- TODO: create a deploy script/workflow that handles migrations

## Readings
Schema design for timeseries data with columnar datastore (BigQuery).

https://cloud.google.com/bigtable/docs/schema-design-time-series

## Migrations
Code generation for database queries requires `sqlc`
- `$ brew install kyleconroy/sqlc/sqlc`
- `$ sqlc generate`

## Installation

Setup
- `$ brew install postgresql`
- `$ brew services restart postgresql`
- `$ brew install vektra/tap/mockery`
- `$ brew upgrade mockery`

