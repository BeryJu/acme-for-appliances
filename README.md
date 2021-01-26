# ACME-for-appliances

ACME for appliances that don't natively support it

## Currently supported

- Netapp ONTAP (tested with 9.8)

  Notes: Currently, the certificate is deleted and re-created. This will fail when the certificate is in use by HTTPS services. The API does not expose where the Certificate is used. I'm currently re-writing the Netapp integration to create certificates with a counter, then modify common HTTPS settings to use the new certificate and attempt to delete the old one.

- Citrix ADC/Netscaler (tested with 13.0)


- VMware vCenter (tested with 7.0u1)

  Notes: After the initial replacement, you might have to accept the new certificate in software that connects to the vCenter, like Veeam.

Supported DNS Providers: https://go-acme.github.io/lego/dns/

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
