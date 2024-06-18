
# Engineering Blog Web Crawler

[![Go](https://img.shields.io/badge/go-1.16-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/yourusername/engineering-blog-web-crawler.svg?branch=main)](https://travis-ci.org/yourusername/engineering-blog-web-crawler)

## Description
This project is a web crawler written in Golang designed to scrape and process web content from engineering blogs efficiently. The crawler can be configured to follow links, extract information, and store the results in a structured format.

## Features
- Concurrent crawling with Goroutines.
- Configurable crawling depth and rate limiting.
- Parsing and extraction using Goquery.
- Support for various output formats (JSON, CSV, etc.).

## Setup and Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/engineering-blog-web-crawler.git
    cd engineering-blog-web-crawler
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    go build -o crawler cmd/main.go
    ```

4. Run the application:
    ```bash
    ./crawler -url=https://example.com -depth=3
    ```

## Usage
Specify the URL and the crawling depth:
```bash
./crawler -url=https://example.com -depth=3
```

## Project Structure
- `cmd/` - Contains the main entry point of the application.
- `internal/` - Contains the core functionality and logic of the crawler.
- `pkg/` - Contains shared libraries and utilities.
- `main.go` - The main file to run the crawler.

## Contribution
Contributions are welcome! Please open an issue or submit a pull request for any changes or improvements.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact
For any questions or inquiries, please contact [Ivan Kwong](mailto:ivankwong22@gmail.com).
