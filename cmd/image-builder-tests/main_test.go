// +build integration

package main

import (
	"context"
	"fmt"
	"encoding/json"
	"net/http"
	"os/exec"
	"testing"
	"time"

	// "github.com/osbuild/image-builder/internal/cloudapi"
	"github.com/osbuild/image-builder/internal/server"

	"github.com/stretchr/testify/require"
)

// TODO build requetss from distro endpoints
// composeRequest := cloudapi.ComposeRequest{
// 	Distribution: "rhel-8",
// 	ImageRequests: []cloudapi.ImageRequest{
// 		cloudapi.ImageRequest{
// 			Architecture: "x86_64",
// 			ImageType: "qcow2",
// 			Repositories: [],
// 		}
// 	},
// }

func AddDummyIdentityHeader(ctx context.Context, request *http.Request) error {
	request.Header.Add("x-rh-identity", "eyJlbnRpdGxlbWVudHMiOnsiaW5zaWdodHMiOnsiaXNfZW50aXRsZWQiOnRydWV9LCJzbWFydF9tYW5hZ2VtZW50Ijp7ImlzX2VudGl0bGVkIjp0cnVlfSwib3BlbnNoaWZ0Ijp7ImlzX2VudGl0bGVkIjp0cnVlfSwiaHlicmlkIjp7ImlzX2VudGl0bGVkIjp0cnVlfSwibWlncmF0aW9ucyI6eyJpc19lbnRpdGxlZCI6dHJ1ZX0sImFuc2libGUiOnsiaXNfZW50aXRsZWQiOnRydWV9fSwiaWRlbnRpdHkiOnsiYWNjb3VudF9udW1iZXIiOiIwMDAwMDAiLCJ0eXBlIjoiVXNlciIsInVzZXIiOnsidXNlcm5hbWUiOiJ1c2VyIiwiZW1haWwiOiJ1c2VyQHVzZXIudXNlciIsImZpcnN0X25hbWUiOiJ1c2VyIiwibGFzdF9uYW1lIjoidXNlciIsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfb3JnX2FkbWluIjp0cnVlLCJpc19pbnRlcm5hbCI6dHJ1ZSwibG9jYWxlIjoiZW4tVVMifSwiaW50ZXJuYWwiOnsib3JnX2lkIjoiMDAwMDAwIn19fQ==")
	return nil
}

func ConfigureClient(client *server.Client) error {
	client.RequestEditor = AddDummyIdentityHeader
	return nil
}

func RunTestWithClient(t *testing.T, ib string)  {
	client, err := server.NewClientWithResponses(ib, ConfigureClient)
	require.NoError(t, err)

	// wait until server is reachable
	ctx := context.Background()
	tries := 0
	for tries  < 5 {
		resp, err := client.GetVersionWithResponse(ctx)
		if err == nil && resp.StatusCode() == http.StatusOK {
			break
		} else if tries == 4 {
			require.NoError(t, err)
		}
		time.Sleep(time.Second)
		tries += 1
	}

	distroResp, err := client.GetDistributionsWithResponse(ctx)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, distroResp.StatusCode(), "Error: got non-200 status. Full response: %s", distroResp.Body)
	require.NotNil(t, distroResp.JSON200)
	require.NotEmpty(t, distroResp.JSON200)

	distro := (*distroResp.JSON200)[0].Name
	archResp, err := client.GetArchitecturesWithResponse(ctx, distro)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, archResp.StatusCode(), "Error: got non-200 status. Full response: %s", archResp.Body)
	require.NotEmpty(t, archResp.JSON200)


	// Build a composerequest
	composeRequest := server.ComposeImageJSONRequestBody{
		Distribution: distro,
		ImageRequests: []server.ImageRequest{
			server.ImageRequest{
				Architecture: (*archResp.JSON200)[0].Arch,
				ImageType: (*archResp.JSON200)[0].ImageTypes[0],
			},
		},
	}

	b, err := json.Marshal(composeRequest)
	require.NoError(t, err)
	fmt.Println(fmt.Sprintf("AAAAAAAAAA %v", string(b)))
	composeResp, err := client.ComposeImageWithResponse(ctx, composeRequest)
	require.NoError(t, err)
	fmt.Println(fmt.Sprintf("BBBBBBB %v", composeResp))

	// Check if we get 404 without the identity header
	client, err = server.NewClientWithResponses(ib)
	versionResp, err := client.GetVersionWithResponse(ctx)
	require.NoError(t, err)
	require.Equalf(t, http.StatusNotFound, versionResp.StatusCode(), "Error: got non-404 status. Full response: %s", versionResp.Body)

}

func TestImageBuilder(t *testing.T) {
	// allow to run against existing instance
	// run image builder
	cmd := exec.Command("/usr/libexec/image-builder/image-builder")
	err := cmd.Start()
	require.NoError(t, err)
	defer cmd.Process.Kill()

	RunTestWithClient(t, "http://127.0.0.1:8086/api/image-builder/v1")
	RunTestWithClient(t, "http://127.0.0.1:8086/api/image-builder/v1.0")
}

// Same test as above but against existing contain on localhost:8087
func TestImageBuilderContainer(t *testing.T) {
	RunTestWithClient(t, "http://127.0.0.1:8087/api/image-builder/v1")
	RunTestWithClient(t, "http://127.0.0.1:8086/api/image-builder/v1.0")
}
