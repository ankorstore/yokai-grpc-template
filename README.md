# Yokai gRPC Template

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)
[![Documentation](https://img.shields.io/badge/Doc-online-cyan)](https://ankorstore.github.io/yokai/)

> gRPC application template based on the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Documentation](#documentation)
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Getting started](#getting-started)
  * [Installation](#installation)
    * [With GitHub](#with-github)
    * [With gonew](#with-gonew)
  * [Usage](#usage)
<!-- TOC -->

## Documentation

For more information about the [Yokai](https://github.com/ankorstore/yokai) framework, you can check its [documentation](https://ankorstore.github.io/yokai).

## Overview

This template provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [gRPC server](https://ankorstore.github.io/yokai/modules/fxgrpcserver/) module installed
- a ready to use [dev environment](docker-compose.yaml), based on [Air](https://github.com/air-verse/air) (for live reloading)
- a ready to use [Dockerfile](Dockerfile) for production
- some examples of [service](internal/service/example.go) and [test](internal/service/example_test.go) to get started

### Layout

This template is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `service/`: gRPC service and test examples
  - `bootstrap.go`: bootstrap
  - `register.go`: dependencies registration
- `proto/`: protobuf definition and stubs

### Makefile

This template provides a [Makefile](Makefile):

```
make up      # start the docker compose stack
make down    # stop the docker compose stack
make logs    # stream the docker compose stack logs
make fresh   # refresh the docker compose stack
make stubs   # generate gRPC stubs with protoc (ex: make stubs from=proto/example.proto)
make test    # run tests
make lint    # run linter
```

## Getting started

### Installation

#### With GitHub

You can create your repository [using the GitHub template](https://github.com/new?template_name=yokai-grpc-template&template_owner=ankorstore).

It will automatically rename your project resources and push them, this operation can take a few minutes.

Once ready, after cloning and going into your repository, simply run:

```shell
make fresh
```

#### With gonew

You can install [gonew](https://go.dev/blog/gonew), and simply run:

```shell
gonew github.com/ankorstore/yokai-grpc-template github.com/foo/bar
cd bar
make fresh
```

### Usage

Once ready, the application will be available on:

- `localhost:50051` for the application gRPC server
- [http://localhost:8081](http://localhost:8081) for the application core dashboard

If you update the [proto definition](proto/example.proto), you can run `make stubs from=proto/example.proto` to regenerate the stubs.

Usage examples with [gRPCurl](https://github.com/fullstorydev/grpcurl):

- with `ExampleService/ExampleUnary`:

```shell
grpcurl -plaintext -d '{"text":"hello"}' localhost:50051 example.ExampleService/ExampleUnary
{
  "text": "response from grpc-app: you sent hello"
}
```

- with `ExampleService/ExampleStreaming`:

```shell
grpcurl -plaintext -d '@' localhost:50051 example.ExampleService/ExampleStreaming <<EOF
{"text":"hello"}
{"text":"world"}
EOF
{
  "text": "response from grpc-app: you sent hello"
}
{
  "text": "response from grpc-app: you sent world"
}
```

You can use any gRPC clients, for example [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).
