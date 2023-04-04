package madmin

import (
	"context"
	"encoding/json"
	"net/http"
)

type VacancyInfo struct {
	Enabled          bool `json:"enabled"`
	CheckInterval    int  `json:"checkInterval"` // day
	VacancyThreshold int  `json:"vacancyThreshold"`
}

// GetVacancyInfo - returns vacancy Info
func (adm *AdminClient) GetVacancyInfo(ctx context.Context) (VacancyInfo, error) {
	var vacancyInfo VacancyInfo

	resp, err := adm.executeMethod(ctx, http.MethodGet, requestData{relPath: adminAPIPrefix + "/vacancyinfo"})
	defer closeResponse(resp)
	if err != nil {
		return vacancyInfo, err
	}

	// Check response http status code
	if resp.StatusCode != http.StatusOK {
		return vacancyInfo, httpRespToErrorResponse(resp)
	}

	// Unmarshal the server's json response
	if err = json.NewDecoder(resp.Body).Decode(&vacancyInfo); err != nil {
		return vacancyInfo, err
	}

	return vacancyInfo, nil
}

// SetVacancy - set vacancy config
func (adm *AdminClient) SetVacancy(ctx context.Context, config VacancyInfo) (err error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	//econfigBytes, err := EncryptData(adm.getSecretKey(), configBytes)
	//if err != nil {
	//	return err
	//}

	reqData := requestData{
		relPath: adminAPIPrefix + "/vacancy",
		content: configBytes,
	}

	// Execute PUT on /minio/admin/v3/config to set config.
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

func (adm *AdminClient) ManualMergeVacancy(ctx context.Context) (err error) {
	reqData := requestData{
		relPath: adminAPIPrefix + "/manual-merge-vacancy",
	}

	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
