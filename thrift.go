package utils

//func GetClient(host, port string) (*msgbus.BaikeMsgbusClient, error) {
//
//	trans, err := thrift.NewTSocket(net.JoinHostPort(host, port))
//
//	if err != nil {
//		return nil, err
//	}
//
//	tf := thrift.NewTFramedTransport(trans)
//
//	if err := tf.Open(); err != nil {
//		panic(err)
//	}
//
//	var protocolFactory thrift.TProtocolFactory
//
//	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
//
//	iprot := protocolFactory.GetProtocol(tf)
//	oprot := protocolFactory.GetProtocol(tf)
//	client := msgbus.NewBaikeMsgbusClient(thrift.NewTStandardClient(iprot, oprot))
//
//	return client, nil
//}
