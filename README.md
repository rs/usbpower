# usbpower

A Go library and command-line tool for reading power samples (voltage and current) from compatible USB power meter devices (e.g., [WITRN C4/C5](https://www.avhzy.com/html/product-detail/usb-c-meter-c4) meters) via USB HID.

## Features
- Reads real-time voltage and current samples from supported USB power meters
- Outputs samples in JSON format (optional)
- Computes and displays statistics (min, max, average, percentiles) for voltage and current
- Supports sampling for a specified duration

## Requirements
- Go 1.18 or newer
- Compatible USB power meter (e.g., WITRN)
- macOS, Linux, or Windows

## Installation

Clone the repository and build the command-line tool:

```sh
git clone <this-repo-url>
cd usbpower
cd cmd/usbpower
go build -o usbpower
```

## Usage

```sh
./usbpower [flags]
```

### Flags
- `-duration <duration>`: Stop sampling after the specified duration (e.g., `10s`, `1m`). If not set, runs until interrupted.
- `-output json`: Output each sample as a JSON object to stdout.

### Example

Sample for 10 seconds and print statistics:

```sh
./usbpower -duration 10s
```

Sample for 5 seconds and output raw JSON samples:

```sh
./usbpower -duration 5s -output json
```

## Output
After sampling, the tool prints statistics for voltage and current:
- Minimum
- Maximum
- Average
- 50th percentile (median)
- 90th percentile

## Library Usage

You can use the `usbpower` package in your own Go programs to access USB power meter data:

```go
import "github.com/rs/usbpower"

device, err := usbpower.OpenDevice()
if err != nil {
    // handle error
}
defer device.Close()
sample, err := device.Read()
// use sample.Voltage, sample.Current, sample.Timestamp
```

## License
MIT
