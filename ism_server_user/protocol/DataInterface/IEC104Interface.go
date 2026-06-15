/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package DataInterface

import (
	"ISMServer/models"
	ISMScript "ISMServer/task/ISMScript/func"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/wendy512/go-iecp5/asdu"
	"github.com/wendy512/iec104/server"
)

var GIEC104InterfaceChan chan bool
var IEC104Wg sync.WaitGroup

type IEC104Ctl struct {
	waitGroup     *sync.WaitGroup
	InterfaceData IEC104InterfaceData
	Conn          asdu.Connect
}
type DataPointStu struct {
	BandData  string
	Type      int
	Point     int
	DataValue float32
	TempValue any
}
type IEC104InterfaceData struct {
	InterfaceName        string
	ProjectUuid          string
	InterfaceUuid        string
	InterfaceType        int
	InterfaceDataUuid    string
	InterfacePort        int
	InterfaceAddr        int
	InterfaceContent     string
	InterfaceDataContent []string
	InterfaceDataPoint   int
	YaoXin               []DataPointStu
	YaoCe                []DataPointStu
	Energy               []DataPointStu
	YaoKong              []DataPointStu
	YaoTiao              []DataPointStu
}

func isIEC104ChanClose() bool {
	select {
	case _, received := <-GIEC104InterfaceChan:
		return !received
	default:
	}
	return false
}
func IEC104InterfaceCloseChan() {
	isOpen := isIEC104ChanClose()
	if !isOpen && GIEC104InterfaceChan != nil {
		close(GIEC104InterfaceChan)
	}
}

