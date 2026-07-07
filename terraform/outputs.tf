output "next_steps" {
  value = <<-EOT
    ArgoCD is installed in the "${kubernetes_namespace.argocd.metadata[0].name}" namespace.

    1. Get admin password:
         kubectl -n argocd get secret argocd-initial-admin-secret \
           -o jsonpath="{.data.password}" | base64 -d
    2. Create app secret (exact command in k8s/secrets.example.yaml)
    3. Bootstrap GitOps app:
         kubectl apply -f ../argocd/application.yaml
    4. Push a commit -> Jenkins builds+pushes images+bumps the tag -> ArgoCD syncs
  EOT
}
