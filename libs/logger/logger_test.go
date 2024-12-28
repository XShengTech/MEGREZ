package logger_test

import (
	"testing"

	"megrez/libs/logger"
)

func TestLogger(t *testing.T) {
	t.Log("TestLogger")
	l, err := logger.NewLogger(logger.DEBUG, "stdout")
	if err != nil {
		t.Fatal(err)
	}
	l.Info("TestInfo")
	l.Warn("TestWarn")
	l.Error("TestError")
	l.Debug("TestDebug")
	// l.Fatal("Fatal")

	t.Log(("SetModel"))
	l.SetModel("Logger")
	l.Info("SetModel TestInfo")
	l.Warn("SetModel TestWarn")
	l.Error("SetModel TestError")
	l.Debug("SetModel TestDebug")

	t.Log(("SetFunction"))
	l.SetFunction("TestLogger")
	l.Info("SetFunction TestInfo")
	l.Warn("SetFunction TestWarn")
	l.Error("SetFunction TestError")
	l.Debug("SetFunction TestDebug")

	l.Close()
}

func BenchmarkLogger(b *testing.B) {
	b.Log("BenchmarkLogger")
	l, err := logger.NewLogger(logger.DEBUG, "stdout")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		go l.Info("BenchmarkLogger: %d", i)
	}
	l.Close()
}
