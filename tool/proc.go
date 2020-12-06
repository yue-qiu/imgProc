package tool

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io/ioutil"
	"math"
	"path"
	"strings"
)

type ImgProcessor struct {

}

const RAW = "raw"
const RESULT = "result"

func NewImgProcessor() (ip ImgProcessor) {
	return ip
}

//input a image matrix as src , return a image matrix by sunseteffect process
func (ip *ImgProcessor)SunsetEffect(il *ImgLoader) [][][]uint8 {
	src := il.GetMatrix()
	height := il.GetDY()
	width := il.GetDX()
	imgMatrix := NewRGBAMatrix(height, width)
	copy(imgMatrix, src)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			imgMatrix[i][j][1] = uint8(float64(imgMatrix[i][j][1]) * 0.7)
			imgMatrix[i][j][2] = uint8(float64(imgMatrix[i][j][2]) * 0.7)
		}
	}

	return imgMatrix
}

// input a image as src , return a image matrix by negativefilmeffect process
func (ip *ImgProcessor)NegativeFilmEffect(il *ImgLoader) [][][]uint8 {
	src := il.GetMatrix()
	height := il.GetDY()
	width := il.GetDX()
	imgMatrix := NewRGBAMatrix(height,width)
	copy(imgMatrix, src)

	for i := 0; i < height; i++{
		for j := 0; j < width; j++{
			imgMatrix[i][j][0] = math.MaxUint8 - imgMatrix[i][j][0]
			imgMatrix[i][j][1] = math.MaxUint8 - imgMatrix[i][j][1]
			imgMatrix[i][j][2] = math.MaxUint8 - imgMatrix[i][j][2]
		}
	}

	return imgMatrix
}

func (ip *ImgProcessor)Rotate(il *ImgLoader) (imgMatrix [][][]uint8) {
	src := il.GetMatrix()
	height := il.GetDY()
	width := il.GetDX()
	imgMatrix = NewRGBAMatrix(width, height)

	for i := 0; i < width; i++{
		for j := 0; j < height; j++{
			imgMatrix[i][j] = src[j][i]
		}
	}

	return imgMatrix
}

// 调整图片亮度，light 最小值为 0
func (ip *ImgProcessor)AdjustBrightness(il *ImgLoader, light float64) (imgMatrix [][][]uint8) {
	src := il.GetMatrix()

	height := len(src)
	width := len(src[0])
	imgMatrix = NewRGBAMatrix(height, width)
	copy(imgMatrix, src)

	for i := 0; i < height; i++{
		for j := 0; j < width; j++{
			for c := 0; c < 3; c++ {
				color := float64(imgMatrix[i][j][c]) * light - 100
				if color < 0 {
					color = 0
				} else if color > 255 {
					color = 255
				}
				imgMatrix[i][j][c] = uint8(color)
			}
		}
	}

	return
}

// 双线性插值法
func (ip *ImgProcessor)Resize(il *ImgLoader, heigth, width int) (imgMatrix [][][]uint8) {
	matrix := il.GetMatrix()

	imgMatrix = NewRGBAMatrix(heigth, width)

	for n := 0; n < 4; n++ {
		for hi := range imgMatrix {
			for wi := range imgMatrix[hi] {
				srcY := (float64(hi) + 0.5) * (float64(il.GetDY()) / float64(heigth)) - 0.5
				srcX := (float64(wi) + 0.5) * (float64(il.GetDX()) / float64(width) ) - 0.5
				srcX0 := int(math.Floor(srcX))
				if srcX0 < 0 {
					srcX0 = 0
				}
				srcX1 := int(math.Min(float64(srcX0 + 1), float64(il.GetDX() - 1)))
				srcY0 := int(math.Floor(srcY))
				if srcY0 < 0 {
					srcY0 = 0
				}
				srcY1 := int(math.Min(float64(srcY0 + 1), float64(il.GetDY() - 1)))

				value0 := (float64(srcX1) - srcX) * float64(matrix[srcY0][srcX1][n]) + (srcX - float64(srcX0)) * float64(matrix[srcY0][srcX0][n])
				value1 := (float64(srcX1) - srcX) * float64(matrix[srcY1][srcX1][n]) + (srcX - float64(srcX0)) * float64(matrix[srcY1][srcX0][n])
				imgMatrix[hi][wi][n] = uint8((float64(srcY1) - srcY) * value1 + (srcY - float64(srcY0)) * value0)
			}
		}
	}

	return
}

