#!/bin/bash

echo "🛠 Creating frontend structure..."

# Frontend
mkdir -p frontend/public
mkdir -p frontend/src/{api,assets,context,hooks,layouts,pages,routes,styles,tests,utils}
mkdir -p frontend/src/components/{Common,Menu,Auth,Orders,Admin,Cart,Payment}

# Frontend Root Files
touch frontend/package.json frontend/vite.config.js
touch frontend/src/{App.js,index.js}
touch frontend/src/styles/index.css

# Pages
touch frontend/src/pages/{Home.js,Menu.js,Orders.js,AdminDashboard.js}

# Routes
touch frontend/src/routes/AppRoutes.js

# Components
touch frontend/src/components/Auth/Auth.js
touch frontend/src/components/Menu/{MenuList.js,MenuItemCard.js}
touch frontend/src/components/Cart/Cart.js
touch frontend/src/components/Orders/{OrderHistory.js,OrderStatus.js}
touch frontend/src/components/Admin/{AdminMenuManager.js,AdminOrderManager.js}
touch frontend/src/components/Payment/PaymentForm.js
touch frontend/src/components/Common/{Header.js,Footer.js}

# Context & Utils
touch frontend/src/context/{AuthContext.js,CartContext.js}
touch frontend/src/utils/{api.js,validators.js}

# Tests
touch frontend/src/tests/{App.test.js,components.test.js}

echo "✅ Frontend folders and files created."


echo "🛠 Creating backend microservices..."

services=("customer-service" "menu-service" "order-service" "payment-service" "admin-service" "gateway-service")

for service in "${services[@]}"; do
  mkdir -p backend/$service/{controllers,models,routes,db,utils,middleware}
  touch backend/$service/{main.go,Dockerfile,go.mod,go.sum}
done

echo "✅ Backend microservices structure ready."


echo "🛠 Creating database folders..."
mkdir -p database/couchbase/buckets
touch database/couchbase/init.sh
touch database/couchbase/buckets/{customers.json,menu.json,orders.json,admin.json}

echo "🛠 Creating deployment configs..."
mkdir -p deploy/{docker,k8s,github-actions,helm}
touch deploy/docker/{docker-compose.yml,couchbase.Dockerfile,nginx.Dockerfile,.env}
touch deploy/k8s/{namespace.yaml,services.yaml,ingress.yaml,secrets.yaml,configmap.yaml}
touch deploy/k8s/{customer-deployment.yaml,menu-deployment.yaml,order-deployment.yaml,payment-deployment.yaml,admin-deployment.yaml,gateway-deployment.yaml,frontend-deployment.yaml}
touch deploy/github-actions/{backend.yml,frontend.yml}

echo "🛠 Creating scripts..."
mkdir -p scripts
touch scripts/{seed_data.sh,create_admin.sh,cleanup.sh,test_endpoints.sh}

echo "🛠 Creating monitoring..."
mkdir -p monitoring/{prometheus,grafana/dashboards,alertmanager}
touch monitoring/prometheus/{prometheus.yml,Dockerfile}
touch monitoring/grafana/Dockerfile

echo "🛠 Creating docs..."
mkdir -p docs/flowcharts
touch docs/{architecture.md,api-reference.md,db-schema.md,README.md}
touch docs/flowcharts/{place-order-flow.png,admin-add-menu.png,structure-diagram.png}

echo "🎉 Project folder structure is fully set up!"

