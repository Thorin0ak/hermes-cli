[![Maintainability](https://img.shields.io/badge/Go_report-A+-success)](https://goreportcard.com/report/github.com/Thorin0ak/hermes-cli)
[![License: MIT](https://img.shields.io/badge/License-AGPL3.0-blue.svg)](https://opensource.org/licenses/AGPL-3.0)
[![codecov](https://codecov.io/gh/Thorin0ak/hermes-cli/branch/main/graph/badge.svg?token=9420L1M2BD)](https://codecov.io/gh/Thorin0ak/hermes-cli)

# HERMES-CLI

A simple CLI tool to publish events to a [Mercure Hub](https://mercure.rocks/), and optionally [subscribe](https://mercure.rocks/spec#subscription) to receive them.

# ⚡️ Getting started

Download the binaries corresponding to your OS [here](https://github.com/Thorin0ak/hermes-cli/releases).

Provide a Yaml file containing the different Mercure configurations per environment you wish to interact with. The
project comes with a [sample config file](https://github.com/Thorin0ak/hermes-cli/blob/main/sample-config.json):

```yaml
{
  "environments": [
    {
      "name": "localhost",
      "url": "https://localhost/.well-known",
      "jwtSecretKey": "!ChangeMe!"
    }
  ]
}
```

![Recording of terminal demo](./myrecord.svg)

# ⚙️ Commands & Options

| Option | Description                                         | Type     | Default                                    | Required? |
|--------|-----------------------------------------------------|----------|--------------------------------------------|-----------|
| `-h`   | Display usage.                                      |          |                                            | No        |
| `-c`   | Load the Mercure Hub configuration per environment. | `path`   | `.`                                        | No        |
| `-n`   | Number of events to publish to the Mercure Hub.     | `int`    | `5`                                        | No        |
| `-t`   | Topic URI used by Mercure to manage pub/sub.        | `string` | `sse://pxc.dev/123456/test_mercure_events` | No        |
| `-p`   | Only publish events, no client subscription.        | `bool`   | `false`                                    | No        |

# TODO

- [ ] publish binaries to allow installation with `brew`
- [ ] encrypt secrets
- [ ] generate different types of mock payload
- [ ] `sub` should be randomly generated, to ensure no collisions
