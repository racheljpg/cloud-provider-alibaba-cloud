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

// DescribeAICImages invokes the ens.DescribeAICImages API synchronously
func (client *Client) DescribeAICImages(request *DescribeAICImagesRequest) (response *DescribeAICImagesResponse, err error) {
	response = CreateDescribeAICImagesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAICImagesWithChan invokes the ens.DescribeAICImages API asynchronously
func (client *Client) DescribeAICImagesWithChan(request *DescribeAICImagesRequest) (<-chan *DescribeAICImagesResponse, <-chan error) {
	responseChan := make(chan *DescribeAICImagesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAICImages(request)
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

// DescribeAICImagesWithCallback invokes the ens.DescribeAICImages API asynchronously
func (client *Client) DescribeAICImagesWithCallback(request *DescribeAICImagesRequest, callback func(response *DescribeAICImagesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAICImagesResponse
		var err error
		defer close(result)
		response, err = client.DescribeAICImages(request)
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

// DescribeAICImagesRequest is the request struct for api DescribeAICImages
type DescribeAICImagesRequest struct {
	*requests.RpcRequest
	ImageId    string `position:"Query" name:"ImageId"`
	PageNumber string `position:"Query" name:"PageNumber"`
	PageSize   string `position:"Query" name:"PageSize"`
	ImageUrl   string `position:"Query" name:"ImageUrl"`
}

// DescribeAICImagesResponse is the response struct for api DescribeAICImages
type DescribeAICImagesResponse struct {
	*responses.BaseResponse
	RequestId string       `json:"RequestId" xml:"RequestId"`
	Images    []ImagesItem `json:"Images" xml:"Images"`
}

// CreateDescribeAICImagesRequest creates a request to invoke DescribeAICImages API
func CreateDescribeAICImagesRequest() (request *DescribeAICImagesRequest) {
	request = &DescribeAICImagesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ens", "2017-11-10", "DescribeAICImages", "ens", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAICImagesResponse creates a response to parse from DescribeAICImages response
func CreateDescribeAICImagesResponse() (response *DescribeAICImagesResponse) {
	response = &DescribeAICImagesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