// fuse two images(filepath) and the size of new image is as il1
func (ip *ImgProcessor)ImageFusion(il1 *ImgLoader, il2 *ImgLoader)(imgMatrix1 [][][]uint8) {
	imgMatrix1 = il1.GetMatrix()

	height := il1.GetDY()
	width := il1.GetDX()

	imgMatrix2 := ip.Resize(il2, height, width)

	for i := 0; i < height; i++{
		for j := 0; j < width; j++{
			imgMatrix1[i][j][0] = (imgMatrix1[i][j][0] >> 1) + (imgMatrix2[i][j][0] >> 1)
			imgMatrix1[i][j][1] = (imgMatrix1[i][j][1] >> 1) + (imgMatrix2[i][j][1] >> 1)
			imgMatrix1[i][j][2] = (imgMatrix1[i][j][2] >> 1) + (imgMatrix2[i][j][2] >> 1)
		}
	}

	return
}

func (ip *ImgProcessor)RGB2Gray(il *ImgLoader) [][][]uint8 {
	src := il.GetMatrix()
	height := il.GetDY()
	width := il.GetDX()
	imgMatrix := NewRGBAMatrix(height, width)
	copy(imgMatrix, src)

	for i := 0; i < height; i++{
		for j := 0;j < width; j++{
			// 平均灰度: avg1 := (imgMatrix[i][j][0] + imgMatrix[i][j][1] + imgMatrix[i][j][3]) / 3
			// 加权灰度
			avg := (uint16(imgMatrix[i][j][0]) * 30 + uint16(imgMatrix[i][j][1]) * 59 + uint16(imgMatrix[i][j][2]) * 11 + 50) / 100
			imgMatrix[i][j][0] = uint8(avg)
			imgMatrix[i][j][1] = uint8(avg)
			imgMatrix[i][j][2] = uint8(avg)
		}
	}
	return imgMatrix
}

func (ip *ImgProcessor)Base64Encode(il *ImgLoader) (err error) {
	var buf bytes.Buffer
	err = png.Encode(&buf, il.img)
	if err != nil {
		return
	}

	pngBytes := buf.Bytes()
	data := make([]byte, base64.StdEncoding.EncodedLen(len(pngBytes)))
	base64.StdEncoding.Encode(data, pngBytes)
	err = ioutil.WriteFile(path.Join(RESULT, "base64-"+il.GetFileName()+ ".txt"), data, 0666)
	if err != nil {
		return
	}
	return
}

func (ip *ImgProcessor)Base642Img(filename string) (err error) {
	encBytes, err := ioutil.ReadFile(path.Join(RAW, filename))
	if err != nil {
		return
	}
	data, err := base64.StdEncoding.DecodeString(string(encBytes))
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path.Join(RESULT, "base64Dec-"+strings.Split(filename, ".")[0]+".png"), data, 0666)
	return
}

func (ip *ImgProcessor)GetFingerPrint(il *ImgLoader) (fp string) {
	matrix := ip.Resize(il, 8, 9)

	height := len(matrix)
	width := len(matrix[0])

	// convert rgb to gray
	gray := make([]uint8, height * width)
	for hi := 0; hi < height; hi++ {
		for wi := 0; wi < width; wi++ {
			gValue := (matrix[hi][wi][0] * 30 + matrix[hi][wi][1] * 59 + matrix[hi][wi][2] * 11) / 100
			gray[9 * hi + wi] = gValue
		}
	}

	var buf bytes.Buffer
	for hi := 0; hi < height; hi++ {
		for wi := 0; wi < width-1; wi++ {
			diff := int16(gray[9 * hi + wi] - gray[9 * hi + wi + 1])
			if diff > 0 {
				buf.WriteByte('1')
			} else {
				buf.WriteByte('0')
			}
		}
	}

	return buf.String()
}

func GetActionList() []string {
	return []string{"Sunset", "NegativeFilm", "Rotate", "AdjustBrightness", "Resize", "Base64Dec", "ToGray",
		"Base64Enc", "Fusion", "FingerPrint"}
}

