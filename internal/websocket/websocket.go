package websocket

import (
    "crypto/tls"
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/yourusername/defi-analytics/internal/config"
)

// PriceUpdate represents a live price update message.
type PriceUpdate struct {
    Symbol string  `json:"symbol"`
    Price  float64 `json:"price"`
}

// Server maintains active WebSocket connections and broadcasts messages.
type Server struct {
    cfg          *config.Config
    upgrader     websocket.Upgrader
    connections  map[*websocket.Conn]bool
    connectionsM sync.Mutex
}

// NewServer creates a new WebSocket server instance.
func NewServer(cfg *config.Config) *Server {
    return &Server{
        cfg: cfg,
        upgrader: websocket.Upgrader{
            // In a production system, you should check the origin.
            CheckOrigin: func(r *http.Request) bool { return true },
        },
        connections: make(map[*websocket.Conn]bool),
    }
}

// HandleConnections upgrades HTTP connections to WebSocket and maintains them.
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
    // Enable TLS upgrade if required (e.g., if using "wss://")
    if s.cfg.ServerAddress[:5] == "wss://" {
        r.TLS = &tls.ConnectionState{}
    }
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    s.connectionsM.Lock()
    s.connections[conn] = true
    s.connectionsM.Unlock()
    log.Printf("New WebSocket connection established")

    // Continuously read messages from the client (e.g., for heartbeat or control)
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("connection closed: %v", err)
            s.connectionsM.Lock()
            delete(s.connections, conn)
            s.connectionsM.Unlock()
            break
        }
    }
}

// BroadcastPriceUpdate sends a live price update to all connected clients.
func (s *Server) BroadcastPriceUpdate(update PriceUpdate) {
    s.connectionsM.Lock()
    defer s.connectionsM.Unlock()

    message, err := json.Marshal(update)
    if err != nil {
        log.Printf("failed to marshal price update: %v", err)
        return
    }

    for conn := range s.connections {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Printf("error sending message: %v", err)
            conn.Close()
            delete(s.connections, conn)
        }
    }
}
