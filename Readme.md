
---

# Revenge

## Introduction
This project implements a Complete Backend Infrastruture Framework in Golang with the ability to deploy Single to Multi system Dsitributed services via Runners and Built in Channel Network Registry with GrPc services for fast communication between Different Runner operated via a Central Revenge Server with features like (GraphQL,Rest,SSE,webSockets,Https and lot more architectures to choose from),. It features a custom message struct for efficiently managing messages and a user registry for storing user context. Additionally, the server utilizes scheduled tasks to maintain and flush user records daily.

## Features
- **SSE Support**: Real-time communication with clients using Server-Sent Events.
- **Custom Message Structs**: Optimized data handling with structured message formats.
- **User Registry**: Dedicated storage for user contexts, improving interaction management.
- **Scheduled Tasks**: Automated daily tasks to maintain and update user records, ensuring data integrity and performance.

## Prerequisites
Before you start, ensure you have the following installed:
- Go (version 1.15 or later)
- Echo Framework
- Cron package for Go

## Installation
To get the server running locally, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/theDewill/SSEbridge.git
   cd golang-sse-server
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the project:**
   ```bash
   go build -o sse-server
   ```

4. **Run the server:**
   ```bash
   ./sse-server
   ```

## Usage
After starting the server, it will listen on port 8080. You can connect to the SSE endpoint at:

```
http://localhost:8080/startSSE
```

Clients connected to this endpoint will receive real-time events sent from the server.

## Scheduled Tasks
The server uses a cron job to call the `Flush()` method on the `user_registry` every day at 12:00 PM. This task ensures that user records are consistently up-to-date and reduces memory usage over time.

## Contributing
Contributions are welcome! Feel free to fork the repository and submit pull requests.

1. **Fork the Project**
2. **Create your Feature Branch** (`git checkout -b feature/AmazingFeature`)
3. **Commit your Changes** (`git commit -m 'Add some AmazingFeature'`)
4. **Push to the Branch** (`git push origin feature/AmazingFeature`)
5. **Open a Pull Request**

## License
Distributed under the MIT License. See `LICENSE` for more information.

## Contact
Nomin Sendinu - nsendinu@gmail.com

Project Link: [https://github.com/theDewill/SSEbridge](https://github.com/theDewill/SSEbridge)

---
