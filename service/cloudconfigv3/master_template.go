package cloudconfigv3

import (
	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/certificatetpr"
	k8scloudconfig "github.com/giantswarm/k8scloudconfig/v_3_0_0"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/randomkeytpr"
)

// NewMasterTemplate generates a new master cloud config template and returns it
// as a base64 encoded string.
func (c *CloudConfig) NewMasterTemplate(customObject v1alpha1.AWSConfig, certs certificatetpr.CompactTLSAssets, keys randomkeytpr.CompactRandomKeyAssets) (string, error) {
	var err error

	var params k8scloudconfig.Params
	{
		params.Cluster = customObject.Spec.Cluster
		params.EtcdPort = customObject.Spec.Cluster.Etcd.Port
		params.Extension = &MasterExtension{
			certs:        certs,
			customObject: customObject,
			keys:         keys,
		}
	}

	var newCloudConfig *k8scloudconfig.CloudConfig
	{
		cloudConfigConfig := k8scloudconfig.DefaultCloudConfigConfig()
		cloudConfigConfig.Params = params
		cloudConfigConfig.Template = k8scloudconfig.MasterTemplate

		newCloudConfig, err = k8scloudconfig.NewCloudConfig(cloudConfigConfig)
		if err != nil {
			return "", microerror.Mask(err)
		}

		err = newCloudConfig.ExecuteTemplate()
		if err != nil {
			return "", microerror.Mask(err)
		}
	}

	return newCloudConfig.Base64(), nil
}

type MasterExtension struct {
	certs        certificatetpr.CompactTLSAssets
	customObject v1alpha1.AWSConfig
	keys         randomkeytpr.CompactRandomKeyAssets
}

func (e *MasterExtension) Files() ([]k8scloudconfig.FileAsset, error) {
	filesMeta := []k8scloudconfig.FileMetadata{
		{
			AssetContent: decryptTLSAssetsScriptTemplate,
			Path:         "/opt/bin/decrypt-tls-assets",
			Owner:        FileOwner,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.APIServerCrt,
			Path:         "/etc/kubernetes/ssl/apiserver-crt.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.APIServerCA,
			Path:         "/etc/kubernetes/ssl/apiserver-ca.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.APIServerKey,
			Path:         "/etc/kubernetes/ssl/apiserver-key.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.ServiceAccountCrt,
			Path:         "/etc/kubernetes/ssl/service-account-crt.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.ServiceAccountCA,
			Path:         "/etc/kubernetes/ssl/service-account-ca.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.ServiceAccountKey,
			Path:         "/etc/kubernetes/ssl/service-account-key.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.CalicoClientCrt,
			Path:         "/etc/kubernetes/ssl/calico/client-crt.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.CalicoClientCA,
			Path:         "/etc/kubernetes/ssl/calico/client-ca.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.CalicoClientKey,
			Path:         "/etc/kubernetes/ssl/calico/client-key.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.EtcdServerCrt,
			Path:         "/etc/kubernetes/ssl/etcd/server-crt.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.EtcdServerCA,
			Path:         "/etc/kubernetes/ssl/etcd/server-ca.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.EtcdServerKey,
			Path:         "/etc/kubernetes/ssl/etcd/server-key.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		// Add second copy of files for etcd client certs. Will be replaced by
		// a separate client cert.
		{
			AssetContent: e.certs.EtcdServerCrt,
			Path:         "/etc/kubernetes/ssl/etcd/client-crt.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.EtcdServerCA,
			Path:         "/etc/kubernetes/ssl/etcd/client-ca.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.certs.EtcdServerKey,
			Path:         "/etc/kubernetes/ssl/etcd/client-key.pem.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  FilePermission,
		},
		{
			AssetContent: waitDockerConfTemplate,
			Path:         "/etc/systemd/system/docker.service.d/01-wait-docker.conf",
			Owner:        FileOwner,
			Permissions:  FilePermission,
		},
		{
			AssetContent: decryptKeysAssetsScriptTemplate,
			Path:         "/opt/bin/decrypt-keys-assets",
			Owner:        FileOwner,
			Permissions:  FilePermission,
		},
		{
			AssetContent: e.keys.APIServerEncryptionKey,
			Path:         "/etc/kubernetes/encryption/k8s-encryption-config.yaml.enc",
			Owner:        FileOwner,
			Encoding:     GzipBase64Encoding,
			Permissions:  0644,
		},
		// Add use-proxy-protocol to ingress-controller ConfigMap, this doesn't work
		// on KVM because of dependencies on hardware LB configuration.
		{
			AssetContent: ingressControllerConfigMapTemplate,
			Path:         "/srv/ingress-controller-cm.yml",
			Owner:        FileOwner,
			Permissions:  0644,
		},
	}

	var newFiles []k8scloudconfig.FileAsset

	for _, fm := range filesMeta {
		c, err := k8scloudconfig.RenderAssetContent(fm.AssetContent, e.customObject.Spec)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		fileAsset := k8scloudconfig.FileAsset{
			Metadata: fm,
			Content:  c,
		}

		newFiles = append(newFiles, fileAsset)
	}

	return newFiles, nil
}

func (e *MasterExtension) Units() ([]k8scloudconfig.UnitAsset, error) {
	unitsMeta := []k8scloudconfig.UnitMetadata{
		{
			AssetContent: decryptTLSAssetsServiceTemplate,
			Name:         "decrypt-tls-assets.service",
			Enable:       true,
			Command:      "start",
		},
		{
			AssetContent: masterFormatVarLibDockerServiceTemplate,
			Name:         "format-var-lib-docker.service",
			Enable:       true,
			Command:      "start",
		},
		{
			AssetContent: ephemeralVarLibDockerMountTemplate,
			Name:         "var-lib-docker.mount",
			Enable:       true,
			Command:      "start",
		},
		{
			AssetContent: decryptKeysServiceTemplate,
			Name:         "decrypt-keys-assets.service",
			Enable:       true,
			Command:      "start",
		},
	}

	var newUnits []k8scloudconfig.UnitAsset

	for _, fm := range unitsMeta {
		c, err := k8scloudconfig.RenderAssetContent(fm.AssetContent, e.customObject.Spec)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		unitAsset := k8scloudconfig.UnitAsset{
			Metadata: fm,
			Content:  c,
		}

		newUnits = append(newUnits, unitAsset)
	}

	return newUnits, nil
}

func (e *MasterExtension) VerbatimSections() []k8scloudconfig.VerbatimSection {
	newSections := []k8scloudconfig.VerbatimSection{
		{
			Name:    "storage",
			Content: instanceStorageTemplate,
		},
		{
			Name:    "storageclass",
			Content: instanceStorageClassTemplate,
		},
	}
	return newSections
}
