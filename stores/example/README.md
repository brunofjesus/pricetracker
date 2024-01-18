# ExampleStore

This Go application is a PriceTracker project that includes functionality for tracking prices of products in an example store.

It utilizes a message queue for communication, allowing for real-time updates on product prices.

## Files

- `config.yaml`: Application configuration file containing the MessageQueue connection string and loop interval settings.

- `config/type.go`: Go file defining the schema for `config.yaml` as a `ApplicationConfiguration` struct.

- `config/loader.go`: Go file responsible for loading the `config.yaml` file into the `ApplicationConfiguration` struct.

- `pkg/crawler/products.go`: Go file containing three product instances representing the products found in the example store.

- `pkg/crawler/crawl.go`: Go file responsible for fetching products, applying random price variations, and calling the function to publish the products.

- `main.go`: Entry point of the application. It loads configurations, establishes a connection to the Message Queue, registers the store, and starts the product update loop.

## Configuration

The `config.yaml` file should be configured with the appropriate values for the MessageQueue connection and loop interval.

If you are running all locally and using the included `docker-compose.yml` from the `catalog` you can leave it as is.


## Build and Usage


```bash
make build
```

To compile for a specific platform, use the following syntax:

```bash
make build system=linux arch=amd64
```

All compiled binaries will be stored in the `dist/$system_arch/store/example` directory.

1. **Build the Project:**

   In order to build the store, ensure you have Golang 1.21 or later installed.

   To build for all more common platforms run:

   ```bash
   make build
   ```

   If you which to target only one platform you can do it like so:
    ```bash
    make build system=linux arch=amd64
    ```

   All compiled binaries will be stored in the `dist/$system_arch/store/example` directory.


2. **Execute**

   The application can now be executed, make sure that the `catalog` and it's dependencies are running first.