func (c *IEC104Ctl) StartIEC104Server() {
	NewSettings := server.NewSettings()
	NewSettings.LogCfg.Enable = false
	NewSettings.Host = "0.0.0.0"
	NewSettings.Port = c.InterfaceData.InterfacePort
	c.Conn = nil
	s := server.New(NewSettings, c)
	s.SetOnConnectionHandler(func(cs asdu.Connect) {})
	s.SetConnectionLostHandler(func(cs asdu.Connect) {
		c.Conn = nil
	})
	s.Start()

	//检测协程是否主动退出
	<-GIEC104InterfaceChan
	c.waitGroup.Done()
	s.Stop()
	logs.Info("重新加载IEC104 Server接口")
}
func (c *IEC104Ctl) IEC104DataMutation() {
	for {
		//检测协程是否主动退出
		select {
		case <-GIEC104InterfaceChan:
			c.waitGroup.Done()
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
		}

		for k, v := range c.InterfaceData.YaoXin {
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}

			if tempValue1, ok1 := value.(int); ok1 {
				if tempValue1 == 1 {
					c.InterfaceData.YaoXin[k].DataValue = 1
				} else {
					c.InterfaceData.YaoXin[k].DataValue = 0
				}
			} else if tempValue2, ok2 := value.(int64); ok2 {
				if tempValue2 == 1 {
					c.InterfaceData.YaoXin[k].DataValue = 1
				} else {
					c.InterfaceData.YaoXin[k].DataValue = 0
				}
			} else if tempValue2, ok2 := value.(float32); ok2 {
				if tempValue2 == 1 {
					c.InterfaceData.YaoXin[k].DataValue = 1
				} else {
					c.InterfaceData.YaoXin[k].DataValue = 0
				}
			} else if tempValue2, ok2 := value.(float64); ok2 {
				if tempValue2 == 1 {
					c.InterfaceData.YaoXin[k].DataValue = 1
				} else {
					c.InterfaceData.YaoXin[k].DataValue = 0
				}
			} else {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}
			if c.InterfaceData.YaoXin[k].TempValue != c.InterfaceData.YaoXin[k].DataValue {
				if c.Conn != nil {
					if v.Type == 1 {
						var SpontaneousValue bool
						if c.InterfaceData.YaoXin[k].DataValue == 1 {
							SpontaneousValue = true
						} else {
							SpontaneousValue = false
						}
						err := asdu.Single(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.SinglePointInfo{
							Ioa:   asdu.InfoObjAddr(v.Point),
							Value: SpontaneousValue,
							Qds:   asdu.QDSGood,
						})
						if err != nil {
							logs.Error(err)
						}
					} else {
						err := asdu.Double(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.DoublePointInfo{
							Ioa:   asdu.InfoObjAddr(v.Point),
							Value: asdu.DoublePoint(c.InterfaceData.YaoXin[k].DataValue),
							Qds:   asdu.QDSGood,
						})
						if err != nil {
							logs.Error(err)
						}
					}
				}
				c.InterfaceData.YaoXin[k].TempValue = c.InterfaceData.YaoXin[k].DataValue
			}
		}

		for k, v := range c.InterfaceData.YaoCe {
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				c.InterfaceData.YaoCe[k].DataValue = 0
			}

			tempValue, ok := value.(float64)
			if ok {
				c.InterfaceData.YaoCe[k].DataValue = float32(tempValue)
			} else if tempValue1, ok1 := value.(float32); ok1 {
				c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
			} else if tempValue1, ok1 := value.(int64); ok1 {
				c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
			} else if tempValue1, ok1 := value.(int); ok1 {
				c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
			} else {
				continue
			}
			if c.InterfaceData.YaoCe[k].TempValue != c.InterfaceData.YaoCe[k].DataValue {
				if c.Conn != nil {
					if v.Type == 3 {
						err := asdu.MeasuredValueFloat(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.MeasuredValueFloatInfo{
							Ioa:   asdu.InfoObjAddr(v.Point),
							Value: c.InterfaceData.YaoCe[k].DataValue,
							Qds:   asdu.QDSGood,
						})
						if err != nil {
							logs.Error(err)
						}
					} else if v.Type == 2 {
						err := asdu.MeasuredValueScaled(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.MeasuredValueScaledInfo{
							Ioa:   asdu.InfoObjAddr(v.Point),
							Value: int16(c.InterfaceData.YaoCe[k].DataValue),
							Qds:   asdu.QDSGood,
						})
						if err != nil {
							logs.Error(err)
						}
					} else if v.Type == 1 {
						err := asdu.MeasuredValueNormal(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.MeasuredValueNormalInfo{
							Ioa:   asdu.InfoObjAddr(v.Point),
							Value: asdu.Normalize(c.InterfaceData.YaoCe[k].DataValue),
							Qds:   asdu.QDSGood,
						})
						if err != nil {
							logs.Error(err)
						}
					}
				}
				c.InterfaceData.YaoCe[k].TempValue = c.InterfaceData.YaoCe[k].DataValue
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// 总招
func (c *IEC104Ctl) OnInterrogation(conn asdu.Connect, pack *asdu.ASDU, quality asdu.QualifierOfInterrogation) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	// TODO
	c.Conn = conn
	SinglePointInfo := make([]asdu.SinglePointInfo, 0, 1000)
	DoublePointInfo := make([]asdu.DoublePointInfo, 0, 1000)
	for k, v := range c.InterfaceData.YaoXin {
		value := ISMScript.GetDeviceData(v.BandData)
		if value == nil {
			c.InterfaceData.YaoXin[k].DataValue = 0
		}
		if tempValue1, ok1 := value.(int); ok1 {
			if tempValue1 == 1 {
				c.InterfaceData.YaoXin[k].DataValue = 1
			} else {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}
		} else if tempValue2, ok2 := value.(int64); ok2 {
			if tempValue2 == 1 {
				c.InterfaceData.YaoXin[k].DataValue = 1
			} else {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}
		} else if tempValue2, ok2 := value.(float32); ok2 {
			if tempValue2 == 1 {
				c.InterfaceData.YaoXin[k].DataValue = 1
			} else {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}
		} else if tempValue2, ok2 := value.(float64); ok2 {
			if tempValue2 == 1 {
				c.InterfaceData.YaoXin[k].DataValue = 1
			} else {
				c.InterfaceData.YaoXin[k].DataValue = 0
			}
		} else {
			c.InterfaceData.YaoXin[k].DataValue = 0
		}
		if v.Type == 1 {
			var TeV bool
			if c.InterfaceData.YaoXin[k].DataValue == 1 {
				TeV = true
			} else {
				TeV = false
			}
			SinglePointInfo = append(SinglePointInfo, asdu.SinglePointInfo{
				Ioa:   asdu.InfoObjAddr(v.Point),
				Value: TeV,
				Qds:   asdu.QDSGood,
			})
		} else {
			DoublePointInfo = append(DoublePointInfo, asdu.DoublePointInfo{
				Ioa:   asdu.InfoObjAddr(v.Point),
				Value: asdu.DoublePoint(c.InterfaceData.YaoXin[k].DataValue),
				Qds:   asdu.QDSGood,
			})
		}
	}

	if len(SinglePointInfo) > 0 {
		_ = asdu.Single(conn, false, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), SinglePointInfo...)
	}

	if len(DoublePointInfo) > 0 {
		_ = asdu.Double(conn, false, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), DoublePointInfo...)
	}

	MeasuredValueFloatInfo := make([]asdu.MeasuredValueFloatInfo, 0, 1000)
	MeasuredValueScaled := make([]asdu.MeasuredValueScaledInfo, 0, 1000)
	MeasuredValueNormalInfo := make([]asdu.MeasuredValueNormalInfo, 0, 1000)
	for k, v := range c.InterfaceData.YaoCe {
		value := ISMScript.GetDeviceData(v.BandData)
		if value == nil {
			c.InterfaceData.YaoCe[k].DataValue = 0
		}

		tempValue, ok := value.(float64)
		if ok {
			c.InterfaceData.YaoCe[k].DataValue = float32(tempValue)
		} else if tempValue1, ok1 := value.(float32); ok1 {
			c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
		} else if tempValue1, ok1 := value.(int64); ok1 {
			c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
		} else if tempValue1, ok1 := value.(int); ok1 {
			c.InterfaceData.YaoCe[k].DataValue = float32(tempValue1)
		} else {
			continue
		}

		if v.Type == 3 {
			MeasuredValueFloatInfo = append(MeasuredValueFloatInfo, asdu.MeasuredValueFloatInfo{
				Ioa:   asdu.InfoObjAddr(v.Point),
				Value: float32(c.InterfaceData.YaoCe[k].DataValue),
				Qds:   asdu.QDSGood,
			})

		} else if v.Type == 2 {
			MeasuredValueScaled = append(MeasuredValueScaled, asdu.MeasuredValueScaledInfo{
				Ioa:   asdu.InfoObjAddr(v.Point),
				Value: int16(c.InterfaceData.YaoCe[k].DataValue),
				Qds:   asdu.QDSGood,
			})
		} else if v.Type == 1 {
			MeasuredValueNormalInfo = append(MeasuredValueNormalInfo, asdu.MeasuredValueNormalInfo{
				Ioa:   asdu.InfoObjAddr(v.Point),
				Value: asdu.Normalize(c.InterfaceData.YaoCe[k].DataValue),
				Qds:   asdu.QDSGood,
			})

		}

	}
	if len(MeasuredValueFloatInfo) > 0 {
		err := asdu.MeasuredValueFloat(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), MeasuredValueFloatInfo...)
		if err != nil {
			logs.Error(err)
		}
	}
	if len(MeasuredValueScaled) > 0 {
		err := asdu.MeasuredValueScaled(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), MeasuredValueScaled...)
		if err != nil {
			logs.Error(err)
		}
	}
	if len(MeasuredValueNormalInfo) > 0 {
		err := asdu.MeasuredValueNormal(c.Conn, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), MeasuredValueNormalInfo...)
		if err != nil {
			logs.Error(err)
		}
	}
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)

	return nil
}

