package apis

import (
	"github.com/Ashutosh-Shukla/csye7374-operator/pkg/apis/csye7374/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
