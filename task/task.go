// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package task

import (
	"fmt"
	"log"
	"reflect"

	"github.com/pivotal-cf/on-demand-service-broker/boshdirector"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter -o fakes/fake_bosh_client.go . BoshClient
type BoshClient interface {
	Deploy(manifest []byte, contextID string, logger *log.Logger) (int, error)
	GetTasks(deploymentName string, logger *log.Logger) (boshdirector.BoshTasks, error)
	GetDeployment(name string, logger *log.Logger) ([]byte, bool, error) // TODO SF found = false => manifest => nil, drop the found flag?
}

// TODO SF previousPlanID is a pointer because it might not exist. Should we have a nil value instead? Should we have a specific type?
//go:generate counterfeiter -o fakes/fake_manifest_generator.go . ManifestGenerator
type ManifestGenerator interface {
	GenerateManifest(deploymentName, planID string, requestParams map[string]interface{}, oldManifest []byte, previousPlanID *string, logger *log.Logger) (BoshManifest, error)
}

type deployer struct {
	boshClient        BoshClient
	manifestGenerator ManifestGenerator
}

func NewDeployer(boshClient BoshClient, manifestGenerator ManifestGenerator) deployer {
	return deployer{
		boshClient:        boshClient,
		manifestGenerator: manifestGenerator,
	}
}

func (d deployer) Create(deploymentName, planID string, requestParams map[string]interface{}, boshContextID string, logger *log.Logger) (int, []byte, error) {
	err := d.assertNoOperationsInProgress(deploymentName, logger)
	if err != nil {
		return 0, nil, err
	}

	return d.doDeploy(deploymentName, planID, "create", requestParams, nil, nil, boshContextID, logger)
}

func (d deployer) Upgrade(deploymentName, planID string, previousPlanID *string, boshContextID string, logger *log.Logger) (int, []byte, error) {
	err := d.assertNoOperationsInProgress(deploymentName, logger)
	if err != nil {
		return 0, nil, err
	}

	oldManifest, err := d.getDeploymentManifest(deploymentName, logger)
	if err != nil {
		return 0, nil, err
	}

	return d.doDeploy(deploymentName, planID, "upgrade", nil, oldManifest, previousPlanID, boshContextID, logger)
}

func (d deployer) Update(
	deploymentName,
	planID string,
	requestParams map[string]interface{},
	previousPlanID *string,
	boshContextID string,
	logger *log.Logger,
) (boshTaskID int, manifest []byte, err error) {
	if err := d.assertNoOperationsInProgress(deploymentName, logger); err != nil {
		return 0, nil, err
	}

	oldManifest, err := d.getDeploymentManifest(deploymentName, logger)
	if err != nil {
		return 0, nil, err
	}

	if err := d.checkForPendingChanges(deploymentName, previousPlanID, oldManifest, logger); err != nil {
		return 0, nil, err
	}

	return d.doDeploy(deploymentName, planID, "update", requestParams, oldManifest, previousPlanID, boshContextID, logger)
}

func (d deployer) getDeploymentManifest(deploymentName string, logger *log.Logger) ([]byte, error) {
	oldManifest, found, err := d.boshClient.GetDeployment(deploymentName, logger)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, NewDeploymentNotFoundError(fmt.Errorf("bosh deployment '%s' not found", deploymentName))
	}

	return oldManifest, nil
}

func (d deployer) assertNoOperationsInProgress(deploymentName string, logger *log.Logger) error {
	clientTasks, err := d.boshClient.GetTasks(deploymentName, logger)
	if err != nil {
		return NewServiceError(fmt.Errorf("error getting tasks for deployment %s: %s\n", deploymentName, err))
	}

	if incompleteTasks := clientTasks.IncompleteTasks(); len(incompleteTasks) != 0 {
		logger.Printf("deployment %s is still in progress: tasks %s\n", deploymentName, incompleteTasks.ToLog())
		return TaskInProgressError{Message: "task in progress"}
	}

	return nil
}

func (d deployer) checkForPendingChanges(
	deploymentName string,
	previousPlanID *string,
	oldManifest BoshManifest,
	logger *log.Logger,
) error {
	regeneratedManifest, err := d.manifestGenerator.GenerateManifest(deploymentName, *previousPlanID, map[string]interface{}{}, oldManifest, previousPlanID, logger)
	if err != nil {
		return err
	}

	convertedRegeneratedManifest, err := ignoreUpdateBlock(regeneratedManifest)
	if err != nil {
		return err
	}

	convertedOldManifest, err := ignoreUpdateBlock(oldManifest)
	if err != nil {
		return err
	}

	manifestsSame := reflect.DeepEqual(convertedRegeneratedManifest, convertedOldManifest)

	pendingChanges := !manifestsSame

	if pendingChanges {
		return PendingChangesNotAppliedError{}
	}

	return nil
}

func (d deployer) doDeploy(
	deploymentName,
	planID string,
	operationType string,
	requestParams map[string]interface{},
	oldManifest []byte,
	previousPlanID *string,
	boshContextID string,
	logger *log.Logger,
) (int, []byte, error) {

	manifest, err := d.manifestGenerator.GenerateManifest(deploymentName, planID, requestParams, oldManifest, previousPlanID, logger)
	if err != nil {
		return 0, nil, err
	}

	boshTaskID, err := d.boshClient.Deploy(manifest, boshContextID, logger)
	if err != nil {
		return 0, nil, fmt.Errorf("error deploying instance: %s\n", err)
	}
	logger.Printf("Bosh task ID for %s deployment %s is %d\n", operationType, deploymentName, boshTaskID)

	return boshTaskID, manifest, nil
}

func ignoreUpdateBlock(rawManifest []byte) (bosh.BoshManifest, error) {
	var boshManifest bosh.BoshManifest
	err := yaml.Unmarshal(rawManifest, &boshManifest)

	if err != nil {
		return bosh.BoshManifest{}, fmt.Errorf("error detecting change in manifest, unable to unmarshal manifest: %s", err)
	}

	boshManifest.Update = bosh.Update{}

	return boshManifest, nil
}
