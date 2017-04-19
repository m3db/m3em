m3em [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]
==============================================================================================

`m3em` (pronounced `meme`) is an acronym for M3DB Environment Manager. Think [CCM](https://github.com/pcmanus/ccm) for M3DB, without the restriction of operating solely on the localhost.

The goal of `m3em` is to make it easy to create, manage and destroy small M3DB clusters across hosts. It is meant for testing a M3DB cluster.
[doc-img]: https://godoc.org/github.com/m3db/m3em?status.svg
[doc]: https://godoc.org/github.com/m3db/m3em
[ci-img]: https://travis-ci.org/m3db/m3em.svg?branch=master
[ci]: https://travis-ci.org/m3db/m3em
[cov-img]: https://coveralls.io/repos/m3db/m3em/badge.svg?branch=master&service=github
[cov]: https://coveralls.io/github/m3db/m3em?branch=master


## Components
There are two primary components in m3em:
(1) `Cluster`: API construct encapsulating placement interactions, along with process orchestration.

(2) `m3em_agent`: process running on remote hosts. It's responsible for process lifecycle, heartbeating back to the coordinating host.

## Usage Example
- TODO(prateek): add usage example

## m3em_agent
```
$ make m3em_agent
$ scp ./out/m3em_agent <remote-host>:<remote-path>
$ ssh <remote-host>
$ cat >m3em.agent.yaml <<EOF
server:
  listenAddress: "0.0.0.0:14541"
  debugAddress: "0.0.0.0:24541"

metrics:
  sampleRate: 0.02
  m3:
    HostPort: "127.0.0.1:9052"
    Service: "m3em"
    IncludeHost: true
    Env: "development"

agent:
  workingDir: /tmp/m3em-agent
  startupCmds:
    - path: /usr/bin/supervisorctl
      args:
        - stop
        - m3dbnode
  releaseCmds:
    - path: /usr/bin/supervisorctl
      args:
        - start
        - m3dbnode
EOF
$ /remote-path/m3em_agent -f m3em.agent.yaml
```