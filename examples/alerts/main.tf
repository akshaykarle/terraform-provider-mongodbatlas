resource "mongodbatlas_alert_configuration" "restarts" {
  group = "${var.group_id}"
  event_type_name = "CLUSTER_MONGOS_IS_MISSING"
  enabled = true

  matchers {
    field_name = "CLUSTER_NAME"
    operator = "EQUALS"
    value = "${var.cluster_name}"
  }

  notifications {
    type_name = "PAGER_DUTY"
    service_key = "${var.pagerduty_service_key}"
    interval_min = 5
    delay_min = 0
  }
}

resource "mongodbatlas_alert_configuration" "high_cpu" {
  group = "${var.group_id}"
  event_type_name = "OUTSIDE_METRIC_THRESHOLD"
  enabled = true

  metric_threshold {
    metric_name = "DISK_PARTITION_UTILIZATION_DATA"
    operator = "GREATER_THAN"
    threshold = 96.0
    mode = "AVERAGE"
    units = "RAW"
  }

  notifications {
    type_name = "PAGER_DUTY"
    service_key = "${var.pagerduty_service_key}"
    interval_min = 5
    delay_min = 0
  }
  notifications {
    type_name = "PAGER_DUTY"
    service_key = "c94636f7fb464065898c6b9ca06192a0"
    interval_min = 5
    delay_min = 0
  }
}
