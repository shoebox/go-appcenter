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

const Default = "application/octet-stream"

func ResolveContentType(e string) string {
	if val, ok := mapping[e]; ok {
		return val
	}

	return Default
}
