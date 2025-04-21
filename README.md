# Conduit Connector Apache Pulsar

The Apache Pulsar connector is one of [Conduit](https://conduit.io) plugins. The connector provides both a source and a destination connector for [Apache Pulsar](https://pulsar.apache.org/).

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests. Tests require Docker to be installed and running. The command will handle starting and stopping docker containers for you.

The Docker compose file at `test/docker-compose.yml` can be used to run the required resource locally.

## Source Configuration

<!-- readmegen:source.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "pulsar"
        settings:
          # The Pulsar topic used by the connector.
          # Type: string
          # Required: yes
          topic: ""
          # The pulsar instance to connect to.
          # Type: string
          # Required: yes
          url: ""
          # The duration for which the client will attempt to establish a
          # connection before timing out.
          # Type: duration
          # Required: no
          connectionTimeout: "0s"
          # Disables pulsar client logs
          # Type: bool
          # Required: no
          disableLogging: "false"
          # Determines if the client should support transactions.
          # Type: bool
          # Required: no
          enableTransaction: "false"
          # Limits the number of connections to each broker.
          # Type: int
          # Required: no
          maxConnectionsPerBroker: "0"
          # Sets the memory limit for the client in bytes. If the limit is
          # exceeded, the client may start to block or fail operations.
          # Type: int
          # Required: no
          memoryLimitBytes: "0"
          # The duration after which an operation is considered to have timed
          # out.
          # Type: duration
          # Required: no
          operationTimeout: "0s"
          # The name of the subscription to be used for consuming messages.
          # Type: string
          # Required: no
          subscriptionName: ""
          # Whether the internal Pulsar client accepts untrusted TLS certificate
          # from broker (default: false)
          # Type: bool
          # Required: no
          tlsAllowInsecureConnection: "false"
          # The path to the TLS certificate file
          # Type: string
          # Required: no
          tlsCertificateFile: ""
          # The path to the TLS key file
          # Type: string
          # Required: no
          tlsKeyFilePath: ""
          # The path to the trusted TLS certificate file
          # Type: string
          # Required: no
          tlsTrustCertsFilePath: ""
          # Whether the pulsar client verifies the validity of the host name
          # from broker (default: false)
          # Type: bool
          # Required: no
          tlsValidateHostname: "false"
          # Maximum delay before an incomplete batch is read from the source.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets read from the source.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Specifies whether to use a schema context name. If set to false, no
          # schema context name will be used, and schemas will be saved with the
          # subject name specified in the connector (not safe because of name
          # conflicts).
          # Type: bool
          # Required: no
          sdk.schema.context.enabled: "true"
          # Schema context name to be used. Used as a prefix for all schema
          # subject names. If empty, defaults to the connector ID.
          # Type: string
          # Required: no
          sdk.schema.context.name: ""
          # Whether to extract and encode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # The subject of the key schema. If the record metadata contains the
          # field "opencdc.collection" it is prepended to the subject name and
          # separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.key.subject: "key"
          # Whether to extract and encode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
          # The subject of the payload schema. If the record metadata contains
          # the field "opencdc.collection" it is prepended to the subject name
          # and separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.payload.subject: "payload"
          # The type of the payload schema.
          # Type: string
          # Required: no
          sdk.schema.extract.type: "avro"
```
<!-- /readmegen:source.parameters.yaml -->

## Destination Configuration

<!-- readmegen:destination.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "pulsar"
        settings:
          # The Pulsar topic used by the connector.
          # Type: string
          # Required: yes
          topic: ""
          # The pulsar instance to connect to.
          # Type: string
          # Required: yes
          url: ""
          # The duration for which the client will attempt to establish a
          # connection before timing out.
          # Type: duration
          # Required: no
          connectionTimeout: "0s"
          # Disables pulsar client logs
          # Type: bool
          # Required: no
          disableLogging: "false"
          # Determines if the client should support transactions.
          # Type: bool
          # Required: no
          enableTransaction: "false"
          # Limits the number of connections to each broker.
          # Type: int
          # Required: no
          maxConnectionsPerBroker: "0"
          # Sets the memory limit for the client in bytes. If the limit is
          # exceeded, the client may start to block or fail operations.
          # Type: int
          # Required: no
          memoryLimitBytes: "0"
          # The duration after which an operation is considered to have timed
          # out.
          # Type: duration
          # Required: no
          operationTimeout: "0s"
          # Whether the internal Pulsar client accepts untrusted TLS certificate
          # from broker (default: false)
          # Type: bool
          # Required: no
          tlsAllowInsecureConnection: "false"
          # The path to the TLS certificate file
          # Type: string
          # Required: no
          tlsCertificateFile: ""
          # The path to the TLS key file
          # Type: string
          # Required: no
          tlsKeyFilePath: ""
          # The path to the trusted TLS certificate file
          # Type: string
          # Required: no
          tlsTrustCertsFilePath: ""
          # Whether the pulsar client verifies the validity of the host name
          # from broker (default: false)
          # Type: bool
          # Required: no
          tlsValidateHostname: "false"
          # Maximum delay before an incomplete batch is written to the
          # destination.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets written to the destination.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Allow bursts of at most X records (0 or less means that bursts are
          # not limited). Only takes effect if a rate limit per second is set.
          # Note that if `sdk.batch.size` is bigger than `sdk.rate.burst`, the
          # effective batch size will be equal to `sdk.rate.burst`.
          # Type: int
          # Required: no
          sdk.rate.burst: "0"
          # Maximum number of records written per second (0 means no rate
          # limit).
          # Type: float
          # Required: no
          sdk.rate.perSecond: "0"
          # The format of the output record. See the Conduit documentation for a
          # full list of supported formats
          # (https://conduit.io/docs/using/connectors/configuration-parameters/output-format).
          # Type: string
          # Required: no
          sdk.record.format: "opencdc/json"
          # Options to configure the chosen output record format. Options are
          # normally key=value pairs separated with comma (e.g.
          # opt1=val2,opt2=val2), except for the `template` record format, where
          # options are a Go template.
          # Type: string
          # Required: no
          sdk.record.format.options: ""
          # Whether to extract and decode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # Whether to extract and decode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
```
<!-- /readmegen:destination.parameters.yaml -->
