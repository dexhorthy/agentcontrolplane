package taskruntoolcall

import (
	"context"

	kubechainv1alpha1 "github.com/humanlayer/smallchain/kubechain/api/v1alpha1"
	. "github.com/humanlayer/smallchain/kubechain/test/utils"
	testutils "github.com/humanlayer/smallchain/kubechain/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("TaskRunToolCall Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-taskruntoolcall"

		var ctx context.Context
		var cancel context.CancelFunc
		var typeNamespacedName = types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}

		var testTool *testutils.TestScopedTool
		var testTaskRunToolCall *kubechainv1alpha1.TaskRunToolCall

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.TODO())

			testTool = &testutils.TestScopedTool{
				Name:        "add",
				Description: "Add two numbers",
				BuiltinName: "add",
				ParametersRaw: `{
					"type": "object",
					"properties": {
						"a": { "type": "number" },
						"b": { "type": "number" }
					},
					"required": ["a", "b"]
				}`,
			}

			By("Setting up the test tool")
			testTool.Setup(k8sClient)

			// Clean up any existing TaskRunToolCall
			existingTRTC := &kubechainv1alpha1.TaskRunToolCall{}
			err := k8sClient.Get(ctx, typeNamespacedName, existingTRTC)
			if err == nil {
				Expect(k8sClient.Delete(ctx, existingTRTC)).To(Succeed())
				Eventually(func() error {
					return k8sClient.Get(ctx, typeNamespacedName, existingTRTC)
				}).Should(HaveOccurred())
			}
		})

		AfterEach(func() {
			cancel()

			By("Cleaning up the test tool")
			testTool.Teardown()

			By("Cleaning up the test TaskRunToolCall")
			if testTaskRunToolCall != nil {
				err := k8sClient.Delete(ctx, testTaskRunToolCall)
				if err == nil {
					Eventually(func() error {
						return k8sClient.Get(ctx, typeNamespacedName, testTaskRunToolCall)
					}).Should(HaveOccurred())
				}
			}
		})

		It("should successfully execute a function tool call", func() {
			testTaskRunToolCall = &kubechainv1alpha1.TaskRunToolCall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.TaskRunToolCallSpec{
					TaskRunRef: kubechainv1alpha1.LocalObjectReference{
						Name: "parent-taskrun",
					},
					ToolRef: kubechainv1alpha1.LocalObjectReference{
						Name: testTool.Name,
					},
					Arguments: `{"a": 2, "b": 3}`,
				},
			}

			By("Creating the TaskRunToolCall")
			Expect(k8sClient.Create(ctx, testTaskRunToolCall)).To(Succeed())

			By("Reconciling the TaskRunToolCall")
			eventRecorder := record.NewFakeRecorder(10)
			reconciler := &TaskRunToolCallReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				recorder: eventRecorder,
			}

			// First reconciliation - should initialize status
			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Second reconciliation executes the tool
			_, err = reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the TaskRunToolCall status")
			updatedTRTC := &kubechainv1alpha1.TaskRunToolCall{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedTRTC)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTRTC.Status.Phase).To(Equal(kubechainv1alpha1.TaskRunToolCallPhaseSucceeded))
			Expect(updatedTRTC.Status.Result).To(Equal("5"))
			Expect(updatedTRTC.Status.Status).To(Equal("Ready"))
			Expect(updatedTRTC.Status.StatusDetail).To(Equal("Tool executed successfully"))

			By("Verifying that execution events were emitted")
			ExpectEvent(eventRecorder).ToEmitEventContaining("ExecutionSucceeded")
		})

		It("should fail with invalid arguments", func() {
			testTaskRunToolCall = &kubechainv1alpha1.TaskRunToolCall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.TaskRunToolCallSpec{
					TaskRunRef: kubechainv1alpha1.LocalObjectReference{
						Name: "parent-taskrun",
					},
					ToolRef: kubechainv1alpha1.LocalObjectReference{
						Name: testTool.Name,
					},
					Arguments: `invalid json`,
				},
			}

			By("Creating the TaskRunToolCall with invalid JSON")
			Expect(k8sClient.Create(ctx, testTaskRunToolCall)).To(Succeed())

			By("Reconciling the TaskRunToolCall")
			eventRecorder := record.NewFakeRecorder(10)
			reconciler := &TaskRunToolCallReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				recorder: eventRecorder,
			}

			// First reconciliation - should initialize status
			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Second reconciliation fails validation
			_, err = reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).To(HaveOccurred())

			By("Verifying the TaskRunToolCall error status")
			updatedTRTC := &kubechainv1alpha1.TaskRunToolCall{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedTRTC)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTRTC.Status.Status).To(Equal("Error"))
			Expect(updatedTRTC.Status.StatusDetail).To(Equal("Invalid arguments JSON"))

			By("Verifying that a validation failed event was created")
			ExpectEvent(eventRecorder).ToEmitEventContaining("ExecutionFailed")
		})
	})
})
