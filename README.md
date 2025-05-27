# 🐳 Bonnes Pratiques Dockerfile

> Guide complet pour créer des Dockerfiles optimisés, sécurisés et performants

[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Best Practices](https://img.shields.io/badge/Best_Practices-✅-green?style=for-the-badge)](https://docs.docker.com/develop/dev-best-practices/)

## 📋 Table des Matières

- [🎯 Objectifs](#-objectifs)
- [🚀 Quick Start](#-quick-start)
- [📏 Règles Fondamentales](#-règles-fondamentales)
- [🏗️ Structure Recommandée](#️-structure-recommandée)
- [⚡ Optimisation des Performances](#-optimisation-des-performances)
- [🔒 Sécurité](#-sécurité)
- [📊 Multi-Stage Builds](#-multi-stage-builds)
- [🛠️ Outils et Validation](#️-outils-et-validation)
- [📈 Métriques](#-métriques)
- [🔗 Ressources](#-ressources)

## 🎯 Objectifs

Ce repository présente les bonnes pratiques pour :
- ✅ **Réduire la taille des images** (jusqu'à -70%)
- ✅ **Accélérer les builds** (jusqu'à -90% avec le cache)
- ✅ **Améliorer la sécurité** (principe du moindre privilège)
- ✅ **Garantir la reproductibilité** (versions fixes)
- ✅ **Optimiser les performances** (temps de démarrage divisé par 4)

## 🚀 Quick Start

### Dockerfile de Base (❌ À éviter)

```dockerfile
FROM node:latest
COPY .. .
RUN npm install
EXPOSE 3000
CMD ["npm", "start"]
```

### Dockerfile Optimisé (✅ Recommandé)
```dockerfile
# Multi-stage build
FROM node:18-alpine AS builder

ENV NODE_ENV=production
WORKDIR /app

# Cache des dépendances
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# Build
COPY . .
RUN npm run build

# Image finale
FROM node:18-alpine
WORKDIR /app

# Sécurité
RUN addgroup -g 1001 appuser && \
    adduser -u 1001 -G appuser -D appuser

# Copie optimisée
COPY --from=builder --chown=appuser:appuser /app/node_modules ./node_modules
COPY --from=builder --chown=appuser:appuser /app/dist ./dist

USER appuser
EXPOSE 3000

# Monitoring
HEALTHCHECK --interval=30s --timeout=5s \
  CMD curl --fail http://localhost:3000/health || exit 1

CMD ["node", "dist/index.js"]
```

## 📏 Règles Fondamentales

### 1. Image de Base
```dockerfile
# ✅ Faire
FROM node:18.17.0-alpine    # Version fixe + image légère
FROM nginx:1.24-alpine      # Image officielle

# ❌ Éviter  
FROM node:latest            # Version flottante
FROM random-user/node       # Image non officielle
```

### 2. Organisation des Couches
```dockerfile
# ✅ Du moins changeant au plus changeant
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./       # Dépendances (change peu)
RUN npm ci
COPY . .                    # Code source (change souvent)
RUN npm run build
```

### 3. Fichier .dockerignore
```bash
# Fichiers à exclure
node_modules/
.git/
.env*
npm-debug.log*
coverage/
.nyc_output/
.DS_Store
```

## 🏗️ Structure Recommandée

### Template Dockerfile Multi-Stage

```dockerfile
# ======================
# Étape 1: Dependencies
# ======================
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# ======================  
# Étape 2: Builder
# ======================
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# ======================
# Étape 3: Runner
# ======================
FROM node:18-alpine AS runner
WORKDIR /app

# Sécurité
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copie sélective
COPY --from=builder /app/dist ./dist
COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder --chown=nextjs:nodejs /app/package.json ./

USER nextjs
EXPOSE 3000

# Métadonnées
LABEL maintainer="team@company.com"
LABEL description="Application Node.js optimisée"
LABEL version="1.0.0"

CMD ["node", "dist/index.js"]
```

## ⚡ Optimisation des Performances

### Cache Docker
```dockerfile
# ✅ Optimisé - Les dépendances changent rarement
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

# ❌ Non optimisé - Tout est recalculé à chaque changement
COPY . .
RUN npm ci && npm run build
```

### Réduction des Couches
```dockerfile
# ✅ Une seule couche
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        git \
        python3 \
    && rm -rf /var/lib/apt/lists/*

# ❌ Plusieurs couches inutiles
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git  
RUN apt-get install -y python3
```

## 🔒 Sécurité

### Utilisateur Non-Root
```dockerfile
# Créer un utilisateur système
RUN groupadd -r appuser && useradd --no-log-init -r -g appuser appuser

# Ou pour Alpine
RUN addgroup -g 1001 appuser && \
    adduser -u 1001 -G appuser -D appuser

# Changer de propriétaire lors de la copie
COPY --from=builder --chown=appuser:appuser /app/dist ./dist

# Basculer vers l'utilisateur
USER appuser
```

### Gestion des Secrets
```dockerfile
# ❌ Jamais de secrets en dur
ENV API_KEY=secret123

# ✅ Utiliser les build secrets ou runtime
RUN --mount=type=secret,id=api_key \
    API_KEY=$(cat /run/secrets/api_key) && \
    # Utiliser la clé sans la persister
```

## 📊 Multi-Stage Builds

### Exemple Python
```dockerfile
# Build stage
FROM python:3.11-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

# Production stage  
FROM python:3.11-slim
WORKDIR /app

# Copier uniquement les packages installés
COPY --from=builder /root/.local /root/.local
COPY . .

# S'assurer que les scripts locaux sont dans PATH
ENV PATH=/root/.local/bin:$PATH

CMD ["python", "app.py"]
```

### Exemple Go
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🛠️ Outils et Validation

### Linting avec Hadolint
```bash
# Installation
brew install hadolint

# Analyse
hadolint Dockerfile
```

### Scan de Sécurité
```bash
# Docker Scout
docker scout cves myapp:latest

# Trivy
trivy image myapp:latest
```

### Extension VS Code
- **Docker Extension (Beta)** : Linting intégré, navigation de code
- **Dockerfile** : Syntaxe highlighting

### Scripts de Build
```bash
#!/bin/bash
# build.sh

# Build avec optimisations
docker build \
  --target production \
  --build-arg NODE_ENV=production \
  --tag myapp:$(git rev-parse --short HEAD) \
  --tag myapp:latest \
  .

# Analyse de sécurité
docker scout cves myapp:latest

# Test de santé
docker run --rm -d --name test-container -p 3000:3000 myapp:latest
sleep 10
curl -f http://localhost:3000/health || exit 1
docker stop test-container
```

## 📈 Métriques

### Comparaison Avant/Après

| Aspect | Avant Optimisation | Après Optimisation | Amélioration |
|--------|-------------------|-------------------|-------------|
| **Taille image** | 1.2 GB | 350 MB | **-71%** |
| **Build initial** | 12 min | 6 min | **-50%** |
| **Build avec cache** | 12 min | 45 sec | **-94%** |
| **Temps démarrage** | 45 sec | 8 sec | **-82%** |
| **Vulnérabilités** | 23 | 4 | **-83%** |
| **Couches Docker** | 15 | 8 | **-47%** |

### Commandes d'Analyse
```bash
# Taille de l'image
docker images myapp:latest

# Historique des couches
docker history myapp:latest

# Analyse détaillée
docker system df -v
```

## 🔗 Ressources

### Documentation Officielle
- [Docker Best Practices](https://docs.docker.com/develop/best-practices/)
- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

### Outils
- [Hadolint](https://github.com/hadolint/hadolint) - Dockerfile linter
- [Docker Scout](https://docs.docker.com/scout/) - Analyse de vulnérabilités
- [Trivy](https://github.com/aquasecurity/trivy) - Scanner de sécurité
- [Dive](https://github.com/wagoodman/dive) - Analyse des couches

### Templates
- [Node.js](examples/nodejs/Dockerfile)
- [Python](examples/python/Dockerfile)
- [Go](examples/golang/Dockerfile)
- [PHP](examples/php/Dockerfile)

---

## 🤝 Contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/amélioration`)
3. Commit (`git commit -am 'Ajout nouvelle pratique'`)
4. Push (`git push origin feature/amélioration`)
5. Créer une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

---

<div align="center">

**Fait avec ❤️ pour la communauté Docker**

[⬆ Retour en haut](#-bonnes-pratiques-dockerfile)

</div>
