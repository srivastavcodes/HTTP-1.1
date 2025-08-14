# ServeLite

A lightweight HTTP server built from scratch in Go without using the builtin server.

## Features

- HTTP/1.1 support
- Static file serving
- Basic routing
- Request/response handling

## Installation

```bash
git clone <repository-url>
cd http-server
go build -o server main.go
```

## Usage

Start the server:

```bash
./server
```

The server will start on `localhost:8080` by default.

## Configuration

```bash
go run main.go
```

Run tests:

```bash
go test ./...
```
