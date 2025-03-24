package aliyunsls

import "github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"

func StartOpenTelemetryProvider(serviceName string, serviceVersion string, project string, instanceID string, traceEndpoint string, metricEndpoint string, accessKeyID string, accessKeySecret string) (func(), error) {
	slsConfig, err := provider.NewConfig(provider.WithServiceName(serviceName),
		provider.WithServiceVersion(serviceVersion),
		provider.WithTraceExporterEndpoint(traceEndpoint),
		provider.WithMetricExporterEndpoint(metricEndpoint),
		provider.WithSLSConfig(project, instanceID, accessKeyID, accessKeySecret))
	// 如果初始化失败则panic，可以替换为其他错误处理方式
	if err != nil {
		return nil, err
	}
	if err := provider.Start(slsConfig); err != nil {
		return nil, err
	}

	return func() { provider.Shutdown(slsConfig) }, nil
}
