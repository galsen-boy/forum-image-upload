#!/bin/bash

# Création des dossiers
mkdir -p forum/cmd
mkdir -p forum/internal/handlers
mkdir -p forum/internal/models
mkdir -p forum/internal/repositories
mkdir -p forum/internal/services
mkdir -p forum/templates
mkdir -p forum/static
mkdir -p forum/config
mkdir -p forum/db

# Création des fichiers
touch forum/cmd/main.go
touch forum/internal/handlers/auth_handler.go
touch forum/internal/handlers/post_handler.go
touch forum/internal/handlers/user_handler.go
touch forum/internal/models/user.go
touch forum/internal/models/post.go
touch forum/internal/repositories/user_repo.go
touch forum/internal/repositories/post_repo.go
touch forum/internal/services/auth_service.go
touch forum/internal/services/post_service.go
touch forum/internal/services/user_service.go
touch forum/main.go
touch forum/go.mod
touch forum/go.sum

echo "Structure du projet créée avec succès!"
