package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// Flags
	reverse := flag.Bool("r", false, "reverse hex dump back to binary")
	output := flag.String("o", "", "output file name (default: auto)")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: xxdtool [options] <file>")
		flag.PrintDefaults()
		return
	}

	inputFile := flag.Arg(0)
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Determine output file
	outFile := *output
	if outFile == "" {
		if *reverse {
			if strings.HasSuffix(inputFile, ".hex") {
				outFile = strings.TrimSuffix(inputFile, ".hex")
			} else {
				outFile = "recovered_" + inputFile
			}
		} else {
			outFile = inputFile + ".hex"
		}
	}

	if *reverse {
		decoded, err := reverseHex(data)
		if err != nil {
			log.Fatalf("Failed to decode hex: %v", err)
		}
		if err := ioutil.WriteFile(outFile, decoded, 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("Reversed hex dump saved to %s\n", outFile)
	} else {
		hexDump := createHexDump(data)
		if err := ioutil.WriteFile(outFile, []byte(hexDump), 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("Hex dump saved to %s\n", outFile)
	}
}

// createHexDump generates a hex dump like xxd
func createHexDump(data []byte) string {
	var builder strings.Builder
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		line := data[i:end]
		builder.WriteString(fmt.Sprintf("%08x: ", i))
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				builder.WriteString(fmt.Sprintf("%02x ", data[i+j]))
			} else {
				builder.WriteString("   ")
			}
		}
		builder.WriteString(" ")
		for _, b := range line {
			if b >= 32 && b <= 126 {
				builder.WriteByte(b)
			} else {
				builder.WriteByte('.')
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// reverseHex converts a hex dump back to binary
func reverseHex(data []byte) ([]byte, error) {
    lines := strings.Split(string(data), "\n")
    var decoded []byte
    for _, line := range lines {
        if len(line) < 10 {
            continue
        }

        // Remove the offset (first 8 chars + colon)
        hexPart := line[9:]
        // Only keep hex bytes before the double space (before ASCII)
        if idx := strings.Index(hexPart, "  "); idx != -1 {
            hexPart = hexPart[:idx]
        }

        // Split hex bytes by spaces
        parts := strings.Fields(hexPart)
        for _, p := range parts {
            if len(p) != 2 {
                continue
            }
            b, err := hex.DecodeString(p)
            if err != nil {
                return nil, err
            }
            decoded = append(decoded, b...)
        }
    }
    return decoded, nil
}


