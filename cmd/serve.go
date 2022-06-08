/*
Copyright © 2022 Remi Ferrand

Contributor(s): Remi Ferrand <riton.github_at_gmail(dot)com>, 2022

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/
package cmd

import (
	"fmt"

	"github.com/riton/dirzipper/fileslist"
	"github.com/riton/dirzipper/httpsrv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serveCmdRun,
}

type serveCmdFlags struct {
	HTTPListen        string
	FilesListFilename string
	ZIPFileUrl        string
	ZIPFilename       string
}

var defaultServeCmdFlags serveCmdFlags

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&defaultServeCmdFlags.HTTPListen, "http-listen", "l", "localhost:8080", "HTTP Listen")
	serveCmd.Flags().StringVarP(&defaultServeCmdFlags.FilesListFilename, "file-list", "f", "dirzipper-filelist.json", "JSON filename with files list")
	serveCmd.Flags().StringVarP(&defaultServeCmdFlags.ZIPFileUrl, "zip-url", "u", "", "URL that triggers the ZIP file download")
	serveCmd.Flags().StringVarP(&defaultServeCmdFlags.ZIPFilename, "zip-filename", "n", "archive", "filename of the generated ZIP archive")
}

func serveCmdRun(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")

	if defaultServeCmdFlags.FilesListFilename == "" {
		logrus.Fatal("empty files list filename")
	}

	httpSrv := httpsrv.NewHTTPServerWithOptions(httpsrv.ServerOptions{
		FilesListProcessor: fileslist.NewJSONFilesListProcessor(defaultServeCmdFlags.FilesListFilename),
		ZIPFileUrl:         defaultServeCmdFlags.ZIPFileUrl,
		ZIPFilename:        defaultServeCmdFlags.ZIPFilename,
		Debug:              debug,
	})

	if err := httpSrv.ListenAndServe(defaultServeCmdFlags.HTTPListen); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("fail to serve HTTP")
	}

	fmt.Println("serve called")
}
