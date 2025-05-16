resource "gitea_user" "example" {
  username   = "bp_example_user"
  login_name = "bp_example_user"
  password   = "Geheim1!"
  email      = "bp_example_user@user.dev"
}

resource "gitea_repository" "repo" {
  username  = gitea_user.example.username
  name      = "bp-branch-protection-example-repo"
  auto_init = false
}

resource "gitea_repository_branch_protection" "main" {
  username = gitea_repository.repo.username
  name     = gitea_repository.repo.name

  rule_name   = "main"
  enable_push = false
}
