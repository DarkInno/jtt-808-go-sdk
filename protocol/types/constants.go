package types

// JTT 808 协议版本
const (
	// Version2011 2011版本
	Version2011 uint8 = 0x00
	// Version2013 2013版本
	Version2013 uint8 = 0x01
	// Version2019 2019版本
	Version2019 uint8 = 0x02
)

// 上行消息ID（终端到平台）
const (
	// MsgIDTerminalCommonResponse 终端通用应答
	MsgIDTerminalCommonResponse uint16 = 0x0001
	// MsgIDTerminalHeartbeat 终端心跳
	MsgIDTerminalHeartbeat uint16 = 0x0002
	// MsgIDTerminalRegister 终端注册
	MsgIDTerminalRegister uint16 = 0x0100
	// MsgIDTerminalAuth 终端鉴权
	MsgIDTerminalAuth uint16 = 0x0102
	// MsgIDTerminalDeregister 终端注销
	MsgIDTerminalDeregister uint16 = 0x0003
	// MsgIDLocationReport 位置信息汇报
	MsgIDLocationReport uint16 = 0x0200
	// MsgIDLocationQueryResponse 位置信息查询应答
	MsgIDLocationQueryResponse uint16 = 0x0201
	// MsgIDEventReport 事件报告
	MsgIDEventReport uint16 = 0x0300
	// MsgIDQuestionResponse 问题答案
	MsgIDQuestionResponse uint16 = 0x0301
	// MsgIDInfoDemandMenu 信息点播/取消
	MsgIDInfoDemandMenu uint16 = 0x0302
	// MsgIDVehicleControlResponse 车辆控制应答
	MsgIDVehicleControlResponse uint16 = 0x0500
	// MsgIDDrivingRecordData 行驶记录数据上传
	MsgIDDrivingRecordData uint16 = 0x0701
	// MsgIDDriverIdentityReport 驾驶员身份信息采集上报
	MsgIDDriverIdentityReport uint16 = 0x0702
	// MsgIDLocationBatchReport 定位数据批量上传
	MsgIDLocationBatchReport uint16 = 0x0704
	// MsgIDCANDataReport CAN总线数据上报
	MsgIDCANDataReport uint16 = 0x0705
	// MsgIDMediaEventReport 多媒体事件信息上传
	MsgIDMediaEventReport uint16 = 0x0800
	// MsgIDMediaDataUpload 多媒体数据上传
	MsgIDMediaDataUpload uint16 = 0x0801
	// MsgIDMediaDataRetransmit 多媒体数据上传
	MsgIDMediaDataRetransmit uint16 = 0x0802
	// MsgIDDataUpTransparent 数据上行透传
	MsgIDDataUpTransparent uint16 = 0x0900
	// MsgIDDataCompressUpload 数据压缩上报
	MsgIDDataCompressUpload uint16 = 0x0901
	// MsgIDPlatformRSA 平台RSA公钥
	MsgIDPlatformRSA uint16 = 0x0A00
)

