# ACME-for-appliances

ACME for appliances that don't natively support it

## Currently supported

- Netapp ONTAP (tested with 9.8)

  The certificate can be changed for either the entire Cluster's Management interface (set the extension `svm_name` to the cluster name), or a SVM's S3 service. Since the S3 update is disruptive, the SVM will be set to down, the cert is replaced and the SVM is started again. In the case of an error, the SVM is started regardless.

- Citrix ADC/Netscaler (tested with 13.0)

  Works pretty much as expected, certificates are updated without any manual actions or workaround required.
  Depending on the Virtual Server setup, you might have to import the Root CA manually, which for Let's Encrypt will be https://www.identrust.com/dst-root-ca-x3.

- VMware vCenter (tested with 7.0u1)

  After the initial replacement, you might have to accept the new certificate in software that connects to the vCenter, like Veeam.

- Synology DSM (tested with 7.2)

Supported DNS Providers: https://go-acme.github.io/lego/dns/

## Running

You can deploy acme-for-appliances in Kubernetes, where it will run as a CronJob, or with docker-compose. For more info on the Kubernetes deployment, check the k8s folder.

For docker-compose, simply download the docker-compose.yml and the example config file. DNS provider-specific configuration can only be set through environment variables, in the docker-compose file.

## Config

Configuration is loaded from `config.toml` if the file exists.

A minimal config looks like this, for a full example/reference, check out `config-example.toml`.

```toml
[acme]
user_email = "jens@beryju.org"
terms_agreed = false

[appliances.my-appliance]
type = "netapp_ontap"
domains = [
    "a.int.domain.tld"
]
url = ""  # Base Connection URL
validate_certs = false  # Validate HTTPS certificates
username = "admin"
password = "admin"
# Alternatively, you can load username and password from environment variables, like so:
# username = "env:MY_ENV_VAR"

[appliances.my-appliance.extension]
cert_name_a = "test-le-cert-a"
cert_name_b = "test-le-cert-b"
svm_name = "cert-test"
```
