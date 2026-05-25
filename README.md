# UUID - Protobuf UUID Type

Wraps [github.com/google/uuid](https://github.com/google/uuid) as a [protobuf](https://developers.google.com/protocol-buffers/) type with implementations of JSON, BSON, and GraphQL serialization.

## Buf Schema Registry

Published on the [Buf Schema Registry](https://buf.build/jpbeatskitano/uuid).

## Usage

### Add the Buf dependency

Add this module to your project's `buf.yaml`:

```yaml
version: v2
deps:
    - buf.build/jpbeatskitano/uuid
```

Then update your Buf lock file:

```sh
buf dep update
```

### Protobuf

The Buf module name is `buf.build/jpbeatskitano/uuid`, but the proto import path is the file path inside that module: `uuid.proto`.

```protobuf
syntax = "proto3";

package models;

import "google/protobuf/timestamp.proto";
import "models/role.proto";
import "uuid.proto";

option go_package = "kefu/proto/generated-go/models";

message User {
    uuid.UUID id = 1;
    google.protobuf.Timestamp created_at = 2;
    google.protobuf.Timestamp updated_at = 3;
}
```

Do not import the module name as a proto path:

```protobuf
// Wrong
import "buf.build/jpbeatskitano/uuid.proto";
import "buf.build/jpbeatskitano/uuid/uuid.proto";
```

### Go

```go
import "github.com/BeatsKitano/uuid"

id := uuid.New()
```

## Publishing

This repository is configured as the Buf module `buf.build/jpbeatskitano/uuid` in `buf.yaml`.

```sh
buf lint
buf breaking --against '.git#branch=main'
buf push
```
