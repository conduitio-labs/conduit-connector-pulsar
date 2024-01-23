# mTLS setup

All `*.pem` files where generated from [this tutorial](https://pulsar.apache.org/docs/3.1.x/security-tls-transport)

`test/setup-tls.sh` are all the commands to generate the keys.


Considerations:

- The generated files may not have the proper file permissions for the docker pulsar service to be able to read them. Use something like `chmod 644 test/*.pem` so solve this.
 
- `standalone.conf`

This is a modified configuration file for the tls setup.

For some reason the pulsar docker service did not pickup the runtime configuration options passed from the environment property. Manually passing the configuration as a flag option circumvents the issue.

You can find the default file contents using this command:

`$ docker exec -it pulsar cat conf/standalone.conf`
