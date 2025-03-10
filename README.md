# Conduit Connector Apache Pulsar

The Apache Pulsar connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [Apache Pulsar](https://pulsar.apache.org/).

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests. Tests require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

The Docker compose file at `test/docker-compose.yml` can be used to run the required resource locally.

## Source Configuration

<!-- readmegen:source.parameters.yaml -->
<!-- /readmegen:source.parameters.yaml -->

## Destination Configuration

<!-- readmegen:destination.parameters.yaml -->
<!-- /readmegen:destination.parameters.yaml -->
