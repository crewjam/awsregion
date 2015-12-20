package awsregion

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

var getAvailabilityZoneOnce sync.Once
var availabilityZone string

// GuessRegion updates the AWS configuration `c` with the
// current region if one is not set by examining EC2 metadata
func GuessRegion(c *aws.Config) {
	if c.Region != nil && *c.Region != "" {
		return
	}

	getAvailabilityZoneOnce.Do(func() {
		req, _ := http.NewRequest("GET",
			"http://169.254.169.254/latest/meta-data/placement/availability-zone",
			nil)
		httpClient := *http.DefaultClient
		httpClient.Timeout = time.Second

		resp, err := httpClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		availabilityZone = string(body)
	})

	if availabilityZone != "" {
		c.Region = aws.String(availabilityZone[:len(availabilityZone)-1])
	}
}
