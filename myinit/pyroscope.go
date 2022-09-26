package myinit

import "github.com/pyroscope-io/pyroscope/pkg/agent/profiler"

func InitPyroscope(appName string, addr string) (*profiler.Profiler, error) {
	return profiler.Start(
		profiler.Config{
			ApplicationName: appName,

			// replace this with the address of pyroscope server
			ServerAddress: addr,

			// by default all profilers are enabled,
			// but you can select the ones you want to use:
			ProfileTypes: []profiler.ProfileType{
				profiler.ProfileCPU,
				profiler.ProfileAllocObjects,
				profiler.ProfileAllocSpace,
				profiler.ProfileInuseObjects,
				profiler.ProfileInuseSpace,
			},
		},
	)

}
