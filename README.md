# authorizer
service for handling application authorizing

## JWT keys (RS256) di Kubernetes

Aplikasi memuat private/public key untuk JWT dengan urutan berikut:

1. **Environment variable (PEM string)**  
   Jika `JWT_PRIVATE_KEY` dan `JWT_PUBLIC_KEY` di-set (isi PEM), key diambil dari env (cocok untuk `secretKeyRef`).

2. **File**  
   Dicoba berurutan: `JWT_PRIVATE_KEY_PATH` / `JWT_PUBLIC_KEY_PATH` → `./private.pem` / `./public.pem` → `/app/secrets/private.pem` / `/app/secrets/public.pem` (mount secret di K8s).

### Opsi 1: Mount secret sebagai volume (disarankan)

Buat secret dari file PEM:
```bash
kubectl create secret generic jwt-keys --from-file=private.pem=./private.pem --from-file=public.pem=./public.pem
```

Di Deployment, mount secret ke `/app/secrets`:
```yaml
volumeMounts:
  - name: jwt-keys
    mountPath: /app/secrets
    readOnly: true
volumes:
  - name: jwt-keys
    secret:
      secretName: jwt-keys
```

Tanpa env tambahan, aplikasi akan otomatis memakai `/app/secrets/private.pem` dan `/app/secrets/public.pem` jika file di path default tidak ada.

Atau set path secara eksplisit:
```yaml
env:
  - name: JWT_PRIVATE_KEY_PATH
    value: /app/secrets/private.pem
  - name: JWT_PUBLIC_KEY_PATH
    value: /app/secrets/public.pem
```

### Opsi 2: Secret sebagai environment variable

Mount key sebagai env (isi raw PEM):
```yaml
env:
  - name: JWT_PRIVATE_KEY
    valueFrom:
      secretKeyRef:
        name: jwt-keys
        key: private.pem
  - name: JWT_PUBLIC_KEY
    valueFrom:
      secretKeyRef:
        name: jwt-keys
        key: public.pem
```
