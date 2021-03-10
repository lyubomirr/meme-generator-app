import { fabric } from "fabric";

const setCanvasBackgroundImage = (canvas, imageUrl) => {
    fabric.Image.fromURL(imageUrl, img => {
        const imgRatio = img.width/img.height;
        canvas.setHeight(canvas.width/imgRatio)

        canvas.setBackgroundImage(img, canvas.renderAll.bind(canvas), {
            scaleX: canvas.width / img.width,
            scaleY: canvas.height / img.height
         })
    })
}

export { setCanvasBackgroundImage }