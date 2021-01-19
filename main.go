package main

import(
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inexio/go-monitoringplugin"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"os"
)

var opts struct{
	URL 		string `short:"u" long:"url" description:"The url for requesting the federation check" required:"true"`
}
type HealthResponse struct {
	Version []string `json:"versions"`
}

type ErrorAndCode struct {
	ExitCode int
	Error    error
}

func CheckHealth(url string, versionArray []string) []ErrorAndCode {
	var errSlice []ErrorAndCode
	var found bool
	if url == "" {
		errSlice = append(errSlice, ErrorAndCode{2, errors.New("URL and server name is required to send GET request")})
		return errSlice
	}
	client := resty.New()
	request := client.SetDebugBodyLimit(1000).R()
	response, err := request.Get(url)
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.Wrap(err, "error during http request")})
		return errSlice
	}
	var resp HealthResponse
	err = json.Unmarshal(response.Body(), &resp)
		for _, s1 := range resp.Version {
			found = false
			for _, s2 := range versionArray {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				errSlice = append(errSlice, ErrorAndCode{2, errors.New("This version could not be found " + s1)})
			}
		}
	return errSlice
}


func OutputMonitoring(errSlice []ErrorAndCode, defaultMessage string) {
	response := monitoringplugin.NewResponse(defaultMessage)
	for i := 0; i < len(errSlice); i++ {
		response.UpdateStatus(errSlice[i].ExitCode, errSlice[i].Error.Error())
	}
	response.OutputAndExit()
}

func fillDiffArray (versionArray []string) []string{
	versionArray[0] = "r0.0.1"
	versionArray[1] = "r0.1.0"
	versionArray[2] = "r0.2.0"
	versionArray[3] = "r0.3.0"
	versionArray[4] = "r0.4.0"
	versionArray[5] = "r0.5.0"
	versionArray[6] = "r0.6.0"

	return versionArray
}

func main() {
 var errSlice []ErrorAndCode
 diffVersions := make([]string, 7)
 var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println("error parsing flags")
		os.Exit(3)
	}
	diffVersions = fillDiffArray(diffVersions)
	errSlice = CheckHealth(opts.URL, diffVersions)
	OutputMonitoring(errSlice, "health checked")


}
