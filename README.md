# Yokai gRPC Template

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.22-blue)](https://go.dev/)

> gRPC application template based on the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
* [Documentation](#documentation)
* [Getting started](#getting-started)
  * [Installation](#installation)
    * [With GitHub](#with-github)
    * [With gonew](#with-gonew)
  * [Usage](#usage)
* [Contents](#contents)
  * [Layout](#layout)
  * [Makefile](#makefile)
<!-- TOC -->

## Overview

This template provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [fxgrpcserver](https://github.com/ankorstore/yokai/tree/main/fxgrpcserver) module installed
- a ready to use [dev environment](docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)
- some examples of [service](internal/service/example.go) and [test](internal/service/example_test.go) to get started

## Documentation

See [Yokai documentation](https://ankorstore.github.io/yokai).

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

If you update the [proto definition](proto/example.proto), you can run `make stubs` to regenerate the stubs.

You can use any gRPC clients, for example [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).

## Contents

### Layout

This template is following the [standard Go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `service/`: gRPC service and test examples
  - `bootstrap.go`: bootstrap (modules, lifecycles, etc)
  - `services.go`: dependency injection
- `proto/`: protobuf definition and stubs

### Makefile

This template provides a [Makefile](Makefile):

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc
make test   # run tests
make lint   # run linter
```
