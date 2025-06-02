# API Documentation

This document provides a human-readable overview of the Andamio Indexer API, based on the OpenAPI (Swagger) specification.

## Base Path

`/api/v1/indexer`

## Authentication

API requests are secured using `ApiKeyAuth`.

## Endpoints

### Addresses

#### Add Address

Adds a new address to the indexer for monitoring transactions and UTxOs.

*   **URL:** `/addresses`
*   **Method:** `POST`
*   **Description:** Adds a new address to the indexer for monitoring transactions and UTxOs.
*   **Request Body:**
    *   `address` (required): The address object containing the address string to be added. Refer to `viewmodel.AddressRequest` schema.
*   **Responses:**
    *   `201 Created`: Successfully added address.
        *   Schema:
            ```json
            {
              "message": "string"
            }
            ```
    *   `400 Bad Request`: Invalid request body or missing address.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Remove Address

Removes an address from the tracking list.

*   **URL:** `/addresses/remove-address`
*   **Method:** `DELETE`
*   **Description:** Removes an address from the tracking list.
*   **Request Body:**
    *   `addressRequest` (required): The address object containing the address string to be removed. Refer to `viewmodel.AddressRequest` schema.
*   **Responses:**
    *   `200 OK`: Successfully removed address.
        *   Schema:
            ```json
            {
              "message": "string"
            }
            ```
    *   `400 Bad Request`: Invalid request body or missing address.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Assets by Address

Retrieve a list of assets held at a specific address, with support for pagination.

*   **URL:** `/addresses/{address}/assets`
*   **Method:** `GET`
*   **Description:** Retrieve a list of assets held at a specific address, with support for pagination.
*   **Parameters:**
    *   `address` (required, path): The address to retrieve assets for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved assets.
        *   Schema: Array of `viewmodel.Asset`
    *   `400 Bad Request`: Invalid address or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Address not found or no assets found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Address

Retrieves transactions associated with a specific address with pagination.

*   **URL:** `/addresses/{address}/transactions`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific address with pagination.
*   **Parameters:**
    *   `address` (required, path): The address to retrieve transactions for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid address or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Address not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

### Assets

#### Get Addresses by Asset Fingerprint

Retrieve a list of addresses that hold a specific asset fingerprint, with support for pagination.

*   **URL:** `/assets/fingerprint/{asset_fingerprint}/addresses`
*   **Method:** `GET`
*   **Description:** Retrieve a list of addresses that hold a specific asset fingerprint, with support for pagination.
*   **Parameters:**
    *   `asset_fingerprint` (required, path): The asset fingerprint (hex-encoded) to retrieve addresses for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved addresses.
        *   Schema: Array of strings
    *   `400 Bad Request`: Invalid asset fingerprint or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Asset fingerprint not found or no addresses found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get UTxOs by Asset Fingerprint

Retrieves UTxOs associated with a specific asset fingerprint with pagination.

*   **URL:** `/assets/fingerprint/{asset_fingerprint}/utxos`
*   **Method:** `GET`
*   **Description:** Retrieves UTxOs associated with a specific asset fingerprint with pagination.
*   **Parameters:**
    *   `asset_fingerprint` (required, path): The asset fingerprint (hex-encoded) to retrieve UTxOs for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved UTxOs.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid asset fingerprint or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Asset fingerprint not found or no UTxOs found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Asset Fingerprint

Retrieves transactions associated with a specific asset fingerprint with pagination.

*   **URL:** `/assets/fingerprint/{asset_fingerprint}/transactions`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific asset fingerprint with pagination.
*   **Parameters:**
    *   `asset_fingerprint` (required, path): The asset fingerprint (hex-encoded) to retrieve transactions for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid asset fingerprint or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Asset fingerprint not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Policy ID and Token Name

Retrieves transactions associated with a specific policy ID and token name with pagination.

*   **URL:** `/assets/policy/{policyId}/token/{tokenname}/transactions`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific policy ID and token name with pagination.
*   **Parameters:**
    *   `policyId` (required, path): The policy ID (hex-encoded) to retrieve transactions for. (string)
    *   `tokenname` (required, path): The token name (hex-encoded) to retrieve transactions for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid policy ID, token name, or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Policy ID and token name not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Policy ID

Retrieves transactions associated with a specific policy ID with pagination.

*   **URL:** `/assets/policy/{policyId}/transactions`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific policy ID with pagination.
*   **Parameters:**
    *   `policyId` (required, path): The policy ID (hex-encoded) to retrieve transactions for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid policy ID or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Policy ID not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Token Name

Retrieves transactions associated with a specific token name with pagination.

*   **URL:** `/assets/token/{tokenname}/transactions`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific token name with pagination.
*   **Parameters:**
    *   `tokenname` (required, path): The token name (hex-encoded) to retrieve transactions for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid token name or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Token name not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

### Transactions

#### Get Transaction by Tx Hash

Retrieves a transaction by its hash.

