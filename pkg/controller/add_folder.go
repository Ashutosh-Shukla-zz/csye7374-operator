package controller

import (
	"github.com/Ashutosh-Shukla/csye7374-operator/pkg/controller/folder"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, folder.Add)
}
