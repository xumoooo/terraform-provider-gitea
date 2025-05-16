resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "gitea_user" "example" {
  username   = "repo_key_example_user"
  login_name = "repo_key_example_user"
  password   = "Geheim1!"
  email      = "repo_key_example_user@user.dev"
}

resource "gitea_repository" "example" {
  name     = "example"
  username = gitea_user.example.username
  private  = true
}

resource "gitea_repository_key" "example" {
  repository = gitea_repository.example.id
  title      = "Example Deploy Key"
  read_only  = true
  key        = tls_private_key.example.public_key_openssh
}