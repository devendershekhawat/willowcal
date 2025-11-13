#!/bin/sh
set -e

echo "🌿 Willowcal Starting..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📡 WebSocket Server: ws://localhost:${PORT}/ws"
echo "🌐 Web Interface: http://localhost:${PORT}"
echo "📂 Workspace: ${WORKSPACE_DIR}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Start the willowcal server
exec /app/willowcal server "${PORT}" "${WORKSPACE_DIR}" "${STATIC_DIR}"
