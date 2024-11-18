package broker_test

import (
	"context"
	"net/http"

	"code.cloudfoundry.org/lager/v3"
	"github.com/cloudfoundry/cloud-service-broker/v2/brokerapi/broker"
	"github.com/cloudfoundry/cloud-service-broker/v2/brokerapi/broker/brokerfakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/internal/storage"
	pkgBroker "github.com/cloudfoundry/cloud-service-broker/v2/pkg/broker"
	pkgBrokerFakes "github.com/cloudfoundry/cloud-service-broker/v2/pkg/broker/brokerfakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/pkg/credstore/credstorefakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/pkg/varcontext"
	"github.com/cloudfoundry/cloud-service-broker/v2/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi/v11/domain"
	"github.com/pivotal-cf/brokerapi/v11/domain/apiresponses"
	"github.com/pivotal-cf/brokerapi/v11/middlewares"
)

var _ = Describe("GetBinding", func() {
	const (
		orgID      = "test-org-id"
		spaceID    = "test-space-id"
		planID     = "test-plan-id"
		serviceID  = "test-service-id"
		offeringID = "test-service-id"
		instanceID = "test-instance-id"
		bindingID  = "test-binding-id"
	)

	var (
		serviceBroker *broker.ServiceBroker

		fakeStorage         *brokerfakes.FakeStorage
		fakeServiceProvider *pkgBrokerFakes.FakeServiceProvider
		fakeCredStore       *credstorefakes.FakeCredStore

		brokerConfig *broker.BrokerConfig

		bindingParams = storage.JSONObject{
			"param1": "value1",
			"param2": 3,
			"param3": true,
			"param4": []string{"a", "b", "c"},
			"param5": map[string]string{"key1": "value", "key2": "value"},
			"param6": struct {
				A string
				B string
			}{"a", "b"},
		}

		fetchBindingID  = bindingID
		fetchInstanceID = instanceID
		fetchServiceID  = serviceID
		fetchPlanID     = planID
	)

	BeforeEach(func() {
		fakeStorage = &brokerfakes.FakeStorage{}
		fakeServiceProvider = &pkgBrokerFakes.FakeServiceProvider{}
		fakeCredStore = &credstorefakes.FakeCredStore{}

		providerBuilder := func(logger lager.Logger, store pkgBroker.ServiceProviderStorage) pkgBroker.ServiceProvider {
			return fakeServiceProvider
		}
		planUpdatable := true

		brokerConfig = &broker.BrokerConfig{
			Registry: pkgBroker.BrokerRegistry{
				"test-service": &pkgBroker.ServiceDefinition{
					ID:       offeringID,
					Name:     "test-service",
					Bindable: true,
					Plans: []pkgBroker.ServicePlan{
						{
							ServicePlan: domain.ServicePlan{
								ID:            planID,
								Name:          "test-plan",
								PlanUpdatable: &planUpdatable,
							},
							ServiceProperties: map[string]any{
								"plan-defined-key":       "plan-defined-value",
								"other-plan-defined-key": "other-plan-defined-value",
							},
						},
					},
					BindComputedVariables: []varcontext.DefaultVariable{
						{Name: "copyOriginatingIdentity", Default: "${json.marshal(request.x_broker_api_originating_identity)}", Overwrite: true},
					},
					ProviderBuilder: providerBuilder,
				},
			},
			Credstore: fakeCredStore,
		}

		serviceBroker = must(broker.New(brokerConfig, fakeStorage, utils.NewLogger("unbind-test-with-credstore")))
	})

	When("binding exsists", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(true, nil)
			fakeStorage.GetBindRequestDetailsReturns(bindingParams, nil)
		})
		It("returns binding details and parameters", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.SyslogDrainURL).To(BeZero())
			Expect(response.RouteServiceURL).To(BeZero())
			Expect(response.VolumeMounts).To(BeZero())
			Expect(response.Parameters).To(BeEquivalentTo(bindingParams))
			Expect(response.Endpoints).To(BeZero())

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(1))

			By("validating storage is asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(1))
		})
		It("does not return binding credentials", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.Credentials).To(BeZero())

			By("validating storage is not asked for binding credentials")
			Expect(fakeStorage.GetServiceBindingCredentialsCallCount()).To(Equal(0))
		})
		It("does not return binding metadata", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.Metadata).To(BeZero())
		})
	})

	When("service is not bindable", func() {
		BeforeEach(func() {
			brokerConfig.Registry["test-service"].Bindable = false
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
		})
		It("returns status code 400 (bad request)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			apiErr, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())
			Expect(apiErr.Error()).To(Equal("bad request"))
			Expect(apiErr.ValidatedStatusCode(nil)).To(Equal(http.StatusBadRequest))

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is not asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(0))

			By("validating storage is not asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(0))
		})
	})

	When("instance does not exist", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(false, nil)
		})
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			apiErr, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                                 // must be a failure response
			Expect(apiErr.Error()).To(Equal("not found"))                          // must contain "Not Found" error message
			Expect(apiErr.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is not asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(0))

			By("validating storage is not asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(0))

			By("validating storage is not asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(0))
		})
	})

	When("binding does not exist", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(false, nil)
		})
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			apiErr, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                                 // must be a failure response
			Expect(apiErr.Error()).To(Equal("not found"))                          // must contain "Not Found" error message
			Expect(apiErr.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(1))

			By("validating storage is not asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(0))
		})
	})

	// with the current implementation, bind is a synchroneous operation which waits for all resources to be created
	//   bindings are stored only in case resource creation is successfull
	//   -> it is not possible to retrieve a binding that is in progress from the store
	// When("binding operation is in progress", func() {
	// 	It("returns status code 404 (not found)", func() {
	// 		Fail("not implemented")
	// 	})
	// })

	When("service_id is not set", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(true, nil)
			fakeStorage.GetBindRequestDetailsReturns(bindingParams, nil)
		})
		It("ignores service_id and returns binding details", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.Parameters).To(BeEquivalentTo(bindingParams))
		})
	})

	When("service_id does not match service for binding", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(true, nil)
			fakeStorage.GetBindRequestDetailsReturns(bindingParams, nil)
		})
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: "otherService", PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			apiErr, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                                 // must be a failure response
			Expect(apiErr.Error()).To(Equal("not found"))                          // must contain "Not Found" error message
			Expect(apiErr.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is not asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(0))

			By("validating storage is not asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(0))
		})
	})

	When("plan_id is not set", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(true, nil)
			fakeStorage.GetBindRequestDetailsReturns(bindingParams, nil)
		})
		It("ignores plan_id and returns binding details", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.Parameters).To(BeEquivalentTo(bindingParams))
		})
	})

	When("plan_id does not match plan for binding", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(true, nil)
			fakeStorage.GetServiceInstanceDetailsReturns(
				storage.ServiceInstanceDetails{
					GUID:             instanceID,
					Name:             "test-instance",
					Outputs:          storage.JSONObject{},
					ServiceGUID:      serviceID,
					PlanGUID:         planID,
					SpaceGUID:        spaceID,
					OrganizationGUID: orgID,
				}, nil)
			fakeStorage.ExistsServiceBindingCredentialsReturns(true, nil)
			fakeStorage.GetBindRequestDetailsReturns(bindingParams, nil)
		})
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetBinding(newContext, fetchInstanceID, fetchBindingID, domain.FetchBindingDetails{ServiceID: fetchServiceID, PlanID: "otherPlan"})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			apiErr, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                                 // must be a failure response
			Expect(apiErr.Error()).To(Equal("not found"))                          // must contain "Not Found" error message
			Expect(apiErr.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is not asked whether binding exists")
			Expect(fakeStorage.ExistsServiceBindingCredentialsCallCount()).To(Equal(0))

			By("validating storage is not asked for bind request details")
			Expect(fakeStorage.GetBindRequestDetailsCallCount()).To(Equal(0))
		})
	})
})
