// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package task_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	. "github.com/pivotal-cf/on-demand-service-broker/task"
	"github.com/pivotal-cf/on-demand-service-broker/task/fakes"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("Manifest Generator", func() {
	var (
		mg              ManifestGenerator
		serviceStemcell serviceadapter.Stemcell
		serviceReleases serviceadapter.ServiceReleases
		serviceAdapter  *fakes.FakeServiceAdapterClient
		serviceCatalog  config.ServiceOffering

		existingPlan config.Plan
		secondPlan   config.Plan

		generatedManifestSecrets serviceadapter.ODBManagedSecrets
	)

	BeforeEach(func() {
		planServiceInstanceLimit := 3
		globalServiceInstanceLimit := 5

		existingPlan = config.Plan{
			ID:   existingPlanID,
			Name: existingPlanName,
			Update: &serviceadapter.Update{
				Canaries:        1,
				CanaryWatchTime: "100-200",
				UpdateWatchTime: "100-200",
				MaxInFlight:     5,
			},
			Quotas: config.Quotas{
				ServiceInstanceLimit: &planServiceInstanceLimit,
			},
			Properties: serviceadapter.Properties{
				"super": "no",
			},
			InstanceGroups: []serviceadapter.InstanceGroup{
				{
					Name:               existingPlanInstanceGroupName,
					VMType:             "vm-type",
					PersistentDiskType: "disk-type",
					Instances:          42,
					Networks:           []string{"networks"},
					AZs:                []string{"my-az1", "my-az2"},
				},
				{
					Name:      "instance-group-name-the-second",
					VMType:    "vm-type",
					Instances: 55,
					Networks:  []string{"networks2"},
				},
			},
		}

		secondPlan = config.Plan{
			ID: secondPlanID,
			Properties: serviceadapter.Properties{
				"super":             "yes",
				"a_global_property": "overrides_global_value",
			},
			InstanceGroups: []serviceadapter.InstanceGroup{
				{
					Name:               existingPlanInstanceGroupName,
					VMType:             "vm-type1",
					PersistentDiskType: "disk-type1",
					Instances:          44,
					Networks:           []string{"networks1"},
					AZs:                []string{"my-az4", "my-az5"},
				},
			},
		}

		serviceCatalog = config.ServiceOffering{
			ID:               serviceOfferingID,
			Name:             "a-cool-redis-service",
			GlobalProperties: serviceadapter.Properties{"a_global_property": "global_value", "some_other_global_property": "other_global_value"},
			GlobalQuotas: config.Quotas{
				ServiceInstanceLimit: &globalServiceInstanceLimit,
			},
			Plans: []config.Plan{
				existingPlan,
				secondPlan,
			},
		}

		serviceReleases = serviceadapter.ServiceReleases{{
			Name:    "name",
			Version: "vers",
			Jobs:    []string{"a", "b"},
		}}

		serviceStemcell = serviceadapter.Stemcell{
			OS:      "ubuntu-trusty",
			Version: "1234",
		}

		serviceAdapter = new(fakes.FakeServiceAdapterClient)

		generatedManifestSecrets = serviceadapter.ODBManagedSecrets{
			"foo":    "bar",
			"secret": "value",
		}

		mg = NewManifestGenerator(
			serviceAdapter,
			serviceCatalog,
			serviceStemcell,
			serviceReleases,
		)
	})

	Describe("GenerateManifest", func() {
		var (
			generateManifestOutput serviceadapter.MarshalledGenerateManifest
			manifest               []byte

			err error

			planGUID       string
			previousPlanID *string
			requestParams  map[string]interface{}
			oldManifest    []byte
		)

		BeforeEach(func() {
			planGUID = existingPlanID
			previousPlanID = nil

			requestParams = map[string]interface{}{"foo": "bar"}

			oldManifest = []byte("oldmanifest")
		})

		JustBeforeEach(func() {
			generateManifestOutput, err = mg.GenerateManifest(deploymentName, planGUID, requestParams, oldManifest, previousPlanID, logger)
			manifest = []byte(generateManifestOutput.Manifest)
		})

		Context("when called with correct arguments", func() {
			generatedManifest := []byte("some manifest")
			BeforeEach(func() {
				serviceAdapter.GenerateManifestReturns(serviceadapter.MarshalledGenerateManifest{Manifest: string(generatedManifest), ODBManagedSecrets: generatedManifestSecrets}, nil)
			})

			It("calls service adapter once", func() {
				Expect(serviceAdapter.GenerateManifestCallCount()).To(Equal(1))
			})

			It("returns result of adapter", func() {
				Expect(manifest).To(Equal(generatedManifest))
				Expect(generateManifestOutput.ODBManagedSecrets).To(Equal(generatedManifestSecrets))
			})

			It("does not return an error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})

			It("logs call to service adapter", func() {
				expectedLog := fmt.Sprintf("service adapter will generate manifest for deployment %s\n", deploymentName)
				Expect(logBuffer.String()).To(ContainSubstring(expectedLog))
			})

			It("calls the service adapter with the service deployment", func() {
				passedServiceDeployment, _, _, _, _, _ := serviceAdapter.GenerateManifestArgsForCall(0)
				expectedServiceDeployment := serviceadapter.ServiceDeployment{
					DeploymentName: deploymentName,
					Releases:       serviceReleases,
					Stemcell:       serviceStemcell,
				}
				Expect(passedServiceDeployment).To(Equal(expectedServiceDeployment))
			})

			It("calls the service adapter with the plan", func() {
				_, passedPlan, _, _, _, _ := serviceAdapter.GenerateManifestArgsForCall(0)
				Expect(passedPlan.InstanceGroups).To(Equal(existingPlan.InstanceGroups))
			})

			It("calls the service adapter with the request params", func() {
				_, _, passedRequestParams, _, _, _ := serviceAdapter.GenerateManifestArgsForCall(0)
				Expect(passedRequestParams).To(Equal(requestParams))
			})

			It("calls the service adapter with the old manifest", func() {
				_, _, _, passedOldManifest, _, _ := serviceAdapter.GenerateManifestArgsForCall(0)
				Expect(passedOldManifest).To(Equal(oldManifest))
			})

			It("merges global and plan properties", func() {
				_, actualPlan, _, _, _, _ := serviceAdapter.GenerateManifestArgsForCall(0)
				expectedProperties := serviceadapter.Properties{
					"a_global_property":          "global_value",
					"some_other_global_property": "other_global_value",
					"super": "no",
				}
				Expect(actualPlan.Properties).To(Equal(expectedProperties))
			})

			Context("when previous plan ID is provided", func() {
				BeforeEach(func() {
					anotherPlan := secondPlanID
					previousPlanID = &anotherPlan
				})

				It("calls the service adapter with the previous plan", func() {
					_, _, _, _, passedPreviousPlan, _ := serviceAdapter.GenerateManifestArgsForCall(0)
					Expect(passedPreviousPlan.InstanceGroups).To(Equal(secondPlan.InstanceGroups))
				})

				It("merges global and previous plan properties, overriding global with plan props", func() {
					_, _, _, _, previousPlan, _ := serviceAdapter.GenerateManifestArgsForCall(0)
					expectedProperties := serviceadapter.Properties{
						"a_global_property":          "overrides_global_value",
						"some_other_global_property": "other_global_value",
						"super": "yes",
					}

					Expect(previousPlan.Properties).To(Equal(expectedProperties))
				})
			})

			Context("when previous plan ID is not provided", func() {
				BeforeEach(func() {
					previousPlanID = nil
				})

				It("calls the service adapter with the nil previous plan", func() {
					_, _, _, _, passedPreviousPlan, _ := serviceAdapter.GenerateManifestArgsForCall(0)
					Expect(passedPreviousPlan).To(BeNil())
				})
			})
		})

		Context("when the plan cannot be found", func() {
			BeforeEach(func() {
				planGUID = "invalid-id"
			})

			It("fails without generating a manifest", func() {
				Expect(serviceAdapter.GenerateManifestCallCount()).To(Equal(0))

				Expect(err).To(Equal(PlanNotFoundError{PlanGUID: planGUID}))
				Expect(logBuffer.String()).To(ContainSubstring(planGUID))
			})
		})

		Context("when the previous plan cannot be found", func() {
			BeforeEach(func() {
				invalidID := "invalid-previous-id"
				previousPlanID = &invalidID
			})

			It("fails without generating a manifest", func() {
				Expect(serviceAdapter.GenerateManifestCallCount()).To(Equal(0))
				Expect(err).To(Equal(PlanNotFoundError{PlanGUID: *previousPlanID}))
				Expect(logBuffer.String()).To(ContainSubstring(*previousPlanID))
			})
		})

		Context("when the adapter returns an error", func() {
			BeforeEach(func() {
				serviceAdapter.GenerateManifestReturns(serviceadapter.MarshalledGenerateManifest{}, errors.New("oops"))
			})

			It("is returned", func() {
				Expect(err).To(MatchError("oops"))
			})
		})
	})

	Describe("GenerateSecretPaths", func() {
		It("generates a list of ManifestSecrets", func() {
			deploymentName := "the-name"
			secretsPath := mg.GenerateSecretPaths(deploymentName, generatedManifestSecrets)
			Expect(secretsPath).To(SatisfyAll(
				ContainElement(ManifestSecret{Name: "foo", Path: fmt.Sprintf("/%s/%s/%s/foo", config.ODBCredhubNamespace, serviceOfferingID, deploymentName), Value: generatedManifestSecrets["foo"]}),
				ContainElement(ManifestSecret{Name: "secret", Path: fmt.Sprintf("/%s/%s/%s/secret", config.ODBCredhubNamespace, serviceOfferingID, deploymentName), Value: generatedManifestSecrets["secret"]}),
			))
		})
	})

	Describe("ReplaceODBRefs", func() {
		It("replaces odb_secret:foo with /odb/<dep-name>/<svc-id>/foo", func() {
			manifest := fmt.Sprintf("name: ((%s:foo))\nsecret: ((%[1]s:bar))", serviceadapter.ODBSecretPrefix)
			secrets := []ManifestSecret{
				{Name: "foo", Value: "something", Path: "/" + config.ODBCredhubNamespace + "/jim/bob/foo"},
				{Name: "bar", Value: "another thing", Path: "/" + config.ODBCredhubNamespace + "/jim/bob/bar"},
			}
			expectedManifest := fmt.Sprintf("name: ((/%s/jim/bob/foo))\nsecret: ((/%[1]s/jim/bob/bar))", config.ODBCredhubNamespace)
			substitutedManifest := mg.ReplaceODBRefs(manifest, secrets)
			Expect(substitutedManifest).To(Equal(expectedManifest))
		})

		It("replaces all occurrences of a managed secret", func() {
			manifest := fmt.Sprintf("name: ((%s:foo))\nsecret: ((%[1]s:foo))", serviceadapter.ODBSecretPrefix)
			secrets := []ManifestSecret{
				{Name: "foo", Value: "something", Path: "/" + config.ODBCredhubNamespace + "/jim/bob/foo"},
			}
			expectedManifest := fmt.Sprintf("name: ((/%s/jim/bob/foo))\nsecret: ((/%[1]s/jim/bob/foo))", config.ODBCredhubNamespace)
			substitutedManifest := mg.ReplaceODBRefs(manifest, secrets)
			Expect(substitutedManifest).To(Equal(expectedManifest))
		})
	})
})
