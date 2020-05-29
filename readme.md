# Notes

# Todo
- [ ] intermittently update default gateway ip
- [ ] super lightweight UI/DB for recording and plotting metrics
- [ ] update deploy actions on github to handle migrations
- [ ] add desktop push notification when pings get terrible
- [ ] bulk upload results
- [ ] update sqlc dependency to support null IPs (track ip from host) 
- [ ] investigate using aggregates for data storage 

# Notes
How to deploy the app to production:
- GO111MODULE=on gcloud app deploy app.yaml

Future work:
- TODO: create a deploy script/workflow that handles migrations

## Readings

### System Design
Schema design for timeseries data with columnar datastore (BigQuery).

https://cloud.google.com/bigtable/docs/schema-design-time-series

### How JWT Signing
- JWT is signed using secret to prevent forgery
- most simple implementation uses secret value + HMAC `HMAC(combine(data, secret))`
- more robust implementation uses public/private keys `PrivateKey(secret).Sign(data)`
    - allows for verification of content without holding the same secret value because the public key can be safely shared between services/systems

Notes:
- keygen: `$ openssl ecparam -genkey -noout -name prime256v1 -out key.pem`
- extract public key: `$ openssl ec -in key.pem -pubout -out public.pem`

  - Elliptical Curve keys have much smaller signatures than RSA. This is optimal for data that is attached to all transmissions.
- [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- [How to provision new key pairs for JWT tokens](https://connect2id.com/products/nimbus-jose-jwt/openssl-key-generation)

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

# Deploying on Balena

`$ balena build --deviceType raspberrypi3 --arch armv7hf --emulated`