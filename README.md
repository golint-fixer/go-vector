## Vector Go

Official golang implementation of the Vector protocol

          | Linux   | OSX | ARM | Windows | Tests
----------|---------|-----|-----|---------|------
develop   | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=Linux%20Go%20develop%20branch)](https://build.vecdev.com/builders/Linux%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=Linux%20Go%20develop%20branch)](https://build.vecdev.com/builders/OSX%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=ARM%20Go%20develop%20branch)](https://build.vecdev.com/builders/ARM%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=Windows%20Go%20develop%20branch)](https://build.vecdev.com/builders/Windows%20Go%20develop%20branch/builds/-1) | [![Buildr+Status](https://travis-ci.org/vector/go-vector.svg?branch=develop)](https://travis-ci.org/vector/go-vector) [![codecov.io](http://codecov.io/github/vector/go-vector/coverage.svg?branch=develop)](http://codecov.io/github/vector/go-vector?branch=develop)
master    | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=Linux%20Go%20master%20branch)](https://build.vecdev.com/builders/Linux%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=OSX%20Go%20master%20branch)](https://build.vecdev.com/builders/OSX%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=ARM%20Go%20master%20branch)](https://build.vecdev.com/builders/ARM%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.vecdev.com/buildstatusimage?builder=Windows%20Go%20master%20branch)](https://build.vecdev.com/builders/Windows%20Go%20master%20branch/builds/-1) | [![Buildr+Status](https://travis-ci.org/vector/go-vector.svg?branch=master)](https://travis-ci.org/vector/go-vector) [![codecov.io](http://codecov.io/github/vector/go-vector/coverage.svg?branch=master)](http://codecov.io/github/vector/go-vector?branch=master)

[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/vector/go-vector) 
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/vector/go-vector?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Automated development builds

The following builds are build automatically by our build servers after each push to the [develop](https://github.com/vectordev1/go-vector/tree/develop) branch.

* [Docker](https://registry.hub.docker.com/u/vector/client-go/)
* [OS X](http://build.vecdev.com/builds/OSX%20Go%20develop%20branch/Mist-OSX-latest.dmg)
* Ubuntu
  [trusty](https://build.vecdev.com/builds/Linux%20Go%20develop%20deb%20i386-trusty/latest/) |
  [utopic](https://build.vecdev.com/builds/Linux%20Go%20develop%20deb%20i386-utopic/latest/)
* [Windows 64-bit](https://build.vecdev.com/builds/Windows%20Go%20develop%20branch/Gvec-Win64-latest.zip)
* [ARM](https://build.vecdev.com/builds/ARM%20Go%20develop%20branch/gvec-ARM-latest.tar.bz2)

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/vectordev1/go-vector/wiki/Building-Vector)
on the wiki.

Building gvec requires both a Go and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gvec

## Executables

Go Vector comes with several wrappers/executables found in 
[the `cmd` directory](https://github.com/vectordev1/go-vector/tree/develop/cmd):

 Command  |         |
----------|---------|
`gvec` | Vector CLI (vector command line interface client) |
`bootnode` | runs a bootstrap node for the Discovery Protocol |
`vectest` | test tool which runs with the [tests](https://github.com/vector/tests) suite: `/path/to/test.json > vectest --test BlockTests --stdin`.
`evm` | is a generic Vector Virtual Machine: `evm -code 60ff60ff -gas 10000 -price 0 -dump`. See `-h` for a detailed description. |
`disasm` | disassembles EVM code: `echo "6001" | disasm` |
`rlpdump` | prints RLP structures |

## Command line options

`gvec` can be configured via command line options, environment variables and config files.

To get the options available:

    gvec help

For further details on options, see the [wiki](https://github.com/vectordev1/go-vector/wiki/Command-Line-Options)

## Contribution

If you'd like to contribute to go-vector please fork, fix, commit and
send a pull request. Commits who do not comply with the coding standards
are ignored (use gofmt!). If you send pull requests make absolute sure that you
commit on the `develop` branch and that you do not merge to master.
Commits that are directly based on master are simply ignored.

See [Developers' Guide](https://github.com/vectordev1/go-vector/wiki/Developers'-Guide)
for more details on configuring your environment, testing, and
dependency management.
