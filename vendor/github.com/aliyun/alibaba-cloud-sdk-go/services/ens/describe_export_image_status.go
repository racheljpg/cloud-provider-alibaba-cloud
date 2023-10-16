package ens

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeExportImageStatus invokes the ens.DescribeExportImageStatus API synchronously
func (client *Client) DescribeExportImageStatus(request *DescribeExportImageStatusRequest) (response *DescribeExportImageStatusResponse, err error) {
	response = CreateDescribeExportImageStatusResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeExportImageStatusWithChan invokes the ens.DescribeExportImageStatus API asynchronously
func (client *Client) DescribeExportImageStatusWithChan(request *DescribeExportImageStatusRequest) (<-chan *DescribeExportImageStatusResponse, <-chan error) {
	responseChan := make(chan *DescribeExportImageStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeExportImageStatus(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeExportImageStatusWithCallback invokes the ens.DescribeExportImageStatus API asynchronously
func (client *Client) DescribeExportImageStatusWithCallback(request *DescribeExportImageStatusRequest, callback func(response *DescribeExportImageStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeExportImageStatusResponse
		var err error
		defer close(result)
		response, err = client.DescribeExportImageStatus(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeExportImageStatusRequest is the request struct for api DescribeExportImageStatus
type DescribeExportImageStatusRequest struct {
	*requests.RpcRequest
	ImageId string `position:"Query" name:"ImageId"`
}

// DescribeExportImageStatusResponse is the response struct for api DescribeExportImageStatus
type DescribeExportImageStatusResponse struct {
	*responses.BaseResponse
	ImageExportStatus string `json:"ImageExportStatus" xml:"ImageExportStatus"`
	RequestId         string `json:"RequestId" xml:"RequestId"`
}

// CreateDescribeExportImageStatusRequest creates a request to invoke DescribeExportImageStatus API
func CreateDescribeExportImageStatusRequest() (request *DescribeExportImageStatusRequest) {
	request = &DescribeExportImageStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ens", "2017-11-10", "DescribeExportImageStatus", "ens", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeExportImageStatusResponse creates a response to parse from DescribeExportImageStatus response
func CreateDescribeExportImageStatusResponse() (response *DescribeExportImageStatusResponse) {
	response = &DescribeExportImageStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
