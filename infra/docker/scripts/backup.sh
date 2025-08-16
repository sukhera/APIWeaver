#!/bin/bash
# MongoDB backup script for APIWeaver
# This script creates compressed backups with rotation

set -euo pipefail

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/backups}"
MONGODB_URI="${MONGODB_URI:-mongodb://admin:apiweaver123@mongodb:27017/apiweaver?authSource=admin}"
BACKUP_RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="apiweaver_backup_${TIMESTAMP}"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}"

# Logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

log "Starting MongoDB backup for APIWeaver..."
log "Backup path: ${BACKUP_PATH}"

# Create the backup
if mongodump --uri="${MONGODB_URI}" --out="${BACKUP_PATH}" --gzip; then
    log "MongoDB backup completed successfully"
else
    log "ERROR: MongoDB backup failed"
    exit 1
fi

# Create archive
log "Creating compressed archive..."
cd "${BACKUP_DIR}"
if tar -czf "${BACKUP_NAME}.tar.gz" "${BACKUP_NAME}"; then
    log "Archive created: ${BACKUP_NAME}.tar.gz"
    rm -rf "${BACKUP_NAME}"
else
    log "ERROR: Failed to create archive"
    exit 1
fi

# Calculate backup size
BACKUP_SIZE=$(du -h "${BACKUP_NAME}.tar.gz" | cut -f1)
log "Backup size: ${BACKUP_SIZE}"

# Clean up old backups
log "Cleaning up backups older than ${BACKUP_RETENTION_DAYS} days..."
find "${BACKUP_DIR}" -name "apiweaver_backup_*.tar.gz" -type f -mtime +${BACKUP_RETENTION_DAYS} -delete

# Show remaining backups
REMAINING_BACKUPS=$(find "${BACKUP_DIR}" -name "apiweaver_backup_*.tar.gz" -type f | wc -l)
log "Backup cleanup completed. Remaining backups: ${REMAINING_BACKUPS}"

# Create backup metadata
cat > "${BACKUP_DIR}/${BACKUP_NAME}.meta" << EOF
{
  "backup_name": "${BACKUP_NAME}",
  "timestamp": "${TIMESTAMP}",
  "date": "$(date -Iseconds)",
  "size": "${BACKUP_SIZE}",
  "mongodb_uri": "$(echo ${MONGODB_URI} | sed 's/:[^@]*@/:***@/')",
  "retention_days": ${BACKUP_RETENTION_DAYS},
  "checksum": "$(sha256sum ${BACKUP_NAME}.tar.gz | cut -d' ' -f1)"
}
EOF

log "Backup completed successfully: ${BACKUP_NAME}.tar.gz"
log "Metadata saved: ${BACKUP_NAME}.meta"