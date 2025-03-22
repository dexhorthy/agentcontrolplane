package taskrun

import (
	"context"

	kubechainv1alpha1 "github.com/humanlayer/smallchain/kubechain/api/v1alpha1"
	testutils "github.com/humanlayer/smallchain/kubechain/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("TaskRun Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-taskrun"
		const taskName = "test-task"
		const agentName = "test-agent"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}

		BeforeEach(func() {
			// Clean up any existing TaskRun
			existingTaskRun := &kubechainv1alpha1.TaskRun{}
			var err error // Declare err at the start
			err = k8sClient.Get(ctx, typeNamespacedName, existingTaskRun)
			if err == nil {
				Expect(k8sClient.Delete(ctx, existingTaskRun)).To(Succeed())
				Eventually(func() error {
					return k8sClient.Get(ctx, typeNamespacedName, existingTaskRun)
				}).Should(HaveOccurred())
			}

			// Create test Agent
			By("Creating a test Agent")
			agent := &kubechainv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.AgentSpec{
					LLMRef: kubechainv1alpha1.LocalObjectReference{
						Name: "test-llm",
					},
					System: "Test agent",
				},
			}
			Expect(k8sClient.Create(ctx, agent)).To(Succeed())

			// Mark Agent as ready
			agent.Status.Ready = true
			agent.Status.Status = "Ready"
			agent.Status.StatusDetail = "Ready for testing"
			Expect(k8sClient.Status().Update(ctx, agent)).To(Succeed())

			// Create test Task
			By("Creating a test Task")
			task := &kubechainv1alpha1.Task{
				ObjectMeta: metav1.ObjectMeta{
					Name:      taskName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.TaskSpec{
					AgentRef: kubechainv1alpha1.LocalObjectReference{
						Name: agentName,
					},
					Message: "Test task",
				},
			}
			Expect(k8sClient.Create(ctx, task)).To(Succeed())

			// Mark Task as ready
			task.Status.Ready = true
			task.Status.Status = "Ready"
			task.Status.StatusDetail = "Ready for testing"
			Expect(k8sClient.Status().Update(ctx, task)).To(Succeed())
		})

		AfterEach(func() {
			// Cleanup test resources
			By("Cleanup the test Agent")
			agent := &kubechainv1alpha1.Agent{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: agentName, Namespace: "default"}, agent)
			if err == nil {
				Expect(k8sClient.Delete(ctx, agent)).To(Succeed())
			}

			By("Cleanup the test Task")
			task := &kubechainv1alpha1.Task{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: taskName, Namespace: "default"}, task)
			if err == nil {
				Expect(k8sClient.Delete(ctx, task)).To(Succeed())
			}

			By("Cleanup the test TaskRun")
			taskRun := &kubechainv1alpha1.TaskRun{}
			err = k8sClient.Get(ctx, typeNamespacedName, taskRun)
			if err == nil {
				Expect(k8sClient.Delete(ctx, taskRun)).To(Succeed())
			}
		})

		It("should successfully validate a taskrun with valid task", func() {
			By("creating the taskrun with valid task reference")
			taskRun := &kubechainv1alpha1.TaskRun{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.TaskRunSpec{
					TaskRef: kubechainv1alpha1.LocalObjectReference{
						Name: taskName,
					},
				},
			}
			Expect(k8sClient.Create(ctx, taskRun)).To(Succeed())

			By("reconciling the taskrun")
			eventRecorder := record.NewFakeRecorder(10)
			reconciler := &TaskRunReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				recorder: eventRecorder,
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("checking the taskrun status")
			updatedTaskRun := &kubechainv1alpha1.TaskRun{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedTaskRun)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTaskRun.Status.Ready).To(BeTrue())
			Expect(updatedTaskRun.Status.Status).To(Equal("Ready"))
			Expect(updatedTaskRun.Status.StatusDetail).To(Equal("Ready to send to LLM"))
			Expect(updatedTaskRun.Status.Phase).To(Equal(kubechainv1alpha1.TaskRunPhaseReadyForLLM))

			By("checking that validation success event was created")
			testutils.ExpectEvent(eventRecorder).ToEmitEventContaining("ValidationSucceeded")
		})

		It("should fail validation with non-existent task", func() {
			By("creating the taskrun with invalid task reference")
			taskRun := &kubechainv1alpha1.TaskRun{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: kubechainv1alpha1.TaskRunSpec{
					TaskRef: kubechainv1alpha1.LocalObjectReference{
						Name: "nonexistent-task",
					},
				},
			}
			Expect(k8sClient.Create(ctx, taskRun)).To(Succeed())

			By("reconciling the taskrun")
			eventRecorder := record.NewFakeRecorder(10)
			reconciler := &TaskRunReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				recorder: eventRecorder,
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`"nonexistent-task" not found`))

			By("checking the taskrun status")
			updatedTaskRun := &kubechainv1alpha1.TaskRun{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedTaskRun)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTaskRun.Status.Ready).To(BeFalse())
			Expect(updatedTaskRun.Status.Status).To(Equal("Error"))
			Expect(updatedTaskRun.Status.StatusDetail).To(ContainSubstring(`"nonexistent-task" not found`))

			By("checking that a failure event was created")
			testutils.ExpectEvent(eventRecorder).ToEmitEventContaining("ValidationFailed")
		})
	})
})
