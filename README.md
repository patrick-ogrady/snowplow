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

<p align="center"><b>
SNOWPLOW IS CONSIDERED ALPHA SOFTWARE. USE AT YOUR OWN RISK!
</b></p>

## Origins
When setting up my own [avalanche](https://www.avax.network) validator,
I couldn't find any simple tools to backup my validator staking
credentials or send simple text message alerts if the validator went haywire.
So, I made my own snowplow to help tame the avalanche...zing.

## Install
To install `snowplow`, you must first install `golang` (to compile the code).
I don't plan on hosting any pre-built binaries because I think it is important
that the users of this tool compile their own code (for their own safety).

```text
git clone https://github.com/patrick-ogrady/snowplow;
make install;
```

## Usage
_For all of the following operations, `snowplow` assumes your staking
credentials are kept in `.avalanchego/staking` (relative directory). Note, this
is different that `avalanchego` which assumes these credentials are in
`$HOME/.avalanchego/staking`._

### Create Staking Credentials
This command will generate new staking credentials in the
`.avalanchego/staking` folder, if staking credentials do not yet exist in that
folder.

```text
snowplow create
```

### View NodeID
This command will print the NodeID associated with the staking credentials in
the `.avalanchego/staking` folder, if they exist.

```text
snowplow view
```

### Encrypt + Backup Staking Credentials
This command encrypts and backs up your staking credentials to the Google Cloud
Storage bucket of your choosing.

```text
export GOOGLE_APPLICATION_CREDENTIALS=path/to/credentials.json
snowplow backup keys [bucket]
```

_Before running this command, make sure to export your
`GOOGLE_APPLICATION_CREDENTIALS` in your terminal. You can learn more about
Google Cloud's authentication mechanism
[here](https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication)._

### Restore + Decrypt Staking Credentials
This command restores and decrypts the staking credentials of the node of your
choosing from the Google Cloud Storage bucket of your choosing.

```text
export GOOGLE_APPLICATION_CREDENTIALS=path/to/credentials.json
snowplow restore keys [bucket] [node ID]
```

_Before running this command, make sure to export your
`GOOGLE_APPLICATION_CREDENTIALS` in your terminal. You can learn more about
Google Cloud's authentication mechanism
[here](https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication)._

## Google Cloud Deployment
### Setup VM
This sequence of commands sets up an Ubuntu 20.04 LTS
OS on Google Cloud to run an avalanche validator.

```text
git clone https://github.com/patrick-ogrady/snowplow;
cd snowplow;
./scripts/setup.sh;
export PATH=$PATH:~/go/bin;
```

### Build Node
This command builds a Docker image containing `avalanchego` and
the health monitoring mechanism from `snowplow`.

```text
make docker-build
```

_To use Docker on Google Cloud, you may need to prepend `sudo` to this command._

### Start Node
This command starts a Docker container that starts `avalanchego` and
the health monitoring mechanism from `snowplow`.

```text
make run-mainnet
```

_To use Docker on Google Cloud, you may need to prepend `sudo` to this command._

#### Twilio Notifications
To enable text message alerts from [Twilio](https://www.twilio.com/), you must
populate a `yaml` file at `.avalanchego/.snowplow.yaml` (in the same directory
containing your staking keys).

```yaml
twilio:
  accountSid: "<accountSid>"
  authToken: "<authToken>"
  sender: "<sender phone number>"
  recipient: "<your phone number>"
```
