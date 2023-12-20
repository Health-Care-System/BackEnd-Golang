package helper

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

func UploadFilesToGCS(c echo.Context, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	bucketName := "dev-healthify"
	key := "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiaGVhbHRoaWZ5LTQwNzMwMCIsCiAgInByaXZhdGVfa2V5X2lkIjogImZlZWZlOGU4ZWYzMGY3OGVkMDk4YWZiZTdlYWFhOThmMmMyNjUyMWEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRHJTN0NRYUQ2L28zN3dcbmpCcm1DeHhUeVdsbFlIMUR6bytnRy8yQjMrNGFNVXhvOElRSnNTRVVSeGlRdGVRTDVPekEvdnVTVGJPVnZIemtcbi9QQWpnVmZCVDliMFdLdUxRMmtLMnBHQms4Y0d2WjU1bEJhODFUenRWbHd2T21nOE93Y3ZpWlpRWWtJaXBwdm5cblA3YlgrTFM3M1p2WWg5ZC9oSUdqRERGcEEyOGNGL05kN3pPSitPNllXc2NQeFRWdkhuR1VlWEloaUxJUWFGdEZcbktEdXJwMVdiWjgyOWNjbzdESkExY1BPQ1ozOTk2WG1PaG1yRy9TdzV3akFKRGhrSFgxaXFCbWxBZzJxcnZvb2lcbk9mdU11aExxZ2dwdk9vejVUaXdFRDR0Y2pnUGJSam42djVzTVg3SVpvelNqMmdBMkRCWDlsc1RyL1l4TmlHb0dcblpwUmpSbjNiQWdNQkFBRUNnZ0VBVFpNejBEaCtUNkp3ekkvMG9sbVJhbEppVEVrbW8zOXJ2T2JkaXI4d2VqUThcbnMvQndKOTNkUVJrN2tTSlc0RkVHZVk3WGxHSEh6cHRKTmhucWRscERlM2wrTGlzNXkzMWJHWGY4TnhOb3IrRStcbmFXa1lvZ05QeGhRTjZvaEFLM051cUE2ZG56ejEvd1NkSy8zR282ZmN0bitXelNua002ZVZNaE5vQ3U1VjZKV2ZcbkVsT1p6Y3VVNUFLa2dSblFtWDJTUVN5bDlNSjArZmdVMWJkSkpEZ2dmREovTGlFakxvYTdhcFZaZ3NqWGcyMTNcblN3Q2lSTkg0S2p6TzhpTDhiUFdYZFQvbjkzOFJLcjU4Zjk2dEZMQTFkNm5iOEFkV1lwWnVGMWl6WEtPUko0QU5cbnlncTFuY0gvdVlEd0xkbGhNVWhqZ1ZkRFh1MDFKNk5kcVhEekxseDl1UUtCZ1FEMkIvWXVsRGtOQ0RqdVZVcWlcbmMvMVZxWEZsajZ1NkJCOWdnLytlNkJpUXdVTDJJZ1NjaVB0clRrYVgxZzZPNEIvRDVjWmRFcTUwOFF0bEc5cXdcbjBwblpVWXRVdkc4R2VuZXpKUWxEQlJqY2h3UlVyNk5NejVpVTVGVTRaaHgvaUt0RWpIRElrTFdEcDZRbCs1d2RcbjZxNlE2eGw2UkdTdzNPUmtlVmJESTJGWXB3S0JnUUQwMUY4SGVwakRQN0NTMFB4NnMyWkJ6NkROMFdBZUltR0VcbjRZeTN5SUhKZXhsd29PNmFQVHdZQ0VXR1ErY3IrTUhZaUwvbGM0VU5GYnQ3ZVhWWXI1WjN2YmU5U1QzMVRCM1hcbndVTFFOYXh4U1YvSlhycmxlcWc2eFhYVmpDVFd5TGMyZWJBY096MW1MQ0h4Q2lpdjk5RnVjaHNLR3IyM0k3ZWRcblg2MDVBc05qclFLQmdRRGIrVnBmWXg5dlMzNjdlWDcxcVFkRTQrOERnMlZqTi9SbDh4OUdFUGFGMW9Yc1U1WVRcbjcxWDhKMHh1elhET3hnMGd0Znlaa3U1d21HUTd6cC9Gb1ViMVN6ZHNWOTVjeUhybHJhT09US3hoNEVZN0FaN2Ncbk9uQ09EMmt5dC9tYS9iTkQ0dDJrTmQ1VkREcHp5M2RXT0ZKRU9DL2JaZk14UHc3bDFxZUFhYzNMMlFLQmdFSVRcbjczSDlUSzJseXVwVkxVK1FpOURIVVFjN3MwMXV5aE1yTE9lTlhqb2ExMHJtcEg5TWQ2T2sxOTdkQk0rQlhCQXRcbkdGMjlSL080SWRtNWRrcHhXWk1IeVVkMU5SdTNaM2FMMnBTSFovdExhbWJYQW1wZEtIcDBRTkZaK2JkZWhOUlhcbkU0a0xGQTgzYWhHOFJDNzZHN1JMWjdEYWRzbXZBaWVmWXdrNERiUWhBb0dCQU9vTnFKejJrM3JCNFZJRVlrWVZcbkFIeTZMcy91OEZxdG51ZTNwS29DcTdGM05PeFVqd2RmdURreFpXZysxZUhhRHB3VDNIMEZlWmZCL1d6Smt2bUNcbnhJUE0zck5vZFpuS1VBOUdONWRXMVRIQ1hkQUw2TVdaM1FHM1N0MC9LbDBVVEZUL1hldW94Z2tSVEYwU0F1Y2JcbmNBQmI3NitRVDBLb3V0aHVUd3pOcXZZV1xuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImRldmVtYmVyQGhlYWx0aGlmeS00MDczMDAuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA0MTgyNTIyMTgwMDU4MzkzNDI5IiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9kZXZlbWJlciU0MGhlYWx0aGlmeS00MDczMDAuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJ1bml2ZXJzZV9kb21haW4iOiAiZ29vZ2xlYXBpcy5jb20iCn0K"

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Fatal("Can't decode service account key")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(keyBytes))
	if err != nil {
		log.Fatal("Can't connect to Google Cloud Storage")
	}

	currentTime := time.Now().UTC()

	year := currentTime.Year()
	month := int(currentTime.Month())
	day := currentTime.Day()
	hour := currentTime.Hour()
	minute := currentTime.Minute()
	second := currentTime.Second()

	formattedTime := fmt.Sprintf("%04d%02d%02d-%d%02d%02d", year, month, day, second, minute, hour)

	filePath := formattedTime + "-" + fileHeader.Filename

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	obj := client.Bucket(bucketName).Object(filePath)

	// create a writer for the object
	wc := obj.NewWriter(ctx)

	// upload
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// generate URL
	objectAttrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	URL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectAttrs.Name)

	return URL, nil
}

