package cmd

import (
	_ "dtool/prober"
	"dtool/utils"

	"github.com/spf13/cobra"
)

var filename string
var thread_num int
var upstreamCmd = &cobra.Command{
	Use: "upstream",
	//Aliases: []string{"up", "stream"},
	Short: "probe upstream recursive resolvers",
	Long: `get the upstream recursive resolvers of the input resolvers
input target can be added as an argument or as a file

	-f	input file with limited ip addresses (limit=50)
	-t	number of goroutine`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//prober.Get_upstream_ip(args[0])
		utils.SendTencentHttpdnsQuery()
	},
}

func init() {
	upstreamCmd.Flags().StringVarP(&filename, "file", "f", "", "input filename")
	upstreamCmd.Flags().IntVarP(&thread_num, "threads", "t", 10, "number of concurrent threads")
	rootCmd.AddCommand(upstreamCmd)
}
