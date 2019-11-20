package k8s

import (
	"fmt"
	"regexp"

	authenticationapi "k8s.io/api/authentication/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

//WhoAmI returns the current user/token subject
func WhoAmI(kubeclient kubernetes.Interface, token string) (string, error) {
	result, err := kubeclient.AuthenticationV1().TokenReviews().Create(&authenticationapi.TokenReview{
		Spec: authenticationapi.TokenReviewSpec{
			Token: token,
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		if k8serrors.IsForbidden(err) {
			return getUsernameFromError(err), nil
		}
		return "", err
	}

	if result.Status.Error != "" {
		return "", fmt.Errorf(result.Status.Error)
	}

	return result.Status.User.Username, nil
}

func getUsernameFromError(err error) string {
	fmt.Println(err.Error())
	re := regexp.MustCompile(`^.* User "(.*)" cannot .*$`)
	return re.ReplaceAllString(err.Error(), "$1")
}
