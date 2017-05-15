// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/pivotal-cf/brokerapi"
	apiauth "github.com/pivotal-cf/brokerapi/auth"
	"github.com/pivotal-cf/on-demand-service-broker/authorizationheader"
	"github.com/pivotal-cf/on-demand-service-broker/boshdirector"
	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/cf"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/credstore"
	"github.com/pivotal-cf/on-demand-service-broker/loggerfactory"
	"github.com/pivotal-cf/on-demand-service-broker/mgmtapi"
	"github.com/pivotal-cf/on-demand-service-broker/network"
	"github.com/pivotal-cf/on-demand-service-broker/serviceadapter"
	"github.com/pivotal-cf/on-demand-service-broker/task"
	"github.com/urfave/negroni"
)

func main() {
	loggerFactory := loggerfactory.New(os.Stdout, "on-demand-service-broker", loggerfactory.Flags)

	logger := loggerFactory.New()
	logger.Println("Starting broker")

	configFilePath := flag.String("configFilePath", "", "path to config file")

	flag.Parse()
	if *configFilePath == "" {
		logger.Fatal("must supply -configFilePath")
	}

	conf, err := config.Parse(*configFilePath)
	if err != nil {
		logger.Fatalf("error parsing config: %s", err)
	}
	startBroker(conf, logger, loggerFactory)
}

func startBroker(conf config.Config, logger *log.Logger, loggerFactory *loggerfactory.LoggerFactory) {
	var (
		boshAuthenticator boshdirector.AuthHeaderBuilder
		err               error
	)

	boshAuthConfig := conf.Bosh.Authentication
	if boshAuthConfig.Basic.IsSet() {
		boshAuthenticator = authorizationheader.NewBasicAuthHeaderBuilder(
			boshAuthConfig.Basic.Username,
			boshAuthConfig.Basic.Password,
		)
	} else if boshAuthConfig.UAA.IsSet() {
		boshAuthenticator, err = authorizationheader.NewClientTokenAuthHeaderBuilder(
			boshAuthConfig.UAA.UAAURL,
			boshAuthConfig.UAA.ID,
			boshAuthConfig.UAA.Secret,
			conf.Broker.TLSCertVerication,
			[]byte(conf.Bosh.TrustedCert),
		)
		if err != nil {
			logger.Fatalf("error creating BOSH authorization header builder: %s", err)
		}
	}

	boshClient, err := boshdirector.New(conf.Bosh.URL, boshAuthenticator, conf.Broker.TLSCertVerication, []byte(conf.Bosh.TrustedCert))
	if err != nil {
		logger.Fatalf("error creating bosh client: %s", err)
	}

	cfAuthenticator, err := conf.CF.NewAuthHeaderBuilder(conf.Broker.TLSCertVerication)
	if err != nil {
		logger.Fatalf("error creating CF authorization header builder: %s", err)
	}

	cfClient, err := cf.New(
		conf.CF.URL,
		cfAuthenticator,
		[]byte(conf.CF.TrustedCert),
		conf.Broker.TLSCertVerication,
	)
	if err != nil {
		logger.Fatalf("error creating Cloud Foundry client: %s", err)
	}

	serviceAdapter := &serviceadapter.Client{
		ExternalBinPath: conf.ServiceAdapter.Path,
		CommandRunner:   serviceadapter.NewCommandRunner(),
	}

	manifestGenerator := task.NewManifestGenerator(
		serviceAdapter,
		conf.ServiceCatalog,
		conf.ServiceDeployment.Stemcell,
		conf.ServiceDeployment.Releases,
	)

	deploymentManager := task.NewDeployer(boshClient, manifestGenerator)
	credStore := credentialStore(conf.Credhub, conf.Broker.TLSCertVerication)

	onDemandBroker, err := broker.New(boshClient, cfClient, credStore, serviceAdapter, deploymentManager, conf.ServiceCatalog, loggerFactory, conf.Features)
	if err != nil {
		logger.Fatalf("error starting broker: %s", err)
	}

	if conf.Broker.StartUpBanner {
		fmt.Println(`
                  .//\
        \\      .+ssso/\     \\
      \---.\  .+ssssssso/.  \----\         ____     ______     ______
    .--------+ssssssssssso+--------\      / __ \   (_  __ \   (_   _ \
  .-------:+ssssssssssssssss+--------\   / /  \ \    ) ) \ \    ) (_) )
 -------./ssssssssssssssssssss:.------- ( ()  () )  ( (   ) )   \   _/
  \--------+ssssssssssssssso+--------/  ( ()  () )   ) )  ) )   /  _ \
    \-------.+osssssssssso/.-------/     \ \__/ /   / /__/ /   _) (_) )
      \---./  ./osssssso/   \.---/        \____/   (______/   (______/
        \/      \/osso/       \/
                  \/:/
								`)
	}

	brokerRouter := mux.NewRouter()
	mgmtapi.AttachRoutes(brokerRouter, onDemandBroker, conf.ServiceCatalog, loggerFactory)

	brokerAPILager := lager.NewLogger("odb-broker-api")
	brokerAPILager.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	brokerapi.AttachRoutes(brokerRouter, onDemandBroker, brokerAPILager)
	authProtectedBrokerAPI := apiauth.NewWrapper(conf.Broker.Username, conf.Broker.Password).Wrap(brokerRouter)

	negroniLogger := &negroni.Logger{ALogger: logger}
	server := negroni.New(negroni.NewRecovery(), negroniLogger, negroni.NewStatic(http.Dir("public")))

	server.UseHandler(authProtectedBrokerAPI)
	server.Run(fmt.Sprintf("0.0.0.0:%d", conf.Broker.Port))
}

func credentialStore(credhub *config.Credhub, tlsCertVerification network.TLSCertVerification) credstore.Client {
	if credhub == nil {
		return credstore.Noop
	}

	return credstore.NewCredhubClient(
		credhub.APIURL,
		credhub.ID,
		credhub.Secret,
		tlsCertVerification,
	)
}
