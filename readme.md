# Auth Service in Go

This is a Go-based authentication service - for demonstration purposes.  
It includes user registration, login, JWT validation, and health checks.  
The service runs on **Kubernetes** using **Helm**, with **PostgreSQL** automatically deployed as a dependency.

---

## **Features**

- `/health` - Health check endpoint  
- `/register` - Register a new user  
- `/login` - Login and receive a JWT token  
- `/validate` - Validate JWT token  
- Automatic PostgreSQL deployment with Helm subchart  
- Environment variables managed via Helm `Secret`

---

## **Prerequisites**

- [Docker](https://www.docker.com/)  
- [kubectl](https://kubernetes.io/docs/tasks/tools/)  
- [k3d](https://k3d.io/) (lightweight Kubernetes)  
- [Helm 3+](https://helm.sh/)  

---

## **Getting Started**

### 1. Clone the Repository

```bash
git clone https://github.com/naphalm/auth-service-cicd.git
cd auth-service-cicd
```

### 2. Create a Local Kubernetes Cluster
```bash
k3d cluster create dev-cluster --agents 1
kubectl get nodes
```
You should see both server and agent nodes in Ready state.

### 3. Deploy Auth Service with PostgreSQL
The Helm chart is configured to automatically deploy PostgreSQL as a subchart.
```bash
helm dependency update ./helm/auth-service
helm install auth-service ./helm/auth-service
```

This will deploy:
- auth-service Deployment and Service
- auth-service-postgresql Deployment, StatefulSet, and Service
- Secret with DATABASE_URL and JWT_SECRET automatically injected

### 4. Verify pods
```bash
kubectl get pods
```
Expected output:
```bash
NAME                             READY   STATUS
auth-service-xxxxxxx             1/1     Running
auth-service-postgresql-0        1/1     Running
```

### 5. Accessing the Service
```bash
kubectl port-forward svc/auth-service 3000:3000
```

Now the endpoints are accessible via:
- Health check: http://localhost:3000/health
- Register: http://localhost:3000/register
- Login: http://localhost:3000/login
- Validate JWT: http://localhost:3000/validate

#### Example API Requests
- Register
```bash
curl -X POST http://localhost:3000/register \
-H "Content-Type: application/json" \
-d '{"email":"test@test.com","password":"Pass1234"}'
```

- Login
```bash
curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email":"test@test.com","password":"Pass1234"}'
```

- Validate JWT
```bash
curl -X GET http://localhost:3000/validate \
-H "Authorization: Bearer <JWT_TOKEN>"
```

### Environment variables
| Name         | Description                        |
| ------------ | ---------------------------------- |
| DATABASE_URL | Connection string for PostgreSQL   |
| JWT_SECRET   | Secret key used to sign JWT tokens |

They are automatically managed by Helm. No .env file is needed.

### Development Notes
Go binary is built using a multi-stage Dockerfile

Migration script (migrate.sh) runs on container start

PostgreSQL credentials are injected into the auth-service via Secret

You can customize CPU/memory limits in values.yaml

### Cleanup
```bash
helm uninstall auth-service
k3d cluster delete dev-cluster
```