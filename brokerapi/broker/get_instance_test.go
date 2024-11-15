package broker_test

import (
	"net/http"

	"code.cloudfoundry.org/lager/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi/v11/domain"
	"github.com/pivotal-cf/brokerapi/v11/domain/apiresponses"
	"github.com/pivotal-cf/brokerapi/v11/middlewares"
	"golang.org/x/net/context"

	"github.com/cloudfoundry/cloud-service-broker/v2/brokerapi/broker/brokerfakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/dbservice/models"
	"github.com/cloudfoundry/cloud-service-broker/v2/internal/storage"
	pkgBroker "github.com/cloudfoundry/cloud-service-broker/v2/pkg/broker"
	pkgBrokerFakes "github.com/cloudfoundry/cloud-service-broker/v2/pkg/broker/brokerfakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/pkg/credstore/credstorefakes"
	"github.com/cloudfoundry/cloud-service-broker/v2/pkg/varcontext"
	"github.com/cloudfoundry/cloud-service-broker/v2/utils"

	"github.com/cloudfoundry/cloud-service-broker/v2/brokerapi/broker"
)

var _ = Describe("GetInstance", func() {

	const (
		orgID      = "test-org-id"
		spaceID    = "test-space-id"
		planID     = "test-plan-id"
		serviceID  = "test-service-id"
		offeringID = "test-service-id"
		instanceID = "test-instance-id"
	)

	var (
		serviceBroker *broker.ServiceBroker

		fakeStorage         *brokerfakes.FakeStorage
		fakeServiceProvider *pkgBrokerFakes.FakeServiceProvider
		fakeCredStore       *credstorefakes.FakeCredStore

		brokerConfig *broker.BrokerConfig

		provisionParams = storage.JSONObject{
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
					ID:   offeringID,
					Name: "test-service",
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

	When("instance exists and provision succeeded", func() {
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
			fakeServiceProvider.PollInstanceReturns(true, "", models.ProvisionOperationType, nil) // Operation status is provision succeeded
			fakeStorage.GetProvisionRequestDetailsReturns(provisionParams, nil)
		})

		It("returns instance details", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())

			By("validating response")
			Expect(response.ServiceID).To(Equal(serviceID))
			Expect(response.PlanID).To(Equal(planID))
			Expect(response.Parameters).To(BeEquivalentTo(provisionParams))
			Expect(response.DashboardURL).To(BeEmpty()) // Broker does not set dashboard URL
			Expect(response.Metadata).To(BeZero())      // Broker does not support instance metadata

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating service provider asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(1))

			By("validating storage is asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(1))
		})
	})

	When("instance does not exist", func() {
		BeforeEach(func() {
			fakeStorage.ExistsServiceInstanceDetailsReturns(false, nil)
		})
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			s, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                            // must be a failure response
			Expect(s.Error()).To(Equal("not found"))                          // must contain "Not Found" error message
			Expect(s.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is not asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(0))

			By("validating service provider is not asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(0))

			By("validating storage is not asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(0))
		})
	})

	When("instance exists and provision is in progress", func() {
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
			fakeServiceProvider.PollInstanceReturns(false, "", models.ProvisionOperationType, nil) // Operation status is provision in progress
		})
		// According to OSB Spec, broker must return 404 in case provision is in progress
		It("returns status code 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			s, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                            // must be a failure response
			Expect(s.Error()).To(Equal("not found"))                          // must contain "not found" error message
			Expect(s.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating service provider is asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(1))

			By("validating storage is not asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(0))
		})

	})

	When("instance exists and update is in progress", func() {
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
			fakeServiceProvider.PollInstanceReturns(false, "", models.UpdateOperationType, nil) // Operation status is update in progress
		})
		It("returns status code 422 (Unprocessable Entity) and error code ConcurrencyError", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID, PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			s, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                                       // must be a failure response
			Expect(s.Error()).To(Equal("ConcurrencyError"))                              // must contain "ConcurrencyError" error message
			Expect(s.ValidatedStatusCode(nil)).To(Equal(http.StatusUnprocessableEntity)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating service provider is asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(1))

			By("validating storage is not asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(0))
		})

	})

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
			fakeServiceProvider.PollInstanceReturns(true, "", models.ProvisionOperationType, nil) // Operation status is provision succeeded
			fakeStorage.GetProvisionRequestDetailsReturns(provisionParams, nil)
		})
		It("ignores service_id and returns instance details", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{PlanID: fetchPlanID})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.ServiceID).To(Equal(serviceID))
		})
	})

	When("service_id does not match service for instance", func() {
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
			fakeServiceProvider.PollInstanceReturns(true, "", models.ProvisionOperationType, nil) // Operation status is provision succeeded
			fakeStorage.GetProvisionRequestDetailsReturns(provisionParams, nil)
		})
		It("returns 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: "otherService", PlanID: fetchPlanID})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			s, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                            // must be a failure response
			Expect(s.Error()).To(Equal("not found"))                          // must contain "not found" error message
			Expect(s.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating service provider is asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(0))

			By("validating storage is not asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(0))
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
			fakeServiceProvider.PollInstanceReturns(true, "", models.ProvisionOperationType, nil) // Operation status is provision succeeded
			fakeStorage.GetProvisionRequestDetailsReturns(provisionParams, nil)
		})
		It("ignores plan_id and returns instance details", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.ServiceID).To(Equal(serviceID))
		})
	})

	When("plan_id does not match plan for instance", func() {
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
			fakeServiceProvider.PollInstanceReturns(true, "", models.ProvisionOperationType, nil) // Operation status is provision succeeded
			fakeStorage.GetProvisionRequestDetailsReturns(provisionParams, nil)
		})
		It("returns 404 (not found)", func() {
			const expectedHeader = "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ=="
			newContext := context.WithValue(context.Background(), middlewares.OriginatingIdentityKey, expectedHeader)

			response, err := serviceBroker.GetInstance(newContext, fetchInstanceID, domain.FetchInstanceDetails{ServiceID: fetchServiceID, PlanID: "otherPlan"})

			By("validating response")
			Expect(response).To(BeZero())

			By("validating error")
			s, isFailureResponse := err.(*apiresponses.FailureResponse)
			Expect(isFailureResponse).To(BeTrue())                            // must be a failure response
			Expect(s.Error()).To(Equal("not found"))                          // must contain "not found" error message
			Expect(s.ValidatedStatusCode(nil)).To(Equal(http.StatusNotFound)) // status code must be 404

			By("validating storage is asked whether instance exists")
			Expect(fakeStorage.ExistsServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating storage is asked for instance details")
			Expect(fakeStorage.GetServiceInstanceDetailsCallCount()).To(Equal(1))

			By("validating service provider is asked for instance status")
			Expect(fakeServiceProvider.PollInstanceCallCount()).To(Equal(0))

			By("validating storage is not asked for provision request details")
			Expect(fakeStorage.GetProvisionRequestDetailsCallCount()).To(Equal(0))
		})
	})
})
