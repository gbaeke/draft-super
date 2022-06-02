# Self-signed cert

openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out draft.crt \
            -keyout draft.key \
            -subj "/C=BE/ST=OV/L=Somewhere/O=Inity/OU=IT/CN=draft.baeke.info"

# Key Vault

Created in AKS resource group
Generate a certificate

# Run az aks draft update

Only updates the service.yaml with annotations and uses ClusterIP:

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubernetes.azure.com/ingress-host: draft.baeke.info
    kubernetes.azure.com/tls-cert-keyvault-uri: https://kvdraft.vault.azure.net/certificates/mycert/7f5607fd5bba4d118ba95a153ed6eb05
  creationTimestamp: null
  name: super-api
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: super-api
  type: ClusterIP
status:
  loadBalancer: {}
```

But you need to follow some instructions first: https://docs.microsoft.com/en-us/azure/aks/web-app-routing

az aks enable-addons --resource-group rg-aks --name clu-git --addons web_application_routing