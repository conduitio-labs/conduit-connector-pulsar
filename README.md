# Conduit Connector Apache Pulsar

The Apache Pulsar connector provides both a source and a destination connector for [Apache Pulsar](https://pulsar.apache.org/).

## How to build?
Run `make build` to build the connector.

## Testing
Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests. Tests require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

The Docker compose file at `test/docker-compose.yml` can be used to run the required resource locally.

## Source
A source connector pulls data from an external resource and pushes it to downstream resources via Conduit.

### Configuration

| name                  | description                           | required | default value |
|-----------------------|---------------------------------------|----------|---------------|
| `URL`                 | The service URL for the Pulsar service                            | true     |           |
| `topic`               | Topic specifies the topic the consumer will subscribe on          | true     |           |
| `subscriptionName`    | SubscriptionName specifies the subscription name for the consumer | true     |           |
| `subscriptionType`    | SubscriptionType specifies the subscription type to be used when subscribing to a topic. | false     | exclusive |


## Destination
A destination connector pushes data from upstream resources to an external resource via Conduit.

### Configuration

| name                  | description                           | required | default value |
|-----------------------|---------------------------------------|----------|---------------|
| `URL`             | The service URL for the Pulsar service                            | true     |           |
| `topic`           | Topic specifies the topic the producer will be publishing on      | true     |           |

