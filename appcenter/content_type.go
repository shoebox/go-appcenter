package appcenter

var mapping = map[string]string{
	"apk":        "application/vnd.android.package-archive",
	"aab":        "application/vnd.android.package-archive",
	"msi":        "application/x-msi",
	"plist":      "application/xml",
	"aetx":       "application/c-x509-ca-cert",
	"cer":        "application/pkix-cert",
	"xap":        "application/x-silverlight-app",
	"appx":       "application/x-appx",
	"appxbundle": "application/x-appxbundle",
	"appxupload": "application/x-appxupload",
	"appxsym":    "application/x-appxupload",
	"msix":       "application/x-msix",
	"msixbundle": "application/x-msixbundle",
	"msixupload": "application/x-msixupload",
	"msixsym":    "application/x-msixupload",
}

const defaultOctetStream = "application/octet-stream"

// ResolveContentType will resolve the right content type for the provided extension name. If the
// provided extension is not of a support format, it will return defaultOctetStream as default
func ResolveContentType(fileExtension string) string {
	if val, ok := mapping[fileExtension]; ok {
		return val
	}

	return defaultOctetStream
}
