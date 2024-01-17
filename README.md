# PriceTracker

PriceTracker is a tool for tracking prices on various online stores and monitoring the price fluctuations of products over time.

This repository also includes integrations with stores Pingo Doce and Worten, communicating via RabbitMQ.

The RabbitMQ integration is utilized for the PriceTracker catalog service to receive product updates and store registrations. 

The catalog service persists store and product information using PostgreSQL.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Build and Usage](#build-and-usage)
- [Development](#development)
- [RabbitMQ Integration](#rabbitmq-integration)
- [PostgreSQL Integration](#postgresql-integration)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Store Price Tracking:** Monitor prices from Pingo Doce and Worten.
- **Product History:** View historical price data to identify trends.
- **RabbitMQ Integration:** Communicate with the PriceTracker catalog for product updates and store registration.
- **PostgreSQL Integration:** Persist store and product information.

## Installation

### Prerequisites
- RabbitMQ server running
- PostgreSQL database running

Both of this applications can be run containerized for development purposes by
using the included `docker-compose.yaml`.

## Build and Usage


```bash
make build
```

To compile for a specific platform, use the following syntax:

```bash
make build system=linux arch=amd64
```

All compiled binaries will be stored in the `dist` directory.

1. **Build the Project:**
   
   In order to build the catalog and stores, ensure you have Golang 1.21 or later installed.

   To build for all more common platforms run:

   ```bash
   make build
   ```
   
    If you which to target only one platform you can do it like so:
    ```bash
    make build system=linux arch=amd64
    ```

    All compiled binaries will be stored in the `dist/$system_$arch` directory.


2. **Configure Catalog:**

   Adapt the `config.yaml` file in the `dist/$system_$arch/catalog` directory to have the correct credentials for both RabbitMQ and PostgreSQL. 
   If you are using the included `docker-compose.yaml` file, you can leave it as is.


3. **Run Catalog:**

   Execute the catalog:

   ```bash
   ./dist/$system_$arch/catalog
   ```

4. **Configure Stores:**

   Go to the desired store directories (e.g., `pingodoce`, `worten`) and ensure the configurations in the `config.yaml` file for the RabbitMQ connection match the ones used in the `catalog`.

5. **Run Stores:**

   Execute the stores from their respective directories:

   ```bash
   ./dist/$system_$arch/store/pingodoce/pingodoce
   ```

   ```bash
   ./dist/$system_$arch/store/worten/worten
   ```

6. **Access Catalog Frontend:**

   The `catalog` frontend is accessible on port `8080`.

   Open your web browser and navigate to [http://localhost:8080](http://localhost:8080) to access the catalog frontend.

**Note:** Ensure that RabbitMQ and PostgreSQL are running and correctly configured.

## Development

### For Catalog:

1. **Go to the catalog directory:**

    ```bash
    cd dist/$system_$arch/catalog
    ```

2. **Create and execute docker containers:**

    ```bash
    docker-compose up
    ```

3. **Install dev_dependencies:**

    ```bash
    make dev_dependencies
    ```

4. **Ensure Golang's bin directory is in your `$PATH`.**

5. **Run using `air` for code hot-reloading:**

    ```bash
    make run
    ```

### For Stores:

1. **Ensure the catalog is running.**

2. **Run the store like any other Golang application:**

    ```bash
    go run main.go
    ```

## RabbitMQ Integration

The RabbitMQ integration is crucial for the PriceTracker catalog service. It facilitates communication for receiving product updates and store registrations. Ensure the RabbitMQ server is running and correctly configured. Update the `config.yaml` file in the `catalog` directory with the RabbitMQ connection details.

## PostgreSQL Integration

The catalog service uses PostgreSQL to persist store and product information. Make sure the PostgreSQL database is running and properly configured. Update the `config.yaml` file in the `catalog` directory with the PostgreSQL connection details.


## Contributing

If you'd like to contribute to PriceTracker, please follow these guidelines:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Make your changes and submit a pull request.
- Provide a detailed description of your changes.

## License

This project is licensed under the [GNU General Public License (GPL)](LICENSE).