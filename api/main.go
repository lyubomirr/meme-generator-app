package main

import "github.com/lyubomirr/meme-generator-app/web"

func main() {
	web.Serve(":8080")
	//t := entities.Template{
	//	Name:          "pesho",
	//	FilePath:      "asd.png",
	//	TextPositions: []entities.TemplateTextPosition{{
	//		TopOffset:  10,
	//		LeftOffset: 20,
	//	}},
	//	MimeType: "alabala",
	//}
	//s := services.NewMemeService(persistence.NewUnitOfWorkFactory())
	//r, err := s.CreateTemplate(context.Background(), []byte{1,2,3,4}, t)
	//fmt.Println(r)
	//
	//err = s.DeleteTemplate(context.Background(), r.ID)
	//fmt.Println(err)
	//_, err := s.Create(context.Background(), []byte {2,3,4,5}, entities.Meme{})
	//fmt.Print(err)
	//m, _ := r.Get(1)
	//m.Comments = append(m.Comments, entities.Comments{
	//	Author:  entities.User{ ID: 1},
	//	Content: "nov komentar",
	//})
	//m, _ = r.Update(m)
	//f,_ := ioutil.ReadFile("download.png")
	//img, _, _ := image.Decode(bytes.NewBuffer(f))
	//a := img.Bounds()
	//width := float64(a.Dx())
	//height := float64(a.Dy())
	//
	//dc := gg.NewContextForImage(img)
	//if err := dc.LoadFontFace("./assets/impact.ttf", 48); err != nil {
	//	panic(err)
	//}
	//dc.SetRGB(0, 0, 0)
	//s := "ONE DOES NOT SIMPLY WALK INTO MORDOR MY FRIENDO"
	//n := 6 // "stroke" size
	//for dy := -n; dy <= n; dy++ {
	//	for dx := -n; dx <= n; dx++ {
	//		if dx*dx+dy*dy >= n*n {
	//			// give it rounded corners
	//			continue
	//		}
	//		x := width/2 + float64(dx)
	//		y := height/5 + float64(dy)
	//		//dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
	//		dc.DrawStringWrapped(s, x, y, 0.5, 0.5, width * 0.8, 1.5, gg.AlignCenter)
	//	}
	//}
	//dc.SetRGB(1, 1, 1)
	//dc.DrawStringWrapped(s, width/2, height/5, 0.5, 0.5, width * 0.8, 1.5, gg.AlignCenter)
	//dc.SavePNG("out.png")

}