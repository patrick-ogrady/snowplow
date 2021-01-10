# Avalanche Runner
a simple repo for creating and managing an avalanche validator

<p align="center">
  <a href="https://goreportcard.com/report/github.com/patrick-ogrady/avalanche-runner"><img src="https://goreportcard.com/badge/github.com/patrick-ogrady/avalanche-runner" /></a>
  <a href="https://github.com/patrick-ogrady/avalanche-runner/blob/master/LICENSE"><img src="https://img.shields.io/github/license/patrick-ogrady/avalanche-runner.svg" /></a>
</p>

## Create Staking Credentials
_Creates in .avalanchego/staking_
```text
avalanche-runner create
```

## Encrypt + Backup Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
avalanche-runner backup [bucket]
```

## Restore + Decrypt Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
avalanche-runner restore [bucket] [node ID]
```

## Start Node
_TODO: add docker cmd_
```text
export TWILIO_TOKEN=twilio_token
avalanche-runner run
```

## TODO
- [x] license generator
- [x] create staking key
- [x] backup staking key
- [x] hardcode directory name of where keys are generated to be
  .avalanchego/staking
- [ ] dockerfile
- [ ] run binary
- [ ] page if stops or unhealthy
- [ ] github workflow tester
- [ ] add sha integrity check on backed up files
- [ ] generate binaries
