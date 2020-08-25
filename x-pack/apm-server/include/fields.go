// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("apm-server", "XPackFields", asset.ModuleFieldsPri, AssetXPackFields); err != nil {
		panic(err)
	}
}

// AssetXPackFields returns asset data.
// This is the base64 encoded gzipped contents of x-pack/apm-server.
func AssetXPackFields() string {
	return "eJycktFunTAMhu95CqvXhQfgYtIeoFKl7b7ykh+OVUgy21l33n6CQsdZ0dSNOxz7i53PLT3j2hOXuXXlZBxccmpnuEqw9mfh8NwQufiEnu4+Pz7Q19959PCad9cQRVhQKUu4p08NEdGSfaDSRr0nTvHsoLWCIIMEKpoL1AV2v5IU36uopJEmCUiGSAPYq8LIargQG/kFdBHzPCrPNAimSH4t6Boiu2T1p5DTIGNPrhUNvaZYv17QUuIZ/bGrNU4roqdRcy1bJF4TzxJ6GngybMEdtv3uvFiVD7Az3G0jx+K3ad5O9vqzkxMB+/eoaHkcFSM74uGV8nDjYe/WuqY5bIYV/vtKfCn8kV1YOLtrYgXVVWTWxc/LYlfZsTS1+Ia5rZsysSOFK32DvwCJJJlrnZGWWQz6QwLsnyxHmEtaR+02wP/Z/tOYwkpOhieXGV3INfk7d1NO47vC7VXgXYFKjh+rur3O6txVO6/8FQAA//8WgjsG"
}
