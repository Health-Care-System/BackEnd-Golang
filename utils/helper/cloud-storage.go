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
	key := "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiaGVhbHRoaWZ5LTQwNzMwMCIsCiAgInByaXZhdGVfa2V5X2lkIjogIjNlMTkzMGI5OGM1MjZjMGU5YWE1ODlkYWQ5OTRjNTFhNTU3ZWVjODQiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ3hNbUYvZVhGaVdkYUlcbmZjNVdJYS9oWkp0M1FNV0M1ZUdieDhDVUhKb0o3WVAxVXZZMEhvY0ZBVTIwOGIwUTdHek5tK0xaQzB0TnFtMmNcbkd6TCtDZTdURGptSkRhYnNFZHlRcEp6ekMrUFFQVjRpU2xGbGE1SUtIV0xkMW0wOFRGd1BKQXl5Ujg4QlVtTDBcbjh0Z0E4K0NCN0ZOSXBwd2xGUWN5dmhaTzFCNjVsVnhhdjByaE9NajBUeXorYVU0U0tOeXpQK21Ja1F2NEoxeUtcbkd3ZkNJVzZySWRtR21zOHVlOE9ReDBCYzdJaXdidnh5T2ExTDR3NDdrMk9rbTVmZHZ1U00wQU9PYlZKZVlTeXBcbjRJY1o1amhST1BSR3JGcTlwQmZJYWRJY0pNVFZPMWNIckhIQjQ3ZS91WGNzOUtNREhBbGRSM3NXUEdQTTRrZXRcbmtMR0hvODZqQWdNQkFBRUNnZ0VBRTVmRkJSSGFFaXhwMXN6R1dHckFKbFY0QlQ4aHQ4QVFhcWV6S1ExaWtPOTBcbjFGYjRoem5LNXloR1BkZ3U5aVNXUVVQWDJQVHhQaVMybXZpODhpOEorOGRLWUZla21keTl2VGFlYmhkbCtMd1lcbllRVkRxOG14cTdHbElXOVl3NzZUdlU3WW9vdmpIVHZOV2xWUndnVTFVVmhldXgxRkRhcWVFYVBJRnE1SXVXOVZcbjNWQW1QTy80MTQyMXN4ajErbWVGOXhzRDVsRTN0MWpsMTJJMGoyWmNibVhLT3ZqUm1oaFQyWVVBOVhUY0J6ckdcbnZWMG0xNERqQlZUTk9oS0lJZGtqbkprcDBJa2FBbXpBbUR4eWd2RzBiazBCSnV2RUZtcWdhVC9CY01PRHkyTk9cbmVLbjVoelUveng2Nlp6UVovOEttUURoL3VvVFlkQ2RKTFJRQnZZcVlYUUtCZ1FEV20zeUNIb2V5ZTZQZDcvbUdcbnZ6TitjNzRkSEd0eTFQTFZZS25McjdURitUTFYvQWYvdWszcnVYUTZFWlI2elA1WFpQaGlTYWxHU3ZjbGtzbUxcblVuOWJPbnFHMUdNSzhiTldKbjkvTVhBSUo5Q1VvTnROelZTeUE1RTNBRXJhVFdJYkV5cHFsWjc3eVpMY2lkbDZcbm40em1ZM1VaOTdCSGd1T0haWWZac2ZLb2h3S0JnUURUWDdNa0prRXRxYk9ub1FXaWdBQTM0TUozRTg1NnlmdFlcbkN2aUYzZ0JaNDEzVTFFaEkzQWJGNGUxK2tWNkprUWN4aXp5RU84MGhlQzk4S1hhTE51NEZGVEN0M0MvaUtGaTdcbnlkYUo2cTM5L3lRSTBNdVhsNkY4b1ZHSzVKREFRODN2TlJhYnQ5cmJmdS9xbW1YUXFIN01IRUhXejlwWUVDTldcbnNVM0JiM3BjQlFLQmdIR2w4YSt1bjBuanRBbktGYWhJQk9zSVBEdUtXMVI2ZFFhT3BCeWJ0ZTNKWkNSeHpZS2RcbmxMb3FnZVJtZnV6eE5oZnQvcU4rUXNoWTFyenRHUkpRNCtUWitSMEJ5Rmw1V2ZGYmZkVkx4dnBxcTBpcVRyaktcbjdmay9ibDFrS0QrbkR6Y3JWU0VRanhyanlvUkQ5QW0rQ0kzUlNhZ3d4UWQ1eHloaW1paXMxY1p0QW9HQkFLY2ZcbmpudHhoNjA3OVNEL3JuM2FLTklGY3B4RjI4YTM5bk9aVVFBL0ZCWCtNRTA3ZnQra24vSkxmTVRLMlcxNWJxK2NcbmdEK3BMTHBlMVdTZFArNDRneDhmcnZwNEVxQUUrSXVadlhnVVJuZUNDSkt6eTVWVFBVcFdIaXZzSmdydVVWL2xcbm9MZUVPWlc4bXFMcWFyLzh5U3hHMTBPcDJlQXcray9zSmlkZ0plV2xBb0dBZmY5RllyYmNRVHJseWwrcmhvZGxcbnArZ0Myd0tLUWd3U1ZWaVhMMWxtekUvSjliaS9mdjhkWGxLVVdSTkoza0Y3SURYQ0VGNW9GTVhCMUJSR0JGYy9cbnZNOXhpV2JrU05rVjB3TTFIR0Z2NVQ2dVdPSWZVMHNySUtTZ1J6US9Wc2lsSXg1TFM5QW1CaUQ3K3crTVNSbGRcbmFJRDlKZEdaRGZkeUFpbFlVQW5oai9zPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImVtYmVyLWRldkBoZWFsdGhpZnktNDA3MzAwLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjExMTU4Mjc0OTQ3NzMzOTM5MDk1OCIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZW1iZXItZGV2JTQwaGVhbHRoaWZ5LTQwNzMwMC5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgInVuaXZlcnNlX2RvbWFpbiI6ICJnb29nbGVhcGlzLmNvbSIKfQo="

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

	key, err := base64.StdEncoding.DecodeString("ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiaGVhbHRoaWZ5LTQwNzMwMCIsCiAgInByaXZhdGVfa2V5X2lkIjogIjNlMTkzMGI5OGM1MjZjMGU5YWE1ODlkYWQ5OTRjNTFhNTU3ZWVjODQiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ3hNbUYvZVhGaVdkYUlcbmZjNVdJYS9oWkp0M1FNV0M1ZUdieDhDVUhKb0o3WVAxVXZZMEhvY0ZBVTIwOGIwUTdHek5tK0xaQzB0TnFtMmNcbkd6TCtDZTdURGptSkRhYnNFZHlRcEp6ekMrUFFQVjRpU2xGbGE1SUtIV0xkMW0wOFRGd1BKQXl5Ujg4QlVtTDBcbjh0Z0E4K0NCN0ZOSXBwd2xGUWN5dmhaTzFCNjVsVnhhdjByaE9NajBUeXorYVU0U0tOeXpQK21Ja1F2NEoxeUtcbkd3ZkNJVzZySWRtR21zOHVlOE9ReDBCYzdJaXdidnh5T2ExTDR3NDdrMk9rbTVmZHZ1U00wQU9PYlZKZVlTeXBcbjRJY1o1amhST1BSR3JGcTlwQmZJYWRJY0pNVFZPMWNIckhIQjQ3ZS91WGNzOUtNREhBbGRSM3NXUEdQTTRrZXRcbmtMR0hvODZqQWdNQkFBRUNnZ0VBRTVmRkJSSGFFaXhwMXN6R1dHckFKbFY0QlQ4aHQ4QVFhcWV6S1ExaWtPOTBcbjFGYjRoem5LNXloR1BkZ3U5aVNXUVVQWDJQVHhQaVMybXZpODhpOEorOGRLWUZla21keTl2VGFlYmhkbCtMd1lcbllRVkRxOG14cTdHbElXOVl3NzZUdlU3WW9vdmpIVHZOV2xWUndnVTFVVmhldXgxRkRhcWVFYVBJRnE1SXVXOVZcbjNWQW1QTy80MTQyMXN4ajErbWVGOXhzRDVsRTN0MWpsMTJJMGoyWmNibVhLT3ZqUm1oaFQyWVVBOVhUY0J6ckdcbnZWMG0xNERqQlZUTk9oS0lJZGtqbkprcDBJa2FBbXpBbUR4eWd2RzBiazBCSnV2RUZtcWdhVC9CY01PRHkyTk9cbmVLbjVoelUveng2Nlp6UVovOEttUURoL3VvVFlkQ2RKTFJRQnZZcVlYUUtCZ1FEV20zeUNIb2V5ZTZQZDcvbUdcbnZ6TitjNzRkSEd0eTFQTFZZS25McjdURitUTFYvQWYvdWszcnVYUTZFWlI2elA1WFpQaGlTYWxHU3ZjbGtzbUxcblVuOWJPbnFHMUdNSzhiTldKbjkvTVhBSUo5Q1VvTnROelZTeUE1RTNBRXJhVFdJYkV5cHFsWjc3eVpMY2lkbDZcbm40em1ZM1VaOTdCSGd1T0haWWZac2ZLb2h3S0JnUURUWDdNa0prRXRxYk9ub1FXaWdBQTM0TUozRTg1NnlmdFlcbkN2aUYzZ0JaNDEzVTFFaEkzQWJGNGUxK2tWNkprUWN4aXp5RU84MGhlQzk4S1hhTE51NEZGVEN0M0MvaUtGaTdcbnlkYUo2cTM5L3lRSTBNdVhsNkY4b1ZHSzVKREFRODN2TlJhYnQ5cmJmdS9xbW1YUXFIN01IRUhXejlwWUVDTldcbnNVM0JiM3BjQlFLQmdIR2w4YSt1bjBuanRBbktGYWhJQk9zSVBEdUtXMVI2ZFFhT3BCeWJ0ZTNKWkNSeHpZS2RcbmxMb3FnZVJtZnV6eE5oZnQvcU4rUXNoWTFyenRHUkpRNCtUWitSMEJ5Rmw1V2ZGYmZkVkx4dnBxcTBpcVRyaktcbjdmay9ibDFrS0QrbkR6Y3JWU0VRanhyanlvUkQ5QW0rQ0kzUlNhZ3d4UWQ1eHloaW1paXMxY1p0QW9HQkFLY2ZcbmpudHhoNjA3OVNEL3JuM2FLTklGY3B4RjI4YTM5bk9aVVFBL0ZCWCtNRTA3ZnQra24vSkxmTVRLMlcxNWJxK2NcbmdEK3BMTHBlMVdTZFArNDRneDhmcnZwNEVxQUUrSXVadlhnVVJuZUNDSkt6eTVWVFBVcFdIaXZzSmdydVVWL2xcbm9MZUVPWlc4bXFMcWFyLzh5U3hHMTBPcDJlQXcray9zSmlkZ0plV2xBb0dBZmY5RllyYmNRVHJseWwrcmhvZGxcbnArZ0Myd0tLUWd3U1ZWaVhMMWxtekUvSjliaS9mdjhkWGxLVVdSTkoza0Y3SURYQ0VGNW9GTVhCMUJSR0JGYy9cbnZNOXhpV2JrU05rVjB3TTFIR0Z2NVQ2dVdPSWZVMHNySUtTZ1J6US9Wc2lsSXg1TFM5QW1CaUQ3K3crTVNSbGRcbmFJRDlKZEdaRGZkeUFpbFlVQW5oai9zPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImVtYmVyLWRldkBoZWFsdGhpZnktNDA3MzAwLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjExMTU4Mjc0OTQ3NzMzOTM5MDk1OCIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZW1iZXItZGV2JTQwaGVhbHRoaWZ5LTQwNzMwMC5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgInVuaXZlcnNlX2RvbWFpbiI6ICJnb29nbGVhcGlzLmNvbSIKfQo=")
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
