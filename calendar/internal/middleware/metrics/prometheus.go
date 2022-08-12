package metrics

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	requestsPerUser = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_requests_total",
		Help: "Number of requests per user",
	}, []string{"user"})
	_ = promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "app_users_count",
		Help: "Number of users",
	}, usersCountDelegate)
	_ = promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "app_events_count",
		Help: "Number of events",
	}, eventsCountDelegate)
)

var usersCountFunc = func() float64 {
	return -1
}

func usersCountDelegate() float64 {
	return usersCountFunc()
}

var eventsCountFunc = func() float64 {
	return -1
}

func eventsCountDelegate() float64 {
	return eventsCountFunc()
}

func IncRequest(user string) {
	requestsPerUser.With(prometheus.Labels{"user": user}).Inc()
}

func Init(usersCount func() (int, error), eventsCount func() (int, error)) {
	usersCountFunc = func() float64 {
		count, err := usersCount()
		if err != nil {
			logging.Logger.Error("failed to get users count", zap.Error(err))
			return -1
		}
		return float64(count)
	}
	eventsCountFunc = func() float64 {
		count, err := eventsCount()
		if err != nil {
			logging.Logger.Error("failed to get events count", zap.Error(err))
			return -1
		}
		return float64(count)
	}
	prometheus.MustRegister(requestsPerUser)
}
