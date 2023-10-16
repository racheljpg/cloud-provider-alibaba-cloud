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

// CreateInstanceActiveOpsTask invokes the ens.CreateInstanceActiveOpsTask API synchronously
func (client *Client) CreateInstanceActiveOpsTask(request *CreateInstanceActiveOpsTaskRequest) (response *CreateInstanceActiveOpsTaskResponse, err error) {
	response = CreateCreateInstanceActiveOpsTaskResponse()
	err = client.DoAction(request, response)
	return
}

// CreateInstanceActiveOpsTaskWithChan invokes the ens.CreateInstanceActiveOpsTask API asynchronously
func (client *Client) CreateInstanceActiveOpsTaskWithChan(request *CreateInstanceActiveOpsTaskRequest) (<-chan *CreateInstanceActiveOpsTaskResponse, <-chan error) {
	responseChan := make(chan *CreateInstanceActiveOpsTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateInstanceActiveOpsTask(request)
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

// CreateInstanceActiveOpsTaskWithCallback invokes the ens.CreateInstanceActiveOpsTask API asynchronously
func (client *Client) CreateInstanceActiveOpsTaskWithCallback(request *CreateInstanceActiveOpsTaskRequest, callback func(response *CreateInstanceActiveOpsTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateInstanceActiveOpsTaskResponse
		var err error
		defer close(result)
		response, err = client.CreateInstanceActiveOpsTask(request)
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

// CreateInstanceActiveOpsTaskRequest is the request struct for api CreateInstanceActiveOpsTask
type CreateInstanceActiveOpsTaskRequest struct {
	*requests.RpcRequest
	InstanceIds *[]string `position:"Query" name:"InstanceIds"  type:"Repeated"`
}

// CreateInstanceActiveOpsTaskResponse is the response struct for api CreateInstanceActiveOpsTask
type CreateInstanceActiveOpsTaskResponse struct {
	*responses.BaseResponse
	RequestId             string                `json:"RequestId" xml:"RequestId"`
	InstanceActiveOpsTask InstanceActiveOpsTask `json:"InstanceActiveOpsTask" xml:"InstanceActiveOpsTask"`
}

// CreateCreateInstanceActiveOpsTaskRequest creates a request to invoke CreateInstanceActiveOpsTask API
func CreateCreateInstanceActiveOpsTaskRequest() (request *CreateInstanceActiveOpsTaskRequest) {
	request = &CreateInstanceActiveOpsTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ens", "2017-11-10", "CreateInstanceActiveOpsTask", "ens", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateInstanceActiveOpsTaskResponse creates a response to parse from CreateInstanceActiveOpsTask response
func CreateCreateInstanceActiveOpsTaskResponse() (response *CreateInstanceActiveOpsTaskResponse) {
	response = &CreateInstanceActiveOpsTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
