package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict in production
		return true
	},
}

// Server represents the WebSocket API server
type Server struct {
	addr     string
	handler  *Handler
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
	shutdown chan bool
}

// NewServer creates a new API server
func NewServer(addr string, handler *Handler) *Server {
	return &Server{
		addr:     addr,
		handler:  handler,
		clients:  make(map[*websocket.Conn]bool),
		shutdown: make(chan bool),
	}
}

// Start starts the WebSocket server
func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/health", s.handleHealth)

	log.Printf("🚀 WebSocket server starting on %s", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

// Broadcast sends a message to all connected clients
func (s *Server) Broadcast(msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Register client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	log.Printf("✅ New WebSocket client connected (total: %d)", len(s.clients))

	// Handle client disconnection
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		log.Printf("❌ Client disconnected (remaining: %d)", len(s.clients))
	}()

	// Read messages from client
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			s.sendError(conn, "", "Invalid message format", err)
			continue
		}

		// Handle message
		response := s.handler.HandleMessage(msg)
		if response != nil {
			s.sendMessage(conn, *response)
		}
	}
}

// sendMessage sends a message to a specific client
func (s *Server) sendMessage(conn *websocket.Conn, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// sendError sends an error message to a client
func (s *Server) sendError(conn *websocket.Conn, requestID string, message string, err error) {
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}

	msg := Message{
		Type: TypeError,
		ID:   requestID,
		Payload: ErrorPayload{
			Message: errMsg,
		},
	}

	s.sendMessage(conn, msg)
}

// GetBroadcaster returns a function that can be used to broadcast messages
func (s *Server) GetBroadcaster() func(Message) {
	return s.Broadcast
}
