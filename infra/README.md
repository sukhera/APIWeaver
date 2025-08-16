# APIWeaver Infrastructure

This directory contains all Docker and infrastructure configuration files organized by environment.

## Directory Structure

```
infra/
├── local/                    # Local development environment
│   ├── docker-compose.yml    # Development services with hot reload
│   └── Dockerfile.dev        # Development container with debugging
├── production/               # Production environment  
│   ├── docker-compose.prod.yml  # Production overrides and scaling
│   └── Dockerfile            # Optimized production container
├── docker/                   # Shared Docker configurations
│   ├── mongodb/             # MongoDB configuration files
│   ├── nginx/               # Nginx configuration files
│   └── scripts/             # Backup and utility scripts
├── Makefile.docker          # Docker management commands
└── .dockerignore            # Docker build context exclusions
```

## Quick Start

### Development Environment

```bash
# Start development environment with hot reload
cd infra
make docker-dev

# View development logs
make docker-logs

# Get shell access to development container
make docker-shell
```

### Production Environment

```bash
# Setup production environment (creates .env from template)
cd infra
make docker-prod-setup

# Start production environment
make docker-prod

# Scale production services
make docker-prod-scale API=3 NGINX=2
```

## Available Commands

Run `make help` in the infra directory to see all available commands:

- **Build Commands**: `docker-build`, `docker-build-no-cache`
- **Development**: `docker-dev`, `docker-dev-logs`, `docker-test`
- **Production**: `docker-prod`, `docker-prod-setup`, `docker-prod-scale`
- **Maintenance**: `docker-backup`, `docker-restore`, `docker-health`
- **Cleanup**: `docker-clean`, `docker-clean-volumes`

## Environment Files

- `.env.example` - Example environment variables
- `.env.production` - Production environment template
- Create `.env` from the appropriate template for your environment

## Best Practices

1. **Development**: Use `local/` configurations for development with hot reload
2. **Production**: Use `production/` configurations with security hardening
3. **Security**: Never commit sensitive environment variables
4. **Scaling**: Use Docker Compose deploy configurations for production scaling
5. **Monitoring**: Enable health checks and logging for production deployments

## Troubleshooting

- Check service health: `make docker-health`
- View logs: `make docker-logs`
- Get container stats: `make docker-stats`
- Clean up resources: `make docker-clean`