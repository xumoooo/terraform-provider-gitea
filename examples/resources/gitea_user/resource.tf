resource "gitea_user" "test" {
  username             = "user-test"
  login_name           = "user-test"
  password             = "Geheim1!"
  email                = "user-test@user.dev"
  must_change_password = false
}