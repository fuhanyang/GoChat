package zap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var sugarLogger *zap.SugaredLogger

//	func InitLogger(path string) {
//		encoder := getEncoder()
//		// test.log记录全量日志
//		c1 := zapcore.NewCore(encoder, getLogWriter(fmt.Sprintf("%s/log/test.log", path)), zapcore.DebugLevel)
//		// test.err.log记录ERROR级别的日志
//		c2 := zapcore.NewCore(encoder, getLogWriter(fmt.Sprintf("%s/log/test.err.log", path)), zap.ErrorLevel)
//		// 使用NewTee将c1和c2合并到core
//		core := zapcore.NewTee(c1, c2)
//
//		logger := zap.New(core)
//		sugarLogger = logger.Sugar()
//		fmt.Println("logger init success at ", path)
//	}
func InitLogger(path string) {
	// 创建一个配置来启用调用者信息
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder // 启用短格式调用者信息 (file:line)

	// 使用配置创建编码器
	encoderWithCaller := getEncoder(config)

	// test.log记录全量日志
	c1 := zapcore.NewCore(encoderWithCaller, getLogWriter(fmt.Sprintf("%s/log/test.log", path)), zapcore.DebugLevel)
	// test.err.log记录ERROR级别的日志
	c2 := zapcore.NewCore(encoderWithCaller, getLogWriter(fmt.Sprintf("%s/log/test.err.log", path)), zap.ErrorLevel)
	// 使用NewTee将c1和c2合并到core
	core := zapcore.NewTee(c1, c2)

	// 创建带有调用者信息的logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	sugarLogger = logger.Sugar()
	fmt.Println("logger init success at ", path)
}
func GetSugarLogger() *zap.SugaredLogger {
	return sugarLogger
}
func getEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewJSONEncoder(config)
}

func getLogWriter(path string) zapcore.WriteSyncer {
	file, _ := os.Create(path)
	return zapcore.AddSync(file)
}
