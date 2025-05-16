resource "gitea_user" "test" {
  username             = "public-key-test"
  login_name           = "public-key-test"
  password             = "Geheim1!"
  email                = "public-key-test@user.dev"
  must_change_password = false
}

resource "gitea_public_key" "test_user_key" {
  title    = "public-key-test"
  key      = file("${path.module}/id_ed25519.pub")
  username = gitea_user.test.username
}
