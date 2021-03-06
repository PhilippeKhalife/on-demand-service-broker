// Copyright (C) 2015-Present Pivotal Software, Inc. All rights reserved.

// This program and the accompanying materials are made available under
// the terms of the under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package credhub_tests

import (
	"fmt"
	"os"

	"testing"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	odbcredhub "github.com/pivotal-cf/on-demand-service-broker/credhub"

	"github.com/totherme/unstructured"
)

const credhubURL = "https://credhub.service.cf.internal:8844"

var (
	dev_env string
	caCerts []string
)

func TestContractTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Credhub Contract Tests Suite")
}

var _ = BeforeSuite(func() {
	dev_env = os.Getenv("DEV_ENV")
	cfUAACACert := os.Getenv("CF_UAA_CA_CERT")
	Expect(cfUAACACert).ToNot(BeEmpty())
	cfCredhubCACert := os.Getenv("CF_CREDHUB_CA_CERT")
	Expect(cfCredhubCACert).ToNot(BeEmpty())
	caCerts = []string{cfUAACACert, cfCredhubCACert}
	ensureCredhubIsClean()
})

var _ = AfterSuite(func() {
	ensureCredhubIsClean()
})

func nameIsCredhubMatcher(data unstructured.Data) bool {
	val, err := data.GetByPointer("/name")
	if err != nil {
		return false
	}
	stringVal, err := val.StringValue()
	if err != nil {
		return false
	}
	return stringVal == "credhub"
}

func testKeyPrefix() string {
	return fmt.Sprintf("/test-%s", dev_env)
}

func makeKeyPath(name string) string {
	return fmt.Sprintf("%s/%s", testKeyPrefix(), name)
}

func ensureCredhubIsClean() {
	credhubClient := underlyingCredhubClient()

	testKeys, err := credhubClient.FindByPath(testKeyPrefix())
	Expect(err).NotTo(HaveOccurred())
	for _, key := range testKeys.Credentials {
		credhubClient.Delete(key.Name)
	}
}

func getCredhubStore() *odbcredhub.Store {
	clientSecret := os.Getenv("TEST_CREDHUB_CLIENT_SECRET")
	Expect(clientSecret).NotTo(BeEmpty(), "Expected TEST_CREDHUB_CLIENT_SECRET to be set")

	credentialStore, err := odbcredhub.Build(
		credhubURL,
		credhub.Auth(auth.UaaClientCredentials("credhub_cli", clientSecret)),
		credhub.CaCerts(caCerts...),
	)
	Expect(err).NotTo(HaveOccurred())
	return credentialStore
}

func underlyingCredhubClient() *credhub.CredHub {
	clientSecret := os.Getenv("TEST_CREDHUB_CLIENT_SECRET")
	Expect(clientSecret).NotTo(BeEmpty(), "Expected TEST_CREDHUB_CLIENT_SECRET to be set")
	Expect(caCerts).ToNot(BeEmpty())

	credhubClient, err := credhub.New(
		credhubURL,
		credhub.Auth(auth.UaaClientCredentials("credhub_cli", clientSecret)),
		credhub.CaCerts(caCerts...),
	)
	Expect(err).NotTo(HaveOccurred())
	return credhubClient
}
