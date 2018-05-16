package asclient

type AerospikeHost struct {
	Host string
	Port int
}

type AerospikeConfig struct {
	Hosts []AerospikeHost
}