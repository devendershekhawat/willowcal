#!/bin/sh
set -e

echo "🌿 Willowcal Starting..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📡 WebSocket Server: ws://localhost:${WILLOWCAL_PORT}/ws"
echo "🌐 Web Interface: http://localhost:${WILLOWCAL_PORT}"
echo "📂 Workspace: ${WORKSPACE_DIR}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Start the willowcal server
exec /app/willowcal server "${WILLOWCAL_PORT}" "${WORKSPACE_DIR}" "${STATIC_DIR}"
