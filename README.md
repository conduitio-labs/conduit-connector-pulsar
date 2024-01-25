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

| name | description | required | default value |
| ---- | ---- | ---- | ---- |
| `url` | url of the Pulsar instance to connect to. | true |  |
| `topic` | topic specifies the Pulsar topic from which the source will consume messages. | true |  |
| `connectionTimeout` | connectionTimeout specifies the duration for which the client will attempt to establish a connection before timing out. | false |  |
| `enableTransaction` | enableTransaction determines if the client should support transactions. | false |  |
| `maxConnectionsPerBroker` | maxConnectionsPerBroker limits the number of connections to each broker. | false |  |
| `memoryLimitBytes` | memoryLimitBytes sets the memory limit for the client in bytes. If the limit is exceeded, the client may start to block or fail operations. | false |  |
| `operationTimeout` | operationTimeout is the duration after which an operation is considered to have timed out. | false |  |
| `subscriptionName` | subscriptionName is the name of the subscription to be used for consuming messages. If none provided, a random uuid will be created as the name. | false |  |
| `subscriptionType` | subscriptionType defines the type of subscription to use. Can be "exclusive", "shared", "failover", "key_shared". Default is "exclusive". | false | exclusive |
| `tlsKeyFilePath` | tlsKeyFilePath sets the path to the TLS key file | false |  |
| `tlsCertificateFile` | tlsCertificateFile sets the path to the TLS certificate file | false |  |
| `tlsTrustCertsFilePath` |tlsTrustCertsFilePath sets the path to the trusted TLS certificate file | false |  |
| `tlsAllowInsecureConnection` | tlsAllowInsecureConnection configures whether the internal Pulsar client accepts untrusted TLS certificate from broker (default: false) | false  |  |
| `tlsValidateHostname` | tlsValidateHostname configures whether the Pulsar client verifies the validity of the host name from broker (default: false) | false |  |


## Destination
A destination connector pushes data from upstream resources to an external resource via Conduit.

### Configuration

| name                     | description                                                                                                                           | required | default value |
|--------------------------|---------------------------------------------------------------------------------------------------------------------------------------|----------|---------------|
| `url`                    | URL of the Pulsar instance to connect to.                                                                                             | true     |               |
| `topic`                  | topic specifies the Pulsar topic to which the destination will produce messages.                                                      | true     |               |
| `connectionTimeout`      | connectionTimeout specifies the duration for which the client will attempt to establish a connection before timing out.              | false    |               |
| `enableTransaction`      | enableTransaction determines if the client should support transactions.                                                               | false    |               |
| `maxConnectionsPerBroker`| maxConnectionsPerBroker limits the number of connections to each broker.                                                              | false    |               |
| `memoryLimitBytes`       | memoryLimitBytes sets the memory limit for the client in bytes. If the limit is exceeded, the client may start to block or fail operations. | false    |               |
| `operationTimeout`       | operationTimeout is the duration after which an operation is considered to have timed out.                                            | false    |               |
| `tlsKeyFilePath` | tlsKeyFilePath sets the path to the TLS key file | false |  |
| `tlsCertificateFile` | tlsCertificateFile sets the path to the TLS certificate file | false |  |
| `tlsTrustCertsFilePath` |tlsTrustCertsFilePath sets the path to the trusted TLS certificate file | false |  |
| `tlsAllowInsecureConnection` | tlsAllowInsecureConnection configures whether the internal Pulsar client accepts untrusted TLS certificate from broker (default: false) | false  |  |
| `tlsValidateHostname` | tlsValidateHostname configures whether the Pulsar client verifies the validity of the host name from broker (default: false) | false |  |
