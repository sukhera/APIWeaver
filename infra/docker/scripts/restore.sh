#!/bin/bash
# MongoDB restore script for APIWeaver
# Usage: ./restore.sh <backup_file.tar.gz>

set -euo pipefail

# Configuration
BACKUP_FILE="${1:-}"
MONGODB_URI="${MONGODB_URI:-}"
if [[ -z "${MONGODB_URI}" ]]; then
    log "ERROR: MONGODB_URI environment variable is not set."
    log "Please set it to your MongoDB connection string (e.g., mongodb://user:pass@host:port/db)."
    exit 1
fi
TEMP_DIR="/tmp/restore_$$"

# Logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Cleanup function
cleanup() {
    log "Cleaning up temporary files..."
    rm -rf "${TEMP_DIR}"
}
trap cleanup EXIT

# Validate input
if [[ -z "${BACKUP_FILE}" ]]; then
    echo "Usage: $0 <backup_file.tar.gz>"
    echo "Available backups:"
    find /backups -name "apiweaver_backup_*.tar.gz" -type f -exec basename {} \; | sort
    exit 1
fi

if [[ ! -f "${BACKUP_FILE}" ]]; then
    log "ERROR: Backup file not found: ${BACKUP_FILE}"
    exit 1
fi

log "Starting MongoDB restore from: ${BACKUP_FILE}"

# Create temporary directory
mkdir -p "${TEMP_DIR}"

# Extract backup
log "Extracting backup archive..."
if tar -xzf "${BACKUP_FILE}" -C "${TEMP_DIR}"; then
    log "Backup extracted successfully"
else
    log "ERROR: Failed to extract backup"
    exit 1
fi

# Find the backup directory
BACKUP_DIR=$(find "${TEMP_DIR}" -maxdepth 1 -type d -name "apiweaver_backup_*" -print -quit)
if [[ -z "${BACKUP_DIR}" ]]; then
    log "ERROR: No backup directory found in archive"
    exit 1
fi

log "Found backup directory: $(basename ${BACKUP_DIR})"

# Confirm restoration
echo "WARNING: This will replace the current database with the backup data."
echo "Current database will be permanently lost."
echo "Backup file: ${BACKUP_FILE}"
echo "Target database: $(echo ${MONGODB_URI} | sed 's/:[^@]*@/:***@/')"
echo ""
read -p "Are you sure you want to proceed? (yes/no): " -r
if [[ ! $REPLY =~ ^yes$ ]]; then
    log "Restore cancelled by user"
    exit 0
fi

# Drop existing database
log "Dropping existing database..."
if mongosh --eval "db.dropDatabase()" "${MONGODB_URI}"; then
    log "Existing database dropped"
else
    log "WARNING: Failed to drop existing database (it may not exist)"
fi

# Restore from backup
log "Restoring from backup..."
if mongorestore --uri="${MONGODB_URI}" --gzip "${BACKUP_DIR}"; then
    log "Database restored successfully"
else
    log "ERROR: Database restore failed"
    exit 1
fi

# Verify restoration
log "Verifying restoration..."
COLLECTIONS=$(mongosh --quiet --eval "db.getCollectionNames().length" "${MONGODB_URI}")
if [[ "${COLLECTIONS}" -gt 0 ]]; then
    log "Verification successful: ${COLLECTIONS} collections restored"
else
    log "WARNING: No collections found after restore"
fi

# Show restoration summary
log "=== Restoration Summary ==="
mongosh --quiet --eval "
    print('Database: ' + db.getName());
    print('Collections: ' + db.getCollectionNames().length);
    db.getCollectionNames().forEach(function(name) {
        var count = db.getCollection(name).countDocuments();
        print('  ' + name + ': ' + count + ' documents');
    });
" "${MONGODB_URI}"

log "Restore completed successfully!"