package network

type (
	// 加密器定义
	Encrypter interface {
		// 加密
		Encrypt([]byte) []byte

		// 解密
		Decrypt([]byte) []byte
	}

	// 加密函数签名
	EncryptionFunc func(uint8, []uint8) []uint8
	// 解密函数签名
	DecryptionFunc func(uint8, []uint8) []uint8

	// 二维空间置换算法
	SETer struct {
		key           string           // 密钥
		commonDivisor int              // 默认公约数，在最大公因数为1和二者本身的时候为此值
		encryptList   []EncryptionFunc // 加密函数组
		decryptList   []DecryptionFunc // 解密函数组
	}
)

const (
	// 默认最大公约数
	DefaultCommonDivisor = 5
	// 默认加密钥匙
	DefaultEncryptionKey = "jarvis"
)

var ()

// 默认加密器
func DefaultEncrypter() Encrypter {
	return NewSETer(DefaultEncryptionKey, DefaultCommonDivisor)
}

// 新建加密器
// key : 密钥
// commonDivisor : 默认公约数，在最大公因数为1和二者本身的时候为此值
func NewSETer(key string, commonDivisor int) Encrypter {
	s := &SETer{
		key:           key,
		commonDivisor: commonDivisor,
	}
	s.encryptList = []EncryptionFunc{
		s.makeUp,
		s.blurValue,
		s.leftMove,
		s.bitOperation,
		s.spaceExchange,
	}
	s.decryptList = []DecryptionFunc{
		s.revMakeUp,
		s.revBlurValue,
		s.revLeftMove,
		s.revBitOperation,
		s.revSpaceExchange,
	}

	return s
}

// 加密
func (s *SETer) Encrypt(v []byte) []byte {
	keys := []uint8(s.key)
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		handler := s.encryptList[i%len(s.encryptList)]

		v = handler(key, v)
	}
	return v
}

// 解密
func (s *SETer) Decrypt(v []byte) []byte {
	keys := []uint8(s.key)
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]
		handler := s.decryptList[i%len(s.decryptList)]

		v = handler(key, v)
	}
	return v
}

// 补齐位📔数，补 0
func (s *SETer) makeUp(n uint8, v []uint8) []uint8 {
	remainder := int(n) - (len(v) % int(n))

	fix := make([]uint8, remainder)

	v = append(v, fix...)

	return v
}

// 去掉补齐位📔数，去掉补位的 0
func (s *SETer) revMakeUp(n uint8, v []uint8) []uint8 {
	remainder := int(n) - (len(v) % int(n))

	a := v
	for i := 0; i < remainder; i++ {
		if len(v)-1-i >= 0 && v[len(v)-1-i] == 0 {
			a = v[:len(v)-1-i]
		}
	}

	return a
}

// 数值模糊
func (s *SETer) blurValue(n uint8, v []uint8) []uint8 {
	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, uint8((int(v[i])+int(n))%256))
	}
	return a
}

// 清晰数值
func (s *SETer) revBlurValue(n uint8, v []uint8) []uint8 {
	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, uint8((int(v[i])+256-int(n))%256))
	}
	return a
}

// 左移数值
func (s *SETer) leftMove(n uint8, v []uint8) []uint8 {
	// 最多左移7位，8位就是原来的值
	n = n % 8

	f := func(i uint8, v uint8) uint8 {
		for i > 0 {
			a := v
			v = v << 1
			if a >= 128 {
				v += 1
			}
			i--
		}
		return v
	}

	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, f(n, v[i]))
	}

	return a
}

// 逆转左移数值
func (s *SETer) revLeftMove(n uint8, v []uint8) []uint8 {
	// 最多左移7位，8位就是原来的值
	n = n % 8

	f := func(i uint8, v uint8) uint8 {
		for i > 0 {
			a := v
			v = v >> 1
			if a%2 != 0 {
				v += 128
			}
			i--
		}
		return v
	}

	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, f(n, v[i]))
	}

	return a
}

// 位运算 是否与 255 非运算
func (s *SETer) bitOperation(n uint8, v []uint8) []uint8 {
	if n%2 == 0 {
		return v
	}

	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, v[i]^255)
	}

	return a
}

// 逆转位运算 是否与 255 非运算
func (s *SETer) revBitOperation(n uint8, v []uint8) []uint8 {
	if n%2 == 0 {
		return v
	}

	a := make([]uint8, 0)
	for i := 0; i < len(v); i++ {
		a = append(a, v[i]^255)
	}

	return a
}

// 空间交换
func (s *SETer) spaceExchange(n uint8, v []uint8) []uint8 {
	// 取得最大公约数
	divisor := maxCommonDivisor(int(n), len(v))
	// 如果最大公约数为 1 ， 默认为 5
	if divisor == 1 {
		divisor = s.commonDivisor
	}

	// 判断 v 的长度取余最大公约数是否为0
	gap := len(v) % divisor

	// 数据的前 gap 位不做置换
	a := v[:gap]

	// 从 gap 位开始进行空间置换
	b := v[gap:]

	// 将 b 以 divisor 切割成二维数组
	c := make([][]uint8, 0)
	for i := 0; i < len(b)/divisor; i++ {
		c = append(c, b[i*divisor:i*divisor+divisor])
	}

	// 交换
	for i := 0; i < len(c)/2; i++ {
		c = exchangeSeat(i, c)
	}

	d := make([]uint8, 0)
	for _, list := range c {
		d = append(d, list...)
	}

	return append(a, d...)
}

// 逆转空间交换
func (s *SETer) revSpaceExchange(n uint8, v []uint8) []uint8 {
	// 取得最大公约数
	divisor := maxCommonDivisor(int(n), len(v))
	// 如果最大公约数为 1 ， 默认为 5
	if divisor == 1 {
		divisor = s.commonDivisor
	}

	// 判断 v 的长度取余最大公约数是否为0
	gap := len(v) % divisor

	// 数据的前 gap 位不做置换
	a := v[:gap]

	// 从 gap 位开始进行空间置换
	b := v[gap:]

	// 将 b 以 divisor 切割成二维数组
	c := make([][]uint8, 0)
	for i := 0; i < len(b)/divisor; i++ {
		c = append(c, b[i*divisor:i*divisor+divisor])
	}

	// 交换
	for i := 0; i < len(c)/2; i++ {
		c = exchangeSeat(i, c)
	}

	d := make([]uint8, 0)
	for _, list := range c {
		d = append(d, list...)
	}

	return append(a, d...)
}

// 最大公约数
func maxCommonDivisor(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	gap := a % b
	smaller := a
	if b < a {
		smaller = b
	}

	return maxCommonDivisor(smaller, gap)
}

// 交换二维数组,交换过程中逆转子数组
func exchangeSeat(i int, v [][]uint8) [][]uint8 {
	if i%2 != 0 {
		return v
	}

	front := reverse(v[i])
	back := reverse(v[len(v)-1-i])

	v[i] = back
	v[len(v)-1-i] = front
	return v
}

// 逆转数组
func reverse(v []uint8) []uint8 {
	a := make([]uint8, 0)
	for i := len(v) - 1; i >= 0; i-- {
		a = append(a, v[i])
	}
	return a
}
