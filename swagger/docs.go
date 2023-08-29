package swagger

import (
	"gitee.com/go-course/go9/projects/devcloud/mcenter/version"
	"github.com/go-openapi/spec"
)

func Docs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "发布中心",
			Description: "Resource for managing k8s Instances",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: version.Short(),
		},
	}
}
