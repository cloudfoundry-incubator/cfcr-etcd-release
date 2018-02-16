# CFCR-ETCD release

## Running the acceptance tests

```
cat > integration_config.json <<EOF
{
	"etcd_endpoint": "http://some-etcd-endpoint:2379"
}
EOF
CONFIG_FILE=$PWD/integration_config.json ginkgo src/acceptance
```
