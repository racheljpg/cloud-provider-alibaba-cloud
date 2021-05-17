package service

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctx2 "k8s.io/cloud-provider-alibaba-cloud/pkg/context"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/context/shared"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/controller/helper"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/model"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/util"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func Add(mgr manager.Manager, ctx *shared.SharedContext) error {
	return add(mgr, newReconciler(mgr, ctx))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, ctx *shared.SharedContext) reconcile.Reconciler {
	recon := &ReconcileService{
		cloud:            ctx.Provider(),
		kubeClient:       mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		record:           mgr.GetEventRecorderFor("service-controller"),
		finalizerManager: helper.NewDefaultFinalizerManager(mgr.GetClient()),
	}
	return recon
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(
		"service-controller", mgr,
		controller.Options{
			Reconciler:              r,
			MaxConcurrentReconciles: 1,
		},
	)
	if err != nil {
		return err
	}
	hand := NewMapHandler()
	if err := c.Watch(&source.Kind{Type: &v1.Service{}}, hand,
		NewPredicateForServiceEvent(mgr.GetEventRecorderFor("service-controller"))); err != nil {
		return fmt.Errorf("watch resource svc error: %s", err.Error())
	}

	if err := c.Watch(&source.Kind{Type: &v1.Endpoints{}}, hand,
		NewPredicateForEndpointEvent(mgr.GetClient())); err != nil {
		return fmt.Errorf("watch resource endpoint error: %s", err.Error())
	}

	if err := c.Watch(&source.Kind{Type: &v1.Node{}}, hand,
		NewPredicateForNodeEvent(mgr.GetEventRecorderFor("service-controller"))); err != nil {
		return fmt.Errorf("watch resource node error: %s", err.Error())
	}
	return nil
}

// ReconcileService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileService{}

// ReconcileService reconciles a AutoRepair object
type ReconcileService struct {
	scheme *runtime.Scheme

	// client
	cloud      prvd.Provider
	kubeClient client.Client

	//record event recorder
	record           record.EventRecorder
	finalizerManager helper.FinalizerManager
}

// TODO
// 是否还需要cache机制？
func (m *ReconcileService) Reconcile(_ context.Context, request reconcile.Request) (reconcile.Result, error) {
	klog.Infof("do reconcile service %s", request.NamespacedName)
	return reconcile.Result{}, m.reconcile(request)
}

type RequestContext struct {
	ctx  context.Context
	svc  *v1.Service
	anno *AnnotationRequest

	// client
	cloud      prvd.Provider
	kubeClient client.Client
}

func (reqCtx *RequestContext) GetAnnotation() *AnnotationRequest { return reqCtx.anno }

func (reqCtx *RequestContext) SetService(svc *v1.Service) { reqCtx.svc = svc }

func NewRequestContext(ctx context.Context, svc *v1.Service, anno *AnnotationRequest,
	cloud prvd.Provider, kubeClient client.Client) *RequestContext {
	return &RequestContext{
		ctx:        ctx,
		svc:        svc,
		anno:       anno,
		cloud:      cloud,
		kubeClient: kubeClient,
	}
}

func (m *ReconcileService) reconcile(request reconcile.Request) error {
	// new context for each request
	ctx := context.Background()

	svc := &v1.Service{}
	err := m.kubeClient.Get(ctx, request.NamespacedName, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("service %s not found, skip", request.NamespacedName)
			// Request object not found, could have been deleted
			// after reconcile request.
			// Owned objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			// Return and don't requeue
			return nil
		}
		return err
	}
	anno := &AnnotationRequest{svc: svc}

	if err := validate(svc, anno); err != nil {
		return fmt.Errorf("validate svc error: %s", err.Error())
	}

	reqContext := &RequestContext{
		ctx:        ctx,
		svc:        svc,
		anno:       anno,
		cloud:      m.cloud,
		kubeClient: m.kubeClient,
	}

	// 1. check to see whither if loadbalancer deletion is needed
	if !isSLBNeeded(svc) {
		return m.cleanupLoadBalancerResources(reqContext)
	}
	err = m.reconcileLoadBalancerResources(reqContext)
	if err != nil {
		klog.Infof("reconcile loadbalancer error: %s", err.Error())
	}
	return err
}

