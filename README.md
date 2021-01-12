# Avalanche Runner
_quick and easy tool for running an avalanche node responsibly_

<p align="center">
  <a href="https://goreportcard.com/report/github.com/patrick-ogrady/avalanche-runner"><img src="https://goreportcard.com/badge/github.com/patrick-ogrady/avalanche-runner" /></a>
  <a href="https://github.com/patrick-ogrady/avalanche-runner/blob/master/LICENSE"><img src="https://img.shields.io/github/license/patrick-ogrady/avalanche-runner.svg" /></a>
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
avalanche-runner create
```

### Encrypt + Backup Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
avalanche-runner backup [bucket]
```

### Restore + Decrypt Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
avalanche-runner restore [bucket] [node ID]
```

### Start Node
_TODO: add docker cmd_
```text
make run-mainnet
```

#### Twilio Notifications
_add .avalanchego/.avalanche-runner.yaml_
```yaml
twilio:
  accountSid: "<accountSid>"
  authToken: "<authToken>"
  sender: "<sender phone number>"
  recipient: "<your phone number>"
```

## TODO
- [x] license generator
- [x] create staking key
- [x] backup staking key
- [x] hardcode directory name of where keys are generated to be
  .avalanchego/staking
- [x] dockerfile
- [x] run binary
- [ ] reorganize to pkg
- [ ] write tests for health with interfaces + mockery
- [ ] send message when node started or node shutdown
- [ ] Config file with phone number + twilio tokens
- [ ] page if stops or unhealthy (only once bootstrapped has gone true)
- [ ] fix liveness check json
```
[NOTIFIER] received error while checking liveness: json: cannot unmarshal object into Go struct field Result.checks.error of type error
```
```
{
    "jsonrpc": "2.0",
    "result": {
        "checks": {
            "P": {
                "message": {
                    "percentConnected": 1
                },
                "timestamp": "2021-01-11T17:57:24.3617665Z",
                "duration": 1646234500,
                "contiguousFailures": 0,
                "timeOfFirstFailure": null
            },
            "chains.default.bootstrapped": {
                "error": {
                    "message": "P-Chain not bootstrapped"
                },
                "timestamp": "2021-01-11T17:57:29.8214535Z",
                "duration": 8000,
                "contiguousFailures": 8,
                "timeOfFirstFailure": "2021-01-11T17:51:20.233096Z"
            },
            "network.validators.heartbeat": {
                "message": {
                    "heartbeat": 1610387849
                },
                "timestamp": "2021-01-11T17:57:29.8214987Z",
                "duration": 26300,
                "contiguousFailures": 0,
                "timeOfFirstFailure": null
            }
        },
        "healthy": false
    },
    "id": 1
}
```
- [ ] github workflow tester
- [ ] add sha integrity check on backed up files
- [ ] generate and host cli binaries
- [ ] setup new host script (install docker, etc)
