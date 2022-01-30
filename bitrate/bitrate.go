package bitrate

import (
	"io"
	"math"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
)

func hann(L int) []float64 {
	r := make([]float64, L)

	if L == 1 {
		r[0] = 1
	} else {
		N := L - 1
		coef := 2 * math.Pi / float64(N)
		for n := 0; n <= N; n++ {
			r[n] = 0.5 * (1 - math.Cos(coef*float64(n)))
		}
	}

	return r
}

func rfft(samples []float64) []float64 {
	fftReal := fft.FFTReal(samples)
	rfft := make([]float64, len(fftReal))

	rfft[0] = real(fftReal[0])

	for i, j := 1, 1; j < len(fftReal)-1; i, j = i+1, j+2 {
		rfft[j] = real(fftReal[i])
		rfft[j+1] = imag(fftReal[i])
	}

	return rfft
}

func findCutoff(a []float64, dx float64, diff float64) int {
	l := len(a)

	for i := int(dx / 4); i < l; i++ {
		m := l - i - int(dx)
		n := l - i
		x := a[m] - a[n]
		if x > diff {
			return l - i - int(dx)
		}
	}
	return l
}

func readWav(wavPath string) (int, []float64, error) {
	var freq int
	var samples []float64

	file, err := os.Open(wavPath)
	if err != nil {
		return freq, samples, err
	}
	reader := wav.NewReader(file)
	defer file.Close()

	format, err := reader.Format()
	if err != nil {
		return freq, samples, err
	}

	freq = int(format.SampleRate)

	for {
		s, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range s {
			samples = append(samples, float64(reader.IntValue(sample, 0)))
		}
	}

	return freq, samples, nil
}

func CheckBitrate(wavPath string) (string, error) {
	var bitrate string

	freq, samples, err := readWav(wavPath)
	if err != nil {
		return bitrate, nil
	}

	seconds := len(samples) / freq
	if seconds > 30 {
		seconds = 30
	}

	spectrum := Fill(0.0, freq)
	window := hann(freq)

	for i := 0; i < seconds-1; i++ {
		audioSecond := Mul(window, samples[(i*freq):(i+1)*freq])
		secondFFT := Abs(rfft(audioSecond))
		spectrum = Add(spectrum, secondFFT)
	}

	spectrum = Divide(spectrum, float64(seconds))
	spectrum = Normalize(spectrum)
	window = Fill(1.0/(float64(freq)/100.0), int(float64(freq)/100.0))
	spectrum = CrossCorrelation(spectrum, window)
	dx := float64(freq / 50.0)
	diff := 1.25
	cutoff := float64((findCutoff(spectrum, dx, diff) + int(dx)/4)) / 2000.0

	kbps := []float64{11.0, 16.0, 19.0, 20.0, 22.0, 24.0}
	closest := ClosestValue(cutoff, kbps)

	switch closest {
	case 11.0:
		bitrate = "64 kbps"
	case 16.0:
		bitrate = "128 kbps"
	case 19.0:
		bitrate = "192 kbps"
	case 20.0:
		bitrate = "320 kbps"
	case 22.0:
		bitrate = "500 kbps"
	case 24.0:
		bitrate = "lossless"
	}

	return bitrate, nil
}
