#!/bin/bash
# Start lfs-test-server with admin interface enabled

export LFS_CONTENTPATH=/opt/lfs-test-server
export LFS_ADMINUSER=admin
export LFS_ADMINPASS=admin123

# Install lfs-test-server if not found
if [ ! `which lfs-test-server` ]; then
  go install github.com/git-lfs/lfs-test-server@latest
fi

export LFS_REPO=/opt/lfs-test-server
export LFS_PORT=8080
export SERVER=`uname -n`
export LOG_FILE="$LFS_REPO/lfs-server.log"

cd "$LFS_REPO" | exit 1

# Kill any existing instance
pkill lfs-test-server

rm "$LOG_FILE" # Not everyone might want this

# Start server in background with verbose logging
nohup ~/go/bin/lfs-test-server -verbose -addr :$LFS_PORT &> "$LOG_FILE" &

echo "LFS Test Server started on $SERVER:$LFS_PORT"
echo "Admin interface: http://$SERVER:$LFS_PORT/mgmt"
echo "Tailing $SERVER:$LOG_FILE..."
tail -f "$LOG_FILE"

