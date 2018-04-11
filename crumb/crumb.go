package crumb

import (
       "net/http"
)

func GenerateCrumb(req *http.Request) (string, error) {

     return "FIXME", nil
}

func ValidateCrumb(crumb_var string) (bool, error) {

     return true, nil
}
