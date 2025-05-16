resource "gitea_org" "example" {
  name = "org-example-org"
}

resource "gitea_repository" "example" {
  username = gitea_org.example.name
  name     = "org-example-repo"
}