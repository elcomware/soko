package actions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"elcomware/app/models"
)

type IssueDigitalLicenseAction struct{}

func NewIssueDigitalLicenseAction() *IssueDigitalLicenseAction {
	return &IssueDigitalLicenseAction{}
}

func (a *IssueDigitalLicenseAction) Execute(product *models.Product) (licenseKey string, downloadUrl string) {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	hexStr := strings.ToUpper(hex.EncodeToString(bytes))
	licenseKey = fmt.Sprintf("SOKO-%s-%s", hexStr[0:8], hexStr[8:16])

	downloadUrl = product.DigitalDownloadUrl
	if downloadUrl == "" {
		downloadUrl = fmt.Sprintf("https://soko-storage.local/downloads/%s.zip", product.ID)
	}

	return licenseKey, downloadUrl
}

