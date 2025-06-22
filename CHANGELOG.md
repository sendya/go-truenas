# Changelog

All notable changes to this project will be documented in this file.

## [0.1.3] 

### Fixed
- Fixed data race in WebSocket client between Close() and Call() methods
- Fixed race condition in connectionManager when sending to closed reconnectCh

## [0.1.2] 

### Fixed
- Handle unmarshaling PoolScan timestamps

## [0.1.0]

### Added
- Initial implementation of TrueNAS WebSocket API client
- Support for most of the TrueNAS API surface including:
  - Authentication (username/password and API key)
  - Pool management
  - Dataset operations
  - Service management
  - System information
  - Network configuration
  - File sharing (SMB, NFS, AFP, WebDAV)
  - User and group management
  - Certificate management
  - Job management
  - Alert system
  - Boot management
  - Disk operations
  - Filesystem operations
  - VM management
  - Smart monitoring
- Thread-safe WebSocket client with automatic reconnection
- Context-based timeout support
- Comprehensive test suite
- Type-safe API clients for all endpoints

[0.1.3]: https://github.com/715d/go-truenas/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/715d/go-truenas/compare/v0.1.1...v0.1.2
[0.1.0]: https://github.com/715d/go-truenas/releases/tag/v0.1.0
