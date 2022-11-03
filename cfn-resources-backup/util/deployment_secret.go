package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type DeploymentSecret struct {
	PublicKey  string              `json:"PublicKey"`
	PrivateKey string              `json:"PrivateKey"`
	ResourceID *ResourceIdentifier `json:"ResourceID"`
	Properties *map[string]string  `json:"Properties"`
}

func CreateDeploymentSecret(req *handler.Request, cfnID *ResourceIdentifier, publicKey string, privateKey string, properties *map[string]string) (*string, error) {
	deploySecret := &DeploymentSecret{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		ResourceID: cfnID,
		Properties: properties,
	}
	fmt.Printf("deploySecret: %v", deploySecret)
	deploySecretString, err := json.Marshal(deploySecret)
	fmt.Printf("deploySecretString: %v", deploySecretString)

	// log.Println("===============================================")
	fmt.Printf("%+v", os.Environ())
	// log.Println("===============================================")

	//sess := credentials.SessionFromCredentialsProvider(creds)
	// create a new secret from this struct with the json string

	// Create service client value configured for credentials
	// from assumed role.
	svc := secretsmanager.New(req.Session)

	//config := &aws.Config{
	//    Region: aws.String("us-east-1"),
	//}
	//svc := secretsmanager.New(session.New(), config)

	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String("MongoDB Atlas Quickstart Deployment Secret"),
		Name:         aws.String(fmt.Sprintf("%v", cfnID)),
		SecretString: aws.String(string(deploySecretString)),
	}

	result, err := svc.CreateSecret(input)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Printf("error create secret: %+v", err.Error())
		return nil, err
		//fmt.Println(err.Error())

	}
	fmt.Printf("Created secret result:%+v", result)
	return result.Name, nil

}

func GetApiKeyFromDeploymentSecret(req *handler.Request, secretName string) (DeploymentSecret, error) {
	fmt.Printf("secretName=%v\n", secretName)
	sm := secretsmanager.New(req.Session)
	output, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		fmt.Printf("Error --- %v", err.Error())
		return DeploymentSecret{}, err
	}
	// fmt.Println(*output.SecretString)
	var key DeploymentSecret
	err = json.Unmarshal([]byte(*output.SecretString), &key)
	if err != nil {
		fmt.Printf("Error --- %v", err.Error())
		return key, err
	}
	fmt.Printf("%v", key)
	return key, nil
}
