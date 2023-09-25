# hny-otel-semantic

Utility to update dataset column descriptions in Honeycomb per OpenTelemetry
semantic conventions. Only columns that do not have a description are updated.
If the dataset column name matches the OpenTelemetry semantic convention
attribute, the description is updated. All datasets in the Honeycomb environment
scoped for the API key are updated.

## Building

This will build the binary: `hny-otel-semantic` 
```shell
make build
```

## Usage

```shell
hny-otel-semantic --honeycomb-api-key <HONEYCOMB_API_KEY> [options]
```

### Options

The following options can be specified on the command line or via environment
variables. The Honeycomb API Key option is required and must be specified on the
or as an environment variable.

| CLI option              | Environment Variable  | Description                                                    | Default   |
|-------------------------|-----------------------|----------------------------------------------------------------|-----------|
| --honeycomb-api-key     | HONEYCOMB_API_KEY     | Honeycomb API Key with permissions to update dataset columns   | `nil`     |
| --model-path            | SEMANTIC_MODEL_PATH   | Path for OpenTelemetry semantic models                         | `model`   |
| --all-metric-attributes |                       | Include all metric attributes, included ones without a prefix  | `false`   |
| --dry-run               |                       | Dry run mode                                                   | `false`   |
| --parse-models-only     |                       | Only parse the semantic models and display details             | `false`   |
| --version               |                       | Display version information                                    | `false`   |

## Semantic Models

The semantic models are stored in the `model` directory. The models are copied
from the OpenTelemetry [semantic conventions](https://github.com/open-telemetry/semantic-conventions)
repository. Models were last copied from the [984079e](https://github.com/open-telemetry/semantic-conventions/commit/984079ee2b98a5700b139989db9737b044ab40e6)
commit.

To sync with the latest models run:

```shell
make sync-models
```

These models contain the definition of the attributes and their
descriptions.

All `.yaml` files in the path are loaded. The attributes are loaded in the order
they are found. If an attribute is found with the same name as a previous
attribute, the new attribute will overwrite the previous attribute.

Custom attribute definitions can be added to the `model` directory. The
attribute definitions must be in the same format as the OpenTelemetry semantic
conventions.