// 下行消息ID（平台到终端）
const (
	// MsgIDPlatformCommonResponse 平台通用应答
	MsgIDPlatformCommonResponse uint16 = 0x8001
	// MsgIDPlatformHeartbeat 平台心跳
	MsgIDPlatformHeartbeat uint16 = 0x8002
	// MsgIDServer补传分包请求 服务器补传分包请求
	MsgIDServerRetransmitRequest uint16 = 0x8003
	// MsgIDTerminalRegisterResponse 终端注册应答
	MsgIDTerminalRegisterResponse uint16 = 0x8100
	// MsgIDSetTerminalParameters 设置终端参数
	MsgIDSetTerminalParameters uint16 = 0x8103
	// MsgIDQueryTerminalParameters 查询终端参数
	MsgIDQueryTerminalParameters uint16 = 0x8104
	// MsgIDQueryTerminalParametersResponse 查询终端参数应答
	MsgIDQueryTerminalParametersResponse uint16 = 0x8104
	// MsgIDTerminalControl 终端控制
	MsgIDTerminalControl uint16 = 0x8105
	// MsgIDQuerySpecificTerminalParameters 查询指定终端参数
	MsgIDQuerySpecificTerminalParameters uint16 = 0x8106
	// MsgIDQueryTerminalAttribute 查询终端属性
	MsgIDQueryTerminalAttribute uint16 = 0x8107
	// MsgIDLocationQuery 下发终端升级包
	MsgIDLocationQuery uint16 = 0x8200
	// MsgIDLocationQueryResponse 位置信息查询
	MsgIDLocationQueryResponse2 uint16 = 0x8201
	// MsgIDTemporaryLocationTracking 临时位置跟踪控制
	MsgIDTemporaryLocationTracking uint16 = 0x8202
	// MsgIDVehicleControl 车辆控制
	MsgIDVehicleControl uint16 = 0x8500
	// MsgIDSetRoundArea 设置圆形区域
	MsgIDSetRoundArea uint16 = 0x8600
	// MsgIDDeleteRoundArea 删除圆形区域
	MsgIDDeleteRoundArea uint16 = 0x8601
	// MsgIDSetRectangularArea 设置矩形区域
	MsgIDSetRectangularArea uint16 = 0x8602
	// MsgIDDeleteRectangularArea 删除矩形区域
	MsgIDDeleteRectangularArea uint16 = 0x8603
	// MsgIDSetPolygonArea 设置多边形区域
	MsgIDSetPolygonArea uint16 = 0x8604
	// MsgIDDeletePolygonArea 删除多边形区域
	MsgIDDeletePolygonArea uint16 = 0x8605
	// MsgIDSetRoute 设置路线
	MsgIDSetRoute uint16 = 0x8606
	// MsgIDDeleteRoute 删除路线
	MsgIDDeleteRoute uint16 = 0x8607
	// MsgIDQueryAreaOrRoute 查询区域或路线数据
	MsgIDQueryAreaOrRoute uint16 = 0x8608
	// MsgIDAreaOrRouteResponse 区域或路线数据应答
	MsgIDAreaOrRouteResponse uint16 = 0x8609
	// MsgIDDrivingRecordParamCommand 行驶记录参数采集命令
	MsgIDDrivingRecordParamCommand uint16 = 0x8700
	// MsgIDDrivingRecordCommandResponse 行驶记录参数命令应答
	MsgIDDrivingRecordCommandResponse uint16 = 0x8700
	// MsgIDReportDrivingRecord 上报行驶记录
	MsgIDReportDrivingRecord uint16 = 0x8701
	// MsgIDReportDrivingRecordResponse 上报行驶记录应答
	MsgIDReportDrivingRecordResponse uint16 = 0x8701
	// MsgIDDriverIdentityRequest 驾驶员身份信息请求
	MsgIDDriverIdentityRequest uint16 = 0x8702
	// MsgIDDriverIdentityResponse 驾驶员身份信息应答
	MsgIDDriverIdentityResponse uint16 = 0x8702
	// MsgIDMediaDataRetransmitRequest 多媒体数据检索应答
	MsgIDMediaDataRetransmitRequest uint16 = 0x8800
	// MsgIDCameraShotRequest 摄像头立即拍摄命令
	MsgIDCameraShotRequest uint16 = 0x8801
	// MsgIDCameraShotResponse 摄像头立即拍摄应答
	MsgIDCameraShotResponse uint16 = 0x8801
	// MsgIDMediaDataRetransmitResponse 存储多媒体数据检索
	MsgIDMediaDataRetransmitResponse uint16 = 0x8802
	// MsgIDDataDownTransparent 数据下行透传
	MsgIDDataDownTransparent uint16 = 0x8900
	// MsgIDPlatformRSARequest 平台RSA公钥
	MsgIDPlatformRSARequest uint16 = 0x8A00
)

// 消息体属性位定义
const (
	// MsgBodyPropertyLengthMask 消息体长度掩码（0-9位）
	MsgBodyPropertyLengthMask uint16 = 0x03FF
	// MsgBodyPropertyEncryptionMask 加密方式掩码（10-12位）
	MsgBodyPropertyEncryptionMask uint16 = 0x1C00
	// MsgBodyPropertySubPackageFlag 分包标志（13位）
	MsgBodyPropertySubPackageFlag uint16 = 0x2000
	// MsgBodyPropertyVersionFlag 协议版本标志（14位，2019版本）
	MsgBodyPropertyVersionFlag uint16 = 0x4000
	// MsgBodyPropertyReserved 保留位（15位）
	MsgBodyPropertyReserved uint16 = 0x8000
)

// 加密方式
const (
	// EncryptionNone 不加密
	EncryptionNone uint8 = 0x00
	// EncryptionRSA RSA加密
	EncryptionRSA uint8 = 0x01
)

// 终端注册结果
const (
	// RegisterResultSuccess 成功
	RegisterResultSuccess uint8 = 0x00
	// RegisterResultVehicleRegistered 车辆已被注册
	RegisterResultVehicleRegistered uint8 = 0x01
	// RegisterResultNoVehicle 找不到数据库中对应的车辆
	RegisterResultNoVehicle uint8 = 0x02
	// RegisterResultTerminalRegistered 终端已被注册
	RegisterResultTerminalRegistered uint8 = 0x03
	// RegisterResultNoTerminal 找不到数据库中对应的终端
	RegisterResultNoTerminal uint8 = 0x04
)

// 通用应答结果
const (
	// CommonResponseSuccess 成功
	CommonResponseSuccess uint8 = 0x00
	// CommonResponseFailure 失败
	CommonResponseFailure uint8 = 0x01
	// CommonResponseMessageNotSupported 消息不支持
	CommonResponseMessageNotSupported uint8 = 0x02
	// CommonResponseAlarmConfirm 报警处理确认
	CommonResponseAlarmConfirm uint8 = 0x03
)

