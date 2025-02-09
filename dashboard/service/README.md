# `dashboard/service`

This module provisions a Google Cloud Monitoring dashboard for a regionalized Cloud Run service.

It assumes the service has the same name in all regions.

```hcl
// Create a network with several regional subnets
module "networking" {
  source = "chainguard-dev/glue/cloudrun//networking"

  name       = "my-networking"
  project_id = var.project_id
  regions    = [...]
}

// Run a regionalized cloud run service "frontend" to serve requests.
resource "google_cloud_run_v2_service" "frontend" {
  for_each = module.networking.regional-networks
  name     = "frontend"

  //...
  template {
    //...
    containers {
      image = "..."
    }
  }
}

// Set up a dashboard for a regionalized service named "frontend".
module "service-dashboard" {
  source       = "chainguard-dev/glue/cloudrun//dashboard/service"
  service_name = "frontend"
}
```

The dashboard it creates includes widgets for service logs, request count, latency (p50,p95,p99), instance count grouped by revision, CPU and memory utilization, startup latency, and sent/received bytes.

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cpu_utilization"></a> [cpu\_utilization](#module\_cpu\_utilization) | ../widgets/xy | n/a |
| <a name="module_incoming_latency"></a> [incoming\_latency](#module\_incoming\_latency) | ../widgets/latency | n/a |
| <a name="module_instance_count"></a> [instance\_count](#module\_instance\_count) | ../widgets/xy | n/a |
| <a name="module_logs"></a> [logs](#module\_logs) | ../widgets/logs | n/a |
| <a name="module_memory_utilization"></a> [memory\_utilization](#module\_memory\_utilization) | ../widgets/xy | n/a |
| <a name="module_received_bytes"></a> [received\_bytes](#module\_received\_bytes) | ../widgets/xy | n/a |
| <a name="module_request_count"></a> [request\_count](#module\_request\_count) | ../widgets/xy | n/a |
| <a name="module_sent_bytes"></a> [sent\_bytes](#module\_sent\_bytes) | ../widgets/xy | n/a |
| <a name="module_startup_latency"></a> [startup\_latency](#module\_startup\_latency) | ../widgets/xy | n/a |

## Resources

| Name | Type |
|------|------|
| [google_monitoring_dashboard.dashboard](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/monitoring_dashboard) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | Name of the service(s) to monitor | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
