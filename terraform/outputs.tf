output "cluster_name" {
  description = "Created kind cluster name"
  value       = kind_cluster.week1.name
}

output "kubeconfig_path" {
  description = "Kubeconfig generated for kubectl"
  value       = abspath("${path.module}/kubeconfig")
}
