package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

var secretsKeyValue map[string]string

func main() {
	secretsKeyValue := make(map[string]string)

	vaultUri := os.Getenv("VAULT_URI")

	GetAllSecrets(vaultUri, secretsKeyValue)
	tmpl, err := template.ParseFiles(os.Getenv("INPUT_TEMPLATE_FILE"))

	if err != nil {
		log.Fatalln(err)
	}

	secretsFilePath := os.Getenv("OUTPUT_FILE_PATH")

	secretsFile, err := os.Create(secretsFilePath)
	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
	}

	err = tmpl.Execute(secretsFile, secretsKeyValue)

	if err != nil {
		log.Fatalln(err)
	}
}

func GetAllSecrets(vaultUri string, secretsKeyValue map[string]string) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	client, err := azsecrets.NewClient(vaultUri, cred, nil)

	if err != nil {
		fmt.Printf("err: %v", err)
	}

	var secretNames []string

	pager := client.ListSecrets(nil)
	for pager.NextPage(context.TODO()) {
		for _, v := range pager.PageResponse().Secrets {

			if *v.Attributes.Enabled {
				fmt.Printf("Secret Name: %s\tSecret Tags: %v\n", *v.ID, v.Tags)

				// split secret url eg. https://contoso.vault.azure.net/secrets/envVar2
				// in order to get the env name eg. 'envVar2'
				secretName := strings.Split(*v.ID, "/")
				secretNames = append(secretNames, secretName[len(secretName)-1])
			}
		}
	}

	for _, secret := range secretNames {
		resp, err := client.GetSecret(context.Background(), secret, nil)
		if err != nil {
			fmt.Println(err)
		}
		secretsKeyValue[secret] = *resp.Value

	}
}
