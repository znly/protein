package main

import (
	"context"
	"fmt"
	"log"
	"time"

	protein_bank "github.com/znly/protein/bank"
	"github.com/znly/protein/protoscan"
	"github.com/znly/protein/protostruct"
	tuyau_client "github.com/znly/tuyauDB/client"
	tuyau_kv "github.com/znly/tuyauDB/kv"
	tuyau_pipe "github.com/znly/tuyauDB/pipe"
	tuyau_service "github.com/znly/tuyauDB/service"
)

// -----------------------------------------------------------------------------

func buildBank() protein_bank.Bank {
	// fetched locally instanciated schemas
	schemas, err := protoscan.ScanSchemas()
	if err != nil {
		log.Fatal(err)
	}

	// build the underlying TuyauDB components: Client{Pipe, KV}
	bufSize := uint(len(schemas) + 1) // cannot block that way
	cs, err := tuyau_client.NewSimple(
		tuyau_pipe.NewRAM(bufSize, bufSize),
		tuyau_kv.NewRAM(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// build a simple TuyauDB Service to sync-up the underlying Pipe & KV
	// components (i.e. what's pushed into the pipe should en up in the kv
	// store)
	ctx, canceller := context.WithCancel(context.Background())
	s, err := tuyau_service.NewSimple(cs)
	if err != nil {
		log.Fatal(err)
	}
	go s.Run(ctx, 0)

	// build the actual Bank that integrates with the TuyauDB Client
	ty := protein_bank.NewTuyau(cs)
	go func() {
		for _, ps := range schemas {
			if err := ty.Put(ps); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Millisecond * 20)
		canceller() // we're done
	}()

	<-ctx.Done()
	// At this point, all the locally-instanciated protobuf schemas should
	// have been Put() into the Bank, which Push()ed them all to its underlying
	// Tuyau Client and, hence, into the RAM-based Tuyau Pipe.
	//
	// Since a Simple Tuyau Service had been running all along, making sure the
	// underlying RAM-based Tuyau KV store was kept in synchronization with
	// the RAM-based Pipe, our Bank should now be able to retrieve any schema
	// directly from its underlying KV store.

	return ty
}

func main() {
	b := buildBank()
	_ = b

	structType, err := protostruct.CreateStructType(
		b, b.FQNameToUID(".protein.TestSchema")[0],
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*structType)
}
