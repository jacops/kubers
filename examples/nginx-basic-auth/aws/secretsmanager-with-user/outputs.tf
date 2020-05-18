output "nginx_identity_access_key_id" {
  value = aws_iam_access_key.nginx_identity.id
  sensitive = true
}

output "nginx_identity_secret_access_key" {
  value = aws_iam_access_key.nginx_identity.secret
  sensitive = true
}

output "nginx_htpasswd_secret_name" {
  value = aws_secretsmanager_secret.nginx_htpasswd.name
}
