package filedeltastore

import (
	"bytes"
	"encoding/json"
	"filedeltastore/src/rsync"
	"fmt"
	"io"
	"os"
	"time"
)

type FileDeltaStore struct {
	r      rsync.RSync
	deltar rsync.RSync
}

// 新建FDS差分存储工具，调用方指定block大小
func NewFileDeltaStore(blocksize int) *FileDeltaStore {
	return &FileDeltaStore{
		r: rsync.RSync{
			BlockSize: blocksize,
		},
		deltar: rsync.RSync{},
	}
}

// 新建签名文件，提供源文件名和输出的签名文件名
func (fds *FileDeltaStore) CreateSigFile(baseFile, sigFile string) error {
	// here we store the whole signature in a byte slice,
	// but it could just as well be sent over a network connection for example
	sig := make([]rsync.BlockHash, 0, 10)
	sigWriter := func(bl rsync.BlockHash) error {
		sig = append(sig, bl)
		return nil
	}

	baseF, err := os.OpenFile(baseFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fds.r.CreateSignature(baseF, sigWriter)

	if err := transBlockHashtoFile(sigFile, sig); err != nil {
		fmt.Println("transBlockHashtoFile in CreateSigFile error, ", err)
		return err
	}
	return nil
}

// 新建差分文件，提供改动过后的新文件名，签名文件名和输出的差分文件名
func (fds *FileDeltaStore) CreateDeltaFile(sigFile, newFile, deltaFile string) error {

	opsOut := make(chan rsync.Operation)
	// var blockCt, blockRangeCt, dataCt, bytes int
	// writeOperation := func(op rsync.Operation) error {
	// 	switch op.Type {
	// 	case rsync.OpBlockRange:
	// 		blockRangeCt++
	// 	case rsync.OpBlock:
	// 		blockCt++
	// 	case rsync.OpData:
	// 		// Copy data buffer so it may be reused in internal buffer.
	// 		b := make([]byte, len(op.Data))
	// 		copy(b, op.Data)
	// 		op.Data = b
	// 		dataCt++
	// 		bytes += len(op.Data)
	// 	}
	// 	opsOut <- op
	// 	return nil
	// }

	newF, err := os.OpenFile(newFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sig, err := transFiletoBlockHash(sigFile)
	if err != nil {
		fmt.Println("transFiletoBlockHash in CreateDeltaFile error, ", err)
		return err
	}

	// if err := fds.r.CreateDelta(newF, sig, writeOperation); err != nil {
	// 	fmt.Println("CreateDelta error, ", err)
	// 	return err
	// }
	go func() {
		var blockCt, blockRangeCt, dataCt, bytes int
		defer close(opsOut)
		err := fds.deltar.CreateDelta(newF, sig, func(op rsync.Operation) error {
			switch op.Type {
			case rsync.OpBlockRange:
				blockRangeCt++
			case rsync.OpBlock:
				blockCt++
			case rsync.OpData:
				// Copy data buffer so it may be reused in internal buffer.
				b := make([]byte, len(op.Data))
				copy(b, op.Data)
				op.Data = b
				dataCt++
				bytes += len(op.Data)
			}
			opsOut <- op
			return nil
		})
		if err != nil {
			fmt.Println("Failed to create delta:", err)
		}
	}()

	if err := transOpstoFile(deltaFile, opsOut); err != nil {
		fmt.Println("transOpstoFile in CreateDeltaFile error, ", err)
		return err
	}

	return nil
}

// 重建文件，提供源文件名，差分文件名，和输出的重建新文件名
func (fds *FileDeltaStore) RebuildNewFile(baseFile, deltaFile, rebuildFile string) error {
	result := new(bytes.Buffer)

	baseF, err := os.OpenFile(baseFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rebuildF, err := os.OpenFile(rebuildFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ops, err := transFiletoOps(deltaFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	opsOut := make(chan rsync.Operation)

	go func() {
		defer close(opsOut)
		for _, op := range ops {
			opsOut <- op
		}
	}()

	err = fds.r.ApplyDelta(result, baseF, opsOut)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err = rebuildF.Write(result.Bytes()); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("RebuildNewFile OK______________________")
	return nil
}

// 将blockhash写入文件
func transBlockHashtoFile(sigFile string, sigs []rsync.BlockHash) error {
	sigF, err := os.OpenFile(sigFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// for _, sig := range sigs {
	// 	sigB, _ := json.Marshal(sig)
	// 	if _, err = sigF.Write(sigB); err != nil {
	// 		fmt.Println("Writefile Error =", err)
	// 		return err
	// 	}
	// }

	sigsB, _ := json.Marshal(sigs)
	if _, err = sigF.Write(sigsB); err != nil {
		fmt.Println("Writefile Error =", err)
		return err
	}
	return nil
}

// 读取sig文件，转换为blockhash
func transFiletoBlockHash(sigFile string) (blockhash []rsync.BlockHash, err error) {
	sigF, err := os.OpenFile(sigFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sigB, err := io.ReadAll(sigF)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := json.Unmarshal(sigB, &blockhash); err != nil {
		fmt.Println("Read file error =", err)
		return nil, err
	} else {
		return blockhash, nil
	}
}

// 将差分操作写入文件
func transOpstoFile(deltaFile string, opsch chan rsync.Operation) error {
	deltaF, err := os.OpenFile(deltaFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// writeBlock := func(op Operation) error {
	// 	if _, err := target.Seek(int64(r.BlockSize*int(op.BlockIndex)), 0); err != nil {
	// 		return err
	// 	}
	// 	n, err := io.ReadAtLeast(target, buffer, r.BlockSize)
	// 	if err != nil && err != io.ErrUnexpectedEOF {
	// 		return err
	// 	}
	// 	block = buffer[:n]
	// 	if _, err = alignedTarget.Write(block); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	ops := make([]rsync.Operation, 0, 10)
	for op := range opsch {
		// switch op.Type {
		// case rsync.OpBlockRange:
		// 	for i := op.BlockIndex; i <= op.BlockIndexEnd; i++ {
		// 		if err := writeBlock(Operation{
		// 			Type:       OpBlock,
		// 			BlockIndex: i,
		// 		}); err != nil {
		// 			if err == io.EOF {
		// 				break
		// 			}
		// 			return err
		// 		}
		// 	}
		// case rsync.OpBlock:
		// 	if err := writeBlock(op); err != nil {
		// 		if err == io.EOF {
		// 			break
		// 		}
		// 		return err
		// 	}
		// case rsync.OpData:
		// 	if _, err := alignedTarget.Write(op.Data); err != nil {
		// 		return err
		// 	}
		// }
		ops = append(ops, op)
	}
	opsB, _ := json.Marshal(ops)
	if _, err = deltaF.Write(opsB); err != nil {
		fmt.Println("Writefile Error =", err)
		return err
	}
	return nil
}

// 读取差分操作文件，转为差分operations
func transFiletoOps(deltaFile string) (ops []rsync.Operation, err error) {
	deltaF, err := os.OpenFile(deltaFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	deltaB, err := io.ReadAll(deltaF)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := json.Unmarshal(deltaB, &ops); err != nil {
		fmt.Println("Read file error =", err)
		return nil, err
	} else {
		return ops, nil
	}
}

func FDSTest() {
	fds := NewFileDeltaStore(16)
	sigFileName := "sig"
	deltaFileName := "delta"
	baseFileName := "base"
	newFileName := "new"
	rebuildFileName := "rebuild"

	fds.CreateSigFile(baseFileName, sigFileName)
	time.Sleep(1 * time.Second)

	fds.CreateDeltaFile(sigFileName, newFileName, deltaFileName)
	time.Sleep(1 * time.Second)

	fds.RebuildNewFile(baseFileName, deltaFileName, rebuildFileName)
}
