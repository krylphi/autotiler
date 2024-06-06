/*
 * MIT License
 *
 * Copyright (c) 2024 The autotiler authors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package unpack

import (
	"image"
	"runtime"
	"sync"
)

// rotateLeft90 rotates the given image 90 degrees counter-clockwise.
//
// Parameters:
// - img: The input image to be rotated.
//
// Returns:
// - A new image that is the result of rotating the input image 90 degrees counter-clockwise.
func rotateLeft90(img *image.NRGBA) *image.NRGBA {
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

// scan copies pixel data from a source image to a destination slice.
// It supports copying a rectangular region of the source image to the destination slice.
// The function uses optimized copying for a single pixel (size == 4) and falls back to a generic copy for larger regions.
//
// Parameters:
// - srcImg: The source image from which to copy pixel data.
// - dstPx: The destination slice to which to copy pixel data.
// - x1, y1: The top-left coordinates of the rectangular region to copy from the source image.
// - x2, y2: The bottom-right coordinates of the rectangular region to copy from the source image.
//
// Returns:
// - None
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

// parallel is a helper function that executes a given function in parallel across multiple goroutines.
// It distributes the work by sending indices from the start to stop (exclusive) to a channel,
// and each goroutine receives an index from the channel and executes the given function.
// The function waits for all goroutines to finish before returning.
//
// Parameters:
// - start: The starting index for the range of indices to be processed.
// - stop: The exclusive ending index for the range of indices to be processed.
// - fn: The function to be executed in parallel. It should accept a channel of integers as its parameter.
//
// Returns:
// - None
func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	// Determine the number of goroutines to use.
	// Use the minimum of the count and the number of available CPUs.
	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	// Create a buffered channel with a capacity equal to the count.
	c := make(chan int, count)

	// Fill the channel with indices from start to stop (exclusive).
	for i := start; i < stop; i++ {
		c <- i
	}

	// Close the channel to indicate that no more indices will be sent.
	close(c)

	var wg sync.WaitGroup

	// Start a goroutine for each processor.
	for i := 0; i < procs; i++ {
		wg.Add(1)

		// Each goroutine receives an index from the channel and executes the given function.
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}
