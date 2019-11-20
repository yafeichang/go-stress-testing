/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"bytes"
	"crypto/tls"
	"go-stress-testing/heper"
	"go-stress-testing/model"
	"go-stress-testing/server/client"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// http go link
func Http(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   request.Timeout,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {
		sleepDur := time.Duration(r.Intn(500)) * time.Millisecond
		if r.Intn(2) == 0 {
			sleepDur = time.Second - sleepDur
		} else {
			sleepDur = time.Second + sleepDur
		}
		time.Sleep(sleepDur)
		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)

		bodyReader := bytes.NewReader(request.Body)

		resp, err := client.HttpRequest(c, request.Method, request.Url, bodyReader, request.Headers, request.Timeout)
		requestTime := uint64(heper.DiffNano(startTime))
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
			resp.Body.Close()
		}

		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}
