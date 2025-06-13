
<h3 align="center">RandomTalk</h3>
<p align="center"> RandomTalk is a simple chat application that connects users to random people to chat with based on their preferences.</p>
<p align="center">
  <a href="https://randomtalk.chat">View Live Demo</a>
  <br/><br/>
  <a href="https://github.com/xfrr/randomtalk/actions/workflows/go.yml">
    <img src="https://github.com/xfrr/randomtalk/actions/workflows/go.yml/badge.svg" alt="Go CI"/>
  </a>
</p>

---

## 🚀 Features <a name = "features"></a>

- In-house **Gale-Shapley algorithm** for matching users based on their preferences.
- Real-time chat using **WebSockets**.
- Monitoring and observability using **Grafana** and **Prometheus**.
- High-concurrency management using Optimistic Locking and NATS JetStream capabilities.
- Scalability using **Kubernetes** and **Helm**.

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Cloning the repository

```bash
git clone https://github.com/xfrr/randomtalk.git
```

### Prerequisites

This project requires the following tools to be installed:

- [Justfile](https://github.com/casey/just) - Task runner.
- [Docker](https://www.docker.com/) - Containerization platform.

### Start the application

To start the application, run the following command:

```bash
just up <streaming_system>
```

For example, to start the application with NATS Jetstream as the messaging and stream processing system, run the following command:

```bash
just up nats
```

> This command will start the application using Docker-Compose.

The following streaming systems are supported:

- `nats` - Start the application with NATS Jetstream as the messaging and stream processing system.
- More coming soon.

### Access the application

Once the application is started, you can access the application at the following URLs:

- **Web Chat Application**: [http://localhost:8081](http://localhost:8081)
- **Grafana**: [http://localhost:3000](http://localhost:3000)
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **NATS UI**: [http://localhost:8222](http://localhost:31311)

## ⛏️ Built Using <a name = "built_using"></a>

- [Go](https://golang.org/) - Programming language.
- [NATS](https://nats.io/) - Messaging and Stream processing system.
- [NATS UI](https://github.com/nats-nui/nui) - NATS Web UI.
- [Grafana](https://grafana.com/) - Monitoring and observability platform.
- [Prometheus](https://prometheus.io/) - Monitoring and alerting toolkit.
- [Docker](https://www.docker.com/) - Containerization platform.
- [Kubernetes](https://kubernetes.io/) - Container orchestration platform.
- [Helm](https://helm.sh/) - Kubernetes package manager.

## 📜 License <a name = "license"></a>

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
