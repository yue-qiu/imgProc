package main

import (
	"fmt"
	"github.com/yue-qiu/imgProc/tool"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type App struct {
	PicList 	[]string
	ActList 	[]string
	Processor 	tool.ImgProcessor
}

func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	app.Run()
}

func NewApp() (App, error) {
	picList, err := loadPics()
	if err != nil {
		return App{}, err
	}

	return App{PicList: picList, ActList: tool.GetActionList()}, err
}

func (app App)Run() {
	for true {
		app.listActions()
		action := strings.ToLower(app.getChoice())
		if isValid := app.checkActChoice(action); !isValid {
			fmt.Println("Error, invalid choice!")
			continue
		}

		switch action {
		case strings.ToLower("Sunset"):
			app.dealWithSunSet()
		case strings.ToLower("Rotate"):
			app.dealWithRotate()
		case strings.ToLower("NegativeFilm"):
			app.dealWithNegativeFilm()
		case strings.ToLower("ToGray"):
			app.dealWithToGray()
		case strings.ToLower("Base64Encode"):
			app.dealWithBase64Enc()
		case strings.ToLower("Base64Decode"):
			app.dealWithBase64Dec()
		case strings.ToLower("FingerPrint"):
			app.dealWithFingerPrint()
		case strings.ToLower("Fusion"):
			app.dealWithFusion()
		case strings.ToLower("Resize"):
			app.dealWithResize()
		case strings.ToLower("AdjustBrightness"):
			app.dealWithAdjBrit()
		}
	}
}

func (app App)listRaw() {
	fmt.Println("There are your raw pictures:")
	col := 0
	for _, v := range app.PicList {
		format := path.Ext(v)
		if format == ".jpg" || format == ".png" {
			fmt.Printf("%s\t", v)
			col++
			if col == 2 {
				col = 0
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func (app App)listActions() {
	for i := 0; i < 40; i++ {
		fmt.Print("*")
	}
	fmt.Println()
	fmt.Println("There are your actions list:")
	col := 0
	for _, v := range app.ActList {
		fmt.Printf("\t%s\t", v)
		col++
		if col == 2 {
			col = 0
			fmt.Println()
		}
	}
	fmt.Println()
}

func (app App)getChoice() string {
	fmt.Print("make your choice: ")
	var choice string
	_, _ = fmt.Scan(&choice)

	if strings.ToLower(choice) == "q" {
		fmt.Println("bye~")
		os.Exit(0)
	}

	return choice
}

func (app App)checkRawChoice(choice string) bool {
	var isValid bool

	for _, v := range app.PicList {
		if strings.ToLower(v) == choice {
			isValid = true
			break
		}
	}

	return isValid
}

func (app App)checkActChoice(choice string) bool {
	var isValid bool

	for _, v := range app.ActList {
		if strings.ToLower(v) == choice {
			isValid = true
			break
		}
	}

	return isValid
}

func (app App)dealWithSunSet()  {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path.Join(tool.RESULT, "Sunset-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.SunsetEffect(&il))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithNegativeFilm() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path.Join(tool.RESULT, "NegativeFilm-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.NegativeFilmEffect(&il))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithRotate() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path.Join(tool.RESULT, "Rotate-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.Rotate(&il))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithToGray() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path.Join(tool.RESULT, "Gray-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.RGB2Gray(&il))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithFingerPrint() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("FingerPrint: ", app.Processor.GetFingerPrint(&il))
	fmt.Println()
}
func (app App)dealWithBase64Enc() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = app.Processor.Base64Encode(&il)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithBase64Dec() {
	files, err := ioutil.ReadDir(tool.RAW)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("choice a base64-txt:")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			fmt.Printf("%s\t", file.Name())
		}
	}
	fmt.Println()

	filename := app.getChoice()
	var isValid bool
	for _, file := range files {
		if filename == file.Name() {
			isValid = true
			break
		}
	}

	if !isValid {
		fmt.Println("inValid input!")
		return
	}

	err = app.Processor.Base642Img(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithResize() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("input height & width, separated by space: ")
	var height, width int
	_, err = fmt.Scan(&height, &width)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if height <= 0 || width <= 0 {
		fmt.Println("height or width has to greater than 0!")
		return
	}

	savePath := path.Join(tool.RESULT, "Resize-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.Resize(&il, height, width))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithFusion() {
	app.listRaw()
	fmt.Printf("make your choices(seperated by space): " )
	var filename1, filename2 string
	_, _ = fmt.Scan(&filename1, &filename2)
	if isValid1, isValid2 := app.checkRawChoice(filename1), app.checkRawChoice(filename2); !isValid1 || !isValid2 {
		fmt.Println("inValid input!")
		return
	}

	il1, err := tool.NewImgLoader(path.Join(tool.RAW, filename1))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	il2, err := tool.NewImgLoader(path.Join(tool.RAW, filename2))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	savePath := path.Join(tool.RESULT, "fusion-"+il1.GetFileName()+il2.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.ImageFusion(&il1, &il2))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func (app App)dealWithAdjBrit() {
	app.listRaw()
	filename := app.getChoice()
	if isValid := app.checkRawChoice(filename); !isValid {
		fmt.Println("inValid input!")
		return
	}

	il, err := tool.NewImgLoader(path.Join(tool.RAW, filename))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("input brightness rate(Negative numbers will be treated as 0): ")
	var light float64
	_, err = fmt.Scan(&light)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if light < 0 {
		light = 0
	}

	savePath := path.Join(tool.RESULT, "AdjBrit-"+il.GetFileName()+".png")
	err = tool.SaveAsPng( savePath, app.Processor.AdjustBrightness(&il, light))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Succeed, enjoy it")
	fmt.Println()
}

func loadPics() ([]string, error) {
	dir ,err := ioutil.ReadDir(tool.RAW)
	if err != nil {
		return []string{}, nil
	}

	picList := make([]string, 0)

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		picList = append(picList, file.Name())
	}

	return picList, nil
}
