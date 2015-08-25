# Recon
Recon is an open source IT infrastructure & service monitoring tool.

#### Goals
* Ease of setup
* Simple UI
* Intelligent configuration
* Smart analytics and alerting

#### Release

We are still in the early stages of development and will release a production ready build by September 2015. Follow the instructions below to install & try Recon out.

#### Installation

If you have Go installed and workspace setup,

```sh
go get github.com/codeignition/recon/...
```

#### Terminology

[`recond`](https://github.com/codeignition/recon/tree/master/cmd/recond) is the daemon (agent) that runs on your target machine (server).

[`marksman`](https://github.com/codeignition/marksman) is the master server that aggregates the metrics from agents and exposes a public HTTP API.

#### Disclaimer

So far the project is tested only on Linux, specifically Ubuntu 14.04.

#### Contributing

Check out the code and jump to the Issues section to join a conversation or to start one. Also, please read the  [CONTRIBUTING](https://github.com/codeignition/recon/blob/master/CONTRIBUTING.md) document.

#### License

BSD 3-clause "New" or "Revised" license