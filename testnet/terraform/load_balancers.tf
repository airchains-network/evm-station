# RPC load balancer
resource "aws_lb" "rpc_alb" {
  name               = "${var.environment}-full-node-rpc-elb"
  internal           = false
  load_balancer_type = "application"

  subnets = [aws_subnet.private_subnet.id, aws_subnet.public_subnet.id]
  security_groups = [aws_security_group.alb_sg.id]

  tags = {
    Name = "${var.environment} Full Nodes RPC alb"
  }
}

resource "aws_lb_listener" "rpc_alb_http" {
  load_balancer_arn = aws_lb.rpc_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "rpc_alb_https" {
  load_balancer_arn = aws_lb.rpc_alb.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate_validation.rpc_acm.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.rpc_alb_tg.arn
  }
}

resource "aws_lb_target_group" "rpc_alb_tg" {
  vpc_id      = aws_vpc.validator_testnet_vpc.id

  name        = "${var.environment}-rpc-lb-tg"
  target_type = "instance"
  port        = 8545
  protocol    = "HTTP"
  health_check {
    protocol = "HTTP"
    path = "/"
    port = 8545
    matcher = "200,405"
  }
}

resource "aws_lb_target_group_attachment" "rpc_alb_tg_attachment" {
  count = var.full_node_count

  target_group_arn = aws_lb_target_group.rpc_alb_tg.arn
  target_id        = aws_instance.validator_testnet_rpc_nodes[count.index].id
  port             = 8545
}

resource "aws_lb" "blockscout_alb" {
  name               = "${var.environment}-blockscout-alb"
  internal           = false
  load_balancer_type = "application"

  subnets = [aws_subnet.private_subnet.id, aws_subnet.public_subnet.id]
  security_groups = [aws_security_group.alb_sg.id]

  tags = {
    Name = "${var.environment} Blockscout alb"
  }
}

resource "aws_lb_listener" "rpc_blockscout_http" {
  load_balancer_arn = aws_lb.blockscout_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "rpc_blockscout_https" {
  load_balancer_arn = aws_lb.blockscout_alb.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate_validation.explorer_acm.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.blockscout_alb_tg.arn
  }
}

resource "aws_lb_target_group" "blockscout_alb_tg" {
  vpc_id      = aws_vpc.validator_testnet_vpc.id

  name        = "${var.environment}-blockscout-lb-tg"
  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  health_check {
    protocol = "HTTP"
    path = "/"
    port = 80
    matcher = "200"
  }
}

resource "aws_lb_target_group_attachment" "blockscout_alb_tg_attachment" {
  target_group_arn = aws_lb_target_group.blockscout_alb_tg.arn
  target_id        = aws_instance.block_explorer.id
  port             = 80
}

# WS load balancer
resource "aws_lb" "ws_alb" {
  name               = "${var.environment}-full-node-ws-elb"
  internal           = false
  load_balancer_type = "application"

  subnets = [aws_subnet.private_subnet.id, aws_subnet.public_subnet.id]
  security_groups = [aws_security_group.alb_sg.id]

  tags = {
    Name = "${var.environment} RPC WS alb"
  }
}

resource "aws_lb_listener" "ws_alb_http" {
  load_balancer_arn = aws_lb.ws_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "ws_alb_https" {
  load_balancer_arn = aws_lb.ws_alb.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate_validation.ws_acm.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ws_alb_tg.arn
  }
}

resource "aws_lb_target_group" "ws_alb_tg" {
  vpc_id      = aws_vpc.validator_testnet_vpc.id

  name        = "${var.environment}-ws-lb-tg"
  target_type = "instance"
  port        = 8546
  protocol    = "HTTP"
  health_check {
    protocol = "HTTP"
    path = "/"
    port = 8545
    matcher = "200,405"
  }
}

resource "aws_lb_target_group_attachment" "ws_alb_tg_attachment" {
  count = var.full_node_count

  target_group_arn = aws_lb_target_group.ws_alb_tg.arn
  target_id        = aws_instance.validator_testnet_rpc_nodes[count.index].id
  port             = 8546
}

# Grafana load balancer
resource "aws_lb" "grafana_alb" {
  name               = "${var.environment}-grafana-elb"
  internal           = false
  load_balancer_type = "application"

  subnets = [aws_subnet.private_subnet.id, aws_subnet.public_subnet.id]
  security_groups = [aws_security_group.alb_sg.id]

  tags = {
    Name = "${var.environment} Grafana ELB"
  }
}

resource "aws_lb_listener" "grafana_alb_http" {
  load_balancer_arn = aws_lb.grafana_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "grafana_alb_https" {
  load_balancer_arn = aws_lb.grafana_alb.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate_validation.grafana_acm.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.grafana_alb_tg.arn
  }
}

resource "aws_lb_target_group" "grafana_alb_tg" {
  vpc_id      = aws_vpc.validator_testnet_vpc.id

  name        = "${var.environment}-grafana-lb-tg"
  target_type = "instance"
  port        = 3000
  protocol    = "HTTP"
  health_check {
    protocol = "HTTP"
    path = "/api/health"
    port = 3000
  }
}

resource "aws_lb_target_group_attachment" "grafana_alb_tg_attachment" {
  target_group_arn = aws_lb_target_group.grafana_alb_tg.arn
  target_id        = aws_instance.metrics.id
  port             = 3000
}
