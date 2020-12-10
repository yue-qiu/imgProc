package tool

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ImgLoader struct {
	filename string
	format   string
	img      image.Image
	matrix   [][][]uint8
}

// 将文件解码为图片对象
func (il *ImgLoader)decode(filePath string) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	il.filename = filepath.Base(filePath)
	il.filename = strings.Split(il.filename, ".")[0]
	var img image.Image
	img, il.format, err = image.Decode(file)
	il.img = convertToNRGBA(img)
	if err != nil {
		return
	}

	return
}

// 构建 ImgLoader 结构体
func NewImgLoader(filePath string) (il ImgLoader, err error) {
	err = il.decode(filePath)
	if err != nil {
		return
	}

	height, width := il.img.Bounds().Dy(), il.img.Bounds().Dx()

	matrix := make([][][]uint8, height)
	for hi := range matrix {
		matrix[hi] = make([][]uint8, width)
		for wi := range matrix[hi] {
			matrix[hi][wi] = make([]uint8, 4)
		    r, g, b ,a := il.img.At(wi, hi).RGBA()
		    matrix[hi][wi][0] = uint8(r)
			matrix[hi][wi][1] = uint8(g)
		    matrix[hi][wi][2] = uint8(b)
		    matrix[hi][wi][3] = uint8(a)
		}
	}

	il.matrix = matrix

	return
}

func (il *ImgLoader)GetImg() image.Image {
	return il.img
}

func (il *ImgLoader)GetMatrix() [][][]uint8 {
	matrix := NewRGBAMatrix(il.GetMY(), il.GetMY())

	copy(matrix, il.matrix)
	//for hi := range matrix {
	//	for wi := range matrix[hi] {
	//		_ = copy(matrix[hi][wi], il.matrix[hi][wi])
	//	}
	//}

	return matrix
}

func (il *ImgLoader)GetMX() int {
	return len(il.matrix[0])
}

func (il *ImgLoader)GetMY() int {
	return len(il.matrix)
}

func (il *ImgLoader)GetFileName() string {
	return il.filename
}

func (il *ImgLoader)GetFormat() string {
	return il.format
}

func SaveAsPng(filename string, matrix [][][]uint8) (err error) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return errors.New("not init yet")
	}
	height, width := len(matrix), len(matrix[0])
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			rgba.Set(j, i, color.RGBA{
				R: matrix[i][j][0],
				G: matrix[i][j][1],
				B: matrix[i][j][2],
				A: matrix[i][j][3]})
		}
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	png.Encode(outfile, rgba)

	return nil
}

// quality 范围 [0, 100]
func SaveAsJpeg(filename string, quality int, matrix [][][]uint8) (err error) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return errors.New("not init yet")
	}
	height, width := len(matrix), len(matrix[0])
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			rgba.Set(j, i, color.RGBA{
				R: matrix[i][j][0],
				G: matrix[i][j][1],
				B: matrix[i][j][2],
				A: matrix[i][j][3]})
		}
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	jpeg.Encode(outfile, rgba, &jpeg.Options{Quality: quality})

	return nil
}

func NewRGBAMatrix(height, width int) [][][]uint8 {
	matrix := make([][][]uint8, height)
	for hi := range matrix {
		matrix[hi] = make([][]uint8, width)
		for wi := range matrix[hi] {
			matrix[hi][wi] = make([]uint8, 4)
		}
	}

	return matrix
}

//convert image to NRGBA
func convertToNRGBA(src image.Image) *image.NRGBA {
	srcBounds := src.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

