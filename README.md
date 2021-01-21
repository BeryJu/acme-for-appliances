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

A minimal config looks like this, for a full example/reference, check out `config-example.toml`.

```toml
[acme]
user_email = "jens@beryju.org"
terms_agreed = false

[appliances.my-appliance]
type = "netapp"
domains = [
    "a.int.domain.tld"
]
url = ""  # Base Connection URL
validate_certs = false  # Validate HTTPS certificates
username = "admin"
password = "admin"

[appliances.my-appliance.extension]
cert_name = "test-le-cert"
svm_name = "cert-test"
```
