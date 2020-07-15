package app

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	Addr              string `envconfig:"ADDR" default:"0.0.0.0:8090"`
	RPCAddr           string `envconfig:"RPC_ADDR" default:"0.0.0.0:5090"`
	UsersRPCAddr      string `envconfig:"USERS_RPC_ADDR" default:"0.0.0.0:5000"`
	StreamsRPCAddr    string `envconfig:"STREAMS_RPC_ADDR" default:"0.0.0.0:5102"`
	DispatcherRPCAddr string `envconfig:"DISPATCHER_RPC_ADDR" default:"0.0.0.0:5003"`
	RedisURI          string `envconfig:"REDISURI" default:"redis://:@127.0.0.1:6379/1"`
	MQURI             string `envconfig:"MQURI" default:"amqp://guest:guest@127.0.0.1:5672"`
	DataDir           string `envconfig:"DATA_DIR" default:"/tmp/data"`
	GDriveKey         string `envconfig:"GDRIVE_KEY" required:"true"`
	AuthTokenSecret   string `envconfig:"AUTH_TOKEN_SECRET" default:"secret"`
	Bucket            string `envconfig:"BUCKET" required:"true"`
	WebRTCServerHost  string `envconfig:"WEBRTC_SERVER_HOST" default:"0.0.0.0"`
	WebRTCServerPort  int    `envconfig:"WEBRTC_SERVER_PORT" default:"6000"`
	RTMPURL           string `envconfig:"RTMP_URL" default:"rtmp://127.0.0.1:1935/live"`
}