// 电量
func (c *IEC104Ctl) OnCounterInterrogation(conn asdu.Connect, pack *asdu.ASDU, quality asdu.QualifierCountCall) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	// TODO
	_ = asdu.CounterInterrogationCmd(conn, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.QualifierCountCall{asdu.QCCGroup1, asdu.QCCFrzRead})
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)
	return nil
}

func (c *IEC104Ctl) OnRead(conn asdu.Connect, pack *asdu.ASDU, addr asdu.InfoObjAddr) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	// TODO
	_ = asdu.Single(conn, false, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.SinglePointInfo{
		Ioa:   addr,
		Value: true,
		Qds:   asdu.QDSGood,
	})
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)
	return nil
}

func (c *IEC104Ctl) OnClockSync(conn asdu.Connect, pack *asdu.ASDU, tm time.Time) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	now := time.Now()
	_ = asdu.ClockSynchronizationCmd(conn, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), now)
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)
	return nil
}

func (c *IEC104Ctl) OnResetProcess(conn asdu.Connect, pack *asdu.ASDU, quality asdu.QualifierOfResetProcessCmd) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	// TODO
	_ = asdu.ResetProcessCmd(conn, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), asdu.QPRGeneralRest)
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)
	return nil
}

func (c *IEC104Ctl) OnDelayAcquisition(conn asdu.Connect, pack *asdu.ASDU, msec uint16) error {
	_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
	// TODO
	_ = asdu.DelayAcquireCommand(conn, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(c.InterfaceData.InterfaceAddr), msec)
	_ = pack.SendReplyMirror(conn, asdu.ActivationTerm)
	return nil
}

func (c *IEC104Ctl) OnASDU(conn asdu.Connect, pack *asdu.ASDU) error {

	// TODO
	cmd := pack.GetSingleCmd()
	for _, v := range c.InterfaceData.YaoKong {
		if asdu.InfoObjAddr(v.Point) == cmd.Ioa {
			_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
			_ = asdu.SingleCmd(conn, pack.Type, pack.Coa, pack.CommonAddr, asdu.SingleCommandInfo{
				Ioa:   cmd.Ioa,
				Value: cmd.Value,
				Qoc:   cmd.Qoc,
			})
			if !cmd.Qoc.InSelect {
				var setValue int
				if cmd.Value {
					setValue = 1
				} else {
					setValue = 0
				}
				ISMScript.SetDeviceData(v.BandData, setValue)
			}
			_ = pack.SendReplyMirror(conn, asdu.ActivationCon)
			break
		}
	}
	return nil
}

