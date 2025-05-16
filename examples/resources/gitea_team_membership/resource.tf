resource "gitea_org" "example_org" {
  name = "membership_example_org"
}

resource "gitea_user" "example_users" {
  count      = 5
  username   = "membership_example_user_${count.index}"
  login_name = "membership_example_user_${count.index}"
  password   = "Geheim1!"
  email      = "membership_example_user_${count.index}@user.dev"
}

resource "gitea_team" "example_team" {
  name         = "membership_example_team"
  organisation = gitea_org.example_org.name
  description  = "An example team for membership testing"
  permission   = "read"
}

resource "gitea_team_membership" "example_team_memberships" {
  for_each = { for user in gitea_user.example_users : user.username => user }
  team_id  = gitea_team.example_team.id
  username = each.value["username"]
}
