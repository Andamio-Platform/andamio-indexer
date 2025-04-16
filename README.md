
# Andamio Indexer

## Description

The Andamio Indexer is a service that indexes blockchain data, making it easily accessible and searchable. It supports various data types, including blocks, transactions, addresses, and UTXOs. The indexer is designed to be modular and extensible, allowing it to be adapted to different blockchain networks and data sources.

## Features

*   **Blockchain Data Indexing:** Indexes various blockchain data types, including blocks, transactions, addresses, and UTXOs.
*   **Modular Design:** Designed to be modular and extensible, allowing it to be adapted to different blockchain networks and data sources.
*   **Configurable:** Supports various configuration options to customize the indexing process.
*   **REST API:** Provides a REST API for accessing indexed data.
*   **SQLite and BadgerDB Support:** Uses SQLite for metadata storage and BadgerDB for blob storage.

## Installation

1.  **Prerequisites:**
    *   Go 1.20 or higher
    *   SQLite
    *   BadgerDB

2.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd andamio-indexer
    ```

3.  **Build the project:**

    ```bash
    go build -o build/andamio-indexer main.go
    ```

## Usage

1.  **Configuration:**

    *   Create a `config.json` file with the desired configuration options.
    *   Refer to the `config/config.go` file for available configuration options.

2.  **Run the indexer:**

    ```bash
    ./build/andamio-indexer --config config/config.json
    ```

## Configuration

The Andamio Indexer can be configured using a `config.json` file. The following configuration options are available:

*   `network`: network configuration options, such as magic and socketPath.
*   `Database`: Database configuration options, such as the database file path.
*   `andamio`: andamio configuration options, such as PolicyID and RefTx.

Refer to the `config/config.go` file for detailed information on each configuration option.

## Contributing

Contributions are welcome! Please follow these guidelines:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and commit them with clear and concise messages.
4.  Submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.