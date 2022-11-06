resource "aws_ecr_repository" "stepfunction_ecs_ecr_repo" {
  name                 = "${var.repo_name}"
  
  tags = {
    Name = "${var.app}-ecr"
  }
}