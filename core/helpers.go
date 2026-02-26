package core

import (
	"log"

	"google.golang.org/protobuf/proto"
)

// MustMarshal marshals protobuf and panic if it's impossible.
func MustMarshal(pb proto.Message) []byte {
	out, err := proto.Marshal(pb)
	if err != nil {
		log.Panicf("couldn't marshal protobuf: %v", err)
	}

	return out
}

// MustUnmarshal marshals protobuf and panic if it's impossible.
func MustUnmarshal(buf []byte, pb proto.Message) {
	if err := proto.Unmarshal(buf, pb); err != nil {
		log.Panicf("couldn't unmarshal protobuf: %v", err)
	}
}
