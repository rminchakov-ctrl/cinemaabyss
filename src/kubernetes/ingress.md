helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.image.repository="registry.k8s.io/ingress-nginx/controller" \
  --set controller.image.tag="v1.8.1" \
  --set controller.resources.requests.cpu=20m \
  --set controller.resources.requests.memory=256Mi \
  --set controller.resources.limits.cpu=500m \
  --set controller.resources.limits.memory=512Mi \
  --set controller.hostNetwork=false \
  --set controller.admissionWebhooks.enabled=true \
  --set controller.admissionWebhooks.createSecret=true \
  --set controller.admissionWebhooks.patch.enabled=true


  --set controller.admissionWebhooks.enabled=false \