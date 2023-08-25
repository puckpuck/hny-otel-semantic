# hny-otel-semantic

Utility to update dataset column descriptions in Honeycomb per OpenTelemetry
semantic conventions. Only columns that do not have a description are updated.
If the dataset column name matches the OpenTelemetry semantic convention
attribute, the description is updated. All datasets in the Honeycomb environment
scoped for the API key are updated.

## Building

```shell
go build -o hny-otel-semantic
```

## Usage

```shell
hny-otel-semantic --honeycomb-api-key <HONEYCOMB_API_KEY> [options]
```

### Options

The following options can be specified on the command line or via environment
variables. The Honeycomb API Key option is required and must be specified on the
or as an environment variable.

| CLI option        | Environment Variable  | Description                                                  | Default   |
|-------------------|-----------------------|--------------------------------------------------------------|-----------|
| dry-run           |                       | Dry run mode                                                 | `false`   |
| honeycomb-api-key | HONEYCOMB_API_KEY     | Honeycomb API Key with permissions to update dataset columns | `nil`     |
| model-path        | SEMANTIC_MODEL_PATH   | Path for OpenTelemetry semantic models                       | `model`   |

## Semantic Models

The semantic models are stored in the `model` directory. The models are copied
from the OpenTelemetry [semantic conventions](https://github.com/open-telemetry/semantic-conventions)
repository. Models last copied from the [4bbb8c9](https://github.com/open-telemetry/semantic-conventions/tree/4bbb8c907402caa90bc077214e8a2c78807c1ab9)
commit.

These models contain the definition of the attributes and their
descriptions.

All `.yaml` files in the path are loaded. The attributes are loaded in the order
they are found. If an attribute is found with the same name as a previous
attribute, the new attribute will overwrite the previous attribute.

Custom attribute definitions can be added to the `model` directory. The
attribute definitions must be in the same format as the OpenTelemetry semantic
conventions.
