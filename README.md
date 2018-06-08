# CFCR-ETCD release


> This project is intended as a sub-component within the [Cloud Foundry Container Runtime](https://docs-cfcr.cfapps.io/) project, usage outside of that are not currently supported by the CFCR team.

## Running the acceptance tests

```
cat > integration_config.json <<EOF
{
	"etcd_endpoint": "http://some-etcd-endpoint:2379"
}
EOF
CONFIG_FILE=$PWD/integration_config.json ginkgo src/acceptance
```