func DeleteFilesFromGCS(filename string) error {
	ctx := context.Background()

	key, err := base64.StdEncoding.DecodeString("ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiaGVhbHRoaWZ5LTQwNzMwMCIsCiAgInByaXZhdGVfa2V5X2lkIjogImZlZWZlOGU4ZWYzMGY3OGVkMDk4YWZiZTdlYWFhOThmMmMyNjUyMWEiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRHJTN0NRYUQ2L28zN3dcbmpCcm1DeHhUeVdsbFlIMUR6bytnRy8yQjMrNGFNVXhvOElRSnNTRVVSeGlRdGVRTDVPekEvdnVTVGJPVnZIemtcbi9QQWpnVmZCVDliMFdLdUxRMmtLMnBHQms4Y0d2WjU1bEJhODFUenRWbHd2T21nOE93Y3ZpWlpRWWtJaXBwdm5cblA3YlgrTFM3M1p2WWg5ZC9oSUdqRERGcEEyOGNGL05kN3pPSitPNllXc2NQeFRWdkhuR1VlWEloaUxJUWFGdEZcbktEdXJwMVdiWjgyOWNjbzdESkExY1BPQ1ozOTk2WG1PaG1yRy9TdzV3akFKRGhrSFgxaXFCbWxBZzJxcnZvb2lcbk9mdU11aExxZ2dwdk9vejVUaXdFRDR0Y2pnUGJSam42djVzTVg3SVpvelNqMmdBMkRCWDlsc1RyL1l4TmlHb0dcblpwUmpSbjNiQWdNQkFBRUNnZ0VBVFpNejBEaCtUNkp3ekkvMG9sbVJhbEppVEVrbW8zOXJ2T2JkaXI4d2VqUThcbnMvQndKOTNkUVJrN2tTSlc0RkVHZVk3WGxHSEh6cHRKTmhucWRscERlM2wrTGlzNXkzMWJHWGY4TnhOb3IrRStcbmFXa1lvZ05QeGhRTjZvaEFLM051cUE2ZG56ejEvd1NkSy8zR282ZmN0bitXelNua002ZVZNaE5vQ3U1VjZKV2ZcbkVsT1p6Y3VVNUFLa2dSblFtWDJTUVN5bDlNSjArZmdVMWJkSkpEZ2dmREovTGlFakxvYTdhcFZaZ3NqWGcyMTNcblN3Q2lSTkg0S2p6TzhpTDhiUFdYZFQvbjkzOFJLcjU4Zjk2dEZMQTFkNm5iOEFkV1lwWnVGMWl6WEtPUko0QU5cbnlncTFuY0gvdVlEd0xkbGhNVWhqZ1ZkRFh1MDFKNk5kcVhEekxseDl1UUtCZ1FEMkIvWXVsRGtOQ0RqdVZVcWlcbmMvMVZxWEZsajZ1NkJCOWdnLytlNkJpUXdVTDJJZ1NjaVB0clRrYVgxZzZPNEIvRDVjWmRFcTUwOFF0bEc5cXdcbjBwblpVWXRVdkc4R2VuZXpKUWxEQlJqY2h3UlVyNk5NejVpVTVGVTRaaHgvaUt0RWpIRElrTFdEcDZRbCs1d2RcbjZxNlE2eGw2UkdTdzNPUmtlVmJESTJGWXB3S0JnUUQwMUY4SGVwakRQN0NTMFB4NnMyWkJ6NkROMFdBZUltR0VcbjRZeTN5SUhKZXhsd29PNmFQVHdZQ0VXR1ErY3IrTUhZaUwvbGM0VU5GYnQ3ZVhWWXI1WjN2YmU5U1QzMVRCM1hcbndVTFFOYXh4U1YvSlhycmxlcWc2eFhYVmpDVFd5TGMyZWJBY096MW1MQ0h4Q2lpdjk5RnVjaHNLR3IyM0k3ZWRcblg2MDVBc05qclFLQmdRRGIrVnBmWXg5dlMzNjdlWDcxcVFkRTQrOERnMlZqTi9SbDh4OUdFUGFGMW9Yc1U1WVRcbjcxWDhKMHh1elhET3hnMGd0Znlaa3U1d21HUTd6cC9Gb1ViMVN6ZHNWOTVjeUhybHJhT09US3hoNEVZN0FaN2Ncbk9uQ09EMmt5dC9tYS9iTkQ0dDJrTmQ1VkREcHp5M2RXT0ZKRU9DL2JaZk14UHc3bDFxZUFhYzNMMlFLQmdFSVRcbjczSDlUSzJseXVwVkxVK1FpOURIVVFjN3MwMXV5aE1yTE9lTlhqb2ExMHJtcEg5TWQ2T2sxOTdkQk0rQlhCQXRcbkdGMjlSL080SWRtNWRrcHhXWk1IeVVkMU5SdTNaM2FMMnBTSFovdExhbWJYQW1wZEtIcDBRTkZaK2JkZWhOUlhcbkU0a0xGQTgzYWhHOFJDNzZHN1JMWjdEYWRzbXZBaWVmWXdrNERiUWhBb0dCQU9vTnFKejJrM3JCNFZJRVlrWVZcbkFIeTZMcy91OEZxdG51ZTNwS29DcTdGM05PeFVqd2RmdURreFpXZysxZUhhRHB3VDNIMEZlWmZCL1d6Smt2bUNcbnhJUE0zck5vZFpuS1VBOUdONWRXMVRIQ1hkQUw2TVdaM1FHM1N0MC9LbDBVVEZUL1hldW94Z2tSVEYwU0F1Y2JcbmNBQmI3NitRVDBLb3V0aHVUd3pOcXZZV1xuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImRldmVtYmVyQGhlYWx0aGlmeS00MDczMDAuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTA0MTgyNTIyMTgwMDU4MzkzNDI5IiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9kZXZlbWJlciU0MGhlYWx0aGlmeS00MDczMDAuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJ1bml2ZXJzZV9kb21haW4iOiAiZ29vZ2xlYXBpcy5jb20iCn0K")
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(key))
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket("dev-healthify")
	object := bucket.Object(filename)

	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}
