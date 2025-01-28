# Certificates for explorer + rpc
data "aws_route53_zone" "root_zone" {
  name         = "zama.ai"
  private_zone = false
}

# block explorer certificates + DNS
resource "aws_acm_certificate" "explorer_acm" {
  domain_name       = var.explorer_dns_name
  validation_method = "DNS"
}

resource "aws_route53_record" "explorer_dns_verification" {
  for_each = {
    for dvo in aws_acm_certificate.explorer_acm.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_route53_record" "explorer_dns_cname" {
  name            = var.explorer_dns_name
  records         = [aws_lb.blockscout_alb.dns_name]
  type            = "CNAME"

  allow_overwrite = true
  ttl             = 60
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_acm_certificate_validation" "explorer_acm" {
  certificate_arn         = aws_acm_certificate.explorer_acm.arn
  validation_record_fqdns = [for record in aws_route53_record.explorer_dns_verification : record.fqdn]
}


# RPC certificates + DNS
resource "aws_acm_certificate" "rpc_acm" {
  domain_name       = var.rpc_dns_name
  validation_method = "DNS"
}

resource "aws_route53_record" "rpc_dns_verification" {
  for_each = {
    for dvo in aws_acm_certificate.rpc_acm.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_route53_record" "rpc_dns_cname" {
  name            = var.rpc_dns_name
  records         = [aws_lb.rpc_alb.dns_name]
  type            = "CNAME"

  allow_overwrite = true
  ttl             = 60
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_acm_certificate_validation" "rpc_acm" {
  certificate_arn         = aws_acm_certificate.rpc_acm.arn
  validation_record_fqdns = [for record in aws_route53_record.rpc_dns_verification : record.fqdn]
}

# WS certificates + DNS
resource "aws_acm_certificate" "ws_acm" {
  domain_name       = var.ws_dns_name
  validation_method = "DNS"
}

resource "aws_route53_record" "ws_dns_verification" {
  for_each = {
    for dvo in aws_acm_certificate.ws_acm.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_route53_record" "ws_dns_cname" {
  name            = var.ws_dns_name
  records         = [aws_lb.ws_alb.dns_name]
  type            = "CNAME"

  allow_overwrite = true
  ttl             = 60
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_acm_certificate_validation" "ws_acm" {
  certificate_arn         = aws_acm_certificate.ws_acm.arn
  validation_record_fqdns = [for record in aws_route53_record.ws_dns_verification : record.fqdn]
}

# Grafana certificates + DNS
resource "aws_acm_certificate" "grafana_acm" {
  domain_name       = var.grafana_dns_name
  validation_method = "DNS"
}

resource "aws_route53_record" "grafana_dns_verification" {
  for_each = {
    for dvo in aws_acm_certificate.grafana_acm.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_route53_record" "grafana_dns_cname" {
  name            = var.grafana_dns_name
  records         = [aws_lb.grafana_alb.dns_name]
  type            = "CNAME"

  allow_overwrite = true
  ttl             = 60
  zone_id         = data.aws_route53_zone.root_zone.zone_id
}

resource "aws_acm_certificate_validation" "grafana_acm" {
  certificate_arn         = aws_acm_certificate.grafana_acm.arn
  validation_record_fqdns = [for record in aws_route53_record.grafana_dns_verification : record.fqdn]
}