*   **URL:** `/transactions/{tx_hash}`
*   **Method:** `GET`
*   **Description:** Retrieves a transaction by its hash.
*   **Parameters:**
    *   `tx_hash` (required, path): The hash of the transaction to retrieve. (string)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transaction.
        *   Schema: `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid transaction hash.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Transaction not found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Block Number

Retrieves transactions associated with a specific block number with pagination.

*   **URL:** `/transactions/by-block-number/{block_number}`
*   **Method:** `GET`
*   **Description:** Retrieves transactions associated with a specific block number with pagination.
*   **Parameters:**
    *   `block_number` (required, path): The block number to retrieve transactions for. (integer)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid block number or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Block number not found or no transactions found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions by Slot Range

Retrieves transactions within a specific slot range with pagination.

*   **URL:** `/transactions/by-slot-range`
*   **Method:** `GET`
*   **Description:** Retrieves transactions within a specific slot range with pagination.
*   **Parameters:**
    *   `from_slot` (required, query): The starting slot number. (integer)
    *   `to_slot` (required, query): The ending slot number. (integer)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions.
        *   Schema: Array of `viewmodel.Transaction`
    *   `400 Bad Request`: Invalid slot range or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: No transactions found in the specified slot range.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get UTxOs by Transaction

Retrieves UTxOs associated with a specific transaction hash.

*   **URL:** `/transactions/{tx_hash}/utxos`
*   **Method:** `GET`
*   **Description:** Retrieves UTxOs associated with a specific transaction hash.
*   **Parameters:**
    *   `tx_hash` (required, path): The hash of the transaction to retrieve UTxOs for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved UTxOs.
        *   Schema: Array of `viewmodel.SimpleUTxO`
    *   `400 Bad Request`: Invalid transaction hash or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Transaction not found or no UTxOs found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get UTxOs Inputs by Transaction

Retrieves transaction inputs associated with a specific transaction hash.

*   **URL:** `/transactions/{tx_hash}/utxos/inputs`
*   **Method:** `GET`
*   **Description:** Retrieves transaction inputs associated with a specific transaction hash.
*   **Parameters:**
    *   `tx_hash` (required, path): The hash of the transaction to retrieve inputs for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transaction inputs.
        *   Schema: Array of `viewmodel.TransactionInput`
    *   `400 Bad Request`: Invalid transaction hash or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Transaction not found or no inputs found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get UTxOs Outputs by Transaction

Retrieves transaction outputs associated with a specific transaction hash.

*   **URL:** `/transactions/{tx_hash}/utxos/outputs`
*   **Method:** `GET`
*   **Description:** Retrieves transaction outputs associated with a specific transaction hash.
*   **Parameters:**
    *   `tx_hash` (required, path): The hash of the transaction to retrieve outputs for. (string)
    *   `limit` (optional, query): Maximum number of results to return. (integer, default: 100)
    *   `offset` (optional, query): Number of results to skip. (integer, default: 0)
*   **Responses:**
    *   `200 OK`: Successfully retrieved transaction outputs.
        *   Schema: Array of `viewmodel.TransactionOutput`
    *   `400 Bad Request`: Invalid transaction hash or pagination parameters.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Transaction not found or no outputs found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

### Metrics

#### Get Addresses Count

Retrieves the total count of addresses in the indexer.

*   **URL:** `/metrics/addresses/count`
*   **Method:** `GET`
*   **Description:** Retrieves the total count of addresses in the indexer.
*   **Responses:**
    *   `200 OK`: Successfully retrieved addresses count.
        *   Schema:
            ```json
            {
              "count": "integer"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Assets Count

Retrieves the total count of assets in the indexer.

*   **URL:** `/metrics/assets/count`
*   **Method:** `GET`
*   **Description:** Retrieves the total count of assets in the indexer.
*   **Responses:**
    *   `200 OK`: Successfully retrieved assets count.
        *   Schema:
            ```json
            {
              "count": "integer"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Latest Block

Retrieves information about the latest block indexed.

*   **URL:** `/metrics/latest-block`
*   **Method:** `GET`
*   **Description:** Retrieves information about the latest block indexed.
*   **Responses:**
    *   `200 OK`: Successfully retrieved latest block information.
        *   Schema:
            ```json
            {
              "block_number": "integer",
              "slot": "integer",
              "hash": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

#### Get Transactions Count

Retrieves the total count of transactions in the indexer.

*   **URL:** `/metrics/transactions/count`
*   **Method:** `GET`
*   **Description:** Retrieves the total count of transactions in the indexer.
*   **Responses:**
    *   `200 OK`: Successfully retrieved transactions count.
        *   Schema:
            ```json
            {
              "count": "integer"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```

### Redeemers

#### Get Redeemer by Tx Hash

Retrieves a redeemer by its transaction hash.

*   **URL:** `/redeemers/{tx_hash}`
*   **Method:** `GET`
*   **Description:** Retrieves a redeemer by its transaction hash.
*   **Parameters:**
    *   `tx_hash` (required, path): The hash of the transaction to retrieve the redeemer for. (string)
*   **Responses:**
    *   `200 OK`: Successfully retrieved redeemer.
        *   Schema: `viewmodel.Redeemer`
    *   `400 Bad Request`: Invalid transaction hash.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `404 Not Found`: Redeemer not found.
        *   Schema:
            ```json
            {
              "error": "string"
            }
            ```
    *   `500 Internal Server Error`: Internal server error.
        *   Schema:
            ```json
            {
              "error": "string"
            }