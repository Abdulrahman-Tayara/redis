# Redis Clone (Golang)

A simple Redis clone implemented in Golang.

## Features

- Supports basic commands.
- Easy to implement any new command.
- Supports two different types of key expiration: passive and random active expiration.
- Tested with the popular Golang Redis client [`go-redis`](https://github.com/redis/go-redis).
- Configuration options via JSON or YAML.
- Supports manual execution and Docker deployment.

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (optional)

### Installation

Clone the repository:

```sh
git clone https://github.com/Abdulrahman-Tayara/redis.git
cd redis
```

Build the project:

```sh
go build -o cmd/server/main.go
```

## Usage

### Running Manually

#### Using Default Configuration

Run the server with the default settings:

```sh
./main
```

#### Using a Custom Configuration File

Specify a JSON or YAML config file:

```sh
./main -config config.json
```

or

```sh
./main -config config.yaml
```

### Running with Docker

Build the Docker image:

```sh
docker build -t redis-clone .
```

Run the container:

```sh
docker run -p 6379:6379 redis-clone
```

To use a custom config file, mount it as a volume:

```sh
docker run -p 6379:6379 -v $(pwd)/config.json:/app/config.json redis-clone --config /app/config.json
```

## Configuration

The server supports configuration through JSON or YAML files. Example configurations:

### JSON (`config.json`)

```json
{
  "port": "6379",
  "version": "6.0.3",
  "proto_version": 3,
  "mode": "standalone" 
}
```

### YAML (`config.yaml`)

```yaml
port: "6379"
version: "6.0.3"
proto_version: 3
mode: "standalone"
```

## Testing

You can test the server using the `go-redis` client:

```go
package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key:", val) // Output: key: value
}
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License.
