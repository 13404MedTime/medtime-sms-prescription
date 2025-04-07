package function

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	fcm "github.com/appleboy/go-fcm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Datas This is response struct from create
type Datas struct {
	Data struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	} `json:"data"`
}

// ClientApiResponse This is get single api response
type ClientApiResponse struct {
	Data ClientApiData `json:"data"`
}

type ClientApiData struct {
	Data ClientApiResp `json:"data"`
}

type ClientApiResp struct {
	Response map[string]interface{} `json:"response"`
}

type Response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// NewRequestBody's Data (map) field will be in this structure
//.   fields
// objects_ids []string
// table_slug string
// object_data map[string]interface
// method string
// app_id string

// but all field will be an interface, you must do type assertion

type HttpRequest struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
	Params  url.Values  `json:"params"`
	Body    []byte      `json:"body"`
}

type AuthData struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type NewRequestBody struct {
	RequestData HttpRequest            `json:"request_data"`
	Auth        AuthData               `json:"auth"`
	Data        map[string]interface{} `json:"data"`
}
type Request struct {
	Data map[string]interface{} `json:"data"`
}

// GetListClientApiResponse This is get list api response
type GetListClientApiResponse struct {
	Data GetListClientApiData `json:"data"`
}

type GetListClientApiData struct {
	Data GetListClientApiResp `json:"data"`
}

type GetListClientApiResp struct {
	Response []map[string]interface{} `json:"response"`
}

type UserNotification struct {
	Title       string
	Description string
	Image       string
	Fcm         string
	Body        string
	Platform    int // 0 - IOS, 1 - Android
}

func DoRequest(url string, method string, body interface{}, appId string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("authorization", "API-KEY")
	request.Header.Add("X-API-KEY", appId)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

// Handle a serverless request
func Handle(req []byte) string {
	Send(string(req))
	var response Response
	var request NewRequestBody
	//const urlConst = ""

	err := json.Unmarshal(req, &request)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling request"}
		response.Status = "error"
		responseByte, _ := json.Marshal(response)
		return string(responseByte)
	}
	if request.Data["app_id"] == nil {
		response.Data = map[string]interface{}{"message": "App id required"}
		response.Status = "error"
		responseByte, _ := json.Marshal(response)
		return string(responseByte)
	}
	// appId := request.Data["app_id"].(string)

	// // you may change table slug  it's related your business logic
	// var tableSlug = ""

	response.Data = map[string]interface{}{}
	response.Status = "done" //if all will be ok else "error"
	responseByte, _ := json.Marshal(response)

	return string(responseByte)
}

func SendNotification(notification UserNotification) {
	msg := &fcm.Message{
		To: notification.Fcm,
	}
	if notification.Platform == 1 {
		msg.Data = map[string]interface{}{
			"title": notification.Title,
			"body":  notification.Body,
		}
	} else if notification.Platform == 0 {
		msg.Notification = &fcm.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		}
	}
	// Create a FCM client to send the message.
	client, _ := fcm.NewClient("AAAA8OkqzfI:APA91bHyBn537ADTKHRwSN_JsjvtaVlY_bJATanjZodV5whU379qKp8M0470kRkeOzVMQxRw1e5vYVta-cy8R1QnQ_y6f_dGDM5eYzEtseB6cxrNFnkDwkGgIZ44jxsoyE6ORUcMtqHF")
	client.Send(msg)

}

func Send(text string) {
	bot, _ := tgbotapi.NewBotAPI("6041044802:AAEDdr0uD4SkxnnGctOOsA2Ua3Ovy-7Sy0A")

	msg := tgbotapi.NewMessage(266798451, text)

	bot.Send(msg)
}

// func GetListObject(url, tableSlug, appId string, request Request) (GetListClientApiResponse, error, Response) {
// 	response := Response{}

// 	getListResponseInByte, err := DoRequest(url+"/v1/object/get-list/{table_slug}?from-ofs=true", "POST", request, appId)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while getting single object"}
// 		response.Status = "error"
// 		return GetListClientApiResponse{}, errors.New("error"), response
// 	}
// 	var getListObject GetListClientApiResponse
// 	err = json.Unmarshal(getListResponseInByte, &getListObject)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while unmarshalling get list object"}
// 		response.Status = "error"
// 		return GetListClientApiResponse{}, errors.New("error"), response
// 	}
// 	return getListObject, nil, response
// }

// func GetSingleObject(url, tableSlug, appId, guid string) (ClientApiResponse, error, Response) {
// 	response := Response{}

// 	var getSingleObject ClientApiResponse
// 	getSingleResponseInByte, err := DoRequest(url+"/v1/object/{table_slug}/{guid}?from-ofs=true", "GET", nil, appId)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while getting single object"}
// 		response.Status = "error"
// 		return ClientApiResponse{}, errors.New("error"), response
// 	}
// 	err = json.Unmarshal(getSingleResponseInByte, &getSingleObject)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while unmarshalling single object"}
// 		response.Status = "error"
// 		return ClientApiResponse{}, errors.New("error"), response
// 	}
// 	return getSingleObject, nil, response
// }

// func CreateObject(url, tableSlug, appId string, request Request) (Datas, error, Response) {
// 	response := Response{}

// 	var createdObject Datas
// 	createObjectResponseInByte, err := DoRequest(url+"/v1/object/{table_slug}?from-ofs=true", "POST", request, appId)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while creating object"}
// 		response.Status = "error"
// 		return Datas{}, errors.New("error"), response
// 	}
// 	err = json.Unmarshal(createObjectResponseInByte, &createdObject)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while unmarshalling create object object"}
// 		response.Status = "error"
// 		return Datas{}, errors.New("error"), response
// 	}
// 	return createdObject, nil, response
// }

// func UpdateObject(url, tableSlug, appId string, request Request) (error, Response) {
// 	response := Response{}

// 	_, err := DoRequest(url+"/v1/object/{table_slug}?from-ofs=true", "PUT", request, appId)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while updating object"}
// 		response.Status = "error"
// 		return errors.New("error"), response
// 	}
// 	return nil, response
// }

// func DeleteObject(url, tableSlug, appId, guid string) (error, Response) {
// 	response := Response{}

// 	_, err := DoRequest(url+"/v1/object/{table_slug}/{guid}?from-ofs=true", "DELETE", Request{}, appId)
// 	if err != nil {
// 		response.Data = map[string]interface{}{"message": "Error while updating object"}
// 		response.Status = "error"
// 		return errors.New("error"), response
// 	}
// 	return nil, response
// }
