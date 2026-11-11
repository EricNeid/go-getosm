<!-- markdownlint-disable MD022 MD032 MD024-->
# Changelog

All notable changes to this project will be documented in this file.

## Unreleased
### Fixed
* Show error message when no bounding box is given

## 0.10.0 - 2026-11-06
### Changed
* Improve naming of output files when being tiled is used
* Add handling for multiple bounding boxes
* Print overall bounding box string as wkt
### Removed (Breaking)
* Removed option -t, use -tx instead

## 0.9.0 - 2026-05-08
### Changed
* Minor improvements in code
### Added
* Argumentes tx and ty to split into columns, rows (and grid)

## 0.8.0 - 2026-04-30
### Added
* Option to configure custom header
### Changed
* Set element limit to 500000000

## 0.7.0 - 2026-01-05
### Added
* Option to configure overpass api url

### 0.6.0 - 2025-11-04
### Added
* Option to skip already downloaded tiles

## 0.5.0 - 2025-10-29
### Added
* Option retryDelaySec to set the delay between 2 retries

## 0.4.0 - 2025-10-29
### Added
* Option to adjust retry limit

## 0.3.2 - 2025-08-29
### Fixed
* Naming of tile files starts with 1
### Removed
* Argument q to print query, because it is obsolete by verbose
### Added
* Print generated files after download
### Misc
* Cleanup

## 0.3.1 - 2025-08-26
### Fixed
* Query not completely printed
* Error in tiling
### Added
* Parameter q to print query

## 0.3.0 - 2025-08-25
### Added
* Option to tile large files
* Automatic retry if error occurs during download
* Options to disable/enable verbose logging
### Fixed
* Wrong description in help page

## 0.2.0 - 2023-01-05
### Added
* Additional log messages for debugging
### Changed
* Disable timeout

## 0.1.0 - 2021-09-29
### Added
* Initial release
