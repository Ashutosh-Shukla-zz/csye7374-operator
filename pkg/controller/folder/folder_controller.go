package folder

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"io/ioutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"

	csye7374v1alpha1 "github.com/Ashutosh-Shukla/csye7374-operator/pkg/apis/csye7374/v1alpha1"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_folder")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Folder Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFolder{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("folder-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}


	err = c.Watch(&source.Kind{Type: &csye7374v1alpha1.Folder{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource Folder
	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Folder
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &csye7374v1alpha1.Folder{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFolder implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFolder{}

// ReconcileFolder reconciles a Folder object
type ReconcileFolder struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

const FolderFinalizerName  = "finalizer.csye7374.com"
// Reconcile reads that state of the cluster for a Folder object and makes changes based on the state read
// and what is in the Folder.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFolder) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Folder")

	// Fetch the Folder instance
	instance := &csye7374v1alpha1.Folder{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//Get Operator Secret from file to be mounted in dir /usr/local/etc/operator/
	awsAccessKeyIDbyte, err := ioutil.ReadFile("/usr/local/etc/operator/aws_access_key_id")
	if err != nil {
		return reconcile.Result{}, err
	}
	awsAccessKeyID := strings.TrimRight(string(awsAccessKeyIDbyte),  "\r\n")
	awsSecretAccessKeybyte,err := ioutil.ReadFile("/usr/local/etc/operator/aws_secret_access_key")
	if err != nil {
		return reconcile.Result{}, err
	}
	awsSecretAccessKey:= strings.TrimRight(string(awsSecretAccessKeybyte),  "\r\n")
	bucketbyte,err := ioutil.ReadFile("/usr/local/etc/operator/bucketname")
	if err != nil {
		return reconcile.Result{}, err
	}
	bucket:=strings.TrimRight(string(bucketbyte),  "\r\n")
	regionbyte,err := ioutil.ReadFile("/usr/local/etc/operator/aws_region")
	if err != nil {
		return reconcile.Result{}, err
	}
	region:=strings.TrimRight(string(regionbyte),  "\r\n")
	token := ""
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)

	_, err = creds.Get()
	if err != nil {
		return reconcile.Result{}, err
	}
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)

	//Create S3 Folder in bucket
	err = createS3Folder(cfg, bucket, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Create Iam User
	createdAwsUser, err := createIamUser(cfg, instance.Spec.Username)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}


func createS3Folder(cfg *aws.Config, bucket string, instance *csye7374v1alpha1.Folder) error {
	svc := s3.New(session.New(), cfg)
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}

	bucketexists := true

	result, err := svc.GetBucketLocation(input)
	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == s3.ErrCodeNoSuchBucket {
		bucketexists = false
		return err
	}
	fmt.Print(result)

	if bucketexists {
		key := instance.Spec.Username + "/"
		_, err := svc.PutObject(&s3.PutObjectInput{
			Body:   strings.NewReader("Hello World!"),
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			return err
		}

	}
	if !bucketexists {
		log.Error(err, "Bucket does not exist")
		return err
	}
	return nil
}

func createIamUser(cfg *aws.Config, username string) (*iam.User, error) {
	svc := iam.New(session.New(), cfg)

	userOutput, err := svc.GetUser(&iam.GetUserInput{
		UserName: aws.String(username),
	})

	if err == nil {
		return userOutput.User, nil
	}

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		userOutputPod, err := svc.CreateUser(&iam.CreateUserInput{
			UserName: aws.String(username),
		})

		if err != nil {
			fmt.Println("CreateUser Error", err)
			return nil, err
		}

		return userOutputPod.User, nil
	} else {
		log.Error(err, "Get User error")
		return nil, err
	}
}


