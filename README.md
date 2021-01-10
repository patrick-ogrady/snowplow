# Avalanche Runner
a simple repo for creating and managing an avalanche validator

<p align="center">
  <a href="https://goreportcard.com/report/github.com/patrick-ogrady/avalanche-runner"><img src="https://goreportcard.com/badge/github.com/patrick-ogrady/avalanche-runner" /></a>
  <a href="https://github.com/patrick-ogrady/avalanche-runner/blob/master/LICENSE"><img src="https://img.shields.io/github/license/patrick-ogrady/avalanche-runner.svg" /></a>
</p>

## Create Staking Credentials
```text
avalanche-runner create [credential directory]
```

## Backup Credentials
_Make sure to set GOOGLE_APPLICATION_CREDENTIALS_
```text
export GOOGLE_APPLICATION_CREDENTIALS=blah
avalance-runner backup [credential directory] [bucket]
```

## Restore Credentials
```text
avalance-runner restore [bucket] [node ID] [credential directory]
```

## Start Node
_TODO: add docker cmd_
```text
export PAGERDUTY_TOKEN=pagerduty_token
avalance-runner run
```

## TODO
- [ ] license generator
- [ ] create staking key
- [ ] run binary
- [ ] page if stops or unhealthy
- [ ] dockerfile
- [ ] github workflow tester
