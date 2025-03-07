# Change Log

## [v0.0.5] - 08.03.2025
### Changed
* Added support last version of 
  * [lib-tinyerrors](https://github.com/crypto-bundle/bc-wallet-common-lib-tinyerrors)
  * [lib-errors](https://github.com/crypto-bundle/bc-wallet-common-lib-errors)
* Added support of Go 1.23
* Added linter and fixed up all linter issues
* Updated License
  * Copyright - new year 2025
  * MIT -> MIT NON-AI
  * Added License banner to *.go files

## [v0.0.4] - 16.04.2024
### Changed
* Bump golang version 1.19 -> 1.22

## [v0.0.3] - 15.04.2024
### Changed
* Removed usage of opentracing. Moved to (go.opentelemetry.io)[https://github.com/open-telemetry/opentelemetry-go-contrib]

## [v0.0.2] - 25.04.2023 14:00 MSK
### Added
* Added gRPC-client roundrobin picker for client-side balancing cryptobundle profile

## [v0.0.1] - 24.03.2023 00:10 MSK
### Changed
* Lib-grpc moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-grpc
    * Removed other packages of old bc-wallet-common repository
* Added MIT license
### Fixed
* Fixed wrong path of dependencies - gRPC dns package