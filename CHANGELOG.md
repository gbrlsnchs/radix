# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.4.5] - 2018-03-02
### Fixed
- Recursively locking for reading in the same goroutine when calling `(*Tree).String`.

## [0.4.4] - 2018-03-01
### Fixed
- `TestRace` now waits all goroutines to finish.

## [0.4.3] - 2018-03-01
### Fixed
- Bounds out of range error when trying to get a node with placeholder.
- Fix concurrency inconsistencies.

## [0.4.2] - 2018-01-24
### Added
- Race condition tests for all operations.

### Fixed
- Race condition when retrieving a node.

## [0.4.1] - 2017-12-17
### Fixed
- Adding an empty string does not enters an infinite loop anymore.
- Getting a node when the tree is empty no longer returns the root.

## [0.4.0] - 2017-12-11
### Added
- Makefile.
- Thread safety (if enabled).
- Thread safety test and example.

### Changed
- Travis CI config script.

### Removed
- Travis CI script for running goimports.

## [0.3.1] - 2017-12-10
### Added
- List of package features in README file.

### Changed
- Package name from `patricia` to `radix`.

## [0.3.0] - 2017-12-04
### Added
- Alphabetical sorting.

### Changed
- Sorting is not made automatically anymore.
- Characters in the printed tree.
- README file.

### Fixed
- CHANGELOG code citations are shown as code instead of ordinary text.
- README typo.
- Tests are now run in the correct order.
- `(*Tree).Add` now correctly adds a node as a prefix.

### Removed
- `(*Tree).WithNode` and `(*Tree).WithoutNode`.

## [0.2.1] - 2017-11-11
### Added
- Benchmark flag.

### Fixed
- `goimports` installation missing.

## [0.2.0] - 2017-11-08
### Added
- Specific methods for inserting and deleting a node and also return the tree.

### Changed
- Enhance basic operations' algorithms.
- Remove basic operations' return value.

## [0.1.3] - 2017-11-01
### Added
- Some use cases to the tests table.

### Changed
- Variable name `Node.Val` to `Node.Value`.

### Fixed
- GetByRune method returning a non-nil map even when not matching the string.

## [0.1.2] - 2017-11-01
### Changed
- Update this file to use "changelog" in lieu of "change log".

### Fixed
- Out of range bug when dynamically looking up a string.

## [0.1.1] - 2017-10-31
### Changed
- README file.

### Fixed
- Output examples.

## 0.1.0 - 2017-10-31
### Added
- This changelog file.
- README file.
- MIT License.
- Travis CI configuration file and scripts.
- Git ignore file.
- Editorconfig file.
- This package's source code, including examples and tests.
- Go dep files.

[0.4.5]: https://github.com/gbrlsnchs/radix/compare/v0.4.4...v0.4.5
[0.4.4]: https://github.com/gbrlsnchs/radix/compare/v0.4.3...v0.4.4
[0.4.3]: https://github.com/gbrlsnchs/radix/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/gbrlsnchs/radix/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/gbrlsnchs/radix/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/gbrlsnchs/radix/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/gbrlsnchs/radix/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/gbrlsnchs/radix/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/gbrlsnchs/radix/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/gbrlsnchs/radix/compare/v0.1.3...v0.2.0
[0.1.3]: https://github.com/gbrlsnchs/radix/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/gbrlsnchs/radix/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/gbrlsnchs/radix/compare/v0.1.0...v0.1.1
