terraform {
  required_version = ">= 0.12.2"

  backend "local" {}
}

provider "aws" {
  version = "~> 2.56"
  region = "eu-west-2"
}

resource "aws_secretsmanager_secret" "nginx_htpasswd" {
    name = "nginx-htpasswd"
}

resource "aws_secretsmanager_secret_version" "nginx_htpasswd" {
  secret_id     = aws_secretsmanager_secret.nginx_htpasswd.id
  secret_string = "test:$2y$10$i1yZDJmNSJkipWW1JDSp6uIpgt12oCdsM5Q8FZcntOBLerl5dFQ56" #user test:test
}

data "aws_iam_policy_document" "read_nginx_htpasswd_secret" {
  statement {
    actions = [
      "secretsmanager:GetResourcePolicy",
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret",
      "secretsmanager:ListSecretVersionIds"
    ]

    resources = [
        aws_secretsmanager_secret.nginx_htpasswd.arn
    ]
  }
}

resource "aws_iam_policy" "read_nginx_htpasswd_secret" {
  name   = "read_nginx_htpasswd_secret"
  path   = "/kubers/"
  policy = data.aws_iam_policy_document.read_nginx_htpasswd_secret.json
}

resource "aws_iam_user" "nginx_identity" {
  name = "kubernetes-nginx"
  path  = "/kubers/"
}

resource "aws_iam_user_policy_attachment" "nginx_identity_can_read_nginx_htpasswd_secret" {
  user       = aws_iam_user.nginx_identity.name
  policy_arn = aws_iam_policy.read_nginx_htpasswd_secret.arn
}

resource "aws_iam_access_key" "nginx_identity" {
  user = aws_iam_user.nginx_identity.name
}
