package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	internalLog "myproject/internal/log"
	pkgLog "myproject/pkg/log"
)

func main() {
	fmt.Println("haha")
	internalLog.Log("haha")
	pkgLog.Log("custom", "haha")
	logrus.Info("haha")
}
