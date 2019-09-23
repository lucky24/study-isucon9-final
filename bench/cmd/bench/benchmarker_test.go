package main

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/chibiegg/isucon9-final/bench/internal/bencherror"
	"github.com/chibiegg/isucon9-final/bench/internal/config"
	"github.com/chibiegg/isucon9-final/bench/internal/endpoint"
	"github.com/chibiegg/isucon9-final/bench/internal/logger"
	"github.com/chibiegg/isucon9-final/bench/isutrain"
	"github.com/chibiegg/isucon9-final/bench/mock"
	"github.com/chibiegg/isucon9-final/bench/payment"
	"github.com/chibiegg/isucon9-final/bench/scenario"
	"github.com/jarcoal/httpmock"
)

// NOTE: go test -login-delay=150 -bench .
var (
	testBenchLoginDelayMsec             = flag.Int("login-delay", 100, "login-delay [msec]")
	testBenchReserveDelayMsec           = flag.Int("reserve-delay", 100, "reserve-delay [msec]")
	testBenchListStationsDelayMsec      = flag.Int("liststations-delay", 100, "liststations-delay [msec]")
	testBenchSearchTrainsDelayMsec      = flag.Int("searchtrains-delay", 100, "searchtrains-delay [msec]")
	testBenchCommitReservationDelayMsec = flag.Int("commitreservation-delay", 100, "commitreservation-delay [msec]")
	testBenchCancelReservationDelayMsec = flag.Int("cancelreservation-delay", 100, "cancelreservation-delay [msec]")
	testBenchListReservationDelayMsec   = flag.Int("listreservation-delay", 100, "listreservation-delay [msec]")
	testBenchListTrainSeatsDelayMsec    = flag.Int("listtrainseats-delay", 100, "listtrainseats-delay [msec]")
)

var testBenchTimeoutSec = flag.Int("timeout-sec", 5, "timeout [sec]")

func setDelay(m *mock.Mock) {
	m.LoginDelay = time.Duration(*testBenchLoginDelayMsec) * time.Millisecond
	m.ReserveDelay = time.Duration(*testBenchReserveDelayMsec) * time.Millisecond
	m.ListStationsDelay = time.Duration(*testBenchListStationsDelayMsec) * time.Millisecond
	m.SearchTrainsDelay = time.Duration(*testBenchSearchTrainsDelayMsec) * time.Millisecond
	m.CommitReservationDelay = time.Duration(*testBenchCommitReservationDelayMsec) * time.Millisecond
	m.CancelReservationDelay = time.Duration(*testBenchCancelReservationDelayMsec) * time.Millisecond
	m.ListReservationDelay = time.Duration(*testBenchListReservationDelayMsec) * time.Millisecond
	m.ListTrainSeatsDelay = time.Duration(*testBenchListTrainSeatsDelayMsec) * time.Millisecond
}

func BenchmarkScore(b *testing.B) {
	flag.Parse()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Debug = true

	logger.InitZapLogger()

	benchmarker := new(benchmarker)
	isutrainClient, _ := isutrain.NewClient()
	paymentClient, _ := payment.NewClient()

	m, _ := mock.Register()
	setDelay(m)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*testBenchTimeoutSec)*time.Second)
	defer cancel()

	benchmarker.run(ctx)

	scenario.FinalCheck(ctx, isutrainClient, paymentClient)

	score := endpoint.CalcFinalScore()
	b.ReportMetric(float64(score), "score")
	b.ReportMetric(float64(bencherror.BenchmarkErrs.Penalty()), "penalty")
}
