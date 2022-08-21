package system

type DynamicInformation struct {
	CPUAvg                  float64
	MemAll                  uint64
	MenFree                 uint64
	MenUsed                 uint64
	MemUsedPercent          float64
	TotalDownStreamDataSize uint64
	TotalUpStreamDataSize   uint64
	NowDownStreamDataSize   int
	NowUpStreamDataSize     int
	DiskInformation         DiskInformation
	CU                      PingStatus
	CT                      PingStatus
	CM                      PingStatus
}

type DiskInformation struct {
	Total   uint64
	Used    uint64
	Percent uint64
}

type PingStatus struct {
	AvgRTT         int64
	PacketsReceive int
	PacketsSent    int
}

var StaticDynamicInformation DynamicInformation
