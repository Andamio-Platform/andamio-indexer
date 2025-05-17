# Andamio Indexer

## Project Overview

The Andamio Indexer is a service designed to efficiently index blockchain data, making it readily available and easily searchable through a REST API. It supports the indexing of various data types, including blocks, transactions, addresses, and UTXOs. With a modular and extensible architecture, the indexer can be adapted to integrate with different blockchain networks and data sources.

## Features

*   **Blockchain Data Indexing:** Indexes various blockchain data types, including blocks, transactions, addresses, and UTXOs.
*   **Modular Design:** Designed to be modular and extensible, allowing it to be adapted to different blockchain networks and data sources.
*   **Configurable:** Supports various configuration options to customize the indexing process.
*   **REST API:** Provides a REST API for accessing indexed data.
*   **SQLite and BadgerDB Support:** Uses SQLite for metadata storage and BadgerDB for blob storage.

## Dependencies and Installation

1.  **Prerequisites:**
    *   Go 1.20 or higher
    *   SQLite
    *   BadgerDB

2.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd andamio-indexer
    ```

3.  **Install Go Modules:**

    ```bash
    go mod download
    ```

## Building and Running

1.  **Build the project:**

    ```bash
    go build -o ./build/andamio-indexer main.go
    ```

2.  **Configuration:**

    *   Create a `config.json` file with the desired configuration options.
    *   Refer to the `config/config.go` file for available configuration options.

3.  **Run the indexer:**

    ```bash
    ./build/andamio-indexer -config config/config.json
    ```

## Architecture and Key Components

The Andamio Indexer follows a modular architecture with key components responsible for different aspects of the indexing process. The main components include:

*   **API Handlers:** Located in the `handlers/v1` directory, these handle incoming REST API requests and interact with the database to retrieve indexed data.
*   **Indexer Core:** Responsible for receiving and processing blockchain transaction events. This includes filtering relevant events, batching transactions, and coordinating the storage of data.
*   **Database Interaction:** The indexer utilizes a database layer to store the indexed blockchain data. It supports SQLite for metadata and BadgerDB for blob data.

For a detailed explanation of the indexer's internal data flow from receiving transactions to storing them, please refer to the [Indexer Internal Workings Documentation](docs/indexer_internals.md).

## API Documentation

For detailed information on the available API endpoints, request parameters, and responses, please refer to the [API Documentation](docs/api.md).

## Contributing

Contributions are welcome! Please follow these guidelines:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and commit them with clear and concise messages.
4.  Submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.