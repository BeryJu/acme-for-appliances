# ACME-for-appliances

ACME for appliances that don't natively support it

## Currently supported

- Netapp Ontap (tested with 9.8)
- Citrix ADC/Netscaler (tested with 13.0)

## Running

```
Use ACME Certificates for appliances which don't natively support them.

Usage:
  acme-for-appliances [flags]

Flags:
  -n, --check-interval int   Interval for infinite mode, in hours (default 24)
  -c, --config string        config file
  -f, --force                force renewal
  -h, --help                 help for acme-for-appliances
  -i, --infinite             Infinite mode, keep running the program infinitley and check every interval.
```

## Config

Configuration is loaded from `config.toml` if the file exists. You can also set settings using environment variables:

`A4A_ACME_USER_EMAIL=foo@bar.baz`

This is the default config:

```toml
# Path where private keys are stored, relative to the executable
# Defaults to storage
# storage = "storage"
[acme]
# ACME Directory URL
# directory_url = "https://acme-staging-v02.api.letsencrypt.org/directory"
# Email for your user account, which will receive renewal notices
user_email = "jens@beryju.org"
# Renewal threshold, certificates with this expiry or less will be replaced
# Default: 15 days
# refresh_threshold = 15
# Resolver used for DNS Challenges, defaults to system nameservers
# resolvers = []
# Need to agree to the terms of service
terms_agreed = false
# All providers from "lego" are supported, see https://go-acme.github.io/lego/dns/
# To Configure the provider, consult the page for your provider
challenge_provider_name = "route53"

# ----- Appliance block
[appliances.my-appliance]
# Appliance type, currently supported: netapp
type = "netapp"
# Domains that the certificate should have
domains = [
    "a.int.domain.tld"
]
# General connection details
url = ""  # Base Connection URL
validate_certs = false  # Validate HTTPS certificates
# Authentication is done by the appliance type
username = "admin"
password = "admin"

# ----- Extension of Appliance block
# Some appliance types require extension configuration,
# which is set using [appliances.<name>.extension]
[appliances.my-appliance.extension]

# ----- NetApp Extension of Appliance block
# [appliances.my-appliance.extension]
# Name of the certificate that is created/update
# cert_name = "test-le-cert"
# SVM for which the Certificate is created. Can be set to the cluster name
# for a cluster certificate
# svm_name = "cert-test"

# ----- Citrix ADC Extension of Appliance block
# [appliances.my-appliance.extension]
# Filename of the certificate and key files
# filename_cert = "le_cert.pem"
# filename_key = "le_key.pem"
# Folder in which the files are uploaded
# path_ssl = "/nsconfig/ssl/"
# Name of the certificate which is created/updated
# cert_name = "test-le-cert"
```
