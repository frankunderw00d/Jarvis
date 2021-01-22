package network

import (
	"jarvis/util/transform"
)

type (
	// 装包者定义
	Packager interface {
		// 克隆一个同类装包者
		Clone() Packager

		// 打包
		Pack([]byte) []byte

		// 解包
		Unpack([]byte) [][]byte
	}

	// 装包者定义实现
	packager struct {
		buffer []byte // 缓存
	}
)

// 此常量组定义了 Packager 定义及实现中可能会发生的错误文本
const (
	DefaultHeadSymbol    = "#*HeadSymbol*#"                            // 默认装包者使用的数据帧头缀
	DefaultHeadSymbolLen = len(DefaultHeadSymbol)                      // 默认装包者使用的数据帧头缀长度
	DefaultRecordDataLen = 4                                           // 默认装包者记录实际数据长度的长度
	DefaultHeaderLen     = DefaultHeadSymbolLen + DefaultRecordDataLen // 默认装包者使用的数据帧头部长度
)

var ()

// 默认装包者
func DefaultPackager() Packager {
	return &packager{
		buffer: make([]byte, 0),
	}
}

// 克隆一个同类装包者
func (p *packager) Clone() Packager {
	return DefaultPackager()
}

// 打包
func (p *packager) Pack(data []byte) []byte {
	return append([]byte(DefaultHeadSymbol), append(transform.IntToBytes(len(data)), data...)...)
}

// 解包
func (p *packager) Unpack(data []byte) [][]byte {
	p.buffer = append(p.buffer, data...)
	length := len(p.buffer)

	comDatas := make([][]byte, 0)

	//	遍历缓存，不得短于可取出头部的长度
	i := 0
	for ; i < length-DefaultHeaderLen; i++ {
		if string(p.buffer[i:i+DefaultHeadSymbolLen]) != DefaultHeadSymbol { // 非头部
			continue
		}
		dataLength := transform.BytesToInt(p.buffer[i+DefaultHeadSymbolLen : i+DefaultHeaderLen]) // 取得数据长度
		if i+DefaultHeaderLen+dataLength > length {                                               // 如果缓存长度不足以获取完整数据，退出
			break
		}

		comDatas = append(comDatas, p.buffer[i:i+DefaultHeaderLen+dataLength]) // 取出数据
		i = i + DefaultHeaderLen + dataLength - 1                              // -1 是为了防止后续的 i++ 导致偏移量错误
	}

	if len(comDatas) != 0 && i != 0 {
		p.buffer = p.buffer[i:]
	}

	if len(comDatas) > 0 {
		finalDatas := make([][]byte, 0)
		for _, d := range comDatas {
			finalDatas = append(finalDatas, d[DefaultHeaderLen:])
		}
		comDatas = finalDatas
	}

	return comDatas
}
