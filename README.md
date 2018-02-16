# CFCR-ETCD release

## Running the acceptance tests

Use `dep ensure` from the `src/acceptance` directory to pull needed dependencies
before running the tests.

```
cat > integration_config.json <<EOF
{
	"etcd_endpoint": "http://some-etcd-endpoint:2379"
}
EOF
CONFIG_FILE=$PWD/integration_config.json ginkgo src/acceptance
```
