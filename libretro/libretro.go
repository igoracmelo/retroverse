package libretro

import (
	"errors"
)

const APIVersion = 1
const DeviceJoypad = 1

var (
	ErrShortBuffer = errors.New("short buffer")
	ErrUnknown     = errors.New("unknown")
)

// Interface for libretro.h
//
// Most comments were adapted from libretro.h
type API interface {
	SetEnvironmentCallback(func(cmd uint, data []byte) bool)
	SetVideoRefreshCallback(func(data []byte, width uint, height uint, pitch uint))
	SetAudioSampleCallback(func(left int16, right int16))
	SetAudioSampleBatchCallback(func(data []int16) uint)
	SetInputPollCallback(func())
	SetInputStateCallback(func(port, device, index, id uint) int16)

	// Library global initialization.
	Init()

	// Library global deinitialization.
	Deinit()

	// Must return APIVersion. Used to validate ABI compatibility when the API is revised.
	APIVersion() uint

	// SystemInfo Gets system info. Can be called at any time, even before Init().
	SystemInfo() SystemInfo

	// Gets information about system audio/video timings and geometry.
	// Can be called only after LoadGame() has successfully completed.
	// NOTE: The implementation of this function might not initialize every
	// variable if needed.
	// E.g. geom.aspect_ratio might not be initialized if core doesn't
	// desire a particular aspect ratio.
	SystemAVInfo() SystemAVInfo

	// Sets device to be used for player 'port'.
	// By default, DeviceJoypad is assumed to be plugged into all
	// available ports.
	// Setting a particular device type is not a guarantee that libretro cores
	// will only poll input based on that particular device type. It is only a
	// hint to the libretro core when a core cannot automatically detect the
	// appropriate input device type on its own. It is also relevant when a
	// core can change its behavior depending on device type.
	SetControllerPortDevice(port uint, device uint)

	// Resets the current game.
	Reset()

	// Runs the game for one video frame.
	// During Run(), InputPollCallback must be called at least once.
	//
	// If a frame is not rendered for reasons where a game "dropped" a frame,
	// this still counts as a frame, and Run() should explicitly dupe
	// a frame if GET_CAN_DUPE* returns true.
	// In this case, the VideoRefreshCallback can take a nil argument for data.
	Run()

	// Returns the amount of data the implementation requires to serialize
	// internal state (save states).
	// Between calls to LoadGame() and UnloadGame(), the
	// returned size is never allowed to be larger than a previous returned
	// value, to ensure that the frontend can allocate a save state buffer once.
	SerializeSize() uint

	// Serializes internal state. If len(data) is lower than SerializeSize(), it
	// returns ErrShortBuffer. Otherwise if the error is unknown it returns ErrUnknown.
	Serialize(data []byte) error

	// Unserialize internal state. If len(data) is lower than SerializeSize(), it
	// returns ErrShortBuffer. Otherwise if the error is unknown it returns ErrUnknown.
	Unserialize(data []byte) error

	SetCheat(index uint, enabled bool, code string)
	ResetCheat()

	// Loads a game and updates GameInfo.Path if GameInfo.NeedFullPath is true.
	LoadGame(*GameInfo) error

	UnloadGame()

	Region() uint

	MemoryData(id uint)
	MemorySize(id uint) uint
}

type SystemInfo struct {
	// Descriptive name of library. Should not contain any version numbers, etc.
	LibraryName string

	// Descriptive version of core.
	LibraryVersion string

	// A string listing probably content extensions the core will be able to load, separated with pipe.
	// I.e. "bin|rom|iso".  Typically used for a GUI to filter out extensions.
	ValidExtensions string

	// If true, LoadGame() is guaranteed to provide a valid pathname in GameInfo::Path.
	// ::Data and ::Size are both invalid.
	//
	// If false, ::Data and ::Size are guaranteed to be valid, but ::Path
	// might not be valid.
	//
	// This is typically set to true for libretro implementations that must
	// load from file.
	// Implementations should strive for setting this to false, as it allows
	// the frontend to perform patching, etc.
	NeedFullPath bool

	// If true, the frontend is not allowed to extract any archives before
	// loading the real content.
	// Necessary for certain libretro implementations that load games
	// from zipped archives.
	BlockExtract bool
}

type GameInfo struct {
	Path string
	Data []byte
	Size uint
	Meta string
}

type SystemAVInfo struct {
	Geometry GameGeometry
	Timing   SystemTiming
}

type GameGeometry struct {
	BaseWidth   uint
	BaseHeight  uint
	MaxWidth    uint
	MaxHeight   uint
	AspectRatio float32
}

type SystemTiming struct {
	FPS        float64
	SampleRate float64
}
