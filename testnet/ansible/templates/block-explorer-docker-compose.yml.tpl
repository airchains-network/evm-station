version: '3.9'

services:
  redis_db:
    extends:
      file: ./services/redis.yml
      service: redis_db
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  db-init:
    extends:
      file: ./services/db.yml
      service: db-init
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  db:
    extends:
      file: ./services/db.yml
      service: db
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  backend:
    depends_on:
      - db
      - redis_db
    extends:
      file: ./services/backend.yml
      service: backend
    build:
      context: ..
      dockerfile: ./docker/Dockerfile
      args:
        CACHE_EXCHANGE_RATES_PERIOD: ""
        API_V1_READ_METHODS_DISABLED: "false"
        DISABLE_WEBAPP: "false"
        API_V1_WRITE_METHODS_DISABLED: "false"
        CACHE_TOTAL_GAS_USAGE_COUNTER_ENABLED: ""
        CACHE_ADDRESS_WITH_BALANCES_UPDATE_INTERVAL: ""
        ADMIN_PANEL_ENABLED: ""
        RELEASE_VERSION: 5.3.3
    links:
      - db:database
    environment:
        ETHEREUM_JSONRPC_HTTP_URL: http://proxy:8545/
        ETHEREUM_JSONRPC_TRACE_URL: http://proxy:8545/
        ETHEREUM_JSONRPC_WS_URL: ws://proxy:8546/
        CHAIN_ID: '{{ numeric_chain_id }}'
        INDEXER_TOKEN_BALANCES_BATCH_SIZE: 10
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  visualizer:
    extends:
      file: ./services/visualizer.yml
      service: visualizer
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  sig-provider:
    extends:
      file: ./services/sig-provider.yml
      service: sig-provider
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  frontend:
    depends_on:
      - backend
    extends:
      file: ./services/frontend.yml
      service: frontend
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  stats-db-init:
    extends:
      file: ./services/stats.yml
      service: stats-db-init
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  stats-db:
    depends_on:
      - backend
    extends:
      file: ./services/stats.yml
      service: stats-db
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  stats:
    depends_on:
      - stats-db
    extends:
      file: ./services/stats.yml
      service: stats
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

  proxy:
    depends_on:
      - backend
      - frontend
      - stats
    extends:
      file: ./services/nginx.yml
      service: proxy
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"
