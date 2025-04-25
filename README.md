# Terraform Provider for Aidbox

This provider allows Terraform to manage Aidbox resources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
```sh
git clone git@github.com:alushanov92/terraform-provider-aidbox.git
```

2. Enter the repository directory
```sh
cd terraform-provider-aidbox
```

3. Build the provider
```sh
go build -o terraform-provider-aidbox
```

## Using the provider

```hcl
terraform {
  required_providers {
    aidbox = {
      source = "alushanov92/aidbox"
    }
  }
}

provider "aidbox" {
  url           = "http://localhost:8080"
  client_id     = "admin"
  client_secret = "password"
}

# Create an Aidbox resource
resource "aidbox_resource" "example" {
  resource_type = "Organization"
  id           = "example-org"
  resource = jsonencode({
    name = "Example Organization"
  })
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is *required*).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

```sh
go build
```

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/alushanov92/aidbox/latest/docs). 