# LoadModule App

This is a simple Go application that sends HTTP requests to a specified URL at a specified rate for a specified amount of time

## Prerequisites

- Go 1.22.x 

## Installation

1. Clone the repository: `git clone https://github.com/your-username/loadmodule.git`
2. Navigate to the project directory: `cd loadmodule`

## Usage

1. Build the application: `ddocker build . -t loadmodule`
2. Run the application with the URL and QPS as command-line arguments: `docker run loadmodule <URL> <QPS> <duration>`

## Example

```bash
docker run loadmodule https://webhook.site/ 5 10s
