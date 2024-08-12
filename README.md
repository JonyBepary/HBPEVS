
# Hybrid BlockChain-based E-Voting System (HBPEVS)

![GitHub](https://img.shields.io/github/license/JonyBepary/HBPEVS)
![GitHub last commit](https://img.shields.io/github/last-commit/JonyBepary/HBPEVS)
![GitHub issues](https://img.shields.io/github/issues/JonyBepary/HBPEVS)
![GitHub stars](https://img.shields.io/github/stars/JonyBepary/HBPEVS)
![GitHub forks](https://img.shields.io/github/forks/JonyBepary/HBPEVS)

<img src="https://your-icon-url.com/icon.png" alt="HBPEVS Icon" width="100"/>

## Overview

The Hybrid BlockChain-based E-Voting System (HBPEVS) is a revolutionary voting system designed to be transparent, permissioned, append-only, and decentralized. This system allows people to vote easily with little to no assistance and verify their votes anonymously without revealing their identity. The system is tamper-proof and does not depend on a single point of failure. It can coexist with the current electoral system and facilitates a complete power structure to prevent any unauthorized activity.

HBPEVS combines the higher security, transparency, anonymity, and immutability of a public blockchain with the high efficiency and regulations of a private blockchain. This hybrid approach ensures compatibility with the public policy of the electoral system.

## Key Features

- **Transparency**: Every vote is recorded on the blockchain, ensuring full transparency and traceability.
- **Decentralized**: No single point of failure; the system is distributed across multiple nodes.
- **Permissioned**: Only authorized users can participate in the voting process.
- **Append-Only**: Once a vote is recorded, it cannot be altered or deleted.
- **Anonymity**: Voters can verify their votes without revealing their identity.
- **Tamper-Proof**: The blockchain ensures that votes cannot be tampered with.
- **High Efficiency**: The private blockchain component ensures efficient processing of votes.
- **Regulations**: Compliance with electoral regulations and public policy.

## Technology Stack

- **Blockchain**: Utilizes a hybrid blockchain approach combining public and private blockchain features.
- **LevelDB**: A fast key-value storage library for storing blockchain data.
- **Libp2p**: A peer-to-peer networking library for decentralized communication.
- **Gin**: A web framework for building the API endpoints.
- **Go**: The programming language used for the implementation.
- **Cryptography**: Uses SHA-256 for hashing and Ed25519 for digital signatures.
- **Base58**: Encoding format used for keys and hashes.
- **Testing**: Benchmarking and testing frameworks to ensure performance and reliability.

## Installation

### Prerequisites

- Go (version 1.16 or later)
- LevelDB
- Libp2p
- Gin

### Steps

1. **Clone the repository**:
   ```sh
   git clone https://github.com/JonyBepary/pqhbevs_hac.git
   cd pqhbevs_hac
   ```

2. **Install dependencies**:
   ```sh
   go mod tidy
   ```

3. **Run the application**:
   ```sh
   go run main.go
   ```

## Usage

### API Endpoints

- **Add Vote**: `POST /add_vote`
  - Adds a new vote to the blockchain.
  - Parameters: `vote`, `pscode`, `digest`, `po_id`, `po_privkey`, `verbose`

- **Generate Genesis Block**: `POST /generate_genesis`
  - Initializes the genesis block.
  - Parameters: `pscode`, `seed`, `start`, `duration`, `config`, `verbose`

### Benchmarking

- **NID Server Benchmark**:
  ```sh
  go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> nid_data.txt
  ```

- **Block Generation Rate Benchmark**:
  ```sh
  go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> newblock_data.txt -timeout 60m
  ```

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

<!-- LinkedIn, X, Mail -->
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/sohel-ahmed-jony/)
[![X](https://img.shields.io/badge/X-000000?style=for-the-badge&logo=x&logoColor=white)](https://twitter.com/jbepary)
[![Email](https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:sohelahmedjony@gmail.com)

