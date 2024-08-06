package rsync

import (
	"fmt"
	"os"
	"strings"
	// "github.com/minio/rsync-go"
)

func RsyncTest() {
	oldReader := strings.NewReader("I am the original content")

	rs := &RSync{}

	// here we store the whole signature in a byte slice,
	// but it could just as well be sent over a network connection for example
	sig := make([]BlockHash, 0, 10)
	writeSignature := func(bl BlockHash) error {
		sig = append(sig, bl)
		return nil
	}
	rs.CreateSignature(oldReader, writeSignature) // CreateSigFile

	currentReader := strings.NewReader("I am the new content")

	opsOut := make(chan Operation)
	writeOperation := func(op Operation) error {
		opsOut <- op
		return nil
	}

	go func() {
		defer close(opsOut)
		rs.CreateDelta(currentReader, sig, writeOperation) // CreateDeltaFile
	}()

	var newWriter strings.Builder
	oldReader.Seek(0, os.SEEK_SET)

	rs.ApplyDelta(&newWriter, oldReader, opsOut) // RebuildNewFile

	fmt.Println(newWriter.String())
	// Output: I am the new content
}
