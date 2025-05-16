terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
  }
}

provider "gitea" {
  base_url = var.gitea_url
  username = var.gitea_username
  password = var.gitea_password
}