func IEC104InterfaceStart() {

	for {
		var is_starting = 0
		type iec104Data struct {
			Port int `json:"port"`
			Addr int `json:"addr"`
		}
		for {
			if is_starting == 1 {
				IEC104Wg.Wait()
			}
			IEC104InterfaceCloseChan()
			GIEC104InterfaceChan = make(chan bool)
			var getData []models.SystemDataInterface
			err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 4").Select("*").Find(&getData).Error
			if err != nil {
				time.Sleep(10 * time.Second)
				continue
			}
			for _, v := range getData {
				if v.InterfaceStatus == 0 {
					continue
				}
				var SIEC104InterfaceData IEC104InterfaceData
				var edata iec104Data
				SIEC104InterfaceData.ProjectUuid = v.InterfaceUuid
				SIEC104InterfaceData.InterfaceName = v.InterfaceName
				SIEC104InterfaceData.InterfaceType = v.InterfaceType
				SIEC104InterfaceData.InterfaceDataUuid = v.InterfaceDataUuid
				SIEC104InterfaceData.InterfaceContent = v.InterfaceContent

				err := json.Unmarshal([]byte(SIEC104InterfaceData.InterfaceContent), &edata)
				if err != nil {
					continue
				}
				SIEC104InterfaceData.InterfacePort = edata.Port
				SIEC104InterfaceData.InterfaceAddr = edata.Addr
				d := strings.Split(v.InterfaceDataUuid, ",")
				if len(d) >= 1 {
					var getDataContent []models.IEC104DataPushModel
					err := models.Db.Model(&models.IEC104DataPushModel{}).Where("muid = ?", d[0]).Select("*").Find(&getDataContent).Error
					if err != nil {
						continue
					}
					for _, v := range getDataContent {
						if v.DataCategory == 1 {
							var tYaoXin DataPointStu
							tYaoXin.BandData = v.BandData
							tYaoXin.Point = v.DataPoint
							tYaoXin.Type = v.DataCategoryYaoKongType
							SIEC104InterfaceData.YaoXin = append(SIEC104InterfaceData.YaoXin, tYaoXin)
						} else if v.DataCategory == 2 {
							var tYaoCe DataPointStu
							tYaoCe.BandData = v.BandData
							tYaoCe.Type = v.DataCategoryYaoTiaoType
							tYaoCe.Point = v.DataPoint
							SIEC104InterfaceData.YaoCe = append(SIEC104InterfaceData.YaoCe, tYaoCe)
						} else if v.DataCategory == 3 {
							var tEnergy DataPointStu
							tEnergy.BandData = v.BandData
							if v.Type == "8" {
								tEnergy.Type = 1
							} else if v.Type == "10" {
								tEnergy.Type = 2
							} else {
								continue
							}
							tEnergy.Point = v.DataPoint
							SIEC104InterfaceData.Energy = append(SIEC104InterfaceData.Energy, tEnergy)
						} else if v.DataCategory == 4 {
							var tYaoKong DataPointStu
							tYaoKong.BandData = v.BandData
							tYaoKong.Type = v.DataCategoryYaoKongType
							tYaoKong.Point = v.DataPoint
							SIEC104InterfaceData.YaoKong = append(SIEC104InterfaceData.YaoKong, tYaoKong)
						} else if v.DataCategory == 5 {
							var tYaoTiao DataPointStu
							tYaoTiao.BandData = v.BandData
							tYaoTiao.Point = v.DataPoint
							tYaoTiao.Type = v.DataCategoryYaoTiaoType
							SIEC104InterfaceData.YaoTiao = append(SIEC104InterfaceData.YaoTiao, tYaoTiao)
						} else {
							continue
						}
					}
				} else {
					continue
				}
				dIec104 := &IEC104Ctl{waitGroup: &IEC104Wg, InterfaceData: SIEC104InterfaceData}
				go dIec104.StartIEC104Server()
				IEC104Wg.Add(1)
				time.Sleep(100 * time.Millisecond)
				go dIec104.IEC104DataMutation()
				IEC104Wg.Add(1)
				is_starting = 1
			}
			if is_starting == 0 {
				time.Sleep(10 * time.Second)
				continue
			}
		}
	}

}
