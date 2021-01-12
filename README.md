<p align="center">
  <a href="https://www.avax.network">
    <img width="25%" alt="snowplow" src="assets/logo.png?raw=true">
  </a>
</p>
<h3 align="center">
  snowplow
</h3>
<p align="center">
quick and easy tool for running and monitoring an <a href="https://www.avax.network">avalanche</a> validator
</p>
<p align="center">
  <a href="https://goreportcard.com/report/github.com/patrick-ogrady/snowplow"><img src="https://goreportcard.com/badge/github.com/patrick-ogrady/snowplow" /></a>
  <a href="https://github.com/patrick-ogrady/snowplow/blob/master/LICENSE"><img src="https://img.shields.io/github/license/patrick-ogrady/snowplow.svg" /></a>
  <a href="https://github.com/patrick-ogrady/snowplow/actions"><img src="https://github.com/patrick-ogrady/snowplow/workflows/go/badge.svg?branch=master" /></a>
  <a href="https://github.com/patrick-ogrady/snowplow/actions"><img src="https://github.com/patrick-ogrady/snowplow/workflows/golangci-lint/badge.svg?branch=master" /></a>
</p>

## Usage
### Install
_must have golang installed_
```text
make install
```

_Creates in .avalanchego/staking (relative directory)_
### Create Staking Credentials
```text
snowplow create
```

### Encrypt + Backup Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
snowplow backup [bucket]
```

### Restore + Decrypt Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
snowplow restore [bucket] [node ID]
```

### Start Node
_TODO: add docker cmd_
```text
make run-mainnet
```

#### Twilio Notifications
_add .avalanchego/.snowplow.yaml_
```yaml
twilio:
  accountSid: "<accountSid>"
  authToken: "<authToken>"
  sender: "<sender phone number>"
  recipient: "<your phone number>"
```

## TODO
* setup new host on google cloud script (install docker, etc)
* cleanup README
* add sha integrity check on backed up files