func validate(svc *v1.Service, anno *AnnotationRequest) error {
	// safety check.
	if svc == nil {
		return fmt.Errorf("service could not be empty")
	}

	// disable public address
	if anno.Get(AddressType) == "" || anno.Get(AddressType) == "internet" {
		if ctx2.CFG.Global.DisablePublicSLB {
			return fmt.Errorf("PublicAddress SLB is Not allowed")
		}
	}
	return nil
}

func (m *ReconcileService) cleanupLoadBalancerResources(req *RequestContext) error {
	if helper.HasFinalizer(req.svc, serviceFinalizer) {
		_, err := m.buildAndApplyModel(req)
		if err != nil {
			m.record.Event(req.svc, v1.EventTypeWarning, helper.ServiceEventReasonFailedReconciled, fmt.Sprintf("Failed reconcile due to %s", err.Error()))
			return err
		}

		if err := m.finalizerManager.RemoveFinalizers(req.ctx, req.svc, serviceFinalizer); err != nil {
			m.record.Eventf(req.svc, v1.EventTypeWarning, helper.ServiceEventReasonFailedRemoveFinalizer, fmt.Sprintf("Failed remove finalizer due to %v", err))
			return err
		}
	}
	return nil
}

func (m *ReconcileService) reconcileLoadBalancerResources(req *RequestContext) error {

	if err := m.finalizerManager.AddFinalizers(req.ctx, req.svc, serviceFinalizer); err != nil {
		m.record.Event(req.svc, v1.EventTypeWarning, helper.ServiceEventReasonFailedAddFinalizer, fmt.Sprintf("Failed add finalizer due to %v", err))
		return err
	}

	lb, err := m.buildAndApplyModel(req)
	if err != nil {
		m.record.Event(req.svc, v1.EventTypeWarning, helper.ServiceEventReasonFailedReconciled, fmt.Sprintf("Failed reconcile due to %s", err.Error()))
		return err
	}

	if err := m.updateServiceStatus(req.ctx, req.svc, lb); err != nil {
		m.record.Event(req.svc, v1.EventTypeWarning, helper.ServiceEventReasonFailedUpdateStatus, fmt.Sprintf("Failed update status due to %v", err))
		return err
	}

	m.record.Event(req.svc, v1.EventTypeNormal, helper.ServiceEventReasonSuccessfullyReconciled, "Successfully reconciled")
	return nil
}

func (m *ReconcileService) buildAndApplyModel(req *RequestContext) (*model.LoadBalancer, error) {
	// build local model
	localModel, err := NewModelBuilder(req, LOCAL_MODEL).Build()
	if err != nil {
		return nil, fmt.Errorf("build slb cluster model error: %s", err.Error())
	}
	// build remote model
	remoteModel, err := NewModelBuilder(req, REMOTE_MODEL).Build()
	if err != nil {
		return nil, fmt.Errorf("build slb cloud model error: %s", err.Error())
	}
	// apply model
	if err := NewLBModelApplier(req).Apply(localModel, remoteModel); err != nil {
		return nil, fmt.Errorf("apply model error: %s", err.Error())
	}
	return remoteModel, nil
}

func (m *ReconcileService) updateServiceStatus(ctx context.Context, svc *v1.Service, lb *model.LoadBalancer) error {

	// TODO
	if len(svc.Status.LoadBalancer.Ingress) != 1 ||
		svc.Status.LoadBalancer.Ingress[0].IP == "" {
		svcOld := svc.DeepCopy()
		svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{
			{
				IP: lb.LoadBalancerAttribute.Address,
			},
		}
		if err := m.kubeClient.Status().Patch(ctx, svc, client.MergeFrom(svcOld)); err != nil {
			return fmt.Errorf("%s failed to update service status:, error: %s", util.Key(svc), err.Error())
		}
	}
	return nil

}
