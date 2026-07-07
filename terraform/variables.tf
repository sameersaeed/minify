variable "kubeconfig_path" {
  description = "Path to the kubeconfig pulled from VM's k3s install"
  type        = string
  default     = "~/.kube/minify.yaml"
}

variable "argocd_chart_version" {
  description = "argo-cd Helm chart version"
  type        = string
  default     = "7.6.12"
}

variable "argocd_namespace" {
  type    = string
  default = "argocd"
}
