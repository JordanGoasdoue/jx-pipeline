package plugins

import (
	"fmt"
	"os"
	"strings"

	jenkinsv1 "github.com/jenkins-x/jx-api/v3/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/extensions"
	"github.com/jenkins-x/jx-helpers/v3/pkg/homedir"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKptBinary returns the path to the locally installed kpt 3 extension
func GetKptBinary(version string) (string, error) {
	if version == "" {
		version = KptVersion
	}
	pluginBinDir, err := homedir.PluginBinDir(os.Getenv("JX_GITOPS_HOME"), ".jx-gitops")
	if err != nil {
		return "", errors.Wrapf(err, "failed to find plugin home dir")
	}
	plugin := CreateKptPlugin(version)
	return extensions.EnsurePluginInstalled(plugin, pluginBinDir)
}

// CreateKptPlugin creates the kpt 3 plugin
func CreateKptPlugin(version string) jenkinsv1.Plugin {
	binaries := extensions.CreateBinaries(func(p extensions.Platform) string {
		return fmt.Sprintf("https://github.com/GoogleContainerTools/kpt/releases/download/v%s/kpt_%s_%s_%s.tar.gz", version, strings.ToLower(p.Goos), strings.ToLower(p.Goarch), version)
	})

	plugin := jenkinsv1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: KptPluginName,
		},
		Spec: jenkinsv1.PluginSpec{
			SubCommand:  "kpt",
			Binaries:    binaries,
			Description: "kpt 3 binary",
			Name:        KptPluginName,
			Version:     version,
		},
	}
	return plugin
}
