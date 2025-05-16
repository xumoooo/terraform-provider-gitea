resource "gitea_org" "example" {
  name = "git-hook-example-org"
}

resource "gitea_repository" "example" {
  username = gitea_org.example.name
  name     = "git-hook-example-repo"
}

resource "gitea_git_hook" "org_repo_post_receive" {
  name    = "post-receive"
  user    = gitea_org.example.name
  repo    = gitea_repository.example.name
  content = file("${path.module}/post-receive.sh")
}