package cloudconfig

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/giantswarm/awstpr"
)

var (
	awsFilesMeta []FileMetadata = []FileMetadata{
		FileMetadata{
			AssetPath:   "templates/operator/aws/decrypt-tls-assets",
			Path:        "/opt/bin/decrypt-tls-assets",
			Owner:       "root:root",
			Permissions: 0700,
		},
	}
	awsUnitsMeta []UnitMetadata = []UnitMetadata{
		UnitMetadata{
			AssetPath: "templates/operator/aws/decrypt-tls-assets.service",
			Name:      "decrypt-tls-certs.service",
			Enable:    true,
			Command:   "start",
		},
	}
	kvmUnitsMeta []UnitMetadata = []UnitMetadata{
		UnitMetadata{
			AssetPath: "templates/operator/kvm/add-route-bridge.service",
			Name:      "add-route-bridge.service",
			Enable:    true,
			Command:   "start",
		},
	}
)

type FileMetadata struct {
	AssetPath   string
	Path        string
	Owner       string
	Permissions int
}

type FileAsset struct {
	Metadata FileMetadata
	Content  []string
}

type UnitMetadata struct {
	AssetPath string
	Name      string
	Enable    bool
	Command   string
}

type UnitAsset struct {
	Metadata UnitMetadata
	Content  []string
}

type OperatorExtension interface {
	GetFiles() ([]FileAsset, error)
	GetUnits() ([]UnitAsset, error)
}

type AWSExtension struct {
	AwsInfo awstpr.Spec
}

func NewAWSExtension() (*AWSExtension, error) {
	return &AWSExtension{}, nil
}

func renderAssetContent(assetPath string, params interface{}) ([]string, error) {
	rawContent, err := Asset(assetPath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(assetPath).Parse(string(rawContent[:]))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, params); err != nil {
		return nil, err
	}

	content := strings.Split(buf.String(), "\n")
	return content, nil
}

func (a *AWSExtension) GetFiles() ([]FileAsset, error) {
	files := make([]FileAsset, 0, len(awsFilesMeta))

	for _, fileMeta := range awsFilesMeta {
		content, err := renderAssetContent(fileMeta.AssetPath, a.AwsInfo)
		if err != nil {
			return nil, err
		}

		fileAsset := FileAsset{
			Metadata: fileMeta,
			Content:  content,
		}

		files = append(files, fileAsset)
	}

	return files, nil
}

func (a *AWSExtension) GetUnits() ([]UnitAsset, error) {
	units := make([]UnitAsset, 0, len(awsUnitsMeta))

	for _, unitMeta := range awsUnitsMeta {
		content, err := renderAssetContent(unitMeta.AssetPath, a.AwsInfo)
		if err != nil {
			return nil, err
		}

		unitAsset := UnitAsset{
			Metadata: unitMeta,
			Content:  content,
		}

		units = append(units, unitAsset)
	}

	return units, nil
}

type KVMExtension struct {
}

func NewKVMExtension() (*KVMExtension, error) {
	return &KVMExtension{}, nil
}

func (k *KVMExtension) GetFiles() ([]FileAsset, error) {
	return nil, nil
}

func (k *KVMExtension) GetUnits() ([]UnitAsset, error) {
	return nil, nil
}
