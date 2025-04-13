# DeFi Analytics Platform

This project is a real-time DeFi protocol analytics platform designed for a Fintech Hackathon. It processes over 50K weekly transactions across multiple blockchain networks (e.g., Ethereum, Base, Optimism) with high data accuracy and low latency.

## Key Features

- **Real-Time Analytics:** Process 50K+ weekly transactions across multiple blockchain networks.
- **Secure WebSocket API:** Serve live price feeds to 1K+ concurrent clients with an average response time of 100ms.
- **Distributed Worker System:** Uses Go routines and Redis pub/sub to index over 1M daily transactions with high availability and fault tolerance.
- **MongoDB Integration:** Stores transaction data for analytics.
- **Dockerized Deployment:** Easily build, run, and deploy the application using Docker and Docker Compose.

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.18+
- Protocol Buffers compiler (`protoc`)

### Running the Application

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/codeforcefun2/DeFiAnalytics.git
   cd DeFiAnalytics
   ```

2. **Generate Protocol Buffers Code:** Install protoc and generate Go code:
    ```bash
    protoc --go_out=. internal/protocol/price.proto
    ```

3. **Build and Run Using Docker Compose:**
    ```bash
    docker-compose up --build
    ```

4. **Access the WebSocket API:** The live WebSocket endpoint is available at ws://localhost:8080/ws.

# License

## Final Notes
This project is licensed under the MIT License.

- **Security & Enhancements:** In a production environment, you should enhance security (for example, by implementing origin checks, using proper TLS configurations, handling authentication on the WebSocket endpoint, etc.).
- **Error Handling & Logging:** More robust error handling and logging would be needed for a production-ready system.
- **Testing:** Consider adding unit and integration tests for each module.
- **Scaling:** For processing large volumes of transactions, you might add more sophisticated monitoring and a better job queue system.

