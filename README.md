# 🚖 GoRide

GoRide is a simple ride-hailing backend built in Go with **MongoDB** and **Kafka**.  
It consists of two microservices:

- **booking_svc** → handles ride bookings  
- **driver_svc** → handles drivers & job acceptance  

---

## ⚡ Setup

### 1. Prerequisites
- Docker & Docker Compose installed
- Go 1.24+ installed (for local builds)

### 2. Project Structure
GoRide/
├── booking_svc/ # booking service
├── driver_svc/ # driver service
├── deploy/ # docker-compose.yml
└── README.md


### 3. Run with Docker Compose
```bash 
cd deploy
docker compose up --build -d
```

### 4. Verify running containers
```bash
docker ps
```

You should see:
- booking_svc (8080)
- driver_svc (8081)
- mongodb (27017)
- kafka (9092)
- zookeeper (2181)

## 🔧 Environment Variables
| Service      | Variable       | Example                   | Description            |
| ------------ | -------------- | ------------------------- | ---------------------- |
| booking\_svc | `MONGO_URI`    | `mongodb://mongodb:27017` | MongoDB connection URI |
| booking\_svc | `KAFKA_BROKER` | `kafka:9092`              | Kafka broker address   |
| driver\_svc  | `MONGO_URI`    | `mongodb://mongodb:27017` | MongoDB connection URI |
| driver\_svc  | `KAFKA_BROKER` | `kafka:9092`              | Kafka broker address   |

## 📌 REST API Examples

# Booking Service (port 8080)

### ➕ Create booking
```bash
curl -X POST localhost:8080/bookings \
-H "Content-Type: application/json" \
-d '{"pickuploc":{"lat":12.9,"lng":77.6},"dropoff":{"lat":12.95,"lng":77.64},"price":220}'
```

### 📋 List bookings
```bash
curl localhost:8080/bookings
```

### 🔍 Get booking by ID (optional feature)
```bash
curl localhost:8080/bookings/<booking_id>
```

# Driver Service (port 8081)

### 👨‍✈️ List drivers
```bash
curl localhost:8081/drivers
```

### 📋 List jobs
```bash
curl localhost:8081/jobs
```

### ✅ Accept job
```bash
curl -X POST localhost:8081/jobs/<booking_id>/accept \
-H "Content-Type: application/json" \
-d '{"driver_id":"d-1"}'
```

## 🔄 Event Flow

- POST /bookings → booking created, saved in MongoDB, event booking.created produced to Kafka

- driver_svc consumes booking.created and shows jobs in /jobs

- Driver accepts via POST /jobs/{id}/accept → produces booking.accepted event

- booking_svc consumes booking.accepted → updates booking in MongoDB with driver_id + ride_status="Accepted"###

# 🛠 Development Notes

- Services are written in Go with chi (router), mongo-driver, and kafka-go
- Multi-stage Docker builds for small, secure images
- MongoDB + Kafka provided via Docker Compose