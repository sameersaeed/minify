resource "kubernetes_namespace" "argocd" {
  metadata {
    name = var.argocd_namespace
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = var.argocd_chart_version
  namespace  = kubernetes_namespace.argocd.metadata[0].name

  # Single small VM to trim replica counts / resource requests on ArgoCD components
  set {
    name  = "redis.resources.requests.cpu"
    value = "25m"
  }
  set {
    name  = "controller.resources.requests.cpu"
    value = "50m"
  }
  set {
    name  = "server.resources.requests.cpu"
    value = "25m"
  }
  set {
    name  = "repoServer.resources.requests.cpu"
    value = "25m"
  }
}
