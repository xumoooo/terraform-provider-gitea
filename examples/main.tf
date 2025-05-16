
module "test_gitea_fork" {
  source = "./resources/gitea_fork"
}

module "test_gitea_git_hook" {
  count  = var.server_side_hooks_enabled ? 1 : 0
  source = "./resources/gitea_git_hook"
}

module "test_gitea_org" {
  source = "./resources/gitea_org"
}

module "test_gitea_public_key" {
  source = "./resources/gitea_public_key"
}

module "test_gitea_repo_branch_protection" {
  source = "./resources/gitea_repo_branch_protection"
}

module "test_gitea_repository" {
  source = "./resources/gitea_repository"
}

module "test_gitea_repository_key" {
  source = "./resources/gitea_repository_key"
}

module "test_gitea_team" {
  source = "./resources/gitea_team"
}

module "test_gitea_team_members" {
  source = "./resources/gitea_team_members"
}

module "test_gitea_team_membership" {
  source = "./resources/gitea_team_membership"
}

module "test_gitea_token" {
  source = "./resources/gitea_token"
}

module "test_gitea_user" {
  source = "./resources/gitea_user"
}
