package unpack

import (
	"image"
	"runtime"
	"sync"
)

func Rotate90(img *image.NRGBA) *image.NRGBA {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	rowSize := width * 4
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	parallel(0, height, func(ys <-chan int) {
		for dstY := range ys {
			i := dstY * dst.Stride
			srcX := height - dstY - 1
			scan(img, dst.Pix[i:i+rowSize], srcX, 0, srcX+1, height)
		}
	})
	return dst
}

func scan(srcImg *image.NRGBA, dstPx []uint8, x1, y1, x2, y2 int) {
	size := (x2 - x1) * 4
	srcStride := y1*srcImg.Stride + x1*4
	dstStride := 0
	if size == 4 { // fast swap
		for y := y1; y < y2; y++ {
			dstPixels := dstPx[dstStride : dstStride+4 : dstStride+4]
			srcPixels := srcImg.Pix[srcStride : srcStride+4 : srcStride+4]
			dstPixels[0] = srcPixels[0]
			dstPixels[1] = srcPixels[1]
			dstPixels[2] = srcPixels[2]
			dstPixels[3] = srcPixels[3]
			srcStride += srcImg.Stride
			dstStride += size
		}
	} else {
		for y := y1; y < y2; y++ {
			copy(dstPx[dstStride:dstStride+size], srcImg.Pix[srcStride:srcStride+size])
			srcStride += srcImg.Stride
			dstStride += size
		}
	}
}

func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}
