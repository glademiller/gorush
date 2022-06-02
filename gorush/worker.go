package gorush

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// InitWorkers for initialize all workers.
func InitWorkers(ctx context.Context, wg *sync.WaitGroup, workerNum int64, queueNum int64) {
	LogAccess.Info("worker number is ", workerNum, ", queue number is ", queueNum)
	QueueNotification = make(chan PushNotification, queueNum)
	for i := int64(0); i < workerNum; i++ {
		go startWorker(ctx, wg, i)
	}
}

// SendNotification is send message to iOS or Android
func SendNotification(ctx context.Context, req PushNotification) {
	if PushConf.Core.Sync {
		defer req.WaitDone()
	}

	switch req.Platform {
	case PlatformIos:
		PushToIOS(req)
	case PlatformAndroid:
		PushToAndroid(req)
	case PlatformWeb:
		PushToWeb(req)
	}
}

func startWorker(ctx context.Context, wg *sync.WaitGroup, num int64) {
	defer wg.Done()
	for notification := range QueueNotification {
		SendNotification(ctx, notification)
	}
	LogAccess.Info("closed the worker num ", num)
}

// markFailedNotification adds failure logs for all tokens in push notification
func markFailedNotification(notification *PushNotification, reason string) {
	LogError.Error(reason)
	for _, token := range notification.Tokens {
		notification.AddLog(getLogPushEntry(FailedPush, token, *notification, errors.New(reason)))
	}
	notification.WaitDone()
}

func waitAndPerformCallback(callbackUrl string, count int, wg *sync.WaitGroup, log *[]LogPushEntry) {
	if wg != nil {
		wg.Wait()
		reqBody, err := json.Marshal(gin.H{
			"success": "ok",
			"counts":  count,
			"logs":    log,
		})
		if err != nil {
			var msg = "Error converting logs to JSON."
			LogAccess.Debug(msg)
			return
		}
		resp, err2 := http.Post(callbackUrl, "application/json", bytes.NewBuffer(reqBody))
		if err2 != nil {
			var msg = "Error posting logs to callback URL."
			LogAccess.Debug(msg)
			return
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
	}
}

// queueNotification add notification to queue list.
func queueNotification(ctx context.Context, req RequestPush) (int, []LogPushEntry) {
	var count int
	var doSync = PushConf.Core.Sync
	if req.Sync != nil {
		doSync = *req.Sync
	}
	var callbackUrl = PushConf.Core.CallbackUrl
	if req.CallbackUrl != nil {
		callbackUrl = *req.CallbackUrl
	}
	wg := sync.WaitGroup{}
	newNotification := []*PushNotification{}
	for i := range req.Notifications {
		notification := &req.Notifications[i]
		switch notification.Platform {
		case PlatformIos:
			if !PushConf.Ios.Enabled && !notification.Voip {
				continue
			}
			if !PushConf.Ios.VoipEnabled && notification.Voip {
				continue
			}
		case PlatformAndroid:
			if !PushConf.Android.Enabled {
				continue
			}
		case PlatformWeb:
			if !PushConf.Web.Enabled {
				continue
			}
		}
		notification.sync = doSync || callbackUrl != ""
		newNotification = append(newNotification, notification)
	}

	log := make([]LogPushEntry, 0, count)
	for _, notification := range newNotification {
		if doSync || callbackUrl != "" {
			notification.wg = &wg
			notification.log = &log
			notification.AddWaitCount()
		}
		if !tryEnqueue(*notification, QueueNotification) {
			markFailedNotification(notification, "max capacity reached")
		}

		switch notification.Platform {
		case PlatformWeb:
			count += len(notification.Subscriptions)
		default:
			count += len(notification.Tokens)
		}

		// Count topic message
		if notification.To != "" {
			count++
		}
	}

	if doSync {
		wg.Wait()
	} else if callbackUrl != "" {
		go waitAndPerformCallback(callbackUrl, count, &wg, &log)
	}

	StatStorage.AddTotalCount(int64(count))

	return count, log
}

// tryEnqueue tries to enqueue a job to the given job channel. Returns true if
// the operation was successful, and false if enqueuing would not have been
// possible without blocking. Job is not enqueued in the latter case.
func tryEnqueue(job PushNotification, jobChan chan<- PushNotification) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}
