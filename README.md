# PromProxy

A daemon that receives Prometheus API calls from a system such as Grafana, sends that request to multiple Prometheus instances, and aggregates the responses.

Systems like Grafana don't work well when your data is split up across multiple non overlapping Prometheus servers. This daemon attempts to trick the requester into thinking that there is just one server, but the requests to go many.
