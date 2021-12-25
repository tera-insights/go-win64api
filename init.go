package winapi

func init() {
	modAdvapi32.Load()
	modKernel32.Load()
	modNetapi32.Load()
	modOffreg.Load()
	modSecur32.Load()
	modUserenv.Load()
	modWtsapi32.Load()
}
