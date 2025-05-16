resource "gitea_user" "example" {
  username   = "repo_example_user"
  login_name = "repo_example_user"
  password   = "Geheim1!"
  email      = "repo_example_user@user.dev"
}

resource "gitea_repository" "test" {
  username     = gitea_user.example.username
  name         = "repository-test"
  private      = true
  issue_labels = "Default"
  license      = "MIT"
  gitignores   = "Go"
}

resource "gitea_repository" "mirror" {
  username                     = gitea_user.example.username
  name                         = "terraform-provider-gitea-mirror"
  description                  = "Mirror of Terraform Provider"
  mirror                       = true
  migration_clone_addresse     = "https://git.uploadfilter24.eu/lerentis/terraform-provider-gitea.git"
  migration_service            = "gitea"
  migration_service_auth_token = var.gitea_mirror_token

  depends_on = [gitea_user.example]
}

resource "gitea_repository" "clone" {
  username                     = gitea_user.example.username
  name                         = "terraform-provider-gitea-clone"
  description                  = "Clone of Terraform Provider"
  mirror                       = false
  migration_clone_address      = "https://git.uploadfilter24.eu/lerentis/terraform-provider-gitea.git"
  migration_service            = "gitea"
  migration_service_auth_token = var.gitea_clone_token

  depends_on = [gitea_user.example]
}

variable "gitea_mirror_token" {
  default = "dummy-token"
}

variable "gitea_clone_token" {
  default = "dummy-token"
}