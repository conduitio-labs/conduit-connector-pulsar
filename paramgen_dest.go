// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-connector-sdk/tree/main/cmd/paramgen

package pulsar

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func (DestinationConfig) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		"URL": {
			Default:     "",
			Description: "URL of the Pulsar instance to connect to.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"connectionTimeout": {
			Default:     "",
			Description: "connectionTimeout specifies the duration for which the client will attempt to establish a connection before timing out.",
			Type:        sdk.ParameterTypeDuration,
			Validations: []sdk.Validation{},
		},
		"enableTransaction": {
			Default:     "",
			Description: "enableTransaction determines if the client should support transactions.",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"maxConnectionsPerBroker": {
			Default:     "",
			Description: "maxConnectionsPerBroker limits the number of connections to each broker.",
			Type:        sdk.ParameterTypeInt,
			Validations: []sdk.Validation{},
		},
		"memoryLimitBytes": {
			Default:     "",
			Description: "memoryLimitBytes sets the memory limit for the client in bytes. If the limit is exceeded, the client may start to block or fail operations.",
			Type:        sdk.ParameterTypeInt,
			Validations: []sdk.Validation{},
		},
		"operationTimeout": {
			Default:     "",
			Description: "operationTimeout is the duration after which an operation is considered to have timed out.",
			Type:        sdk.ParameterTypeDuration,
			Validations: []sdk.Validation{},
		},
		"topic": {
			Default:     "",
			Description: "topic specifies the Pulsar topic to which the destination will produce messages.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
	}
}
