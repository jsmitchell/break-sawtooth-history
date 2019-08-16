package main

import (
	"fmt"
	"os"
	"time"
	"math/rand"
	"github.com/google/uuid"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
	"github.com/taekion-org/break-sawtooth-history/client"
	"github.com/taekion-org/break-sawtooth-history/tp/handler"
	flag "github.com/spf13/pflag"
)

const DEFAULT_REST_URL = "http://localhost:8008"
const DEFAULT_WAIT_TIME = 10
const DEFAULT_BLOCK_SIZE = 32768
const DEFAULT_ITERATIONS = 1000

func main() {
	// Set up params
	var err error
	var url *string = flag.String("url", DEFAULT_REST_URL, "Sawtooth REST API URL")
	var keyFile *string = flag.String("keyfile", "", "Sawtooth Private Key File")
	var wait *uint = flag.Uint("wait", DEFAULT_WAIT_TIME, "Time to wait for blocking commits")
	var blockSize *int = flag.Int("blocksize", DEFAULT_BLOCK_SIZE, "Block size in bytes")
	var iterations *int = flag.Int("iterations", DEFAULT_ITERATIONS, "Number of iterations")
	flag.Parse()

	// Output the params
	fmt.Printf("Taekion 'Break Sawtooth History' Demo\n")
	fmt.Printf("REST URL: %s\n", *url)
	fmt.Printf("Wait Time: %ds\n", *wait)
	fmt.Printf("Block Size: %d bytes\n", *blockSize)
	fmt.Printf("Iterations: %d\n", *iterations)
	fmt.Println("")

	// Create a client instance
	cli, err := client.NewBreakSawtoothClient(*url, *keyFile)
	if err != nil {
		handleError(err)
	}

	// Create and commit the original block (we will test against this one).
	idOriginal := uuid.New()
	buf := make([]byte, *blockSize)
	rand.Read(buf)
	sha256Original := handler.Hexdigest(string(buf))

	_, err = cli.SetBlock(idOriginal, buf, *wait)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("UUID Original: %s\n", idOriginal)
	fmt.Printf("Address Original: %s\n", handler.GetAddress(idOriginal))
	fmt.Printf("SHA256 Original: %s\n", sha256Original)
	fmt.Println("")

	// Make sure that we can read the block back, save the current head, and check the checksum
	dataCheck, headCheck, err := cli.GetBlock(idOriginal, "")
	if err != nil {
		handleError(err)
	}

	sha256Check := handler.Hexdigest(string(dataCheck))
	fmt.Printf("Head Check: %s\n", headCheck)
	fmt.Printf("SHA256 Check: %s\n", sha256Check)
	fmt.Println("")

	if sha256Original != sha256Check {
		handleError(fmt.Errorf("SHA256 doesn't match..."))
	}

	// Do the writes
	fmt.Println("Starting Updates:")
	for i:=0; i<*iterations-1; i++ {
		rand.Read(buf)
		if i < *iterations-1 {
			id := uuid.New()
			_, err = cli.SetBlock(id, buf, 0)

		} else {
			_, err = cli.SetBlock(idOriginal, buf, *wait)
		}

		if err != nil {
			switch err.(type) {
			case *errors.SawtoothClientTransportError:
				sawtoothError := err.(*errors.SawtoothClientTransportError)
				if sawtoothError.ErrorCode == 31 {
					fmt.Printf("(r%-d)", i)
					time.Sleep(time.Millisecond*500)
					i -= 1
					continue
				}
			default:
				handleError(err)
			}
		}
		fmt.Printf(".")
	}
	fmt.Printf("\n\n")

	fmt.Println("Final Check:")
	fmt.Printf("Requesting address %s at head %s\n", handler.GetAddress(idOriginal), headCheck)
	dataFinal, headFinal, err := cli.GetBlock(idOriginal, headCheck)
	if err != nil {
		handleError(err)
	}

	fmt.Println()
	fmt.Println("If you see this message, the failure has not happened.")
	fmt.Printf("Head Final: %s\n", headFinal)
	sha256Final := handler.Hexdigest(string(dataFinal))

	fmt.Printf("SHA256 Final: %s\n", sha256Final)
	if sha256Check != sha256Final {
		handleError(fmt.Errorf("SHA256 doesn't match..."))
	}
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
