### README.md

# Conduit Connector Apache Pulsar

The Apache Pulsar connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [Apache Pulsar](https://pulsar.apache.org/).

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests. Tests require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

The Docker compose file at `test/docker-compose.yml` can be used to run the required resource locally.

## Shared Configuration

The following table lists configuration options common to both source and destination connectors.

| name                         | description                                                                                                                                 | required | default value |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `url`                        | URL of the Pulsar instance to connect to.                                                                                                   | true     |               |
| `topic`                      | Topic specifies the Pulsar topic to which the source / destination will interact with.                                                      | true     |               |
| `connectionTimeout`          | ConnectionTimeout specifies the duration for which the client will attempt to establish a connection before timing out.                     | false    |               |
| `operationTimeout`           | OperationTimeout is the duration after which an operation is considered to have timed out.                                                  | false    |               |
| `maxConnectionsPerBroker`    | MaxConnectionsPerBroker limits the number of connections to each broker.                                                                    | false    |               |
| `memoryLimitBytes`           | MemoryLimitBytes sets the memory limit for the client in bytes. If the limit is exceeded, the client may start to block or fail operations. | false    |               |
| `enableTransaction`          | EnableTransaction determines if the client should support transactions.                                                                     | false    |               |
| `tlsKeyFilePath`             | TLSKeyFilePath sets the path to the TLS key file                                                                                            | false    |               |
| `tlsCertificateFile`         | TLSCertificateFile sets the path to the TLS certificate file                                                                                | false    |               |
| `tlsTrustCertsFilePath`      | TLSTrustCertsFilePath sets the path to the trusted TLS certificate file                                                                     | false    |               |
| `tlsAllowInsecureConnection` | TLSAllowInsecureConnection configures whether the internal Pulsar client accepts untrusted TLS certificate from broker (default: false)     | false    |               |
| `tlsValidateHostname`        | TLSValidateHostname configures whether the Pulsar client verifies the validity of the host name from broker (default: false)                | false    |               |

## Destination Configuration

The destination connector uses only the shared configuration options listed above.

## Source Configuration

Additional to the shared configuration, the source connector has the following configurations.

| name               | description                                                                                                                                      | required | default value |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | ------------- |
| `subscriptionName` | SubscriptionName is the name of the subscription to be used for consuming messages. If none provided, a random uuid will be created as the name. | false    |               |
| `subscriptionType` | SubscriptionType defines the type of subscription to use. Can be "exclusive", "shared", "failover", "key_shared". Default is "exclusive".        | false    | exclusive     |

## Example pipeline.yml

Example of a [pipeline.yml](https://conduit.io/docs/pipeline-configuration-files/getting-started) file using `file to apache pulsar` and `apache pulsar to file` pipelines:

```yaml
version: 2.0
pipelines:
  - id: file-to-pulsar
    status: running
    connectors:
      - id: file.in
        type: source
        plugin: builtin:file
        name: file-destination
        settings:
          path: ./file.in
      - id: pulsar.out
        type: destination
        plugin: standalone:pulsar
        name: pulsar-source
        settings:
          url: pulsar://localhost:6650
          topic: demo-topic
          sdk.record.format: template
          sdk.record.format.options: '{{ printf "%s" .Payload.After }}'

  - id: pulsar-to-file
    status: running
    connectors:
      - id: pulsar.in
        type: source
        plugin: standalone:pulsar
        name: pulsar-source
        settings:
          url: pulsar://localhost:6650
          topic: demo-topic
          subscriptionName: demo-topic-subscription

      - id: file.out
        type: destination
        plugin: builtin:file
        name: file-destination
        settings:
          path: ./file.out
          sdk.record.format: template
          sdk.record.format.options: '{{ printf "%s" .Payload.After }}'
```