// 报警标志位
const (
	// AlarmFlagSOS 紧急报警
	AlarmFlagSOS uint32 = 0x00000001
	// AlarmFlagOverSpeed 超速报警
	AlarmFlagOverSpeed uint32 = 0x00000002
	// AlarmFlagFatigue 疲劳驾驶
	AlarmFlagFatigue uint32 = 0x00000004
	// AlarmFlagDanger 危险预警
	AlarmFlagDanger uint32 = 0x00000008
	// AlarmFlagGNSSFault GNSS模块故障
	AlarmFlagGNSSFault uint32 = 0x00000010
	// AlarmFlagGNSSAntenna GNSS天线未接或被剪断
	AlarmFlagGNSSAntenna uint32 = 0x00000020
	// AlarmFlagGNSSShortCircuit GNSS天线短路
	AlarmFlagGNSSShortCircuit uint32 = 0x00000040
	// AlarmFlagPowerLow 终端主电源欠压
	AlarmFlagPowerLow uint32 = 0x00000080
	// AlarmFlagPowerOff 终端主电源掉电
	AlarmFlagPowerOff uint32 = 0x00000100
	// AlarmFlagLCDFault 终端LCD或显示器故障
	AlarmFlagLCDFault uint32 = 0x00000200
	// AlarmFlagTTSSpeechFault TTS模块故障
	AlarmFlagTTSSpeechFault uint32 = 0x00000400
	// AlarmFlagCameraFault 摄像头故障
	AlarmFlagCameraFault uint32 = 0x00000800
	// AlarmFlagICCardFault 道路运输证IC卡模块故障
	AlarmFlagICCardFault uint32 = 0x00001000
	// AlarmFlagOverSpeedWarning 超速预警
	AlarmFlagOverSpeedWarning uint32 = 0x00002000
	// AlarmFlagFatigueWarning 疲劳驾驶预警
	AlarmFlagFatigueWarning uint32 = 0x00004000
	// AlarmFlagReserved1 保留
	AlarmFlagReserved1 uint32 = 0x00008000
	// AlarmFlagDriveTimeout 当天累计驾驶超时
	AlarmFlagDriveTimeout uint32 = 0x00010000
	// AlarmFlagStopTimeout 停车超时
	AlarmFlagStopTimeout uint32 = 0x00020000
	// AlarmFlagAreaIn 进出区域
	AlarmFlagAreaIn uint32 = 0x00040000
	// AlarmFlagRouteDeviation 进出路线
	AlarmFlagRouteDeviation uint32 = 0x00080000
	// AlarmFlagRoadTimeShort 路段行驶时间不足/过长
	AlarmFlagRoadTimeShort uint32 = 0x00100000
	// AlarmFlagRouteDeviation2 路线偏离报警
	AlarmFlagRouteDeviation2 uint32 = 0x00200000
	// AlarmFlagVSSFault 车辆VSS故障
	AlarmFlagVSSFault uint32 = 0x00400000
	// AlarmFlagOilMass abnormal 车辆油量异常
	AlarmFlagOilMassAbnormal uint32 = 0x00800000
	// AlarmFlagVehicleStolen 车辆被盗
	AlarmFlagVehicleStolen uint32 = 0x01000000
	// AlarmFlagVehicleIllegalMove 车辆非法点火
	AlarmFlagVehicleIllegalMove uint32 = 0x02000000
	// AlarmFlagVehicleIllegalMove2 车辆非法位移
	AlarmFlagVehicleIllegalMove2 uint32 = 0x04000000
	// AlarmFlagCollision 碰撞预警
	AlarmFlagCollision uint32 = 0x08000000
	// AlarmFlagRollover 侧翻预警
	AlarmFlagRollover uint32 = 0x10000000
	// AlarmFlagIllegalOpenDoor 非法开门报警
	AlarmFlagIllegalOpenDoor uint32 = 0x20000000
	// AlarmFlagReserved2 保留
	AlarmFlagReserved2 uint32 = 0x40000000
	// AlarmFlagReserved3 保留
	AlarmFlagReserved3 uint32 = 0x80000000
)

// 状态位定义
const (
	// StatusACC ACC状态
	StatusACC uint32 = 0x00000001
	// Status定位状态
	StatusPositioning uint32 = 0x00000002
	// Status南纬
	StatusSouthLatitude uint32 = 0x00000004
	// Status西经
	StatusWestLongitude uint32 = 0x00000008
	// Status停运状态
	StatusOutOfService uint32 = 0x00000010
	// Status经纬度保密
	StatusEncrypted uint32 = 0x00000020
	// Status1 北斗卫星定位
	StatusBeiDou uint32 = 0x00000040
	// Status1 GPS卫星定位
	StatusGPS uint32 = 0x00000080
	// Status1 GLONASS卫星定位
	StatusGLONASS uint32 = 0x00000100
	// Status1 Galileo卫星定位
	StatusGalileo uint32 = 0x00000200
)

// 协议标志位
const (
	// ProtocolFlag 标志位
	ProtocolFlag byte = 0x7E
	// ProtocolEscape 转义字符
	ProtocolEscape byte = 0x7D
)

// 转义映射
var (
	// EscapeMap 转义映射
	EscapeMap = map[byte]byte{
		0x7E: 0x02,
		0x7D: 0x01,
	}
	// UnescapeMap 反转义映射
	UnescapeMap = map[byte]byte{
		0x02: 0x7E,
		0x01: 0x7D,
	}
)
