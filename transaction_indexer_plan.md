# Andamio Transaction Indexer Plan

This document outlines the plan for designing the database and API for the Andamio transaction indexer, extending its capabilities to include comprehensive transaction details from the Cardano blockchain and providing a rich set of API endpoints for frontend consumption.

## 1. Database Design

The existing database schema will be extended to include detailed transaction information. This will involve adding new tables and modifying existing ones to store data such as block information, full transaction details, input/output specifics (including datum and redeemer references), asset information, and witnesses.

Here is the proposed database schema:

```mermaid
erDiagram
    blocks {
        INTEGER id PK
        TEXT block_hash UK "Indexed"
        INTEGER block_number "Indexed"
        INTEGER slot "Indexed"
        DATETIME block_time "Indexed"
        DATETIME created_at
        DATETIME updated_at
    }

    transactions {
        INTEGER id PK
        TEXT tx_hash UK "Indexed"
        INTEGER block_id FK "Indexed"
        INTEGER slot "Indexed"
        INTEGER fee
        INTEGER validity_interval_start
        INTEGER validity_interval_end
        TEXT metadata "JSON or BLOB"
        DATETIME created_at
        DATETIME updated_at
    }

    transaction_inputs {
        INTEGER id PK
        INTEGER transaction_id FK "Indexed"
        INTEGER utxo_id FK "Indexed"
        INTEGER redeemer_id FK "Optional, Indexed"
        INTEGER input_index
        TEXT script_hash "Optional, Indexed"
    }

    transaction_outputs {
        INTEGER id PK
        INTEGER transaction_id FK "Indexed"
        INTEGER utxo_id FK "Indexed"
        INTEGER address_id FK "Indexed"
        INTEGER datum_id FK "Optional, Indexed"
        INTEGER output_index
        TEXT script_hash "Optional, Indexed"
    }

    utxos {
        INTEGER ID PK
        BLOB TxId "Indexed: tx_id_output_idx"
        INTEGER OutputIdx "Indexed: tx_id_output_idx"
        INTEGER AddedSlot "Indexed"
        INTEGER DeletedSlot "Indexed"
        BLOB PaymentKey "Indexed"
        BLOB StakingKey "Indexed"
    }

    addresses {
        INTEGER ID PK
        TEXT Address UK "Indexed"
        DATETIME CreatedAt
        DATETIME UpdatedAt
    }

    assets {
        INTEGER id PK
        TEXT policy_id "Indexed"
        TEXT asset_name "Indexed"
        TEXT fingerprint UK "Indexed"
        DATETIME created_at
        DATETIME updated_at
    }

    utxo_assets {
        INTEGER utxo_id FK "Indexed"
        INTEGER asset_id FK "Indexed"
        INTEGER quantity
        PRIMARY KEY (utxo_id, asset_id)
    }

    datum {
        INTEGER id PK
        TEXT datum_hash UK "Indexed"
        BLOB datum_cbor
    }

    redeemers {
        INTEGER id PK
        TEXT redeemer_hash UK "Indexed"
        BLOB redeemer_cbor
    }

    witnesses {
        INTEGER id PK
        INTEGER transaction_id FK "Indexed"
        BLOB witness_data
    }

    blocks ||--o{ transactions : "contains"
    transactions ||--o{ transaction_inputs : "has"
    transactions ||--o{ transaction_outputs : "has"
    transaction_inputs }o--|| utxos : "spends"
    transaction_inputs }o--o| redeemers : "uses"
    transaction_outputs }o--|| utxos : "creates"
    transaction_outputs }o--|| addresses : "to"
    transaction_outputs }o--o| datum : "with"
    utxos ||--o{ utxo_assets : "has"
    assets ||--o{ utxo_assets : "is included in"
    transactions ||--o{ witnesses : "has"
```

## 2. API Design

The following API endpoints will be implemented to support the required query capabilities:

*   `GET /api/v1/transactions/{tx_hash}`: Retrieve a single transaction by its hash, including all its details (inputs, outputs, metadata, witnesses, associated datum and redeemer content).
*   `GET /api/v1/addresses/{address}/transactions`: Retrieve a list of transactions associated with a specific address (either as an input or an output). This endpoint will support pagination and filtering by time range.
*   `GET /api/v1/assets/policy/{policyId}/transactions`: Retrieve a list of transactions that involve assets with a specific policy ID. This endpoint will support pagination and filtering by time range.
*   `GET /api/v1/assets/token/{tokenname}/transactions`: Retrieve a list of transactions that involve assets with a specific token name. This endpoint will support pagination and filtering by time range.
*   `GET /api/v1/assets/fingerprint/{asset_fingerprint}/transactions`: Retrieve a list of transactions that involve a specific asset identified by its fingerprint. This endpoint will support pagination and filtering by time range.
*   `GET /api/v1/assets/policy/{policyId}/token/{tokenname}/transactions`: Retrieve a list of transactions that involve a specific asset identified by its policy ID and token name. This endpoint will support pagination and filtering by time range.
*   `GET /api/v1/addresses/{address}/utxos`: Retrieve the list of unspent transaction outputs (UTXOs) for a given address. (Enhancement of existing functionality).
*   `GET /api/v1/transactions/{tx_hash}/utxos`: Retrieve the list of UTXOs created by a specific transaction.