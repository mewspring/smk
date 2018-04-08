package smk

import (
	"github.com/lunixbochs/struc"
	"github.com/pkg/errors"
)

// parseFileHeader parses the file header of the Smacker file.
func (f *File) parseFileHeader() error {
	if err := struc.Unpack(f.r, &f.FileHeader); err != nil {
		return errors.WithStack(err)
	}
	// Verify Smacker signature.
	switch f.Signature {
	case "SMK2", "SMK4":
		// Smacker version 2 and 4, respectively.
	default:
		return errors.Errorf(`invalid Smacker signature; got %q, want "SMK2" or "SMK4"`, f.Signature)
	}
	return nil
}

// FileHeader is a general file description header.
type FileHeader struct {
	// File signature; "SMK2" or "SMK4".
	Signature string `struc:"[4]byte"`
	// Frame width.
	Width int `struc:"uint32,little"`
	// Frame height.
	Height int `struc:"uint32,little"`
	// Number of frames. Excluding "ring" frames.
	NFrames int `struc:"uint32,little"`
	// Frame rate.
	FrameRate FrameRate `struc:"int32,little"`
	// Video flags.
	Flags Flag
	// Size of the largest unpacked audio data buffer in bytes; one per channel.
	AudioSize []int `struc:"[7]uint32,little"`
	// Total size in bytes of Huffman trees stored in file.
	TreesSize int `struc:"uint32,little"`
	// Allocation size for the mono blocks maps Huffman tree.
	MMapSize int `struc:"uint32,little"`
	// Allocation size for the mono blocks colours Huffman tree.
	MClrSize int `struc:"uint32,little"`
	// Allocation size for the full blocks Huffman tree.
	FullSize int `struc:"uint32,little"`
	// Allocation size for the block type descriptors Huffman tree.
	TypeSize int `struc:"uint32,little"`
	// Frequency and format information for each sound track; one per channel.
	// TODO: Verify if little or big endian encoding.
	AudioRate []AudioRate `struc:"[7]uint32,little"`
	// Unused.
	_ uint32
	// Frame size in number of bytes. Bit 0 determines if the frame is a key
	// frame. The purpose of bit 1 is unknown. Note, to get the proper length,
	// clear bit 0 and 1.
	FrameSizes []int `struc:"[]uint32,little,sizefrom=NFrames"`
	// Frame types.
	FrameTypes []FrameType `struc:"sizefrom=NFrames"`
}

// FrameRate specifies the number of frames per second.
//
// The frame rate can be determined as follows:
//
//    if (FrameRate > 0) {
//       fps = 1000 / FrameRate
//    } else if (FrameRate < 0) {
//       fps = 100000 / (-FrameRate)
//    } else {
//       fps = 10
//    }
type FrameRate int32

// FPS returns the frame rate in frames per second.
func (rate FrameRate) FPS() float64 {
	switch {
	case rate > 0:
		return 1000 / float64(rate)
	case rate < 0:
		return 100000 / float64(-rate)
	default:
		return 10
	}
}

// Flag specifies a set of video flags.
type Flag uint32

// Video flags.
const ()

// AudioRate describes the frequency and format information of a sound track.
//
// The 32 constituent bits have the following meaning:
//
//    bit 31 - data is compressed
//    bit 30 - indicates that audio data is present for this track
//    bit 29 - 1 = 16-bit audio; 0 = 8-bit audio
//    bit 28 - 1 = stereo audio; 0 = mono audio
//    bits 27-26 - if both set to zero - use v2 sound decompression
//    bits 25-24 - unused
//    bits 23-0 - audio sample rate
type AudioRate uint32

// Audio rates.
const ()

// FrameType describes the contents of the corresponding frame.
//
// The 8 bits have the following meaning when set:
//
//    7 - frame contains audio data corresponding to track 6
//    6 - frame contains audio data corresponding to track 5
//    5 - frame contains audio data corresponding to track 4
//    4 - frame contains audio data corresponding to track 3
//    3 - frame contains audio data corresponding to track 2
//    2 - frame contains audio data corresponding to track 1
//    1 - frame contains audio data corresponding to track 0
//    0 - frame contains a palette record
type FrameType uint8

// Frame types.
const ()
