package cobratest

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// ./sync master ./files_master --http-listen-addr=127.0.0.1 --http-listen-port=9988 --slave-grpc-addr=127.0.0.1 --slave-grpc-port=9989 --server-side-encryption false
var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "master run a master server, with all function enabled",
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		fmt.Println("files store path is:", filepath)
		encS, err := cmd.Flags().GetString("server-side-encryption")
		if err != nil {
			fmt.Println("server-side-encryption flag error", err)
			return
		}
		encS = strings.ToLower(encS)
		enc := false
		switch encS {
		case "false":
			enc = false
		case "0":
			enc = false
		case "true":
			enc = true
		case "1":
			enc = true

		}
		httpIP, err := cmd.Flags().GetString("http-listen-addr")
		if err != nil {
			fmt.Println("http-listen-addr flag error", err)
			return
		}
		httpPort, err := cmd.Flags().GetInt("http-listen-port")
		if err != nil {
			fmt.Println("http-listen-port flag error", err)
			return
		}
		rpcIP, err := cmd.Flags().GetString("slave-grpc-addr")
		if err != nil {
			fmt.Println("slave-grpc-addr flag error", err)
			return
		}
		rpcPort, err := cmd.Flags().GetInt("slave-grpc-port")
		if err != nil {
			fmt.Println("slave-grpc-port flag error", err)
			return
		}
		// server, err := server.NewServer(def.ROLE_MASTER, httpIP, httpPort, rpcIP, rpcPort, filepath, enc)
		// if err != nil {
		// 	fmt.Println("NewServer error:", err.Error())
		// 	return
		// }
		fmt.Printf("server-side-encryption: %v, http-listen-addr: %v, http-listen-port: %v, slave-grpc-addr: %v, slave-grpc-port: %v, filepath: %v\n", enc, httpIP, httpPort, rpcIP, rpcPort, filepath)
		// server.Run()
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	masterCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	masterCmd.Flags().StringP("server-side-encryption", "", "", "server side encryption, use 'false' or 'true'")
	masterCmd.Flags().StringP("http-listen-addr", "", "", "http listen addr")
	masterCmd.Flags().IntP("http-listen-port", "", -1, "http listen port")
	masterCmd.Flags().StringP("slave-grpc-addr", "", "", "slave grpc addr")
	masterCmd.Flags().IntP("slave-grpc-port", "", -1, "slave grpc port")
	rootCmd.AddCommand(masterCmd)
}
