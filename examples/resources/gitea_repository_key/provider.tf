terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.0.4"
    }
  }
}
